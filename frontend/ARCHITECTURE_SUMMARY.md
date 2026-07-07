# Frontend Architecture Summary

## Modular Monolith: The Best of Both Worlds

This frontend follows a **modular monolith** architecture pattern:

- Single application (not microservices)
- Organized into independent, self-contained modules
- Clear boundaries and explicit dependencies
- Easy to test, maintain, and scale

## Why Modular Monolith?

| Aspect      | Monolith | Microservices | Modular Monolith |
| ----------- | -------- | ------------- | ---------------- |
| Simplicity  | ⭐⭐⭐   | ⭐            | ⭐⭐⭐           |
| Scalability | ⭐       | ⭐⭐⭐        | ⭐⭐             |
| Testing     | ⭐⭐     | ⭐⭐⭐        | ⭐⭐⭐           |
| Development | ⭐⭐⭐   | ⭐⭐          | ⭐⭐⭐           |
| Maintenance | ⭐⭐     | ⭐            | ⭐⭐⭐           |
| Deployment  | ⭐⭐⭐   | ⭐            | ⭐⭐⭐           |

## Project Structure

```
frontend/
│
├── app/                              # Next.js app routes
│   ├── layout.tsx                   # Root layout
│   ├── page.tsx                     # Home page
│   ├── (auth)/                      # Auth routes
│   ├── jobs/                        # Job routes
│   ├── ideas/                       # Ideas routes
│   └── globals.css                  # Global styles
│
├── src/
│   │
│   ├── core/                        # Shared utilities & configuration
│   │   ├── api/
│   │   │   └── client.ts            # Centralized HTTP client
│   │   ├── config/
│   │   │   ├── theme.ts             # MUI Material 3 theme
│   │   │   └── index.ts             # App config & constants
│   │   ├── types/
│   │   │   └── index.ts             # Global types (User, Error, etc)
│   │   └── index.ts                 # Core exports
│   │
│   └── modules/                     # Feature modules
│       │
│       ├── auth/                    # Authentication
│       │   ├── components/
│       │   ├── context/
│       │   ├── hooks/
│       │   ├── services/
│       │   ├── types/
│       │   ├── README.md
│       │   └── index.ts
│       │
│       ├── jobs/                    # Job management
│       │   ├── components/
│       │   ├── hooks/
│       │   ├── services/
│       │   ├── types/
│       │   ├── README.md
│       │   └── index.ts
│       │
│       ├── ideas/                   # Ideas & collaboration
│       ├── profile/                 # User profile
│       ├── notifications/           # Notifications
│       └── linkedin/                # LinkedIn integration
│
├── config/
│   └── theme.ts                     # (Moved to src/core/config/)
│
├── MODULAR_ARCHITECTURE.md          # Architecture guide
├── MIGRATION_GUIDE.md               # How to migrate existing code
├── ARCHITECTURE_SUMMARY.md          # This file
├── MATERIAL3_MIGRATION.md           # Material UI 3 guide
├── package.json
├── tsconfig.json
└── next.config.js
```

## Core Module

The `core` module provides shared functionality:

### Config

```typescript
import { APP_CONFIG, API_ENDPOINTS } from "@/core/config";

// Usage
const apiUrl = `${APP_CONFIG.API_BASE_URL}${API_ENDPOINTS.JOBS_LIST}`;
```

### API Client

```typescript
import { apiClient } from "@/core/api/client";

// Usage
const data = await apiClient.get("/endpoint");
const posted = await apiClient.post("/endpoint", payload);
```

### Types

```typescript
import { User, AppError, ErrorType, Notification } from "@/core/types";

// All global types defined here
```

### No Direct Imports!

Core module does NOT export:

- ❌ Providers (use AppProvider in app/layout.tsx)
- ❌ Context (access via hooks from modules)
- ❌ Hooks (modules export their own)

## Feature Modules

Each module is self-contained:

### Module Structure

```
module/
├── components/     # UI components
├── hooks/         # Custom hooks (data fetching, state)
├── services/      # API calls
├── types/         # Module types
├── context/       # Local state management
├── constants/     # Module constants
├── README.md      # Documentation
└── index.ts       # Public API (barrel export)
```

### Module Communication

✅ **Within Module**

```typescript
import { useJobs } from "./hooks";
import { Job } from "./types";
import { jobsService } from "./services";
```

✅ **From Core**

```typescript
import { apiClient, User, useAuth } from "@/core";
```

✅ **From Module to Module** (via public API)

```typescript
import { jobsService, Job } from "@/modules/jobs";
```

❌ **NOT ALLOWED**

```typescript
// Direct internal imports
import { jobsService } from "@/modules/jobs/services/jobsService";
```

## Data Flow

### Typical Feature Request

```
Component (React)
    ↓
Hook (useJobs, useJobDetail)
    ↓
Service (jobsService)
    ↓
Core API Client (apiClient)
    ↓
Backend API
```

### Example: Load Jobs

