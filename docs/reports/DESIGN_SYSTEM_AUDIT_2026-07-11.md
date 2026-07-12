# Kirmya Design System Audit — Theme Migration Status

**Date:** 2026-07-11
**Scope:** `frontend/` — all routes, shared components, profile module
**Method:** Git-history cross-reference (which files the redesign commits actually touched) + static grep for legacy signals (MUI imports, hardcoded hex colors, inline `style={{}}`, orphaned theme files) across all 29 routes and every shared/profile component.

## 1. Ground truth: what "the new design system" is

`frontend/app/globals.css` is the actual source of truth — Tailwind v4 `@theme` block mapping semantic CSS variables (`--background`, `--primary`, `--card`, `--radius: 0.75rem`, etc.) with a `.dark` class override, plus "Kirmya premium" glassmorphism utility classes (`.glass-nav`, `.glass-card`, `.glow-orb`, `.animate-shimmer`). Primary color is blue (`#2563eb` light / `#60a5fa` dark). `components.json` is configured for shadcn (`new-york` style, `@/components/ui` alias) but **no `components/ui/*` primitives have ever been generated** — there is no shared Button/Input/Card/Dialog kit, which is a root cause of the drift below.

Two commits are the actual redesign:
- `5809fed` "Redesign Kirmya platform into a premium, theme-aware Career Operating System" (2026-07-06)
- `6dd4325` "Redesign and build Kirmya Profile Workspace with 15 sections..." (2026-07-06)

Everything those commits touched is the canonical new design language. Everything else is untouched legacy.

## 2. Pages reviewed (29 routes + shared shell + profile module)

All 29 `app/**/page.tsx` routes, `components/shared/*`, and `components/profile/*` were checked for: SiteNav/SiteFooter shell usage, MUI imports, hardcoded hex colors (vs. CSS-variable tokens), inline `style={{}}`, and orphaned/dead theming code.

### Already on the new design system (compliant)
| Page/Component | Notes |
|---|---|
| `app/page.tsx` (home) | Redesigned, token-driven |
| `app/dashboard/page.tsx` | Redesigned |
| `app/coach/page.tsx` | Redesigned |
| `app/mentorship/page.tsx` | Redesigned |
| `app/jobs/page.tsx` | Redesigned |
| `app/referrals/page.tsx` | Redesigned |
| `app/profile/page.tsx`, `app/profile/[id]/page.tsx`, `app/profile/edit/page.tsx` | Redesigned |
| `app/sign-in/page.tsx`, `app/sign-up/page.tsx`, `app/forgot-password/page.tsx` | Redesigned, 0 hardcoded colors |
| `components/shared/SiteNav.tsx`, `SiteFooter.tsx` | Redesigned shell, used by 24/29 pages |
| `components/profile/*` (ProfileWorkspace, ProfileCenterWorkspace, ProfileLeftSidebar, ProfileOnboarding, ProfileAiAssistant, mockData, types) | Built fresh in the redesign, 0 hardcoded colors, token-driven — best reference implementation in the codebase |

### Still on legacy styling (never touched by the redesign)
Ranked by hardcoded-hex-color count (each one is a spot bypassing the token system):

| Page | Hardcoded hex colors | Notes |
|---|---:|---|
| `app/recruiter/page.tsx` | 128 | Worst offender |
| `app/settings/page.tsx` | 121 | |
| `app/directions/page.tsx` | 92 | |
| `app/communities/page.tsx` | 81 | |
| `app/resume/builder/page.tsx` | 74 | |
| `app/career-paths/page.tsx` | 73 | |
| `app/jobs/detail/page.tsx` | 69 | |
| `app/inbox/page.tsx` | 62 | Messaging UI |
| `app/mentors/page.tsx` | 57 | |
| `app/faq/page.tsx` | 55 | |
| `app/about/page.tsx` | 52 | |
| `app/resume/page.tsx` | 22 | |
| `app/admin/page.tsx` | 21 | |
| `app/verify-email/page.tsx` | 17 | Also missing SiteNav/SiteFooter shell |
| `app/reset-password/page.tsx` | 17 | Also missing SiteNav/SiteFooter shell |
| `app/pricing/page.tsx` | 8 | |
| `app/search/page.tsx` | 0 hex, but **still built on `@mui/material`** — the only page still using MUI components directly |

