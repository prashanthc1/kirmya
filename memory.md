# Kirmya Project Memory

Permanent source of truth, architectural standard, and operational guideline for Kirmya.com.

---

## 1. Executive Summary & Core Mission

### 1.1 Mission Statement
Kirmya is an AI-first career recovery platform designed for professionals navigating transitions, layoffs, or employment gaps. The platform is engineered to serve as a secure, premium recovery workspace where users optimize resumes, practice mock interview scenarios, secure internal referrals, and connect with peer groups or industry mentors.

### 1.2 Platform Identity
Unlike copycat networks of LinkedIn, Indeed, or Glassdoor, Kirmya maintains a focused utility design:
* **No Social Feeds**: Excludes viral posting, vanity updates, or engagement loops to eliminate visual noise and anxiety.
* **Direct Career Utility**: Built to provide practical, high-efficiency tools—resume analyzers, skill-gap training roadmaps, verified peer referral networks, and secure mentorship logs.

---

## 2. Business Model & Billing Restrictions

### 2.1 Strictly Free Operation
Kirmya is currently a completely free platform. All features are accessible to authenticated users without cost or restriction.

### 2.2 UI & Code Restrictions
* **Zero Subscriptions or Payments**: Do not implement pricing pages, payment screens, upgrade triggers, premium badges, or payment gateways (such as Stripe, PayPal, or Razorpay).
* **UI Suppression**: Hide all "Premium", "Pro", "Business", "Upgrade", or locked-feature overlays and upsell banners.
* **Platform Information Segment**: In Settings, replace the billing page with a "Platform Information" section containing:
  * Current Version & Release Notes
  * Feature Roadmap
  * Open Source Licenses & Attributions
  * Terms of Service & Privacy Policies
* **Modular Code Readiness**: Keep database models and service architectures modular so that subscription hooks can be integrated in future phases without core refactoring, but ensure all payment code remains disabled and excluded from the current codebase.

---

## 3. Technology Stack & Coding Standards

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           KIRMYA ARCHITECTURE                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   Frontend: Next.js (App Router) • React 19 • TypeScript • MUI v6        │
│                                                                         │
│                                  │ ▲                                    │
│                                  ▼ │ (REST / SSE / WebSockets)          │
│                                                                         │
│   Backend Modular Monolith: Go 1.26 • Gin Framework • ServeMux          │
│                                                                         │
│           │                      │                     │                │
│           ▼                      ▼                     ▼                │
│     PostgreSQL (pgx)       Redis Cache           OpenSearch             │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.1 Frontend Stack
* **Core**: Next.js (App Router), React 19, strict TypeScript.
* **UI Library**: **MUI v6 only** with Emotion for styling.
* **Visual Theme**: Curated Glassmorphism styling, native Light/Dark modes, responsive mobile-first grids, custom CSS, and Framer Motion micro-interactions.
* **Framework Restriction**: **Never use Tailwind CSS, Bootstrap, Chakra UI, Ant Design, or competing CSS frameworks** for any new components. Maintain custom CSS and MUI styled components for maximum design control.
* **Accessibility (WCAG AA)**: Semantic HTML5 tags, strict keyboard navigation (`tabIndex`), explicit focus indicators, and descriptive `aria-label` tags on all interactive elements.
* **Performance**: Lazy loading for modal overlays and dropdowns, code splitting, virtual lists (`react-window`), and responsive image optimization.

### 3.2 Backend Stack
* **Core**: Go 1.26 modular monolith (Clean Architecture / Domain-Driven Design per module).
* **API Delivery**: Gin Framework and standard `http.ServeMux` endpoints under `/api/v1`.
* **Database**: PostgreSQL (pgx connection pool) with raw SQL and SQL-ready ORM mapping.
* **Caching**: Redis.
* **Event Broker**: NATS (in-process event bus for local dev).
* **Search Engine**: OpenSearch (with direct SQL `ILIKE` fallback).
* **Observability**: OpenTelemetry, Prometheus metrics (`/metrics`), structured JSON logging (`slog`), request correlation IDs, and deep health check status checks.

### 3.3 Coding Principles
* **SOLID, Clean Code, DRY, KISS, YAGNI**: No over-engineered layers.
* **Zero technical debt**: Zero console errors, zero linting warnings, zero TypeScript compiler errors, and zero `TODO` comments in production code.
* **Zero Mock APIs**: All endpoints must call verified backend routes; mocks are restricted to unit test suites only.

---

## 4. Platform Modules & Features Directory

