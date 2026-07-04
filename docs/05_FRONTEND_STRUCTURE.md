# Kirmya — Frontend Folder Structure

> Next.js (App Router, latest) · TypeScript · TailwindCSS · ShadCN UI · mobile-first.

## 1. Design Principles
- Clean, professional, modern-SaaS aesthetic. Mobile-first responsive.
- Feature-modular: code organized by domain feature, mirroring backend modules.
- Server Components by default; Client Components only where interactivity is needed.
- Typed API client; no `any` at module boundaries. Zod for runtime validation of API responses.

## 2. Layout

```
frontend/
├── app/                                  # App Router (routes = URLs)
│   ├── (marketing)/                      # public: landing, pricing, about
│   │   └── page.tsx
│   ├── (auth)/                           # auth route group (no app chrome)
│   │   ├── login/page.tsx
│   │   ├── register/page.tsx
│   │   ├── verify-email/page.tsx
│   │   ├── forgot-password/page.tsx
│   │   ├── reset-password/page.tsx
│   │   └── oauth/[provider]/callback/page.tsx
│   ├── (app)/                            # authenticated app shell (sidebar + topbar)
│   │   ├── layout.tsx                    # guards session, renders nav
│   │   ├── dashboard/page.tsx
│   │   ├── profile/[username]/page.tsx
│   │   ├── profile/edit/page.tsx
│   │   ├── resume/page.tsx
│   │   ├── career/page.tsx               # skill-gap, paths, salary insights
│   │   ├── jobs/page.tsx
│   │   ├── jobs/[id]/page.tsx
│   │   ├── jobs/applications/page.tsx
│   │   ├── referrals/page.tsx
│   │   ├── communities/page.tsx
│   │   ├── communities/[slug]/page.tsx
│   │   ├── mentorship/page.tsx
│   │   ├── mentorship/[mentorId]/page.tsx
│   │   ├── messages/page.tsx
│   │   ├── messages/[conversationId]/page.tsx
│   │   ├── notifications/page.tsx
│   │   ├── coach/page.tsx                # AI Career Coach chat
│   │   └── settings/page.tsx
│   ├── (admin)/admin/...                 # admin console (RBAC-gated)
│   ├── layout.tsx                        # root layout, providers, fonts
│   └── globals.css
├── src/
│   ├── components/
│   │   ├── ui/                           # ShadCN primitives (button, input, card, dialog, ...)
│   │   └── shared/                       # composed app components (Navbar, Sidebar, EmptyState, ...)
│   ├── features/                         # one folder per domain (mirrors backend)
│   │   ├── auth/                         # components, hooks, api, schemas
│   │   ├── profile/
│   │   ├── resume/
│   │   ├── career/
│   │   ├── jobs/
│   │   ├── referrals/
│   │   ├── communities/
│   │   ├── mentorship/
│   │   ├── messaging/
│   │   ├── notifications/
│   │   └── coach/
│   ├── lib/
│   │   ├── api/
│   │   │   ├── client.ts                 # fetch wrapper: base URL, auth header, refresh-on-401, CSRF
│   │   │   └── endpoints.ts              # typed endpoint map
│   │   ├── auth/                         # session helpers, token storage strategy
│   │   ├── hooks/                        # shared hooks (useToast, useDebounce, ...)
│   │   ├── utils/                        # cn(), formatters, date utils
│   │   └── config.ts                     # public env (NEXT_PUBLIC_*)
│   ├── types/                            # shared TS types / generated from OpenAPI
│   └── styles/                           # tailwind layers, tokens
├── public/                               # static assets
├── tests/
│   ├── components/                       # component tests (Vitest + Testing Library)
│   └── e2e/                              # Playwright specs
├── tailwind.config.ts
├── components.json                       # ShadCN config
├── next.config.ts
├── tsconfig.json
└── package.json
```

> Note: the current frontend has a partial `src/` and `app/` already (auth, profile, jobs, ideas pages). New work follows the `features/` + route-group structure above; existing pages are migrated as touched.

## 3. Feature folder convention

```
src/features/<feature>/
├── api.ts          # typed calls using lib/api/client
├── schemas.ts      # Zod schemas for requests/responses
├── hooks.ts        # data hooks (React Query/SWR) for the feature
├── components/     # feature-specific UI
└── types.ts
```

## 4. API Client & Auth
- `lib/api/client.ts`: attaches `Authorization: Bearer <access>`, sends CSRF token for unsafe methods, and on `401` attempts a silent refresh (`/auth/refresh`) once, then redirects to login.
- Access token in memory; refresh token in an httpOnly cookie (set by backend). This avoids XSS token theft.
- Data fetching via React Query (caching, retries, optimistic updates).

## 5. Pages required by spec → routes

| Spec page | Route |
|---|---|
| Dashboard | `/(app)/dashboard` |
| Profile | `/(app)/profile/[username]` + `/edit` |
| Resume | `/(app)/resume` |
| Community | `/(app)/communities` + `/[slug]` |
| Referral | `/(app)/referrals` |
| Mentorship | `/(app)/mentorship` |
| Jobs | `/(app)/jobs` |
| Settings | `/(app)/settings` |

## 6. Testing
- **Component tests:** Vitest + React Testing Library (`tests/components`).
- **E2E:** Playwright (`tests/e2e`) covering critical journeys: register→verify→onboard, login, request referral, apply to job, book mentorship.
- Run: `npm run test`, `npm run test:e2e`.
