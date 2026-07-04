/**
 * Runtime configuration for the browser client.
 *
 * The browser always talks to the Next.js origin under `NEXT_PUBLIC_API_BASE`
 * (default `/api/v1`); Next proxies those requests to the Go backend
 * (`API_PROXY_TARGET`, server-side only — see next.config.ts rewrites). Keeping
 * the browser on a same-origin path means the httpOnly refresh cookie and the
 * `csrf_token` double-submit cookie are first-party.
 */
export const config = {
  /** Base path the browser client prefixes onto every API route. */
  apiBase: process.env.NEXT_PUBLIC_API_BASE ?? "/api/v1",
} as const;

export type AppConfig = typeof config;
