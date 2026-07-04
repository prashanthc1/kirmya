# Migration Guide: Modular Monolith Architecture

This guide explains how to migrate existing code to the new modular monolith architecture.

## Phase 1: Understand the Structure

The project is now organized as:

```
src/
├── core/          # Shared across all modules
├── modules/       # Feature-specific modules
│   ├── auth/
│   ├── jobs/
│   ├── ideas/
│   ├── profile/
│   ├── notifications/
│   └── linkedin/
```

## Phase 2: Move Files

### Example: Move authentication code

**Before:**
```
lib/auth-context.tsx
lib/auth.ts
components/Navigation.tsx
components/LinkedInLoginButton.tsx
```

**After:**
```
src/modules/auth/
├── context/AuthContext.tsx
├── services/index.ts
├── components/
│   ├── Navigation.tsx
│   └── LinkedInLoginButton.tsx
├── types/index.ts
└── index.ts (exports)
```

### Gradual Migration Steps

1. Create target module folder
2. Move related files into module
3. Update imports in moved files
4. Create index.ts with public API
5. Update imports in consuming code
6. Test thoroughly

## Phase 3: Update Imports

### Old Imports (to be removed)

```typescript
import { useAuth } from '@/lib/auth-context';
import { api } from '@/lib/api';
import { Navigation } from '@/components/Navigation';
```

### New Imports

```typescript
import { useAuth } from '@/core/providers';  // From core AppProvider
import { apiClient } from '@/core/api/client';  // Core API client
import { Navigation } from '@/core/components';  // Shared component

// Or from specific modules:
import { LinkedInLoginButton } from '@/modules/auth';
import { jobsService } from '@/modules/jobs';
```

## Phase 4: Module-Specific Patterns

### Creating a Service in a Module

```typescript
// src/modules/jobs/services/index.ts
import { apiClient } from '@/core/api/client';
import { API_ENDPOINTS } from '@/core/config';

export const jobsService = {
  searchJobs: (filters) => 
    apiClient.get(`${API_ENDPOINTS.JOBS_LIST}?...`),
  
  getJob: (id) =>
    apiClient.get(API_ENDPOINTS.JOBS_DETAIL(id)),
};
```

### Creating a Hook in a Module

```typescript
// src/modules/jobs/hooks/useJobs.ts
import { useState, useEffect } from 'react';
import { jobsService } from '../services';

export function useJobs(filters) {
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    jobsService.searchJobs(filters)
      .then(setJobs)
      .catch(setError)
      .finally(() => setLoading(false));
  }, [filters]);

  return { jobs, loading, error };
}
```

### Creating a Component in a Module

```typescript
// src/modules/jobs/components/JobCard.tsx
import { Card, CardContent, Button } from '@mui/material';
import { Job } from '../types';

interface JobCardProps {
  job: Job;
  onApply?: () => void;
}

export function JobCard({ job, onApply }: JobCardProps) {
  return (
    <Card>
      <CardContent>
        <h3>{job.title}</h3>
        <p>{job.company}</p>
        <Button onClick={onApply}>Apply</Button>
      </CardContent>
    </Card>
  );
}
```

## Phase 5: Update Page Components

### Before (Old Structure)

```typescript
// app/jobs/page.tsx
import { api } from '@/lib/api';
import { Navigation } from '@/components/Navigation';
import styles from './page.module.css';

export default function JobsPage() {
  const [jobs, setJobs] = useState([]);
  
  useEffect(() => {
    api.searchJobs().then(setJobs);
  }, []);

  return (
    <div className={styles.container}>
      {jobs.map(j => <div key={j.id}>{j.title}</div>)}
    </div>
  );
}
```

### After (New Structure)

```typescript
// app/jobs/page.tsx
import { Container, Box, CircularProgress } from '@mui/material';
import { Navigation } from '@/core/components';
import { JobList } from '@/modules/jobs/components';
import { useJobs } from '@/modules/jobs/hooks';

export default function JobsPage() {
  const { jobs, loading, error } = useJobs();

  if (loading) return <CircularProgress />;
  if (error) return <Box>Error loading jobs</Box>;

  return (
    <Container>
      <Navigation />
      <JobList jobs={jobs} />
    </Container>
  );
}
```

## Phase 6: TypeScript Configuration

Update `tsconfig.json` for clean imports:

