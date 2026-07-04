"use client";

/**
 * Client-side authentication state for the app.
 *
 * The access token lives only in memory (in the API client), so on a fresh page
 * load there is no token yet. On mount the provider seeds the CSRF cookie and
 * then asks GET /users/me: the API client transparently uses the httpOnly
 * refresh cookie to mint a new access token on the 401, so a returning user with
 * a valid session is restored without any visible login step. If there is no
 * session, `user` stays null and the logged-out UI is shown.
 */
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import {
  api,
  fetchCsrfToken,
  logout as apiLogout,
  setAccessToken,
} from "@/lib/api/client";

export interface AuthUser {
  id: string;
  email: string;
  full_name: string;
  email_verified: boolean;
  mfa_enabled: boolean;
  roles: string[];
}

interface AuthContextValue {
  user: AuthUser | null;
  /** True until the initial session-restore attempt has finished. */
  loading: boolean;
  /** Replace the current user (e.g. right after register/login). */
  setUser: (user: AuthUser | null) => void;
  /** Re-fetch the current user from the server. */
  refreshUser: () => Promise<void>;
  /** Clear the session on the server and locally. */
  signOut: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [loading, setLoading] = useState(true);

  const refreshUser = useCallback(async () => {
    try {
      const me = await api.get<AuthUser>("/users/me");
      setUser(me ?? null);
    } catch {
      setUser(null);
    }
  }, []);

  useEffect(() => {
    let active = true;
    (async () => {
      // Seed the csrf_token cookie so cookie-auth calls (refresh/logout) carry a
      // valid double-submit token where it is enforced. Best-effort.
      try {
        await fetchCsrfToken();
      } catch {
        /* ignore */
      }
      try {
        const me = await api.get<AuthUser>("/users/me");
        if (active) setUser(me ?? null);
      } catch {
        if (active) setUser(null);
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  const signOut = useCallback(async () => {
    try {
      await apiLogout();
    } finally {
      setAccessToken(null);
      setUser(null);
    }
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({ user, loading, setUser, refreshUser, signOut }),
    [user, loading, refreshUser, signOut],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return ctx;
}
