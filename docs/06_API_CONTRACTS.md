# Kirmya — API Contracts

> REST/JSON over HTTPS · base path `/api/v1` · auth via `Authorization: Bearer <accessToken>`.
> This is the human-readable contract; the machine-readable source of truth is `backend/docs/openapi.yaml` (served at `/swagger-ui/`).

## 1. Conventions

- **Content type:** `application/json; charset=utf-8`.
- **IDs:** UUID strings.
- **Timestamps:** RFC 3339 / ISO-8601 UTC.
- **Auth:** short-lived JWT access token in the `Authorization: Bearer` header
  (CSRF-immune). The refresh token lives in an httpOnly `SameSite=Strict` cookie
  used only by `/auth/refresh` and `/auth/logout`. A `GET /auth/csrf` endpoint
  and a `csrf_token` cookie exist for an optional double-submit pattern, but the
  current client does not send `X-CSRF-Token` and the server does not yet enforce
  it — do not rely on it as a guarantee (see `CSRF_SECURITY.md`).
- **Pagination:** there is **no cursor pagination**. List endpoints return a
  named array under the resource key (e.g. `{ "jobs": [...] }`,
  `{ "notifications": [...], "unread": N }`). The admin user list uses
  `?limit=&offset=` and returns `{ "users": [...], "total", "limit", "offset" }`.
- **Errors:** consistent envelope:
  ```json
  { "error": { "code": "validation_error", "message": "Email is required", "details": { "field": "email" } } }
  ```
  HTTP codes: 400 validation, 401 unauthenticated, 403 forbidden (RBAC / unverified email), 404 not found, 409 conflict (incl. optimistic-lock), 422 unprocessable, 429 rate-limited, 500 server.
- **Success envelope:** every success body is
  `{ "success": true, "data": <payload>, "meta": { "timestamp": "<RFC3339>" } }`.
  Clients read `data` (and `error` on failure).

