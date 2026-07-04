# CSRF Protection — Kirmya

This describes the **CSRF model as implemented**. It supersedes the earlier draft
that described localStorage-stored tokens, a length-only token check, and an
Origin whitelist as the primary defense — that is not the running design.

## Threat model and primary defenses

CSRF only affects requests the browser authenticates *ambiently* (cookies). In
Kirmya:

1. **Bearer-token endpoints are CSRF-immune.** Every authenticated feature
   endpoint authenticates from the `Authorization: Bearer <jwt>` header, which a
   cross-site page cannot set on a cross-origin request. No cookie is consulted,
   so there is nothing for a forged request to ride.

2. **Cookie-authenticated endpoints rely on `SameSite=Strict`.** Only
   `/api/v1/auth/refresh` and `/api/v1/auth/logout` authenticate from the
   httpOnly `refresh_token` cookie. That cookie is `HttpOnly`, `SameSite=Strict`,
   path-scoped to `/api/v1/auth`, and `Secure` in production **and** whenever the
   request arrives over HTTPS (direct TLS or `X-Forwarded-Proto: https`). With
   `SameSite=Strict` the browser does not attach the cookie to cross-site
   requests, which is the CSRF defense for these two routes.

3. **OAuth callback uses cookie-bound `state`.** `GET /auth/oauth/{provider}`
   sets an httpOnly `oauth_state` cookie; the callback must echo the same value,
   compared in constant time. This prevents login-CSRF / account fixation on the
   OAuth boundary.

## Enforced double-submit check on cookie routes

As of the latest change, `/api/v1/auth/refresh` and `/api/v1/auth/logout`
**enforce** the double-submit cookie pattern in addition to `SameSite=Strict`.
The handler (`verifyDoubleSubmitCSRF`) requires the request to carry an
`X-CSRF-Token` header equal (constant-time compare) to the non-httpOnly
`csrf_token` cookie issued by `GET /auth/csrf`. A cross-site attacker can ride
the victim's cookies but cannot read `csrf_token` to copy it into the header, so
a match proves the request came from our own first-party JavaScript.

It is **opt-in** (`CSRF_DOUBLE_SUBMIT=true`) and is **enabled in production**
(`docker-compose.prod.yml`). The default is off to match the conservative
posture of the Origin check below — enabling it requires every cookie-auth
client to first call `GET /auth/csrf` and echo the token, so it is turned on per
environment rather than globally. The frontend client
(`frontend/lib/api/client.ts`) reads the `csrf_token` cookie and sends the
`X-CSRF-Token` header automatically on every state-changing request, and calls
`GET /auth/csrf` to seed the cookie.

## Origin check (defense-in-depth)

`platform/middleware.VerifyOrigin` adds an Origin allowlist on state-changing
methods. It is off unless `CSRF_VERIFY_ORIGIN=true`, and is **enabled in
production** (`docker-compose.prod.yml`) where `APP_URL` is pinned to the single
canonical domain. It stays off by default elsewhere because it rejects
legitimate traffic when the app is reached via multiple hosts or an
origin-rewriting proxy. Requests with no `Origin` header (curl, server-to-server)
are allowed.

## The `/auth/csrf` endpoint and `csrf_token` cookie

`GET /auth/csrf` issues a non-httpOnly `csrf_token` cookie and returns the same
value in the body. This is the seed for the enforced double-submit check above.

## What changed vs. the old doc

- No SQLite, no bcrypt, no localStorage-stored CSRF token, no length-only token
  check, and no Origin *whitelist as primary defense*.
- Error responses are generic (`writeError` maps known errors to safe messages
  and falls back to "something went wrong"); they do **not** leak raw DB driver
  errors. The previous "Duplicate entry" example was inaccurate and is removed.
