# Security Review — Kirmya Auth/Security Surface

> ## ✅ Remediation status — 2026-06-15 (completion pass)
>
> All 11 findings were actioned in code/docs. **Note:** changes were made with
> the editor file tools; the Go toolchain was unavailable in this environment, so
> the backend was **not compiled or unit-tested** — a `go build ./... && go test
> ./...` run is still required before merge (Definition of Done).
>
> | ID | Status | What changed |
> |---|---|---|
> | H1 OAuth state | ✅ Fixed | `OAuthStart` sets an httpOnly `oauth_state` cookie; `OAuthCallback` requires `state` in the body and constant-time compares it (`crypto/subtle`). DTO gained `state`. PKCE noted as follow-up. |
> | H2 Mailer logs tokens | ✅ Fixed | `LogMailer` no longer logs raw tokens (only a non-reversible ref) and **fails closed in production**. |
> | H3 Email verify at login | ✅ Fixed | `Login` returns `ErrEmailNotVerified` for unverified password accounts (toggle `EMAIL_VERIFICATION_REQUIRED`, default true); mapped to 403. |
> | M1 CSRF/Secure cookie | ◑ Partial | Cookies now `Secure` over HTTPS (not only prod). Double-submit enforcement deferred (client doesn't participate yet) — documented in `CSRF_SECURITY.md`. |
> | M2 XFF spoofing | ✅ Fixed | `clientIP` honours `X-Forwarded-For` only when `TRUST_PROXY=true`, else uses `RemoteAddr`. |
> | M3 Token revocation lag | ◑ Mitigated | `JWT_ACCESS_TTL` now capped at 1h. Durable `token_version` revocation needs a migration (tracked). |
> | M4 Email PII in directory | ✅ Fixed | `email` removed from `directoryEntryDTO`, the OpenAPI `DirectoryUser` schema, and the TS type. |
> | L1 JWT_SECRET | ✅ Fixed | Production requires `JWT_SECRET` ≥ 32 bytes (process aborts otherwise). |
> | L2 MFA key | ✅ Fixed | Production aborts if neither `MFA_ENC_KEY` nor `JWT_SECRET` is set (no hard-coded prod key). |
> | L3 MFA throttle | ✅ Fixed | Per-account in-memory attempt limiter (5/15 min) on the login MFA step. Same-window replay prevention noted (needs persisted counter). |
> | L4 Stale docs | ✅ Fixed | `AUTHENTICATION.md` and `CSRF_SECURITY.md` rewritten to match the real implementation. |

**Reviewer:** security-auth-reviewer (OWASP-oriented, static review only)
**Date:** 2026-06-15
**Scope:** AuthN, AuthZ/RBAC, CSRF, input handling (SQL), secrets/data exposure
**Grounding docs:** `AUTHENTICATION.md`, `CSRF_SECURITY.md`
**Method:** Read-only static analysis. Go toolchain not run. Line numbers reference the files at review time.

> Note on grounding docs: `AUTHENTICATION.md` and `CSRF_SECURITY.md` describe an **older SQLite/bcrypt/localStorage-CSRF design** that does **not** match the current code. The implemented identity module uses PostgreSQL, Argon2id, HS256 JWT access tokens, an httpOnly SameSite=Strict refresh cookie, and refresh-token rotation with reuse detection. Several "documented controls" (length-only CSRF token check, localStorage CSRF storage, bcrypt cost 10, origin whitelist) are **stale documentation**, not the running implementation. Findings below are graded against the **actual code**, with doc drift flagged as its own Low finding.

---

## Severity Summary

| Severity | Count |
|----------|-------|
| Critical | 0 |
| High     | 3 |
| Medium   | 4 |
| Low      | 4 |
| **Total**| **11** |

**No Critical findings.** The highest-impact issues are the missing OAuth `state` validation (login-CSRF / account fixation), account-takeover tokens written to logs by the default mailer, and email-verification not being enforced at login.

---

## High

### H1 — OAuth `state` parameter is generated but never validated (login CSRF / account linking)
**Location:** `backend/internal/identity/api/handlers.go:166-195` (`OAuthStart`, `OAuthCallback`); `backend/internal/identity/application/oauth.go:11-67`; `backend/internal/identity/infrastructure/oauth/provider.go:32-40`

**Risk:** `OAuthStart` mints a `state` value and returns it to the client, but `OAuthCallback` accepts only `{code}` (`oauthCallbackRequest`) and **never receives or verifies `state`**. There is no server-side persistence of the issued state and no PKCE. The OAuth authorization-code flow therefore has no CSRF protection on the callback and no proof that the code was redeemed by the same browser that began the flow.

**Exploit sketch:** Attacker completes an OAuth flow up to obtaining a `code` for the *attacker's* Google account, then tricks a victim (or victim's already-authenticated app session) into POSTing that `code` to `/oauth/{provider}/callback`, linking/logging the victim into the attacker-controlled identity (or vice-versa, fixating the victim's session to attacker's account).

