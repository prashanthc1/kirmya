# Authentication — Kirmya

This document describes the **authentication system as implemented** in
`backend/internal/identity`. It supersedes earlier drafts that described a
SQLite / bcrypt / localStorage-token design — none of which is the running code.

## Stack at a glance

| Concern | Implementation |
|---|---|
| User store | PostgreSQL (pgx), `users` + role/credential tables |
| Password hashing | **Argon2id** (m=64 MiB, t=3, p=2, 16-byte salt, 32-byte key), PHC-encoded, constant-time verify |
| Access token | **HS256 JWT**, 15-min default TTL (`JWT_ACCESS_TTL`, capped at 1h), claims: `sub`, `email`, `roles`, `iss=kirmya` |
| Refresh token | Opaque 256-bit random; only a SHA-256 hash is stored; httpOnly `SameSite=Strict` cookie scoped to `/api/v1/auth` |
| Refresh rotation | Single-use with **reuse detection** — replaying a rotated token revokes the whole family |
| MFA | TOTP (RFC 6238), secret AES-256-GCM encrypted at rest; per-account attempt throttle on the login step |
| OAuth | Google + LinkedIn (OIDC authorization-code), **cookie-bound `state`** validated on callback |
| Email/reset tokens | Opaque, hashed at rest, one-time, expiring (verify 24h, reset 1h) |

## Token model

- **Access token (JWT, Bearer).** Held in memory on the client (mirrored to
  `localStorage` only for reload continuity) and sent as
  `Authorization: Bearer <token>`. Stateless: authorization derives from the
  signed `roles`/`sub` claims. The signing method is pinned to HMAC in
  `Parse` (no `alg` confusion). Because access-protected endpoints read a
  custom header, they are inherently immune to CSRF.
- **Refresh token (opaque cookie).** `POST /api/v1/auth/refresh` reads the
  httpOnly `refresh_token` cookie and rotates it. The raw token is never stored;
  lookups use its SHA-256 hash. Rotation marks the old token replaced and issues
  a new one in the same family. Re-use of an already-replaced token triggers
  `RevokeFamily` (theft response).

## Login flow

1. `POST /api/v1/auth/register` → Argon2id hash, default role `job_seeker`,
   verification email queued. `EmailVerified=false`.
2. `POST /api/v1/auth/login` → verifies password (generic
   `invalid credentials` for unknown user / no password / bad password —
   enumeration-resistant).
   - **Email verification is enforced**: unverified password accounts are
     rejected with `email not verified` (403). Set
     `EMAIL_VERIFICATION_REQUIRED=false` to relax this in local/dev/seed.
   - If MFA is enabled and no `code` is supplied, the response is
     `{ "mfa_required": true }`; the client re-submits `login` with `code`.
     MFA code attempts are rate-limited per account (5 / 15 min).
3. On success the server returns `{ access_token, token_type, expires_in, user }`
   and sets the `refresh_token` cookie.

## OAuth flow (Google / LinkedIn)

1. `GET /api/v1/auth/oauth/{provider}` mints an unpredictable `state`, sets it in
   an httpOnly `oauth_state` cookie (Lax, 10-min), and returns the provider
   authorization URL + `state`.
2. The browser completes consent and the client calls
   `POST /api/v1/auth/oauth/{provider}/callback` with `{ code, state }`.
3. The server **constant-time compares** `state` against the `oauth_state`
   cookie before exchanging the code (login-CSRF / account-fixation defense),
   then provisions/links the account (`EmailVerified=true` for OAuth) and issues
   a session. The state cookie is single-use.

> Roadmap: add PKCE (`code_challenge`/`code_verifier`, S256) on top of the
> cookie-bound state — both providers support it.

## Sessions, logout, password reset

- `POST /auth/logout` revokes the presented refresh token and clears the cookie.
- `POST /auth/forgot-password` / `reset-password` use one-time, hashed,
  1-hour tokens. `forgot-password` always returns success (no enumeration).
- Email delivery: the default `LogMailer` is **dev-only**. It never logs raw
  tokens (only a non-reversible reference) and **fails closed in production**
  (`APP_ENV=production`) until a real mailer is configured.

## Required configuration

| Env var | Purpose | Production requirement |
|---|---|---|
| `JWT_SECRET` | HS256 signing key | Required, **≥ 32 bytes** (process aborts otherwise) |
| `JWT_ACCESS_TTL` | Access-token TTL (seconds) | Optional; clamped to ≤ 1h |
| `JWT_REFRESH_TTL` | Refresh-token TTL (seconds) | Optional (default 30d) |
| `MFA_ENC_KEY` | AES key source for TOTP secrets | Required if MFA used (falls back to `JWT_SECRET`; aborts in prod if neither set) |
| `EMAIL_VERIFICATION_REQUIRED` | Enforce verified email at login | Optional (default `true`) |
| `TRUST_PROXY` | Honour `X-Forwarded-For` for client IP | Set `true` only behind a trusted proxy |
| `APP_ENV` | `production` enables fail-closed checks + `Secure` cookies | Set `production` in prod |

## Known residual risks (tracked)

- **Access-token revocation lag (M3).** A suspended user / revoked role remains
  valid until the access token expires. Mitigated by the short TTL + 1h cap;
  a `token_version`/revocation check is the durable fix (needs a migration).
- **MFA replay within a code window (L3).** The per-account throttle limits
  brute force; same-window replay prevention needs a persisted last-used counter.

See `CSRF_SECURITY.md` for the CSRF model and `docs/reports/SECURITY_REVIEW.md`
for the full findings and remediation status.