## 2. Identity & Auth  (`/api/v1/auth`, `/api/v1/users`)  — built first

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/auth/register` | public | Register with email + password (optional `role`, defaults `job_seeker`). Sends verification email and **auto-logs in** (201 → access token + refresh cookie, same shape as login). If `EMAIL_VERIFICATION_REQUIRED=true`, returns `{ user, verification_required: true }` with no tokens instead. |
| POST | `/auth/login` | public | Email/password login → access token + sets refresh cookie. Email verification is **not required by default** (`EMAIL_VERIFICATION_REQUIRED=false`); when enabled, unverified email → 403 `email_not_verified`. If MFA enabled and no `code` → `{ "mfa_required": true }`; re-POST `login` with `code`. |
| POST | `/auth/logout` | refresh cookie | Revoke current refresh token + clear cookie. |
| POST | `/auth/refresh` | refresh cookie | Rotate refresh token, return new access token. |
| POST | `/auth/verify-email` | public | Confirm email with token. |
| POST | `/auth/resend-verification` | public | Resend verification email. |
| POST | `/auth/forgot-password` | public | Send password-reset email. |
| POST | `/auth/reset-password` | public | Reset password with token. |
| GET | `/auth/csrf` | public | Issue a `csrf_token` cookie + value. |
| GET | `/auth/oauth/{provider}` | public | Start Google/LinkedIn OAuth → `{ url, state }`; sets httpOnly `oauth_state` cookie. |
| POST | `/auth/oauth/{provider}/callback` | public | OAuth callback. Body `{ code, state }`; `state` is verified against the cookie → access token + refresh cookie. |
| POST | `/auth/mfa/setup` | bearer | Begin TOTP enrollment → `{ otpauth_url }`. |
| POST | `/auth/mfa/verify` | bearer | Confirm TOTP enrollment with `{ code }`. |
| POST | `/auth/mfa/disable` | bearer | Disable TOTP after validating a current `{ code }` → `{ mfa_enabled: false }`. No-op if already off. |
| POST | `/auth/change-password` | bearer | Change password with `{ current_password, new_password }`. Verifies the current password, enforces complexity, and revokes **all** refresh tokens (other devices must re-authenticate) → `{ changed: true }`. |
| POST | `/auth/logout-all` | bearer | Revoke **all** of the caller's refresh tokens (sign out on every device) → `{ signed_out: true }`. |
| DELETE | `/users/me` | bearer | Deactivate (soft-close) the caller's account and revoke all sessions → `{ deactivated: true }`. Deactivated accounts can no longer log in. |
| GET | `/users/me` | bearer | Current user (id, email, full_name, roles, email_verified, **mfa_enabled**). |
| PUT | `/users/me/roles` | bearer | Reconcile the caller's self-assignable roles to `{ roles: [...] }` (subset of `job_seeker, referrer, mentor, recruiter`). Admin is not self-assignable and an existing admin role is preserved. Emits `UserRolesUpdated` → updated user. |
| GET | `/users/search?q=` | bearer | Directory search → `{ users: DirectoryUser[] }` (no email field). |
| GET | `/users/{id}` | bearer | Directory entry by id (no email field). |
| GET | `/me/dashboard` | bearer | Role-segmented per-user summary counts (jobs, referrals, mentorship, notifications) → `DashboardSummary`. |

> MFA-during-login is **not** a separate endpoint: `login` returns
> `mfa_required`, and the client re-calls `login` with `code`. There is no
> `/auth/mfa/challenge`. The directory DTO intentionally omits `email` (PII).

### Example: register
```http
POST /api/v1/auth/register
{ "email": "alex@example.com", "password": "S3cure!pass", "full_name": "Alex Doe" }
```
```json
201 {
  "data": {
    "access_token": "eyJ…",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": { "id": "…", "email": "alex@example.com", "full_name": "Alex Doe", "email_verified": false, "roles": ["job_seeker"] }
  }
}
// Set-Cookie: refresh_token=…; HttpOnly; Secure; SameSite=Strict; Path=/api/v1/auth
// When EMAIL_VERIFICATION_REQUIRED=true instead: 201 { "data": { "user": {…}, "verification_required": true } } (no tokens/cookie)
```

### Example: login
```http
POST /api/v1/auth/login
{ "email": "alex@example.com", "password": "S3cure!pass" }
```
```json
200 {
  "data": {
    "access_token": "eyJ…",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": { "id": "…", "email": "alex@example.com", "roles": ["job_seeker"] }
  }
}
// Set-Cookie: refresh_token=…; HttpOnly; Secure; SameSite=Strict; Path=/api/v1/auth
```

### Example: refresh (rotation)
```http
POST /api/v1/auth/refresh        // refresh_token cookie sent automatically
```
```json
200 { "data": { "access_token": "eyJ…", "expires_in": 900 } }
// new refresh cookie set; old token invalidated. Reuse of an old token revokes the whole family (401).
```

## 3. Profile  (`/api/v1/profiles`)

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/profiles/me` | bearer | Own full profile. |
| PUT | `/profiles/me` | bearer | Update headline/about/photo/location/website. |
| GET | `/profiles/{id}` | bearer | Public profile of another user, by user id. |
| POST/PUT/DELETE | `/profiles/me/experiences[/:id]` | bearer | CRUD work experience. |
| POST/PUT/DELETE | `/profiles/me/educations[/:id]` | bearer | CRUD education. |
| POST/PUT/DELETE | `/profiles/me/certifications[/:id]` | bearer | CRUD certifications. |
| PUT | `/profiles/me/skills` | bearer | Set skills. |
| PUT | `/profiles/me/languages` | bearer | Set languages. |
| PUT | `/profiles/me/portfolio` | bearer | Set portfolio links. |

> The profile response includes both `about` and `bio` (currently duplicative;
> `about` is the canonical field) plus `headline`, `photo_url`, `location`,
> `website`.