**Fix:** Persist the issued `state` server-side (httpOnly cookie or short-lived store) at `OAuthStart`; require `state` in the callback payload and compare in constant time before calling `Exchange`. Add PKCE (`code_challenge`/`code_verifier`, S256) since both Google and LinkedIn support it. Reject callbacks whose `state` is missing/unknown.

### H2 — Account-takeover tokens written to application logs by default mailer
**Location:** `backend/internal/identity/infrastructure/mailer/log_mailer.go:24-32`; wired unconditionally in `backend/internal/identity/module.go:44` (`Mailer: mailer.NewLogMailer()`)

**Risk:** The default and only-wired `Mailer` logs the **raw** email-verification and password-reset tokens to stdout (`log.Printf("[mailer] password reset ... token=%s", rawToken)`). These raw tokens are exactly the secrets that grant email verification and password reset (= account takeover). There is no env guard switching to a real mailer; whatever ships uses this.

**Exploit sketch:** Anyone with read access to application logs (aggregator, container stdout, shared ops dashboard) calls `ForgotPassword(victim)`, reads the logged `reset-password?token=...`, and resets the victim's password.

**Fix:** Do not log raw tokens. For dev, log only that an email was "sent" plus a redacted/hashed reference. Gate the LogMailer behind `APP_ENV != production` and fail closed (or refuse to start) if no real mailer is configured in production. Treat reset/verify tokens as secrets in all logging.

### H3 — Email verification is never enforced at login
**Location:** `backend/internal/identity/application/service.go:154-185` (`Login`); `ErrEmailNotVerified` defined at `service.go:22` but never returned anywhere

**Risk:** `Register` sets `EmailVerified=false` and sends a verification email, but `Login` checks only `IsActive()`, `HasPassword()`, and the password hash. It never checks `u.EmailVerified`. The `ErrEmailNotVerified` sentinel exists but is dead code. Unverified, attacker-controlled email addresses get full sessions.

