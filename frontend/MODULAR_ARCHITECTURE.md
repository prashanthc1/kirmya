# Modular Monolith Architecture

## Overview

The frontend has been refactored into a **modular monolith** architecture:
- Single application (not microservices)
- Organized into independent feature modules
- Clear module boundaries and dependencies
- Shared core utilities and components
- Easy to scale and maintain

## Architecture Principles

### 1. Module Independence
- Each module is self-contained
- Modules don't directly import from other modules
- Communication through well-defined interfaces

### 2. Clear Boundaries
- Each module owns its components, hooks, services, types
- Public API defined in `index.ts`
- Internal implementation details hidden

### 3. Reduced Coupling
- Minimize module-to-module dependencies
- Share only through core/common modules
- Use dependency injection where needed

### 4. Scalability
- Easy to add new modules
- New features don't require restructuring
- Team members can work on separate modules independently

## Directory Structure

```
frontend/
│
├── app/                              # Next.js app directory (entry points)
│   ├── layout.tsx                   # Root layout
│   ├── page.tsx                     # Home page
│   ├── globals.css                  # Global styles
│   └── (auth)/                      # Auth-related pages
│
├── src/
│   │
│   ├── core/                        # Shared across all modules
│   │   ├── api/                     # API client configuration
│   │   │   ├── client.ts            # Axios/fetch config
│   │   │   └── interceptors.ts      # Request/response handlers
│   │   │
│   │   ├── providers/               # Context providers
│   │   │   ├── AppProvider.tsx      # Root provider wrapper
│   │   │   ├── ThemeProvider.tsx    # MUI theme provider
│   │   │   └── index.ts
│   │   │
│   │   ├── hooks/                   # Global reusable hooks
│   │   │   ├── useAuth.ts           # Authentication hook
│   │   │   ├── useNotifications.ts  # Notifications hook
│   │   │   ├── useApi.ts            # API calling hook
│   │   │   └── index.ts
│   │   │
│   │   ├── types/                   # Global types
│   │   │   ├── api.ts               # API response types
│   │   │   ├── common.ts            # Common types
│   │   │   └── index.ts
│   │   │
│   │   ├── utils/                   # Utility functions
│   │   │   ├── formatters.ts        # Date, currency, etc
│   │   │   ├── validators.ts        # Form validation
│   │   │   ├── constants.ts         # App constants
│   │   │   └── index.ts
│   │   │
│   │   ├── components/              # Shared UI components
│   │   │   ├── Layout/
│   │   │   │   ├── Header.tsx
│   │   │   │   ├── Footer.tsx
│   │   │   │   └── index.ts
│   │   │   ├── ErrorBoundary/
│   │   │   ├── LoadingSpinner/
│   │   │   └── index.ts
│   │   │
│   │   ├── config/                  # Core config
│   │   │   ├── theme.ts             # MUI theme
│   │   │   └── constants.ts         # App constants
│   │   │
│   │   └── index.ts                 # Core barrel export
│   │
│   ├── modules/                     # Feature modules (isolated)
│   │
│   ├── modules/auth/                # Authentication Module
│   │   ├── components/              # Auth-specific components
│   │   │   ├── LoginForm.tsx
│   │   │   ├── RegisterForm.tsx
│   │   │   ├── LinkedInButton.tsx
│   │   │   └── index.ts
│   │   │
│   │   ├── pages/                   # Auth pages
│   │   │   ├── LoginPage.tsx
│   │   │   ├── RegisterPage.tsx
│   │   │   └── index.ts
│   │   │
│   │   ├── hooks/                   # Auth-specific hooks
│   │   │   ├── useLogin.ts
│   │   │   ├── useRegister.ts
│   │   │   ├── useLinkedInAuth.ts
│   │   │   └── index.ts
│   │   │
│   │   ├── services/                # Auth API calls
│   │   │   ├── authService.ts
│   │   │   ├── linkedinService.ts
│   │   │   └── index.ts
│   │   │
│   │   ├── types/                   # Auth types
│   │   │   ├── auth.types.ts
│   │   │   └── index.ts
│   │   │
│   │   ├── context/                 # Auth context (optional)
│   │   │   ├── AuthContext.tsx
│   │   │   └── index.ts
│   │   │
│   │   ├── README.md                # Module documentation
│   │   └── index.ts                 # Module barrel export
│   │
│   ├── modules/jobs/                # Jobs Module
│   │   ├── components/
│   │   │   ├── JobCard.tsx
│   │   │   ├── JobList.tsx
│   │   │   ├── JobFilters.tsx
│   │   │   ├── PostJobForm.tsx
│   │   │   ├── ApplyJobForm.tsx
│   │   │   └── index.ts
│   │   │
│   │   ├── pages/
│   │   │   ├── JobsPage.tsx
│   │   │   ├── JobDetailPage.tsx
│   │   │   └── index.ts
│   │   │
│   │   ├── hooks/
│   │   │   ├── useJobs.ts           # Fetch jobs list
│   │   │   ├── useJobDetail.ts      # Fetch single job
│   │   │   ├── usePostJob.ts        # Create job
│   │   │   ├── useApplyJob.ts       # Apply to job
│   │   │   └── index.ts
│   │   │
│   │   ├── services/
│   │   │   ├── jobsService.ts       # Jobs API calls
│   │   │   ├── applicationsService.ts
│   │   │   └── index.ts
│   │   │
│   │   ├── types/
│   │   │   ├── job.types.ts
│   │   │   ├── application.types.ts
│   │   │   └── index.ts
│   │   │
│   │   ├── constants/
│   │   │   ├── jobTypes.ts
│   │   │   └── index.ts
│   │   │
│   │   ├── README.md
│   │   └── index.ts
│   │
│   ├── modules/ideas/               # Ideas Module (similar structure)
│   ├── modules/profile/             # Profile Module
│   ├── modules/notifications/       # Notifications Module
│   └── modules/linkedin/            # LinkedIn Integration Module
│
├── lib/                             # Legacy (can be migrated to core/)
│   └── (existing files for migration)
│
├── config/
│   └── theme.ts
│
├── package.json
├── tsconfig.json
├── next.config.js
└── README.md
```

