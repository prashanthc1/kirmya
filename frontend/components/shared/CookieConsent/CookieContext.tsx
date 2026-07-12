"use client";

import React, { createContext, useContext, useState, useEffect, useCallback } from "react";
import { api } from "@/lib/api/client";
import { useAuth } from "@/lib/auth/auth-context";

export interface CookiePreferences {
  essential: boolean;
  functional: boolean;
  analytics: boolean;
  marketing: boolean;
  performance: boolean;
  personalization: boolean;
  ai_preferences: boolean;
  consent_version: string;
}

export interface CookieContextType {
  preferences: CookiePreferences;
  hasChoiceBeenMade: boolean;
  showModal: boolean;
  setShowModal: (show: boolean) => void;
  acceptAll: () => Promise<void>;
  rejectNonEssential: () => Promise<void>;
  saveCustom: (prefs: Partial<CookiePreferences>) => Promise<void>;
  anonymousId: string | null;
}

const DEFAULT_PREFERENCES: CookiePreferences = {
  essential: true,
  functional: false,
  analytics: false,
  marketing: false,
  performance: false,
  personalization: false,
  ai_preferences: false,
  consent_version: "1.0",
};

const CookieContext = createContext<CookieContextType | undefined>(undefined);

export function getCookie(name: string): string | null {
  if (typeof document === "undefined") return null;
  const match = document.cookie.match(new RegExp("(^| )" + name + "=([^;]+)"));
  return match ? decodeURIComponent(match[2]) : null;
}

export function setCookie(name: string, value: string, days: number) {
  if (typeof document === "undefined") return;
  const date = new Date();
  date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
  const expires = "; expires=" + date.toUTCString();
  const secure = window.location.protocol === "https:" ? "; Secure" : "";
  document.cookie = `${name}=${encodeURIComponent(value)}${expires}; path=/; SameSite=Lax${secure}`;
}

export function eraseCookie(name: string) {
  setCookie(name, "", -1);
}

// Generate anonymous UUIDv4
function generateUUID(): string {
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, function (c) {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

// Global script loader registry to prevent double-loading
const loadedScripts = new Set<string>();

export function loadConditionalScript(src: string, category: keyof CookiePreferences, preferences: CookiePreferences) {
  if (typeof document === "undefined") return;
  if (!preferences[category]) return;
  if (loadedScripts.has(src)) return;

  loadedScripts.add(src);
  const script = document.createElement("script");
  script.src = src;
  script.async = true;
  document.head.appendChild(script);
}

export function CookieProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const [preferences, setPreferences] = useState<CookiePreferences>(DEFAULT_PREFERENCES);
  const [hasChoiceBeenMade, setHasChoiceBeenMade] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [anonymousId, setAnonymousId] = useState<string | null>(null);

  // Initialize
  useEffect(() => {
    if (typeof window === "undefined") return;

    // Manage Anonymous ID
    let anonId = localStorage.getItem("kirmya_anonymous_id");
    if (!anonId) {
      anonId = generateUUID();
      localStorage.setItem("kirmya_anonymous_id", anonId);
    }
    setAnonymousId(anonId);

    // Read local cookie preferences
    const storedPrefsString = getCookie("kirmya_cookie_preferences") || localStorage.getItem("kirmya_cookie_preferences");
    if (storedPrefsString) {
      try {
        const parsed = JSON.parse(storedPrefsString) as CookiePreferences;
        // Check Consent Versioning
        if (parsed.consent_version === DEFAULT_PREFERENCES.consent_version) {
          setPreferences(parsed);
          setHasChoiceBeenMade(true);
        } else {
          // Consent version mismatch: require re-consent
          setHasChoiceBeenMade(false);
        }
      } catch (_) {
        setHasChoiceBeenMade(false);
      }
    }
  }, []);

  // Synchronize when user logs in (or whenever user changes)
  useEffect(() => {
    if (typeof window === "undefined") return;
    if (!user) return;

    // Fetch from backend
    api.get<CookiePreferences>("/cookies/preferences")
      .then((data) => {
        if (data && data.consent_version === DEFAULT_PREFERENCES.consent_version) {
          setPreferences(data);
          setHasChoiceBeenMade(true);
          setCookie("kirmya_cookie_preferences", JSON.stringify(data), 365);
          localStorage.setItem("kirmya_cookie_preferences", JSON.stringify(data));
        } else {
          // No settings on backend or version mismatch: save current preferences to backend
          const localPrefs = getCookie("kirmya_cookie_preferences") || localStorage.getItem("kirmya_cookie_preferences");
          const payload = localPrefs ? JSON.parse(localPrefs) : preferences;
          
          const anonId = localStorage.getItem("kirmya_anonymous_id");
          api.post("/cookies/preferences", {
            ...payload,
            anonymous_id: anonId || undefined,
          }).catch(() => {});
        }
      })
      .catch(() => {});
  }, [user]);

  const savePreferencesState = useCallback(async (newPrefs: CookiePreferences) => {
    setPreferences(newPrefs);
    setHasChoiceBeenMade(true);

    // 1. Save in Cookies
    setCookie("kirmya_cookie_preferences", JSON.stringify(newPrefs), 365);

    // 2. Save in Local Storage
    localStorage.setItem("kirmya_cookie_preferences", JSON.stringify(newPrefs));

    // 3. Save to database if logged in or anonymous
    const anonId = localStorage.getItem("kirmya_anonymous_id");
    const payload = {
      ...newPrefs,
      anonymous_id: anonId || undefined,
    };

    try {
      await api.post("/cookies/preferences", payload);
    } catch (_) {
      // Fail-safe: preferences are stored in cookies/localStorage even if backend fails
    }
  }, []);

  const acceptAll = useCallback(async () => {
    const allPrefs: CookiePreferences = {
      essential: true,
      functional: true,
      analytics: true,
      marketing: true,
      performance: true,
      personalization: true,
      ai_preferences: true,
      consent_version: DEFAULT_PREFERENCES.consent_version,
    };
    await savePreferencesState(allPrefs);
  }, [savePreferencesState]);

  const rejectNonEssential = useCallback(async () => {
    const essentialOnly: CookiePreferences = {
      ...DEFAULT_PREFERENCES,
      consent_version: DEFAULT_PREFERENCES.consent_version,
    };
    await savePreferencesState(essentialOnly);
  }, [savePreferencesState]);

  const saveCustom = useCallback(async (custom: Partial<CookiePreferences>) => {
    const combined: CookiePreferences = {
      ...preferences,
      ...custom,
      essential: true, // Always true
      consent_version: DEFAULT_PREFERENCES.consent_version,
    };
    await savePreferencesState(combined);
  }, [preferences, savePreferencesState]);

  return (
    <CookieContext.Provider
      value={{
        preferences,
        hasChoiceBeenMade,
        showModal,
        setShowModal,
        acceptAll,
        rejectNonEssential,
        saveCustom,
        anonymousId,
      }}
    >
      {children}
    </CookieContext.Provider>
  );
}

export function useCookieConsent() {
  const context = useContext(CookieContext);
  if (!context) {
    throw new Error("useCookieConsent must be used within a CookieProvider");
  }
  return context;
}
