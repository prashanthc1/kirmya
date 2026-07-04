# Developer Quick Start Guide

Welcome to the Modular Monolith frontend! This guide helps you get started.

## Understanding the Architecture (5 min)

Read in this order:
1. [ARCHITECTURE_SUMMARY.md](./ARCHITECTURE_SUMMARY.md) - Overview (5 min)
2. [MODULAR_ARCHITECTURE.md](./MODULAR_ARCHITECTURE.md) - Detailed structure (10 min)

## Setting Up Your Environment

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Open browser
# http://localhost:3000
```

## Project Structure Quick Reference

```
src/
├── core/              # Shared config, types, API client
├── modules/
│   ├── auth/         # Authentication features
│   ├── jobs/         # Job management
│   ├── ideas/        # Ideas & collaboration
│   ├── profile/      # User profile
│   └── ...
```

## Adding a Feature (Complete Example)

### Scenario: Add "Favorite Jobs" Feature

#### Step 1: Plan Module Structure

```
modules/jobs/
├── components/
│   ├── FavoriteButton.tsx      # ← NEW
│   └── index.ts
├── hooks/
│   ├── useFavorite.ts          # ← NEW
│   └── index.ts
├── services/
│   └── index.ts                # Update: add favorite API calls
└── types/
    └── index.ts                # Update: add Favorite type
```

#### Step 2: Add Type

```typescript
// modules/jobs/types/index.ts
export interface Favorite {
  job_id: string;
  user_id: string;
  created_at: string;
}
```

#### Step 3: Add Service

```typescript
// modules/jobs/services/index.ts
export const jobsService = {
  // ... existing code ...
  
  favoriteJob: (token: string, jobId: string): Promise<Favorite> => {
    return apiClient.post(
      `/jobs/${jobId}/favorite`,
      {},
      { headers: { Authorization: `Bearer ${token}` } }
    );
  },

  unfavoriteJob: (token: string, jobId: string): Promise<void> => {
    return apiClient.delete(
      `/jobs/${jobId}/favorite`,
      { headers: { Authorization: `Bearer ${token}` } }
    );
  },
};
```

#### Step 4: Add Hook

```typescript
// modules/jobs/hooks/useFavorite.ts
import { useState } from 'react';
import { jobsService } from '../services';
import { useAuth } from '@/core';
import { useNotifications } from '@/modules/notifications';

export function useFavorite(jobId: string) {
  const [isFavorited, setIsFavorited] = useState(false);
  const [loading, setLoading] = useState(false);
  const { token } = useAuth();
  const { notify } = useNotifications();

  const toggleFavorite = async () => {
    if (!token) {
      notify('Please login first', 'warning');
      return;
    }

    setLoading(true);
    try {
      if (isFavorited) {
        await jobsService.unfavoriteJob(token, jobId);
        setIsFavorited(false);
        notify('Removed from favorites', 'info');
      } else {
        await jobsService.favoriteJob(token, jobId);
        setIsFavorited(true);
        notify('Added to favorites', 'success');
      }
    } catch (error) {
      notify('Failed to update favorite', 'error');
    } finally {
      setLoading(false);
    }
  };

  return { isFavorited, loading, toggleFavorite };
}
```

#### Step 5: Add Component

```typescript
// modules/jobs/components/FavoriteButton.tsx
import { IconButton, CircularProgress } from '@mui/material';
import { useFavorite } from '../hooks';

interface FavoriteButtonProps {
  jobId: string;
  onlyIcon?: boolean;
}

export function FavoriteButton({ jobId, onlyIcon }: FavoriteButtonProps) {
  const { isFavorited, loading, toggleFavorite } = useFavorite(jobId);

  if (onlyIcon) {
    return (
      <IconButton
        onClick={toggleFavorite}
        disabled={loading}
        color={isFavorited ? 'error' : 'default'}
      >
        {loading ? <CircularProgress size={24} /> : (isFavorited ? '❤️' : '🤍')}
      </IconButton>
    );
  }

  return (
    <Button
      onClick={toggleFavorite}
      disabled={loading}
      variant={isFavorited ? 'contained' : 'outlined'}
      color={isFavorited ? 'error' : 'primary'}
    >
      {isFavorited ? 'Favorited' : 'Add to Favorites'}
    </Button>
  );
}
```

#### Step 6: Export from Module

```typescript
// modules/jobs/components/index.ts
export { FavoriteButton } from './FavoriteButton';  // ← ADD THIS

// modules/jobs/hooks/index.ts
export { useFavorite } from './useFavorite';  // ← ADD THIS
```

#### Step 7: Use in Component

```typescript
// In JobCard component
import { FavoriteButton } from '@/modules/jobs';