```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/core": ["src/core/index.ts"],
      "@/core/*": ["src/core/*"],
      "@/modules/*": ["src/modules/*"],
      "@/config/*": ["config/*"],
      "@/*": ["app/*", "src/*"]
    }
  }
}
```

## Phase 7: Testing Each Module

### Test Structure

```
modules/jobs/
├── __tests__/
│   ├── services.test.ts
│   ├── hooks.test.ts
│   └── components.test.ts
├── components/
├── hooks/
└── services/
```

### Service Tests

```typescript
// modules/jobs/__tests__/services.test.ts
import { jobsService } from '../services';

describe('jobsService', () => {
  it('should fetch jobs', async () => {
    const jobs = await jobsService.searchJobs();
    expect(Array.isArray(jobs)).toBe(true);
  });
});
```

### Hook Tests

```typescript
// modules/jobs/__tests__/hooks.test.ts
import { renderHook, waitFor } from '@testing-library/react';
import { useJobs } from '../hooks';

describe('useJobs', () => {
  it('should load jobs', async () => {
    const { result } = renderHook(() => useJobs());
    
    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });
    
    expect(result.current.jobs.length).toBeGreaterThan(0);
  });
});
```

## Phase 8: Dependency Management

### Check Module Dependencies

Create a dependency matrix:

```
Module Dependencies:
├── auth
│   └── depends on: core
├── jobs
│   ├── depends on: core
│   └── used by: profile
├── profile
│   ├── depends on: core, jobs
│   └── used by: (none)
├── ideas
│   └── depends on: core
└── notifications
    └── depends on: core
```

### Avoid Circular Dependencies

❌ Bad:
```typescript
// modules/jobs/components/JobCard.tsx
import { ProfileLink } from '@/modules/profile'; // profile depends on jobs!
```

✅ Good:
```typescript
// modules/jobs/components/JobCard.tsx
import { useAuth } from '@/core'; // Only use from core
```

## Phase 9: Public API Definition

Each module must have a clear public API in `index.ts`:

```typescript
// modules/jobs/index.ts
export * from './components';  // All components
export * from './hooks';       // All hooks
export * from './services';    // Services
export * from './types';       // Types

// NOT exported:
// - internals from services/
// - API client directly
// - unused utility functions
```

## Migration Checklist

- [ ] Create core module structure
- [ ] Create auth module (security-critical)
- [ ] Create jobs module
- [ ] Create ideas module
- [ ] Create profile module
- [ ] Update tsconfig.json with paths
- [ ] Update app/layout.tsx for providers
- [ ] Update app/page.tsx for home
- [ ] Update app/jobs/page.tsx
- [ ] Update app/ideas/page.tsx
- [ ] Update app/profile/page.tsx
- [ ] Remove old lib/ directory
- [ ] Remove old components/ directory
- [ ] Update all imports
- [ ] Run full test suite
- [ ] Build and verify
- [ ] Deploy

## Common Issues & Solutions

### Issue: Circular Dependencies

**Problem**: Module A imports from Module B, which imports from Module A

**Solution**: Move shared code to core module

### Issue: Missing Type Exports

**Problem**: Type not exported from module index

**Solution**: Add to `index.ts`:
```typescript
export * from './types';
```

### Issue: Direct Service Import

**Problem**: Importing internal service:
```typescript
import { jobsService } from '@/modules/jobs/services/jobsService';
```

**Solution**: Import from module public API:
```typescript
import { jobsService } from '@/modules/jobs';
```

## Performance Considerations

- Module bundling: Each module loads only when needed
- Code splitting: Route-based code splitting still works
- Lazy loading: Import modules dynamically:

```typescript
const JobsPage = dynamic(() => import('@/modules/jobs/pages/JobsPage'), {
  loading: () => <LoadingSpinner />,
});
```

## Future: Converting to Micro Modules

If modules become independent services later:

```bash
# Publish to private npm registry
npm publish @company/auth-module
npm publish @company/jobs-module

# Use in other projects
npm install @company/auth-module
```

The structure makes this transition seamless!

## Summary

The migration transforms the codebase from a flat structure to a modular monolith:
- **Better organization**: Related code in one module
- **Easier scaling**: Add features as new modules
- **Independent testing**: Test modules in isolation
- **Clear dependencies**: Explicit module relationships
- **Future proof**: Ready for microservices evolution

Take time to understand the structure before starting migration!