**Exploit sketch:** Register with `victim@corp.com` (an address you don't control), never verify, log in immediately — you now hold a session whose email claim is an address you do not own, which downstream features (directory, messaging, referrals keyed on email/identity, recruiter trust) may treat as authentic.

**Fix:** If product policy requires verification, return `ErrEmailNotVerified` from `Login` (and map it in `writeError`) when `!u.EmailVerified` for password accounts. If self-serve-before-verify is intended, document it and ensure no feature treats the email claim as proof of ownership. OAuth accounts (set `EmailVerified=true`) are fine.

---

## Medium

### M1 — Refresh endpoint relies solely on SameSite=Strict; no anti-CSRF token; `csrf_token` cookie is issued but never validated
**Location:** `backend/internal/identity/api/handlers.go:81-95` (`Refresh`), `267-279` (`CSRF` issues a token), `293-303` (refresh cookie); `backend/internal/platform/middleware/csrf.go:26-49` (`VerifyOrigin`, **off** unless `CSRF_VERIFY_ORIGIN=true`)

**Risk:** `POST /auth/refresh` authenticates purely from the `refresh_token` cookie. Its CSRF defense is the cookie's `SameSite=Strict` attribute. That is a reasonable primary control, **but**: (a) the `Secure` flag is only set when `APP_ENV==production` (`secureCookies()`, `handlers.go:317`), so non-prod/misconfigured deploys send the refresh cookie over plain HTTP; (b) the app also issues a `csrf_token` cookie (`HttpOnly:false`) that is **never validated** on any route — it is dead defense-in-depth that may give a false sense of protection; (c) the Origin check (`VerifyOrigin`) is disabled by default and allows requests with no `Origin` header outright.

**Exploit sketch:** On a deployment where the app is reachable over HTTP (or a browser that has historically downgraded SameSite handling), a cross-site `<form>`/`fetch` auto-replays the refresh cookie to mint a fresh access token. Low likelihood on a correctly-configured modern browser/HTTPS, hence Medium.

**Fix:** Either enforce a double-submit CSRF token on `/auth/refresh` and `/auth/logout` (validate the issued `csrf_token` cookie against an `X-CSRF-Token` header) **or** remove the unused `csrf_token` cookie to avoid implying protection that isn't enforced. Set `Secure` on cookies whenever the request is HTTPS (not only when `APP_ENV==production`). Consider enabling `VerifyOrigin` by default.

### M2 — `clientIP` blindly trusts `X-Forwarded-For` (audit/log spoofing, rate-limit evasion)
**Location:** `backend/internal/identity/api/handlers.go:32-37`; consumed by audit records (`service.go:288-290`, persisted in `oauth_mfa_audit.go:84-94`) and login/register IP fields

**Risk:** `clientIP` returns the raw `X-Forwarded-For` header if present, with no trusted-proxy allowlist and taking the whole (potentially attacker-supplied, comma-joined) value. An attacker sets `X-Forwarded-For: <anything>` to forge the IP recorded in `audit_logs` and (if a future rate limiter keys on it) to evade throttling.

**Exploit sketch:** Send login attempts with rotating spoofed `X-Forwarded-For` values to poison the audit trail and frustrate IP-based investigation/blocking.

**Fix:** Only honour `X-Forwarded-For` from known proxy hops; parse the left-most untrusted-but-real client per a configured trusted-proxy count, else fall back to `RemoteAddr`. Document the deployment's proxy assumption.

### M3 — JWT access tokens cannot be revoked; role/status changes do not take effect until expiry
**Location:** `backend/internal/identity/infrastructure/jwtauth/token_factory.go:69-105`; middleware `backend/internal/identity/api/middleware.go:19-68`; admin status/role changes in `backend/internal/admin/...`

**Risk:** Authorization is derived entirely from the signed JWT claims (`Roles`, `Subject`). There is no per-request check against current account status or a revocation list. When an admin suspends a user (`SetUserStatus`) or revokes a role (`RevokeRole`), the user's existing access token (TTL 15 min default, but `JWT_ACCESS_TTL` is operator-configurable and could be large) remains fully valid until it expires. A reset password / logout revokes refresh tokens but not the outstanding access token.

**Exploit sketch:** A user about to be banned keeps acting (posting, messaging, escalated-role actions) for the remainder of the access-token lifetime after the ban is applied.

**Fix:** Keep access-token TTL short (the 15-min default is good — enforce a sane upper bound on `JWT_ACCESS_TTL`). For sensitive actions, re-check account status/roles server-side against the DB, or add a lightweight revocation/`token_version` check. At minimum, document the revocation-lag window.

### M4 — Directory/search endpoints expose every user's email to any authenticated user (PII over-exposure)
**Location:** `backend/internal/identity/api/handlers.go:223-255` (`directoryEntryDTO` includes `Email`; `SearchUsers`, `GetUser`); repo `user_repository.go:148-166` (`directoryCols` selects `u.email`)

**Risk:** `GET /users/search?q=` and `GET /users/{id}` return each matched user's email address to any authenticated caller. Search matches on email too (`... OR u.email ILIKE $1`), so an attacker can enumerate/confirm emails and harvest a member email list — useful for phishing/spam and a privacy concern (GDPR-style PII minimisation).

**Exploit sketch:** Authenticated attacker iterates the search endpoint with common name/domain fragments and scrapes a directory of `{name, email}` for the whole user base.

**Fix:** Drop `email` from the public directory DTO (return it only on `/users/me` / the owner's own record). If email is needed for a specific connected relationship, gate it behind that relationship. Avoid matching search queries against email unless intentionally a feature.

---

## Low

### L1 — `JWT_SECRET` falls back to a random per-process key; multi-instance/dev tokens silently break or are weakly keyed
**Location:** `backend/internal/identity/infrastructure/jwtauth/token_factory.go:45-58`

**Risk:** If `JWT_SECRET` is unset, a random per-process key is generated (logged as a WARNING). Across multiple instances each gets a different key (tokens issued by one are rejected by another), and a restart invalidates all sessions. This is safe-by-accident in the sense it isn't a *predictable* key, but it is an operational footgun and there is no minimum-length enforcement when the secret **is** set (a 4-char `JWT_SECRET` is accepted for HS256).

**Fix:** In production (`APP_ENV==production`) refuse to start without a `JWT_SECRET` of sufficient entropy (e.g. >= 32 bytes). Keep the random-key dev fallback but log loudly.

### L2 — TOTP MFA encryption key falls back to a hard-coded constant
**Location:** `backend/internal/identity/infrastructure/crypto/totp.go:24-33`

**Risk:** `NewTOTPService` derives the AES-256-GCM key from `MFA_ENC_KEY`, else `JWT_SECRET`, else the literal `"kirmya-dev-mfa-key"`. If both env vars are unset in a real deployment, all TOTP secrets at rest are encrypted under a publicly-known key (i.e. effectively plaintext to anyone with the source).

**Fix:** Fail closed in production if neither `MFA_ENC_KEY` nor a strong `JWT_SECRET` is set; never ship the hard-coded fallback as a usable production path.

### L3 — TOTP validation has no replay/rate-limit window and MFA-challenge step is unthrottled
**Location:** `backend/internal/identity/infrastructure/crypto/totp.go:50-56` (default skew, no used-code tracking); `service.go:173-181` (MFA branch in `Login`)

**Risk:** `totp.Validate` accepts a code within the default window with no record of already-consumed codes, and there is no brute-force throttle on the MFA `code` submission in `Login`. A 6-digit code is brute-forceable if unlimited attempts are allowed; a captured code can be replayed within its window.

**Fix:** Add per-account attempt rate-limiting/lockout on the MFA step, and track consumed codes (or a last-used counter) to prevent same-window replay.

### L4 — Documentation drift: `AUTHENTICATION.md` / `CSRF_SECURITY.md` describe controls that don't match the code
**Location:** `AUTHENTICATION.md`, `CSRF_SECURITY.md`

**Risk:** The docs claim bcrypt cost 10 (code uses Argon2id), SQLite (code uses PostgreSQL), localStorage-stored JWT and CSRF tokens with a length-only CSRF check and an Origin whitelist (none of which is the implemented flow), and an error-response example that leaks raw DB driver errors (`Error 1062 ... Duplicate entry`). Stale security docs cause reviewers/operators to trust controls that aren't there (and to overlook the real ones). The leaked-DB-error example, if ever implemented as shown, would be an internals-disclosure issue — the actual `writeError` (`handlers.go:326-344`) correctly returns generic messages, which is the right behaviour.

**Fix:** Rewrite both docs to describe the current implementation (Argon2id params, HS256 + refresh rotation/reuse detection, SameSite=Strict httpOnly refresh cookie, Bearer-token CSRF immunity). Remove the raw-DB-error example.

---

## Controls Verified CORRECT (do not change)

- **Password hashing — Argon2id, OWASP baseline.** `crypto/password.go:25-31` uses m=64 MiB, t=3, p=2, 16-byte salt, 32-byte key, with constant-time verify (`subtle.ConstantTimeCompare`, line 62) and proper PHC encode/decode. Correct.
- **JWT algorithm pinning — no `alg` confusion.** `token_factory.go:92-97` rejects any non-HMAC signing method in the keyfunc and checks `parsed.Valid`. Correct.
- **Refresh-token rotation + reuse detection (theft response).** `service.go:207-243`: tokens are single-use; reuse of an already-`ReplacedBy` token triggers `RevokeFamily` (whole-family revocation). Only SHA-256 hashes are stored (`token_factory.go:107-122`, repo `token_repository.go`). Correct and well-designed.
- **Refresh cookie hardening.** `handlers.go:293-303`: `HttpOnly`, `SameSite=Strict`, path-scoped to `/api/v1/auth`. Correct (see M1 re: `Secure` only in prod).
- **Account-enumeration resistance.** `ForgotPassword`/`ResendVerification` always return success regardless of account existence (`auth_flows.go:51-83`); login returns a single generic `ErrInvalidCredentials` for unknown user, no password, and bad password. Correct.
- **One-time, hashed, expiring verification/reset tokens.** `token_repository.go:84-104`: consume is an atomic `UPDATE ... WHERE used_at IS NULL AND expires_at > now() RETURNING user_id`. Correct.
- **Parameterized SQL throughout.** Every `infrastructure/postgres/*.go` query uses `$N` placeholders with user values passed as args; the dynamic-WHERE builders in `admin/.../repository.go:52-83`, `jobs/.../repository.go:50-70`, and identity search interpolate only placeholder indices/column constants, never user values. No string-built SQL with user input found. Correct.
- **IDOR / ownership checks.**
  - Resume: `requireOwner` on every read/mutate (`resume/application/service.go:143-218`). Correct.
  - Profile sub-resources: every update/delete scoped `WHERE id=$1 AND user_id=$2` with `owned()` returning `ErrNotFound` on 0 rows (`profile/.../repository.go:112-219`). Correct.
  - Messaging: `requireParticipant` on list/send/read/typing (`messaging/application/service.go:86-176`). Correct.
  - Mentorship: status changes gated to the owning mentor (`UpdateStatus`, mentor.UserID check), reviews gated to the mentee on a completed session, self-booking blocked (`mentorship/application/service.go:96-137,59-66`). Correct.
  - Referrals: seeker-cannot-review-own, directed-referral ownership, both-participant outcome checks (`referrals/application/service.go:90-130`). Correct.
- **RBAC wiring.** All feature routes wrapped by `AuthMiddleware`; admin routes gated by `AdminMiddleware` (RequireRole admin), recruiter writes by `RoleMiddleware(RoleRecruiter)` (`platform/router.go:54-68`, `admin/api/routes.go`). Role check is server-side in middleware (`identity/api/middleware.go:36-56`). Correct.
- **Error responses don't leak internals.** `writeError` maps known errors to generic messages and falls back to `"something went wrong"` for everything else (`handlers.go:326-344`). Correct (contradicts the stale doc example).
- **Security headers + request body limits.** `SecurityHeaders` sets `nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy`, `Permissions-Policy`, HSTS, path-aware CSP (`platform/middleware/security.go`). JSON bodies capped at 1 MiB via `http.MaxBytesReader` in each `decode` (`handlers.go:28-30`, etc.). Server has read/write/idle timeouts (`server.go:29-35`). Correct.
- **MFA secret encryption at rest.** AES-256-GCM with random nonce per secret (`crypto/totp.go:58-97`). Correct construction (see L2 re: key sourcing).

---

## Recommended Remediation Order

1. **H1** — Validate OAuth `state` (+ add PKCE). Login-CSRF/account-fixation on the auth boundary.
2. **H2** — Stop logging raw verify/reset tokens; gate LogMailer to non-prod.
3. **H3** — Decide and enforce the email-verification policy at `Login`.
4. **M2 / M1** — Trusted-proxy `X-Forwarded-For` handling; resolve the unused `csrf_token` cookie and `Secure`-over-HTTPS.
5. **M4 / M3** — Trim email from directory DTOs; document/limit access-token revocation lag.
6. **L1–L4** — Production fail-closed checks for `JWT_SECRET` / `MFA_ENC_KEY`, MFA brute-force throttle, and rewrite the two grounding docs to match reality.