## 4. Resume  (`/api/v1/resumes`)

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/resumes` | bearer | Upload PDF/DOCX (multipart). Triggers async parse+score. Emits `ResumeUploaded`. |
| GET | `/resumes` | bearer | List own resumes + versions. |
| GET | `/resumes/{id}` | bearer | Resume detail + latest score. |
| DELETE | `/resumes/{id}` | bearer | Delete a resume → `{ "status": "deleted" }`. |
| POST | `/resumes/{id}/versions` | bearer | Upload a new version (multipart). |
| GET | `/resumes/{id}/versions` | bearer | Version history. |
| GET | `/resumes/{id}/score` | bearer | ATS/keyword/format score + suggestions. |
| POST | `/resumes/{id}/review` | bearer | (Re)compute the deterministic score. Returns a **`ResumeScore`** `{ overall, formatting, keywords, ats, suggestions }`. |

> Note the two distinct "review" shapes: `POST /resumes/{id}/review` returns a
> numeric **`ResumeScore`**, whereas the LLM endpoint `POST /ai/resume-review`
> (§12) returns a richer **`ResumeReview`** `{ summary, ats_score,
> keyword_feedback, formatting_feedback, strengths, improvements }`.

## 5. Career Intelligence  (`/api/v1/career`)

| Method | Path | Auth | Notes |
|---|---|---|---|
| GET | `/career/paths?from=<role>` | bearer | Career ladder reachable from `from`: rungs with pay bands (first rung marked `current`), the headline `target` rung, and `gap_skills` to reach it → `CareerPath`. Curated reference data with a generic fallback for roles outside the library. |

This complements skill-gap analysis (`POST /api/v1/ai/skill-gap`, see §12). The
ladder/pay-band data is curated reference data owned by the `career` module and
needs no database. Other career routes floated in earlier drafts
(`/career/salary`, `/career/market-demand`, `/career/learning`) are still not
implemented.

## 6. Jobs  (`/api/v1/jobs`)

| Method | Path | Auth | Role | Description |
|---|---|---|---|---|
| GET | `/jobs` | bearer | any | Search jobs (`q`, `location`, `type`) → `{ jobs: [...] }`. `posted_by=<id>` / `posted_by=me` / `mine=true` filters to a poster's own jobs (recruiter dashboard). |
| GET | `/jobs/{id}` | bearer | any | Job detail. |
| GET | `/jobs/saved` | bearer | any | The caller's saved jobs → `{ jobs: [...] }`. |
| GET | `/jobs/applications` | bearer | any | The caller's applications → `{ applications: [...] }`. |
| GET | `/jobs/matches` | bearer | any | AI-matched jobs → `{ matches: JobMatch[] }`. |
| POST | `/jobs/{id}/apply` | bearer | any | Apply. Emits `JobApplied`. |
| POST | `/jobs/{id}/save` | bearer | any | Save/unsave (toggle) → `{ saved: bool }`. |
| POST | `/jobs` | bearer | recruiter | Post a job. Emits `JobPosted`. |
| PUT | `/jobs/{id}` | bearer | recruiter | Update a job. |
| DELETE | `/jobs/{id}` | bearer | recruiter | Delete a job. |
| GET | `/jobs/{id}/applicants` | bearer | recruiter | List applicants for a job. |
| PATCH | `/applications/{id}` | bearer | recruiter | Update an application's status. |

> Seeker actions require only authentication (not a `job_seeker` role); only the
> posting/managing routes are recruiter-gated. Recommendations live at
> `/jobs/matches` (there is no `/jobs/recommended`).

## 7. Referrals  (`/api/v1/referrals`)

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/referrals` | bearer | Create referral request (job/company). Emits `ReferralRequested`. |
| GET | `/referrals/incoming` | bearer | As referrer: requests to review. |
| GET | `/referrals/outgoing` | bearer | As seeker: my requests + status. |
| POST | `/referrals/:id/accept` | bearer | Referrer accepts. Emits `ReferralAccepted`. |
| POST | `/referrals/:id/decline` | bearer | Referrer declines. |
| PATCH | `/referrals/:id/outcome` | bearer | Update outcome (interview/offer/hired). |