### 4.1 Authentication & Security
* **JWT Engine**: Short-lived memory access tokens paired with rotated, family-revocable, `httpOnly` secure Refresh Cookies.
* **MFA Readiness**: Fully engineered structure to support future multi-factor authentication.
* **Double-Submit CSRF Protection**: CSRF cookie matching with headers on all unsafe state-changing operations.
* **RBAC Controls**: Strict role validations (Guest, Job Seeker, Recruiter, Admin, Moderator, Support Agent, Super Admin).

### 4.2 User Profile Workspace
Consists of fifteen unified profile sections:
1. **Basic Information**: Name, photo, location, contact.
2. **Professional Headline**: Single sentence career elevator pitch.
3. **About / Bio**: Detailed professional narrative.
4. **Work Experience**: Roles, companies, dates, accomplishments.
5. **Education**: Degrees, institutions, years.
6. **Certifications**: Title, issuer, date, verification URL.
7. **Technical Skills**: Interactive tagged skill badges.
8. **Projects**: Name, description, tech stack, github link.
9. **Portfolio**: Live links, screenshots, media assets.
10. **Achievements**: Awards, publications, honors.
11. **Languages**: Spoken/written languages with proficiency scale.
12. **Resume & Documents**: Uploaded ATS PDF documents.
13. **Preferences**: Remote vs. hybrid, location radius, timezone, salary expectations.
14. **Social Links**: Verified GitHub, LinkedIn, and personal site links.
15. **Professional References**: Recommendations and contact credentials.

### 4.3 Jobs & Recruitment
* **Application Lifecycle**: Easy Apply, Cover Letter uploads, and status pipelines.
* **AI Match Score**: Scans job parameters against profile vectors to calculate match index.
* **Resume Parsing**: ATS compatibility checker scoring missing industry keywords.
* **Referral Requests**: Connects applicants with verified internal employees for warm submissions.
* **Job Attributes**: Visa sponsorship toggles, salary ranges, remote status, location radius, and job alerts.

### 4.4 Professional Networking & Messaging
* **Directory Search**: Direct search matching across users, mentors, and companies.
* **Messaging Features**: Typings indicators, read confirmations, attachments, emojis, pinned history, group messaging, and dedicated recruiter/mentor channels.
* **Anti-Spam Controls**: No automated outreach. Message connections require mutual authorization.

### 4.5 Mentorship Hub
* **Mentor Profiles**: Display areas of focus (Resume Review, Mock Interviews, Pivots) and experience.
* **Booking Workflow**: Dynamic calendar appointment scheduler with timezone synchronization.

### 4.6 Command Palette & Search
* **Trigger**: Global `Ctrl + K` or `⌘K` overlay with full backdrop blurring.
* **Targets**: Instant matching across People, Jobs, Companies, Communities, Mentors, Skills, and Courses.

### 4.7 Cookie Consent Management (CMP)
* **GDPR Dialog**: Glassmorphic, non-blocking first-visit popup with background blur.
* **Customization Settings**: Tabbed control modal mapping choices (Essential, Functional, Analytics, AI Personalization).

### 4.8 Admin Console
* RBAC-gated dashboard managing users, roles, jobs, communities, audit logs, feature flags, email templates, cookie policies, and server health.

---

## 5. Operations, Database & Infrastructure Guidelines

### 5.1 Database Conventions
* **Primary Keys**: Always use UUID v4 values.
* **Index Strategy**: Force indexes on foreign keys, email fields, username search indexes, and search tokens.
* **Migration Rules**: Forward-only migration files named `NNN_title.sql` in `backend/migrations/`. Automated rollbacks are prohibited; schema updates require a new forward migration script.

### 5.2 Caching Strategy
* Cache active user sessions, configuration settings, and heavy database lookups in Redis.
* Clear or rotate cache records on data updates to prevent stale UI states.

### 5.3 Git Workflow & Branching
* **Branches**: Use `feature/name`, `bugfix/name`, or `hotfix/name`.
* **Commits**: Follow Conventional Commits: `feat(profile): ...`, `fix(auth): ...`.
* **Definition of Done**: 100% test coverage for new components, successful production builds (`npm run build`), no ESLint warnings, and zero broken E2E integration specs.

---

## 6. AI Prompt & Generation Directives

### 6.1 Unified Generation Rules
When coding Kirmya features:
1. **Never write partial code**: Always deliver complete, production-ready, compiler-passing code. Do not output `// TODO` or `// ... remaining code` blocks.
2. **Synchronized Updates**: Modify frontend components, backend endpoints, test suites, Swagger docs, migrations, and `.env.example` in a single coordinated task.
3. **No Tailwind**: For all new visual elements, write clean Emotion code or styled components using the defined color palette tokens.