That's **17 of 29 pages (59%)** still on legacy styling — the redesign covers less than half the app. No employer-specific or separate admin-suite pages exist beyond `app/admin/page.tsx` and `app/recruiter/page.tsx`.

## 3. Inconsistent / duplicate component layers

The app currently runs **two parallel UI systems at once**:

1. **Legacy MUI layer** (still wired into the global tree):
   - `components/shared/ThemeProvider.tsx` wraps the *entire app* in `MuiThemeProvider` + `CssBaseline`, even though only one page (`search`) still renders MUI components. CssBaseline applies a global browser-style reset alongside Tailwind's own reset — redundant and a latent source of subtle spacing/typography drift.
   - `components/shared/MuiModal.tsx`, `components/shared/Notifications.tsx` — MUI `Dialog`/`Snackbar`-based, not token-driven, don't match the glass-card/radius/shadow language used everywhere else.
   - `app/search/page.tsx` — the only page still rendering raw MUI components (`TextField`, `Card`, etc.).
   - `lib/theme.ts` — a **fully orphaned, unused file** (confirmed via grep — nothing imports it) defining a *third, completely different* color palette ("Sunset Coral" `#D66838` / "Eucalyptus Green" `#37614D`, "Organic Premium Tech") that matches neither `globals.css` nor the inline MUI theme in `ThemeProvider.tsx`. Dead code, but risky to leave — anyone touching it would restyle against the wrong palette.

2. **No shared component kit**: despite `components.json` declaring shadcn with the `@/components/ui` alias, that directory doesn't exist. Every page hand-rolls its own buttons, cards, inputs, badges — which is *why* radius/shadow/spacing drift page to page even among pages that don't use raw hex (Tailwind utility classes typed by hand, not a shared primitive).

## 4. Concrete cross-cutting bugs found

- **Font mismatch**: `app/layout.tsx` loads `Public_Sans` via `next/font/google` and exposes `--font-public-sans`, but `app/globals.css` sets `font-family: "Geist", "Inter", "SF Pro Display", ...` on `html, body` — none of which are ever loaded. `--font-public-sans` is never referenced. Net effect: the site renders on system-font fallback everywhere, not the intended typeface. One-line fix once a direction is picked (load Geist/Inter, or point globals.css at `--font-public-sans`).
- Dark mode is implemented correctly at the token layer (`.dark` class + CSS vars, driven by `next-themes`), and the MUI shim in `ThemeProvider.tsx` does sync its palette to `resolvedTheme` — so dark mode itself isn't broken. But every hardcoded hex color in the 17 legacy pages above **does not respond to theme changes at all**, since a literal `#fff`/`#111` etc. ignores the `.dark` class entirely. This is the most likely source of "breaks when switching themes" bugs the audit brief called out.

## 5. Scale reality check

This is not a small cleanup: ~1,000 individual hardcoded-color occurrences across 17 pages, one of which (`components/profile/ProfileCenterWorkspace.tsx`, already compliant) is 93KB/1,479 lines on its own — indicating the legacy pages are of similar size. A safe migration means judgment per color (some hex values may be intentional one-off illustration/chart/brand-icon colors that should *not* become semantic tokens), not a blind find-and-replace.

**Sandbox verification limits** (this environment, not the app): `next build` and `vitest` cannot run here (native binary crashes / missing Linux bindings) and the Go toolchain isn't installed. `tsc --noEmit` and `eslint` do work reliably and will be used to verify every edit, but final confirmation of `npm run build`, `npm run test`, and visual QA needs to happen on your machine.