## 8. Communities  (`/api/v1/communities`)

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/communities` | bearer | List/search communities. |
| GET | `/communities/:slug` | bearer | Community detail. |
| POST | `/communities/:slug/join` | bearer | Join/leave (toggle). |
| GET | `/communities/:slug/posts` | bearer | Feed. Optional `?tag=` filter. |
| POST | `/communities/:slug/posts` | bearer | Create post (optional `tags[]`). Emits `CommunityPostCreated`. |
| GET | `/communities/:slug/tags` | bearer | Tags in use + post counts. |
| GET | `/communities/:slug/reports` | bearer | Open moderation reports (moderator only; 403 otherwise). |
| DELETE | `/communities/:slug/posts/:id` | bearer | Hide/remove a post (moderator only). Emits `CommunityPostHidden`. |
| GET | `/posts/:id/comments` | bearer | List comments. |
| POST | `/posts/:id/comments` | bearer | Comment. Emits `CommunityCommentAdded`. |
| POST | `/posts/:id/reactions` | bearer | React (toggle). |
| POST | `/posts/:id/polls` | bearer | Attach a poll (author only, 2–6 options; one per post). |
| POST | `/posts/:id/report` | bearer | Report a post → moderation queue. Emits `CommunityPostReported`. |
| GET | `/polls/:id` | bearer | Poll with options + vote counts. |
| POST | `/polls/:id/vote` | bearer | Vote for an option (one per user; re-voting moves it). Emits `CommunityPollVoted`. |

## 9. Mentorship  (`/api/v1/mentorship`)

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/mentorship/mentors` | bearer | Create/update the caller's mentor profile. |
| GET | `/mentorship/mentors` | bearer | Browse mentors. |
| GET | `/mentorship/mentors/:id` | bearer | Mentor profile. |
| GET | `/mentorship/mentors/:id/reviews` | bearer | List a mentor's reviews + `average_rating` and `count`. 404 if mentor unknown. |
| GET | `/mentorship/mentors/:id/availability` | bearer | A mentor's open (unbooked) slots. |
| POST | `/mentorship/availability` | bearer | Mentor opens an availability slot (`starts_at`, `ends_at`). |
| POST | `/mentorship/sessions` | bearer | Book session; optional `slot_id` consumes a slot. Emits `MentorshipBooked`. |
| GET | `/mentorship/sessions` | bearer | My sessions (as mentee and mentor). |
| PATCH | `/mentorship/sessions/:id` | bearer | Mentor advances status. Emits `SessionConfirmed`/`SessionCompleted`. |
| POST | `/mentorship/sessions/:id/review` | bearer | Rate + review (after completion). Emits `ReviewLeft`. |

## 10. Messaging  (`/api/v1/conversations`)

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/conversations` | bearer | List conversations → `{ conversations: [...] }`. |
| POST | `/conversations` | bearer | Start DM or group (`participant_ids[]`, optional `title`). |
| GET | `/conversations/stream` | bearer | **SSE** stream of message/typing/read events. |
| GET | `/conversations/{id}/messages` | bearer | Message history → `{ messages: [...] }`. |
| POST | `/conversations/{id}/messages` | bearer | Send a message. Emits `MessageSent`. |
| POST | `/conversations/{id}/read` | bearer | Mark read (receipts). |
| POST | `/conversations/{id}/typing` | bearer | Broadcast a typing indicator. |

> Real-time transport is **Server-Sent Events**, not WebSocket. There is no
> `/ws` endpoint. The client uses `fetch`-based SSE so the Bearer token rides the
> `Authorization` header.

## 11. Notifications  (`/api/v1/notifications`)

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/notifications?limit=&offset=` | bearer | List, newest/unread first (bounded; default limit 50, max 100) → `{ notifications: [...], unread: N, limit, offset }`. |
| GET | `/notifications/stream` | bearer | **SSE** stream of pushed notifications. |
| POST | `/notifications/{id}/read` | bearer | Mark read. |
| POST | `/notifications/read-all` | bearer | Mark all read. |

