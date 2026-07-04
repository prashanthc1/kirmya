# Frontend ⇄ Backend API Coverage

_Audit of the redesigned frontend pages against the backend's registered `/api/v1`
routes. Generated 2026-06-29._

## Summary

- **Backend coverage is strong.** Every data-driven page has matching endpoints in
  one of the 16 backend modules (identity, profile, jobs, mentorship, messaging,
  community, referrals, resume, ai, notifications, settings, admin, search, …).
- **Nothing is wired yet.** The redesigned frontend is currently 100% static: `lib/`
  is empty, there is no API client, and there are **zero** `/api/v1` / `fetch` calls in
  `app/` or `components/`. `NEXT_PUBLIC_API_BASE=/api/v1` and the proxy target are still
  configured in `.env`, so the plumbing exists — the pages just don't call it.
- A few real gaps and naming mismatches are listed at the bottom.

## Per-page mapping

| Page (route) | Data it shows | Backend endpoint(s) | Status |
|---|---|---|---|
| Home `/` | Marketing | — | n/a (static) |
| About `/about` | Marketing | — | n/a |
| Pricing `/pricing` | Marketing | — | n/a |
| FAQ `/faq` | Marketing | — | n/a |
| Directions `/directions` | Marketing audience split (candidates/recruiters) | — | n/a |
| Sign In `/sign-in` | Login / register / OAuth / forgot | `POST /auth/login`, `/auth/register`, `GET/POST /auth/oauth/{provider}`, `/auth/forgot-password`, `/auth/refresh`, `/auth/csrf` | ✅ full |
| Jobs `/jobs` | Listing, matches, saved, filters | `GET /jobs`, `/jobs/matches`, `/jobs/saved`, `POST /jobs/{id}/save` | ✅ full |
| Job Detail `/jobs/detail` | One job, apply, save | `GET /jobs/{id}`, `POST /jobs/{id}/apply`, `/jobs/{id}/save` | ✅ full (route naming, see below) |
| Mentors `/mentors` | Mentor directory, reviews, availability | `GET /mentorship/mentors`, `/mentors/{id}`, `/mentors/{id}/availability`, `/mentors/{id}/reviews` | ✅ full |
| Mentorship `/mentorship` | My sessions, booking | `GET /mentorship/sessions`, `POST /mentorship/sessions`, `PATCH /sessions/{id}`, `POST /sessions/{id}/review` | ✅ full |
| Profile `/profile` | Profile view | `GET /profiles/me`, `/profiles/{id}` | ✅ full |
| Edit Profile `/profile/edit` | Basics, summary, experience, skills, availability, roles | `PUT /profiles/me`, `…/experiences`, `…/educations`, `…/certifications`, `…/skills`, `…/languages`, `…/portfolio`; mentor availability `POST /mentorship/availability` | ✅ mostly (roles tab, see gaps) |
| Settings `/settings` | Account, privacy, notifications, security/MFA, roles | `GET /settings`, `PATCH /settings/{general,privacy,notifications,security}`, `POST /auth/mfa/{setup,verify,disable}`, `/auth/change-password` | ✅ mostly (roles toggle, see gaps) |
| Inbox `/inbox` | Conversations + messages | `GET /conversations`, `/conversations/{id}/messages`, `POST …/messages`, `…/read`, `…/typing`, `GET /conversations/stream` | ✅ full (incl. realtime) |
| Coach `/coach` | AI coach chat + threads | `POST /ai/coach`, `GET /ai/coach/threads`, `/threads/{id}` | ✅ full |
| Communities `/communities` | Communities, posts, polls, reactions, comments | `GET /communities`, `/{slug}`, `/{slug}/posts`, `/{slug}/tags`, `POST /{slug}/join`, `/{slug}/posts`, `/posts/{id}/{comments,reactions,polls,report}`, `/polls/{id}/vote` | ✅ full |
| Referrals `/referrals` | Incoming/outgoing, accept/decline/outcome | `GET /referrals/{incoming,outgoing}`, `POST /referrals`, `/{id}/accept`, `/{id}/decline`, `PATCH /{id}/outcome` | ✅ full |
| Resume Check `/resume` | Resume list + score + AI review | `GET /resumes`, `/resumes/{id}/score`, `POST /ai/resume-review`, `/resumes/{id}/review` | ✅ full |
| Resume Builder `/resume/builder` | Create/edit resume + versions | `POST /resumes`, `GET/POST /resumes/{id}/versions`, `DELETE /resumes/{id}` | ✅ full |
| Career Paths `/career-paths` | Role ladder + skill gap → coach/mentor | `POST /ai/skill-gap` (+ links to `/ai/coach`, `/mentorship/mentors`) | ⚠️ partial — see gaps |
| Dashboard `/dashboard` | Role-based summary (job-seeker / recruiter / mentor) | aggregates `GET /jobs/matches`, `/jobs/applications`, `/mentorship/sessions`, `/notifications`, `/referrals/*` | ⚠️ no single endpoint — see gaps |
| Recruiter `/recruiter` | Posted roles, applicants, candidate search | `POST /jobs`, `PUT/DELETE /jobs/{id}`, `GET /jobs/{id}/applicants`, `PATCH /applications/{id}`, `GET /users/search` | ⚠️ "my posted roles" filter — see gaps |
| Admin `/admin` | Stats, users, reports moderation | `GET /admin/stats`, `/admin/users`, `/admin/users/{id}`, `PATCH /admin/users/{id}/status`, `GET /admin/reports`, `PATCH /admin/reports/{id}`, `DELETE /admin/{posts,comments}/{id}` | ✅ full |

## Gaps & mismatches to resolve before wiring

1. **No self-serve "roles" endpoint.** Both Settings and Edit Profile let a user toggle
   their job-seeker / recruiter / mentor roles, but the only role-mutation routes are
   admin-only (`POST`/`DELETE /admin/users/{id}/roles`). A user-scoped endpoint
   (e.g. `PUT /users/me/roles` or `PATCH /settings/roles`) is needed.

2. **No user dashboard/summary endpoint.** `/admin/stats` is admin-only. The user
   Dashboard must either fan out to several list endpoints client-side or get a new
   `GET /me/dashboard` (or `/me/summary`) aggregate.

3. **Career Paths needs a role-ladder data source.** `POST /ai/skill-gap` covers the
   "skills standing between you and the offer" portion, but the ladder of next roles
   with pay bands has no endpoint. Either extend the AI skill-gap response or add a
   `GET /career/paths` endpoint.

4. **Job Detail route naming.** Frontend route is static `/jobs/detail`; backend is
   `GET /jobs/{id}`. Wire it as a dynamic `app/jobs/[id]/page.tsx` to match.

5. **Recruiter "my posted roles".** Recruiter dashboard lists the recruiter's own jobs;
   confirm `GET /jobs` supports an owner filter (e.g. `?postedBy=me` / `?mine=true`),
   otherwise add one.

6. **Frontend is unwired.** No API client exists in the redesign. Re-introduce a
   `lib/api/client.ts` (auth header + httpOnly refresh-on-401 + CSRF per the contract in
   `docs/06_API_CONTRACTS.md`) and per-feature `api.ts` modules before the static pages
   can show live data.

> Source of truth for request/response shapes: `docs/06_API_CONTRACTS.md` (human) and
> `backend/docs/openapi.yaml` (machine, served at `/swagger-ui/`).
