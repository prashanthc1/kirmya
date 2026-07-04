# Frontend Error Handling Guide

## Overview

The frontend implements a comprehensive error handling system with automatic classification, retry logic, user-friendly messages, and consistent error display across all pages.

## Components

### 1. Error Classification System (`lib/error-handling.ts`)

Categorizes errors into types for appropriate handling:

- **NETWORK**: Connection failures, timeouts
- **AUTH**: Authentication/authorization failures (401, 403)
- **VALIDATION**: Input validation errors (4xx except 401/403)
- **NOT_FOUND**: Resource not found (404)
- **SERVER**: Server errors (5xx)
- **TIMEOUT**: Request timeouts
- **FORBIDDEN**: Permission denied
- **UNKNOWN**: Uncategorized errors

### 2. Error Severity Levels

- **INFO**: Informational messages
- **WARNING**: Non-critical issues
- **ERROR**: Important errors requiring user action
- **CRITICAL**: System-level failures

### 3. Enhanced API Client (`lib/api-enhanced.ts`)

Provides automatic retry logic with:

- Exponential backoff (1s, 2s, 4s, etc.)
- Configurable retry attempts (default: 3)
- Request timeout handling (default: 10s)
- Automatic request cancellation for duplicate calls
- Retry callbacks for UI updates

Usage:
```typescript
import { apiEnhanced } from "@/lib/api-enhanced";

const jobs = await apiEnhanced.searchJobs(keyword, location, type);
```

### 4. Error Handler Hook (`lib/use-error-handler.ts`)

Provides consistent error handling in components:

```typescript
const { error, isRetrying, handleError, handleErrorWithRetry, clearError } = 
  useErrorHandler({ context: "feature-name" });

// Handle basic errors
try {
  await api.someCall();
} catch (err) {
  handleError(err, "Custom message");
}

// Handle errors with automatic retry
try {
  await api.someCall();
} catch (err) {
  await handleErrorWithRetry(err, () => api.someCall());
}
```

### 5. Error Boundary Component (`components/ErrorBoundary.tsx`)

Catches React component errors and displays a user-friendly error UI with:

- Error message display
- Development-only error stack trace
- "Try Again" button (resets error boundary)
- "Go Home" button (navigates to home)

The ErrorBoundary is wrapped around the entire app via LayoutWrapper.

## Implementation in Pages

All major pages have been updated to use the error handling system:

### Jobs Page (`app/jobs/page.tsx`)

- Uses `useErrorHandler` hook for consistent error handling
- Notifies users of success/errors via notifications
- Automatic error classification and retry logic for network errors

Handlers:
- `loadJobs()`: Search and filter jobs
- `handlePostJob()`: Post a new job
- `handleApplyJob()`: Apply for a job

### Ideas Page (`app/ideas/page.tsx`)

- Loads and filters community ideas
- Creates new ideas with validation
- Integrated error handling with retry support

### Idea Detail Page (`app/ideas/[id]/page.tsx`)

- Loads full idea details with discussions and tasks
- Adds new discussion comments
- Error context includes idea ID for debugging

### Profile Page (`app/profile/page.tsx`)

- Loads user applications
- Edits user profile
- Displays statistics with error fallbacks

## Testing Error Scenarios

### 1. Network Error Simulation

Edit the API endpoint temporarily:
```typescript
// In lib/api.ts
const API_BASE = "/api/invalid-endpoint"; // Triggers network error
```

Expected behavior:
- Error is classified as NETWORK type
- Automatic retry after 1s, 2s, 4s delays
- User sees retry notifications
- Error displays after max retries exceeded

### 2. Validation Error (400)

Send invalid data to endpoints:
```typescript
await api.postJob(token, "", "", "", "x", "", ""); // Title too short
```

Expected behavior:
- Error classified as VALIDATION
- Not retried (client error)
- User sees validation error message
- Notification displays with error severity

### 3. Authentication Error (401)

Use expired/invalid token:
```typescript
const invalidToken = "invalid.token.here";
await api.postJob(invalidToken, "Title", "Company", ...);
```

Expected behavior:
- Error classified as AUTH
- Not retried
- User should be redirected to login (implement in future)

### 4. Timeout Error

Configure timeout in api-enhanced.ts:
```typescript
await apiCallWithRetry("/jobs", { timeout: 100 }); // Very short timeout
```

Expected behavior:
- Request aborts after 100ms
- Classified as TIMEOUT type
- Automatic retry with exponential backoff

### 5. Server Error (500)

Server returns 500 status:
```typescript
// Backend simulates error
```

Expected behavior:
- Error classified as SERVER
- Automatically retried 3 times
- User sees retry notifications
- Error displayed if all retries fail

## Error Display Patterns

### Inline Errors (Per-Page)

```tsx
{pageError && (
  <div className={styles.error}>
    <button onClick={clearError}>×</button>
    {pageError.message}
  </div>
)}
```

### Toast Notifications (Global)

```tsx
import { useNotifications } from "@/components/Notifications";

const { notify } = useNotifications();
notify("Success message", "success");
notify("Error message", "error");
```

### Error Boundary (Component Level)

Wraps entire app to catch React component rendering errors.

## Best Practices

1. **Always provide context**: `useErrorHandler({ context: "feature-name" })`
2. **Use specific error messages**: `handleError(err, "Failed to post job")`
3. **Notify users of success**: `notify("Job posted!", "success")`
4. **Clear errors when retrying**: `clearError()`
5. **Let API client handle retries**: Don't manually retry in handlers
6. **Use AppError type**: Consistent error structure across app

## Production Monitoring Hook

The error system includes a production monitoring hook in `logError()`:

```typescript
if (process.env.NODE_ENV === "production") {
  // sendToErrorTrackingService(logData)
  // Example: Send to Sentry, LogRocket, or custom monitoring service
}
```

Implement your error tracking service here for production monitoring.

## Error Logging

All errors are logged with context:

```
🚨 Error: {
  timestamp: "2025-06-13T10:30:00Z",
  context: "jobs",
  type: "NETWORK",
  severity: "ERROR",
  message: "Network error. Please check your connection and try again.",
  retryable: true,
  retryCount: 1
}
```

Check browser console for error logs in development.

## Future Enhancements

1. **Error recovery strategies**: Implement offline mode for retryable errors
2. **Analytics**: Track error rates by type and page
3. **User feedback**: Add error reporting UI for critical failures
4. **Retry UI**: Show user when automatic retries are happening
5. **Error persistence**: Store recent errors for debugging