## 12. AI  (`/api/v1/ai`)

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/ai/resume-review` | bearer | LLM resume review → `ResumeReview` (see §4). Body `{ resume_text }`. |
| POST | `/ai/skill-gap` | bearer | Skill-gap engine → `SkillGap`. Body `{ current_role, target_role, current_skills[] }`. |
| POST | `/ai/coach` | bearer | Career coach chat. Body `{ message, thread_id? }` → `{ thread_id, reply }`. |
| GET | `/ai/coach/threads` | bearer | List coach threads → `{ threads: [...] }`. |
| GET | `/ai/coach/threads/{id}` | bearer | Thread detail incl. `messages[]`. |

## 13. Admin  (`/api/v1/admin`)  — RBAC: admin

| Method | Path | Description |
|---|---|---|
| GET | `/admin/stats` | Platform analytics overview (users/jobs/referrals/communities/reports). |
| GET | `/admin/users?q=&status=&role=&limit=&offset=` | List/search/filter users (paginated). |
| GET | `/admin/users/:id` | User detail (roles, status, profile headline). |
| PATCH | `/admin/users/:id/status` | Set account status (`active`/`suspended`/`deactivated`). |
| POST | `/admin/users/:id/roles` | Grant an RBAC role (`{role}`). |
| DELETE | `/admin/users/:id/roles/:role` | Revoke an RBAC role. |
| DELETE | `/admin/posts/:id` | Remove a community post (moderation). |
| DELETE | `/admin/comments/:id` | Remove a community comment (moderation). |
| GET | `/admin/reports?status=` | Moderation report queue. |
| PATCH | `/admin/reports/:id` | Triage a report (`{status: resolved\|dismissed, action_taken}`). |

Report filing is open to any authenticated user:

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/reports` | bearer | File a content report (`{target_type, target_id, reason}`). |

All admin mutations are written to `audit_logs` with the acting admin's id. An
admin cannot suspend their own account or revoke their own admin role.

## 14. Search  (`/api/v1/search`)

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/search?q=&type=&limit=` | bearer | Fuzzy full-text search over users, jobs, communities, skills → `{ results: SearchHit[], engine }`. `type` may be repeated or comma-separated (`user`/`job`/`community`/`skill`); omit for all. |
| GET | `/search/autocomplete?q=&limit=` | bearer | Prefix autocomplete (title match) → `{ results: SearchHit[] }`. |

Backed by **OpenSearch** when `OPENSEARCH_URL` is set; otherwise the API falls
back to PostgreSQL `ILIKE` queries so search always works. The index is
backfilled on startup and kept fresh by subscribing to `UserRegistered`,
`ProfileUpdated`, and `JobPosted` events. Each result is
`{type, ref_id, title, subtitle, url, score}`; the `/search` response also
reports which engine served it (`opensearch` | `database`).

## 15. Settings  (`/api/v1/settings`)

Per-user preferences spanning general, privacy, notification and security
sections. The row is created lazily on first read (defaults applied), so `GET`
always returns a complete object. Each `PATCH` replaces one whole section and
returns the full updated settings object. Writes are optimistic-locked on
`version` (stale write → 409 `conflict`).

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/settings` | bearer | Full settings: `{ general, privacy, notifications, security, version, updated_at }`. Materialises defaults on first access. |
| PATCH | `/settings/general` | bearer | Replace general section: `{ language, timezone, theme, email_digest }`. `theme∈{light,dark,system}`, `email_digest∈{off,daily,weekly}`. |
| PATCH | `/settings/privacy` | bearer | Replace privacy section: `{ profile_visibility, show_email, discoverable, allow_messages }`. `profile_visibility∈{public,network,private}`, `allow_messages∈{everyone,network,none}`. |
| PATCH | `/settings/notifications` | bearer | Replace notification toggles: `{ email_jobs, email_mentorship, email_messages, email_referrals, inapp_jobs, inapp_mentorship, inapp_messages, inapp_referrals }` (all booleans). |
| PATCH | `/settings/security` | bearer | Replace security preferences: `{ login_alerts }` (boolean). Password and MFA are managed via the `/auth/*` endpoints in §2. |

> Section-level events are published on the bus: `SettingsUpdated`,
> `PrivacySettingsChanged`, `NotificationSettingsChanged`.
>
> **Enforcement:** `discoverable` filters people-search results; `profile_visibility:
> private` hides `GET /profiles/{id}` from other users; `allow_messages: none` blocks
> starting new conversations; and the in-app notification toggles suppress matching
> notifications. (`network` scopes await a connections graph.)

## 16. Rate Limits (defaults)
- Auth endpoints: 10 req/min/IP (login/register/reset stricter: 5/min).
- General authenticated: 120 req/min/user.
- Search: 60 req/min/user.
- 429 returns `Retry-After`.