## Module Template

Each module follows this structure:

```
module/
├── components/          # UI components (most specific to least)
├── pages/              # Page components
├── hooks/              # Custom hooks
├── services/           # API calls & business logic
├── types/              # TypeScript types
├── constants/          # Module constants
├── context/            # State management (optional)
├── README.md           # Module documentation
└── index.ts            # Public API (barrel export)
```

## Module Communication

### ✅ Allowed (Good)
```typescript
// Within same module
import { useJobs } from './hooks';
import { Job } from './types';

// From core
import { useAuth, useNotifications } from '@/core';
import { API_BASE_URL } from '@/core/config';
```

### ❌ Not Allowed (Bad)
```typescript
// Jobs module importing from Ideas module directly
import { useIdeas } from '@/modules/ideas';  // NO!

// Accessing internal files
import { jobsService } from '@/modules/jobs/services/jobsService';  // NO!
// Use public API instead:
import { jobsService } from '@/modules/jobs';  // OK!
```

## Module Index (Barrel Export)

Each module exports its public API through `index.ts`:

```typescript
// modules/jobs/index.ts
export * from './components';
export * from './pages';
export * from './hooks';
export * from './types';
export * from './services';
```

This ensures:
- Clean imports: `import { JobCard } from '@/modules/jobs'`
- Hidden internals: Users can't access internal services directly
- Easy refactoring: Move files without breaking imports

## Core Module

The `core` module contains:
- **API Client**: Centralized HTTP requests
- **Providers**: Auth, Theme, Notifications
- **Hooks**: Global utilities (useAuth, useNotifications)
- **Types**: Shared data structures
- **Utils**: Formatting, validation, constants
- **Components**: Layout, ErrorBoundary, LoadingSpinner

## Benefits

### Scalability
✅ Easy to add new features (create new module)  
✅ Independent development (teams work on separate modules)  
✅ Clear code organization  

### Maintainability
✅ Find related code easily (everything in one module)  
✅ Understand dependencies at a glance  
✅ Easier testing (modules are isolated)  

### Reusability
✅ Share common code through core/  
✅ Publish modules independently (future microservices)  
✅ Module dependencies are explicit  

### Testability
✅ Mock entire module (unit tests)  
✅ Test modules independently  
✅ Clear test boundaries  

## Path Aliases

TypeScript paths configured for clean imports:

```json
{
  "compilerOptions": {
    "paths": {
      "@/core": ["src/core/index.ts"],
      "@/modules/*": ["src/modules/*/index.ts"],
      "@/modules/*/components": ["src/modules/*/components/index.ts"],
      "@/config/*": ["config/*"]
    }
  }
}
```

## Adding a New Module

1. Create module folder: `src/modules/newFeature/`
2. Create standard folders:
   ```
   components/
   pages/
   hooks/
   services/
   types/
   constants/
   ```
3. Create `index.ts` with barrel exports
4. Create `README.md` documenting the module
5. Implement features within module
6. Export through `index.ts` for use in other modules

## Module Dependencies

### auth module
- Depends on: core
- Used by: all other modules (for user context)

### jobs module
- Depends on: core
- Used by: profile module (for user's applications)

### ideas module
- Depends on: core
- Used by: none

### profile module
- Depends on: core, jobs (for applications list)
- Used by: none

### notifications module
- Depends on: core
- Used by: all modules (for user feedback)

### linkedin module
- Depends on: core, auth
- Used by: auth module (for social login)

## Migration Roadmap

- [x] Create core module structure
- [ ] Migrate auth module
- [ ] Migrate jobs module
- [ ] Migrate ideas module
- [ ] Migrate profile module
- [ ] Create notifications module
- [ ] Create linkedin module
- [ ] Update app/ routes to use new structure
- [ ] Remove old lib/ directory
- [ ] Update documentation

## Deployment

The modular monolith structure doesn't affect deployment:
- Still builds as single Next.js application
- Code splitting works as expected
- Bundle size optimizations remain the same
- Deploy as usual: `npm run build && npm start`

## Future: Micro Modules

If needed in the future, individual modules can be extracted:
- Publish to npm (private registry)
- Share between projects
- Evolve into microservices
- Current structure makes this easy

## Testing Strategy

Each module should include tests:

```
modules/jobs/
├── __tests__/
│   ├── hooks.test.ts
│   ├── services.test.ts
│   └── components.test.ts
├── components/
├── hooks/
└── services/
```

## Documentation

Every module has `README.md`:
```markdown
# Jobs Module

## Overview
Handles job listing, posting, and applications.

## Public API
- `useJobs()` - Fetch jobs list
- `JobCard` - Display job card
- `JobsPage` - Main jobs page

## Dependencies
- core
- profile (imports applications)

## Adding Features
[Instructions for developers]
```

## Conclusion

This modular monolith architecture provides:
- ✅ Clear code organization
- ✅ Independent modules
- ✅ Explicit dependencies
- ✅ Easy to scale
- ✅ Easy to test
- ✅ Easy to refactor
- ✅ Foundation for future microservices

It's the sweet spot between monolith simplicity and microservices complexity.
