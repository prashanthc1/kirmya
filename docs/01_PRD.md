# Kirmya — Product Requirements Document (PRD)

> Status: Draft v1 · Owner: Product · Last updated: 2026-06-14
> Codebase note: Kirmya is an evolution of the existing `workspace-app` ("Recession Recovery Workspace"). It is built by extending and rebranding that modular monolith, not as a greenfield rewrite.

## 1. Problem Statement

> "I lost my job. Help me get my next opportunity faster."

Job loss is financially and emotionally destabilizing, and the existing tools (large social networks, job boards) optimize for engagement and volume, not for *recovery*. Kirmya is a **career-recovery ecosystem** that compresses the time-to-next-opportunity by combining referrals, mentorship, community, skill-gap closing, and AI guidance into one focused workflow.

**Kirmya is explicitly NOT a LinkedIn clone.** It does not optimize for feed engagement or vanity metrics. It optimizes for one outcome: *the user gets hired (or transitions) faster.*

## 2. Goals & Non-Goals

### Goals (MVP)
- Reduce a job seeker's time-to-first-interview through warm referrals and AI-targeted applications.
- Make skill gaps explicit and give a concrete, time-boxed learning path to close them.
- Provide structured mentorship and supportive communities, not an open social feed.
- Give recruiters a high-signal candidate pool with context (referrals, verified skills).

### Non-Goals (MVP)
- Public social feed / influencer mechanics.
- Native mobile apps (mobile-first responsive web only for MVP).
- Payments / paid mentorship marketplace (mentorship is free in MVP; monetization later).
- Multi-tenant / white-label (architecture must *allow* it later; not built in MVP).

## 3. Target Users & Personas

| Persona | Primary need | Key actions |
|---|---|---|
| **Job Seeker** | Get the next role fast | Build profile, upload resume, request referrals, search/apply jobs, get AI guidance, find mentors |
| **Referrer** (current employee) | Help others / earn referral bonus | Offer referrals, review requests, track outcomes |
| **Mentor** | Give back, build reputation | Provide guidance, schedule sessions, share resources |
| **Recruiter** | Fill roles with high-signal candidates | Post jobs, search candidates, manage applicants |
| **Admin** | Keep platform healthy & safe | Manage users, moderate content, handle reports, view analytics |

A single account may hold multiple roles (e.g. a Referrer is also a Job Seeker). Roles are RBAC grants, not account types.

## 4. Core User Journeys

1. **Recovery onboarding:** Register → verify email → import/upload resume → AI parses & scores it → skill-gap analysis → personalized dashboard with next actions.
2. **Referral loop:** Find target company/role → request referral → employee reviews → accepts → application submitted with referral attached → outcome tracked (interview/offer/hired).
3. **Skill closing:** AI compares current skills vs target-role requirements → generates learning path → user tracks progress → re-score.
4. **Mentorship:** Browse mentors → book session → meet → rate/review.
5. **Community:** Join domain community (e.g. Facilities Management) → post/ask → get replies, polls, resources.

## 5. MVP Feature Scope

### 5.1 Identity & Auth (build first)
Email registration, login, logout, password reset, email verification, Google OAuth, LinkedIn OAuth, JWT access tokens, refresh-token rotation, MFA-ready (TOTP), RBAC.

### 5.2 Professional Profiles
Photo, headline, about, work experience, education, certifications, skills, languages, resume link, portfolio links.

### 5.3 Resume Module
Upload PDF/DOCX, parsing → structured data, scoring (formatting/keywords/ATS), version history, improvement suggestions.

### 5.4 Career Intelligence (AI)
Skill-gap analysis, career-path suggestions, salary insights, market-demand analysis, learning recommendations.

### 5.5 Referral Marketplace
Request referrals, employee review, status tracking, connect with employees.
Workflow: `Request → Review → Accepted → Application Submitted → Hired`.

### 5.6 Communities
Posts, comments, reactions, polls, tags, moderation. Domain communities: Facilities Management, Construction, Logistics, Technology, HR, Operations.

### 5.7 Jobs
Job posting, search, AI job matching, saved jobs, application tracking.

### 5.8 Mentorship
Mentor profiles, session booking, ratings, reviews.

### 5.9 Messaging
Direct messages, group chats, attachments, read receipts.

### 5.10 Notifications
Real-time: referral requests, messages, community activity, job matches.

### 5.11 AI Features
- **AI Resume Reviewer** — formatting, keywords, ATS compatibility.
- **AI Career Coach** — career planning, interview prep, skill recommendations.
- **AI Skill-Gap Engine** — current vs target-role skills → learning path.
Providers: Claude (primary) + OpenAI (fallback/specialized), behind a provider-agnostic interface.

### 5.12 Admin
User management, content moderation, report queue, analytics.

## 6. Success Metrics (North Star + supporting)
- **North Star:** median *time-to-first-interview* for active job seekers.
- Referral acceptance rate; % of applications with a warm referral attached.
- Skill-gap → learning-path completion rate.
- Resume score improvement (before/after suggestions).
- Mentorship sessions completed; mentor rating.
- 30-day retention of job seekers.

## 7. Functional Requirements (selected, MVP)
- FR-1: A user can register with email+password (Argon2id hashed) and must verify email before accessing protected features.
- FR-2: A user can authenticate via Google or LinkedIn OAuth; first OAuth login auto-provisions an account.
- FR-3: Access tokens are short-lived JWTs (~15 min); refresh tokens rotate on use and are revocable.
- FR-4: A job seeker can upload a resume (PDF/DOCX ≤ 10 MB), which is parsed and scored asynchronously.
- FR-5: A job seeker can request a referral for a specific job; the referrer can accept/decline with a note; status transitions are auditable.
- FR-6: Search returns users/jobs/communities/skills/companies with autocomplete, fuzzy matching, and filters.
- FR-7: All state-changing actions are recorded in an audit log with actor, action, target, and timestamp.

## 8. Non-Functional Requirements
- **Scale:** 10M users, 100M messages; horizontally scalable stateless API; read replicas + Redis caching.
- **Performance:** p95 API < 300 ms for cached reads; search autocomplete < 100 ms.
- **Availability:** 99.9% MVP target; multi-region ready.
- **Security:** OWASP Top 10, JWT + refresh rotation, RBAC, rate limiting, CSRF, XSS hardening, security headers, audit logging, MFA-ready.
- **Privacy:** soft deletes + data export/erasure path (GDPR-minded).
- **Observability:** OpenTelemetry traces, Prometheus metrics, structured logs, Grafana dashboards.
- **Extractability:** modules must be extractable into Identity / Jobs / Messaging / AI / Community services without rewriting business logic (interface-based module boundaries, event bus).

## 9. Release Plan (MVP roadmap summary)
See [10_MVP_ROADMAP.md](10_MVP_ROADMAP.md). Phase 1 = Identity + Profiles. Phase 2 = Resume + Jobs + Referrals. Phase 3 = AI Career Intelligence + Communities + Mentorship. Phase 4 = Messaging + Notifications + Admin + hardening.

## 10. Open Questions
- Salary insights data source (3rd-party API vs aggregated internal)?
- Resume parsing: in-house (Go + heuristics) vs LLM-only vs hybrid? (MVP: hybrid — text extraction + LLM structuring.)
- Email/SMS provider for verification & notifications (e.g. SES, Postmark, Twilio)?
