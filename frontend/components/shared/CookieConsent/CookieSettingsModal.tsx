"use client";

import React, { useEffect, useRef, useState } from "react";
import { useCookieConsent } from "./CookieContext";

export default function CookieSettingsModal() {
  const { showModal, setShowModal, preferences, saveCustom, acceptAll, rejectNonEssential } = useCookieConsent();
  const [localPrefs, setLocalPrefs] = useState({
    functional: false,
    analytics: false,
    marketing: false,
    performance: false,
    personalization: false,
    ai_preferences: false,
  });

  const modalRef = useRef<HTMLDivElement>(null);
  const closeBtnRef = useRef<HTMLButtonElement>(null);

  // Sync state when modal opens
  useEffect(() => {
    if (showModal) {
      setLocalPrefs({
        functional: preferences.functional,
        analytics: preferences.analytics,
        marketing: preferences.marketing,
        performance: preferences.performance || false,
        personalization: preferences.personalization || false,
        ai_preferences: preferences.ai_preferences || false,
      });
      // Focus the modal for accessibility
      setTimeout(() => {
        closeBtnRef.current?.focus();
      }, 50);
    }
  }, [showModal, preferences]);

  // Trap Focus & Close on ESC
  useEffect(() => {
    if (!showModal) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        setShowModal(false);
        return;
      }

      if (e.key === "Tab" && modalRef.current) {
        const focusableElements = modalRef.current.querySelectorAll(
          'button, [href], input, select, textarea, [tabindex="0"]'
        );
        const firstElement = focusableElements[0] as HTMLElement;
        const lastElement = focusableElements[focusableElements.length - 1] as HTMLElement;

        if (e.shiftKey) {
          if (document.activeElement === firstElement) {
            lastElement.focus();
            e.preventDefault();
          }
        } else {
          if (document.activeElement === lastElement) {
            firstElement.focus();
            e.preventDefault();
          }
        }
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [showModal, setShowModal]);

  if (!showModal) return null;

  const handleToggle = (key: keyof typeof localPrefs) => {
    setLocalPrefs((prev) => ({ ...prev, [key]: !prev[key] }));
  };

  const handleSave = async () => {
    await saveCustom(localPrefs);
    setShowModal(false);
  };

  const handleAcceptAll = async () => {
    await acceptAll();
    setShowModal(false);
  };

  const handleRejectAll = async () => {
    await rejectNonEssential();
    setShowModal(false);
  };

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-labelledby="cookie-settings-title"
      style={{
        position: "fixed",
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        zIndex: 10000,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        padding: "20px",
      }}
    >
      {/* Dark Overlay */}
      <div
        onClick={() => setShowModal(false)}
        style={{
          position: "fixed",
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          background: "rgba(43, 38, 32, 0.4)",
          backdropFilter: "blur(4px)",
          zIndex: -1,
        }}
      />

      {/* Modal Container */}
      <div
        ref={modalRef}
        style={{
          background: "#FFFFFF",
          border: "1px solid #EFE7DC",
          borderRadius: "24px",
          maxWidth: "680px",
          width: "100%",
          maxHeight: "85vh",
          display: "flex",
          flexDirection: "column",
          boxShadow: "0 24px 60px rgba(43, 38, 32, 0.16)",
          overflow: "hidden",
          fontFamily: "'Public Sans', sans-serif",
          color: "#2B2620",
          animation: "scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1) forwards",
        }}
      >
        {/* Header */}
        <div
          style={{
            padding: "24px 28px",
            borderBottom: "1px solid #EFE7DC",
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <h2 id="cookie-settings-title" style={{ fontSize: "20px", fontWeight: 700, margin: 0 }}>
            Cookie Preference Manager
          </h2>
          <button
            ref={closeBtnRef}
            onClick={() => setShowModal(false)}
            aria-label="Close settings modal"
            style={{
              border: "none",
              background: "transparent",
              fontSize: "24px",
              cursor: "pointer",
              color: "#8A8175",
              padding: "4px",
              lineHeight: 1,
            }}
          >
            ×
          </button>
        </div>

        {/* Scrollable Categories List */}
        <div style={{ padding: "20px 28px", overflowY: "auto", flex: 1, display: "flex", flexDirection: "column", gap: "20px" }}>
          
          {/* ESSENTIAL */}
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: "16px" }}>
            <div style={{ flex: 1 }}>
              <h3 style={{ margin: "0 0 4px 0", fontSize: "15px", fontWeight: 600 }}>Essential Cookies</h3>
              <p style={{ margin: 0, fontSize: "13px", color: "#5B554C", lineHeight: "1.45" }}>
                Necessary for the site to function properly. Includes secure sessions, user sign-in tokens, CSRF validation, and security monitors. Cannot be disabled.
              </p>
            </div>
            <span style={{ fontSize: "13px", fontWeight: 700, color: "#4F7C6A", background: "#E8F0EC", padding: "4px 10px", borderRadius: "100px", alignSelf: "center" }}>
              Always Active
            </span>
          </div>

          {/* FUNCTIONAL */}
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: "16px", borderTop: "1px solid #FCFAF7", paddingTop: "16px" }}>
            <div style={{ flex: 1 }}>
              <h3 style={{ margin: "0 0 4px 0", fontSize: "15px", fontWeight: 600 }}>Functional Cookies</h3>
              <p style={{ margin: 0, fontSize: "13px", color: "#5B554C", lineHeight: "1.45" }}>
                Enables advanced settings like language preferences, custom theme selections (dark/light mode), dashboard configurations, and saved workspace filters.
              </p>
            </div>
            <label style={{ position: "relative", display: "inline-block", width: "44px", height: "24px", alignSelf: "center", cursor: "pointer" }}>
              <input
                type="checkbox"
                checked={localPrefs.functional}
                onChange={() => handleToggle("functional")}
                style={{ opacity: 0, width: 0, height: 0 }}
              />
              <span style={{
                position: "absolute",
                top: 0, left: 0, right: 0, bottom: 0,
                backgroundColor: localPrefs.functional ? "#C2683C" : "#E2D9CC",
                borderRadius: "34px",
                transition: "0.2s",
              }} />
              <span style={{
                position: "absolute",
                height: "18px", width: "18px",
                left: "3px", bottom: "3px",
                backgroundColor: "white",
                borderRadius: "50%",
                transition: "0.2s",
                transform: localPrefs.functional ? "translateX(20px)" : "none",
              }} />
            </label>
          </div>

          {/* ANALYTICS */}
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: "16px", borderTop: "1px solid #FCFAF7", paddingTop: "16px" }}>
            <div style={{ flex: 1 }}>
              <h3 style={{ margin: "0 0 4px 0", fontSize: "15px", fontWeight: 600 }}>Analytics &amp; Performance</h3>
              <p style={{ margin: 0, fontSize: "13px", color: "#5B554C", lineHeight: "1.45" }}>
                Helps us gather aggregated, anonymous telemetry of usage patterns, error rates, and load performance to build a faster and cleaner interface.
              </p>
            </div>
            <label style={{ position: "relative", display: "inline-block", width: "44px", height: "24px", alignSelf: "center", cursor: "pointer" }}>
              <input
                type="checkbox"
                checked={localPrefs.analytics}
                onChange={() => handleToggle("analytics")}
                style={{ opacity: 0, width: 0, height: 0 }}
              />
              <span style={{
                position: "absolute",
                top: 0, left: 0, right: 0, bottom: 0,
                backgroundColor: localPrefs.analytics ? "#C2683C" : "#E2D9CC",
                borderRadius: "34px",
                transition: "0.2s",
              }} />
              <span style={{
                position: "absolute",
                height: "18px", width: "18px",
                left: "3px", bottom: "3px",
                backgroundColor: "white",
                borderRadius: "50%",
                transition: "0.2s",
                transform: localPrefs.analytics ? "translateX(20px)" : "none",
              }} />
            </label>
          </div>

          {/* MARKETING */}
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: "16px", borderTop: "1px solid #FCFAF7", paddingTop: "16px" }}>
            <div style={{ flex: 1 }}>
              <h3 style={{ margin: "0 0 4px 0", fontSize: "15px", fontWeight: 600 }}>Marketing Cookies</h3>
              <p style={{ margin: 0, fontSize: "13px", color: "#5B554C", lineHeight: "1.45" }}>
                Enables targeted recommendations and sponsors. Currently disabled by default; kept for future advertising/sponsorship integration only.
              </p>
            </div>
            <label style={{ position: "relative", display: "inline-block", width: "44px", height: "24px", alignSelf: "center", cursor: "pointer" }}>
              <input
                type="checkbox"
                checked={localPrefs.marketing}
                onChange={() => handleToggle("marketing")}
                style={{ opacity: 0, width: 0, height: 0 }}
              />
              <span style={{
                position: "absolute",
                top: 0, left: 0, right: 0, bottom: 0,
                backgroundColor: localPrefs.marketing ? "#C2683C" : "#E2D9CC",
                borderRadius: "34px",
                transition: "0.2s",
              }} />
              <span style={{
                position: "absolute",
                height: "18px", width: "18px",
                left: "3px", bottom: "3px",
                backgroundColor: "white",
                borderRadius: "50%",
                transition: "0.2s",
                transform: localPrefs.marketing ? "translateX(20px)" : "none",
              }} />
            </label>
          </div>

          {/* AI PREFERENCES */}
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: "16px", borderTop: "1px solid #FCFAF7", paddingTop: "16px" }}>
            <div style={{ flex: 1 }}>
              <h3 style={{ margin: "0 0 4px 0", fontSize: "15px", fontWeight: 600 }}>AI Preferences &amp; Copilot</h3>
              <p style={{ margin: 0, fontSize: "13px", color: "#5B554C", lineHeight: "1.45" }}>
                Saves chat history contextual prompts and allows real-time coaching suggestions locally on your browser to optimize AI tokens.
              </p>
            </div>
            <label style={{ position: "relative", display: "inline-block", width: "44px", height: "24px", alignSelf: "center", cursor: "pointer" }}>
              <input
                type="checkbox"
                checked={localPrefs.ai_preferences}
                onChange={() => handleToggle("ai_preferences")}
                style={{ opacity: 0, width: 0, height: 0 }}
              />
              <span style={{
                position: "absolute",
                top: 0, left: 0, right: 0, bottom: 0,
                backgroundColor: localPrefs.ai_preferences ? "#C2683C" : "#E2D9CC",
                borderRadius: "34px",
                transition: "0.2s",
              }} />
              <span style={{
                position: "absolute",
                height: "18px", width: "18px",
                left: "3px", bottom: "3px",
                backgroundColor: "white",
                borderRadius: "50%",
                transition: "0.2s",
                transform: localPrefs.ai_preferences ? "translateX(20px)" : "none",
              }} />
            </label>
          </div>

          {/* FUTURE & GENERAL PERSISTENCE */}
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: "16px", borderTop: "1px solid #FCFAF7", paddingTop: "16px" }}>
            <div style={{ flex: 1 }}>
              <h3 style={{ margin: "0 0 4px 0", fontSize: "15px", fontWeight: 600 }}>Performance &amp; Personalization</h3>
              <p style={{ margin: 0, fontSize: "13px", color: "#5B554C", lineHeight: "1.45" }}>
                Saves personalized interface modules, page caching preferences, and experimental feature flags (A/B testing).
              </p>
            </div>
            <label style={{ position: "relative", display: "inline-block", width: "44px", height: "24px", alignSelf: "center", cursor: "pointer" }}>
              <input
                type="checkbox"
                checked={localPrefs.performance}
                onChange={() => handleToggle("performance")}
                style={{ opacity: 0, width: 0, height: 0 }}
              />
              <span style={{
                position: "absolute",
                top: 0, left: 0, right: 0, bottom: 0,
                backgroundColor: localPrefs.performance ? "#C2683C" : "#E2D9CC",
                borderRadius: "34px",
                transition: "0.2s",
              }} />
              <span style={{
                position: "absolute",
                height: "18px", width: "18px",
                left: "3px", bottom: "3px",
                backgroundColor: "white",
                borderRadius: "50%",
                transition: "0.2s",
                transform: localPrefs.performance ? "translateX(20px)" : "none",
              }} />
            </label>
          </div>

        </div>

        {/* Footer Actions */}
        <div
          style={{
            padding: "18px 28px",
            borderTop: "1px solid #EFE7DC",
            background: "#FCFAF7",
            display: "flex",
            gap: "12px",
            justifyContent: "space-between",
            flexWrap: "wrap",
          }}
        >
          <div style={{ display: "flex", gap: "10px" }}>
            <button
              onClick={handleRejectAll}
              style={{
                border: "1px solid #E2D9CC",
                background: "#FFFFFF",
                color: "#5B554C",
                padding: "10px 18px",
                borderRadius: "100px",
                fontSize: "13px",
                fontWeight: 600,
                cursor: "pointer",
              }}
            >
              Reject All
            </button>
            <button
              onClick={handleAcceptAll}
              style={{
                border: "1px solid #E2D9CC",
                background: "#FFFFFF",
                color: "#5B554C",
                padding: "10px 18px",
                borderRadius: "100px",
                fontSize: "13px",
                fontWeight: 600,
                cursor: "pointer",
              }}
            >
              Accept All
            </button>
          </div>
          <div style={{ display: "flex", gap: "10px" }}>
            <button
              onClick={() => setShowModal(false)}
              style={{
                border: "1px solid #E2D9CC",
                background: "transparent",
                color: "#5B554C",
                padding: "10px 18px",
                borderRadius: "100px",
                fontSize: "13px",
                fontWeight: 600,
                cursor: "pointer",
              }}
            >
              Cancel
            </button>
            <button
              onClick={handleSave}
              style={{
                border: "none",
                background: "#C2683C",
                color: "#FFFFFF",
                padding: "10px 22px",
                borderRadius: "100px",
                fontSize: "13px",
                fontWeight: 600,
                cursor: "pointer",
              }}
            >
              Save Preferences
            </button>
          </div>
        </div>

        <style dangerouslySetInnerHTML={{ __html: `
          @keyframes scaleIn {
            from { transform: scale(0.92); opacity: 0; }
            to { transform: scale(1); opacity: 1; }
          }
        `}} />
      </div>
    </div>
  );
}