```typescript
// In JobsPage component
import { useJobs } from "@/modules/jobs";

function JobsPage() {
  // Hook handles fetching, loading, errors
  const { jobs, loading, error } = useJobs({
    keyword: "engineer",
    location: "remote",
  });

  // useJobs internally uses jobsService
  // jobsService uses core apiClient
}
```

## Module Dependencies

```
notification ← core
│
auth ← core
│   ├→ components (LoginForm, etc)
│   └→ pages (LoginPage)
│
jobs ← core
│   ├→ components (JobCard, etc)
│   └→ hooks (useJobs, etc)
│
ideas ← core
│
profile ← core + jobs (for applications)
│
linkedin ← core + auth (OAuth)
```

## Design System: Material 3

All UI uses Material UI components:

```typescript
import {
  Box, // Layout
  Container, // Width container
  Button, // Buttons
  TextField, // Forms
  Card, // Card layout
  Typography, // Text
} from "@mui/material";

import material3Theme from "@/core/config/theme";
```

### Key Colors (Material 3)

- **Primary**: Indigo (#6366f1)
- **Secondary**: Violet (#8b5cf6)
- **Success**: Emerald (#10b981)
- **Error**: Red (#ef4444)

## Error Handling

Comprehensive error handling system:

```typescript
import { ErrorType, ErrorSeverity, AppError } from "@/core/types";

// Errors are classified automatically
// Network errors → automatic retry
// Auth errors → redirect to login
// Validation errors → show to user
```

## Authentication

Authentication managed in auth module:

```typescript
import { useAuth } from "@/core/providers"; // Global hook
import { authService } from "@/modules/auth"; // Service

// Usage
const { user, token, login, logout } = useAuth();
```

## Performance Optimizations

### Code Splitting

- Each route loads only necessary modules
- Unused modules excluded from bundle

### Lazy Loading

```typescript
const JobsPage = dynamic(() =>
  import('@/modules/jobs/pages/JobsPage'),
  { loading: () => <Spinner /> }
);
```

### Caching

- API responses cached in React Query (future)
- Service workers for offline support (future)

## Testing Strategy

Each module includes tests:

```
modules/jobs/__tests__/
├── services.test.ts    # API calls
├── hooks.test.ts       # React hooks
└── components.test.ts  # Components
```

Run tests:

```bash
npm test                 # All tests
npm test -- jobs        # Single module
```

## TypeScript Configuration

Path aliases for clean imports:

```json
{
  "paths": {
    "@/core": ["src/core/index.ts"],
    "@/modules/*": ["src/modules/*/index.ts"],
    "@/*": ["app/*", "src/*"]
  }
}
```

## Adding New Features

1. **Create module folder**

   ```bash
   mkdir src/modules/newfeature
   ```

2. **Create structure**

   ```
   components/
   hooks/
   services/
   types/
   index.ts
   ```

3. **Implement features**
   - Keep module independent
   - Export public API in index.ts

4. **Test thoroughly**
   - Unit tests for services
   - Component tests
   - Integration tests

5. **Document**
   - Create README.md
   - Add JSDoc comments

## Deployment

No changes to deployment:

- Single Next.js build
- Code splitting works automatically
- Deploy as usual: `npm run build && npm start`

## Future Evolution

This architecture supports evolution:

### Phase 1: Single App (Current)

- Modular monolith in one repo
- Clear module boundaries

### Phase 2: Multiple Apps (Future)

- Extract modules to npm packages
- Share across projects

### Phase 3: Microservices (Future)

- Convert modules to separate services
- Use same types/interfaces

## Key Benefits

✅ **Scalability**: Easy to add new features  
✅ **Maintainability**: Related code organized together  
✅ **Testability**: Test modules independently  
✅ **Reusability**: Share code between modules  
✅ **Clarity**: Explicit dependencies  
✅ **Flexibility**: Evolve toward microservices

## Quick Reference

### Import from Core

```typescript
import {
  apiClient, // API client
  APP_CONFIG, // Configuration
  API_ENDPOINTS, // API paths
  User, // Types
} from "@/core";
```

### Import from Modules

```typescript
import {
  useJobs, // Hooks
  jobsService, // Services
  JobCard, // Components
  Job, // Types
} from "@/modules/jobs";
```

### Import from MUI

```typescript
import { Box, Container, Button, TextField, Card } from "@mui/material";
```

## Resources

- [MODULAR_ARCHITECTURE.md](./MODULAR_ARCHITECTURE.md) - Detailed architecture guide
- [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) - How to migrate existing code
- [MATERIAL3_MIGRATION.md](./MATERIAL3_MIGRATION.md) - Material UI 3 guide
- [Each Module README.md](./src/modules/) - Module-specific documentation

## Conclusion

This modular monolith architecture provides:

- Professional code organization
- Clear separation of concerns
- Easy to test and maintain
- Foundation for future microservices
- Best practices from both monoliths and microservices

**Happy coding! 🚀**
