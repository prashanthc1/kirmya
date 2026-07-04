/**
 * lib/api/client.ts — the single HTTP client every feature uses to reach the
 * Kirmya backend. It encodes the auth + contract conventions from CLAUDE.md §5:
 *
 *  - Bearer access token in the `Authorization` header (held in memory, never
 *    in localStorage — XSS-safe).
 *  - httpOnly, SameSite=Strict refresh cookie sent automatically via
 *    `credentials: "include"`; used transparently to refresh on a 401.
 *  - Double-submit CSRF: the `csrf_token` cookie (set by GET /auth/csrf) is
 *    echoed in the `X-CSRF-Token` header on every state-changing request, which
 *    is what protects the cookie-authenticated endpoints (refresh, logout).
 *  - Response envelopes: success `{ data: ... }`, error
 *    `{ error: { code, message, details } }`.
 *
 * Refresh-on-401 is single-flight: concurrent requests that all see a 401 share
 * one refresh round-trip and then retry once.
 */
import { config } from "@/lib/config";

/** Shape of the error envelope returned by the backend. */
export interface ApiErrorBody {
  code: string;
  message: string;
  details?: unknown;
}

/** Thrown for any non-2xx response (after a failed refresh, for 401s). */
export class ApiError extends Error {
  readonly status: number;
  readonly code: string;
  readonly details?: unknown;

  constructor(status: number, body: ApiErrorBody) {
    super(body.message || `request failed with status ${status}`);
    this.name = "ApiError";
    this.status = status;
    this.code = body.code || "unknown";
    this.details = body.details;
  }
}

// --- in-memory access token -------------------------------------------------
// Kept in a module variable rather than storage so it never survives a full
// reload (the refresh cookie re-establishes the session) and is unreachable to
// injected scripts reading storage.
let accessToken: string | null = null;

export function setAccessToken(token: string | null): void {
  accessToken = token;
}

export function getAccessToken(): string | null {
  return accessToken;
}

// --- helpers ----------------------------------------------------------------
const SAFE_METHODS = new Set(["GET", "HEAD", "OPTIONS"]);

/** Reads a cookie value by name in the browser; null on the server. */
function readCookie(name: string): string | null {
  if (typeof document === "undefined") return null;
  const match = document.cookie.match(
    new RegExp("(?:^|; )" + name.replace(/([.$?*|{}()[\]\\/+^])/g, "\\$1") + "=([^;]*)"),
  );
  return match ? decodeURIComponent(match[1]) : null;
}

export interface RequestOptions extends Omit<RequestInit, "body"> {
  /** JSON-serializable body; set automatically with the JSON content-type. */
  json?: unknown;
  /** Pre-encoded body (FormData, string, etc.). Takes precedence over `json`. */
  body?: BodyInit | null;
  /** Skip the refresh-on-401 retry (used internally by the refresh call). */
  skipAuthRefresh?: boolean;
}

function buildUrl(path: string): string {
  if (/^https?:\/\//.test(path)) return path;
  const base = config.apiBase.replace(/\/$/, "");
  return path.startsWith("/") ? base + path : `${base}/${path}`;
}

function buildHeaders(method: string, opts: RequestOptions): Headers {
  const headers = new Headers(opts.headers);
  if (opts.json !== undefined && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (accessToken && !headers.has("Authorization")) {
    headers.set("Authorization", `Bearer ${accessToken}`);
  }
  // Double-submit CSRF token for cookie-authenticated, state-changing requests.
  if (!SAFE_METHODS.has(method)) {
    const csrf = readCookie("csrf_token");
    if (csrf && !headers.has("X-CSRF-Token")) {
      headers.set("X-CSRF-Token", csrf);
    }
  }
  return headers;
}

// --- single-flight refresh --------------------------------------------------
let refreshInFlight: Promise<boolean> | null = null;

/**
 * Attempts to mint a new access token from the refresh cookie. Returns true on
 * success (and updates the in-memory token). Concurrent callers share one call.
 */
async function refreshAccessToken(): Promise<boolean> {
  if (!refreshInFlight) {
    refreshInFlight = (async () => {
      try {
        const res = await fetch(buildUrl("/auth/refresh"), {
          method: "POST",
          credentials: "include",
          headers: (() => {
            const h = new Headers();
            const csrf = readCookie("csrf_token");
            if (csrf) h.set("X-CSRF-Token", csrf);
            return h;
          })(),
        });
        if (!res.ok) return false;
        const payload = (await res.json()) as { data?: { access_token?: string } };
        const token = payload?.data?.access_token;
        if (!token) return false;
        setAccessToken(token);
        return true;
      } catch {
        return false;
      } finally {
        // Reset as soon as this refresh settles. Concurrent callers already hold
        // a reference to this same promise (so they still observe its result);
        // clearing synchronously prevents a *later*, independent request from
        // reusing an already-resolved refresh outcome.
        refreshInFlight = null;
      }
    })();
  }
  return refreshInFlight;
}

async function parseError(res: Response): Promise<ApiError> {
  let body: ApiErrorBody = { code: "unknown", message: res.statusText };
  try {
    const payload = (await res.json()) as { error?: ApiErrorBody };
    if (payload?.error) body = payload.error;
  } catch {
    /* non-JSON error body — keep the status-text fallback */
  }
  return new ApiError(res.status, body);
}

/**
 * Core request method. Returns the unwrapped `data` payload on success and
 * throws {@link ApiError} otherwise. On a 401 it transparently refreshes once
 * and retries.
 */
export async function request<T = unknown>(path: string, opts: RequestOptions = {}): Promise<T> {
  const method = (opts.method ?? "GET").toUpperCase();

  const send = (): Promise<Response> => {
    const headers = buildHeaders(method, opts);
    const body = opts.body ?? (opts.json !== undefined ? JSON.stringify(opts.json) : undefined);
    return fetch(buildUrl(path), {
      ...opts,
      method,
      headers,
      body,
      credentials: "include",
    });
  };

  let res = await send();

  if (res.status === 401 && !opts.skipAuthRefresh) {
    const refreshed = await refreshAccessToken();
    if (refreshed) {
      res = await send();
    } else {
      setAccessToken(null);
    }
  }

  if (!res.ok) throw await parseError(res);

  if (res.status === 204) return undefined as T;
  const text = await res.text();
  if (!text) return undefined as T;
  const payload = JSON.parse(text) as { data?: T };
  // Endpoints wrap success in { data }, but tolerate bare bodies too.
  return (payload && "data" in payload ? payload.data : (payload as unknown)) as T;
}

/** Convenience verbs. */
export const api = {
  get: <T = unknown>(path: string, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: "GET" }),
  post: <T = unknown>(path: string, json?: unknown, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: "POST", json }),
  put: <T = unknown>(path: string, json?: unknown, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: "PUT", json }),
  patch: <T = unknown>(path: string, json?: unknown, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: "PATCH", json }),
  delete: <T = unknown>(path: string, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: "DELETE" }),
};

// --- auth helpers built on the client --------------------------------------

/** Fetches a fresh CSRF token (also sets the `csrf_token` cookie). */
export async function fetchCsrfToken(): Promise<string> {
  const data = await api.get<{ csrf_token: string }>("/auth/csrf");
  return data.csrf_token;
}

/** Clears the in-memory token after a logout round-trip (best-effort). */
export async function logout(): Promise<void> {
  try {
    await api.post("/auth/logout");
  } finally {
    setAccessToken(null);
  }
}