function JobCard({ job }) {
  return (
    <Card>
      <CardContent>
        <h3>{job.title}</h3>
        <FavoriteButton jobId={job.id} onlyIcon />
      </CardContent>
    </Card>
  );
}
```

## Common Patterns

### Pattern 1: Fetch Data with Hook

```typescript
// hooks/useJobs.ts
import { useState, useEffect } from 'react';
import { jobsService } from '../services';

export function useJobs(filters) {
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchJobs = async () => {
      setLoading(true);
      try {
        const data = await jobsService.searchJobs(filters);
        setJobs(data);
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    };

    fetchJobs();
  }, [filters]);

  return { jobs, loading, error };
}
```

### Pattern 2: Mutation Hook

```typescript
// hooks/usePostJob.ts
import { useState } from 'react';
import { jobsService } from '../services';
import { useAuth } from '@/core';

export function usePostJob() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const { token } = useAuth();

  const postJob = async (data) => {
    setLoading(true);
    setError(null);
    try {
      const result = await jobsService.postJob(token, data);
      return result;
    } catch (err) {
      setError(err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { postJob, loading, error };
}
```

### Pattern 3: Component with Form

```typescript
// components/PostJobForm.tsx
import { useState } from 'react';
import { Box, TextField, Button } from '@mui/material';
import { usePostJob } from '../hooks';

export function PostJobForm({ onSuccess }) {
  const [formData, setFormData] = useState({
    title: '',
    company: '',
    location: '',
  });
  const { postJob, loading, error } = usePostJob();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await postJob(formData);
      setFormData({ title: '', company: '', location: '' });
      onSuccess?.();
    } catch (err) {
      // Error already handled in hook
    }
  };

  return (
    <Box component="form" onSubmit={handleSubmit}>
      <TextField
        label="Job Title"
        value={formData.title}
        onChange={(e) =>
          setFormData({ ...formData, title: e.target.value })
        }
        fullWidth
        required
      />
      <TextField
        label="Company"
        value={formData.company}
        onChange={(e) =>
          setFormData({ ...formData, company: e.target.value })
        }
        fullWidth
        required
      />
      <Button type="submit" variant="contained" disabled={loading}>
        {loading ? 'Posting...' : 'Post Job'}
      </Button>
      {error && <Box color="error.main">{error.message}</Box>}
    </Box>
  );
}
```

## Debugging

### Check Module Exports

```typescript
// In any component, you can import and console.log
import * as Jobs from '@/modules/jobs';
console.log('Jobs exports:', Jobs);
```

### Verify Service Call

```typescript
// In browser console
const { jobsService } = await import('/src/modules/jobs/services/index.ts');
jobsService.searchJobs();
```

### API Network Tab

1. Open DevTools → Network tab
2. Filter by "api"
3. Check request/response

## Testing

### Run Tests

```bash
# All tests
npm test

# Single module
npm test -- jobs

# Watch mode
npm test -- --watch
```

### Write a Test

```typescript
// modules/jobs/__tests__/services.test.ts
import { jobsService } from '../services';

describe('jobsService', () => {
  it('should search jobs', async () => {
    const jobs = await jobsService.searchJobs({ keyword: 'engineer' });
    expect(Array.isArray(jobs)).toBe(true);
  });

  it('should get single job', async () => {
    const job = await jobsService.getJob('job-123');
    expect(job.id).toBe('job-123');
  });
});
```

## Performance Tips

1. **Use dynamic imports for heavy modules**
   ```typescript
   const JobsPage = dynamic(() => import('@/modules/jobs'));
   ```

2. **Memoize components**
   ```typescript
   export const JobCard = React.memo(({ job }) => (
     // ...
   ));
   ```

3. **Use useCallback for functions**
   ```typescript
   const handleClick = useCallback(() => {
     // ...
   }, [dependencies]);
   ```

## Useful Commands

```bash
npm run dev       # Start dev server
npm run build     # Build for production
npm start         # Run production build
npm test          # Run tests
npm run lint      # Lint code
```

## Getting Help

1. **Read module README.md** - Each module documents its features
2. **Check examples** - Look at similar components
3. **Read MODULAR_ARCHITECTURE.md** - Architecture details
4. **Ask the team** - Discuss patterns and best practices

## Next Steps

1. ✅ Read ARCHITECTURE_SUMMARY.md
2. ✅ Start dev server: `npm run dev`
3. ✅ Explore a module: `src/modules/jobs/`
4. ✅ Try the feature example above
5. ✅ Build something cool!

## Tips for Success

✅ **Keep modules independent** - Don't import between feature modules  
✅ **Use public APIs** - Import from module index.ts, not internal files  
✅ **Test everything** - Write tests as you code  
✅ **Document your code** - Add README to new modules  
✅ **Follow patterns** - Use existing patterns as templates  
✅ **Ask for review** - Code review helps catch issues early  

---

**You're ready to develop! Happy coding! 🚀**
