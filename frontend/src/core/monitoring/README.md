# API Monitoring Architecture

## Overview

This directory contains a clean, separated implementation of API monitoring using Sentry. The architecture follows the **Single Responsibility Principle** by separating HTTP logic from monitoring logic.

## Architecture

### 1. Axios Instance (`../api/client/axios-instance.ts`)
**Responsibility**: Handle HTTP requests and authentication

- Makes API calls to the backend
- Handles authentication (adds tokens to requests)
- Handles token refresh on 401 errors
- **Calls monitoring service** when errors occur

### 2. Sentry API Monitor (`sentry-api-monitor.ts`)
**Responsibility**: Handle all Sentry logging and monitoring

- Logs API errors with context
- Logs authentication failures
- Provides consistent error reporting
- **Pure monitoring logic** - no HTTP concerns

## How It Works

### Clean Separation of Concerns

```typescript
// 1. Axios interceptor handles HTTP logic
axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    // Handle HTTP concerns (token refresh, etc.)

    // Call monitoring service for logging
    logApiError(error, {
      url: originalRequest?.url,
      method: originalRequest?.method,
      status: error.response?.status,
      errorMessage: error.message,
    });
  }
);

// 2. Monitoring service handles Sentry logic
export function logApiError(error: Error, context: ApiErrorContext): void {
  Sentry.captureException(error, {
    tags: { source: 'api' },
    contexts: { api_request: context }
  });
}
```

### Benefits of This Architecture

1. **Single Responsibility**: Each file has one clear purpose
2. **Easy to Test**: Monitoring functions can be tested independently
3. **Easy to Maintain**: Changes to monitoring don't affect HTTP logic
4. **Easy to Understand**: Clear separation makes code readable
5. **Reusable**: Monitoring functions can be used elsewhere

## Usage

### Making API Calls
```typescript
import { axiosInstance } from '@/core/api/client/axios-instance';

// Just use normally - monitoring happens automatically
const users = await axiosInstance.get('/api/users');
const newUser = await axiosInstance.post('/api/users', { name: 'John' });
```

### Manual Monitoring (if needed)
```typescript
import { logApiError, logAuthError } from '@/core/monitoring/sentry-api-monitor';

// Log custom errors
logApiError(new Error('Custom error'), {
  url: '/api/custom',
  method: 'POST',
  status: 500,
  errorMessage: 'Something went wrong'
});
```

## What Gets Logged

### API Errors
- URL, method, status code
- Error message and response data
- User context (automatically from Sentry)

### Authentication Errors
- Token refresh failures
- Original request details
- Error context

### Performance (Optional)
- Slow response warnings
- Success response tracking

## Debugging

### Check Sentry for:
- **API errors**: Look for `source: 'api'` tag
- **Auth errors**: Look for `source: 'auth'` tag
- **Performance**: Look for `performance` breadcrumbs

### Common Issues:
1. **401 errors**: Check if token refresh is working
2. **Slow responses**: Check for performance breadcrumbs
3. **Missing context**: Verify user is logged in

## Testing

Run the test file to verify everything works:
```typescript
import { testApiMonitoringFunctions } from './sentry-api-monitor.test';
testApiMonitoringFunctions(); // Should log success
```

## Maintenance

### Adding New Monitoring Features
1. Add new functions to `sentry-api-monitor.ts`
2. Call them from axios interceptors as needed
3. Keep HTTP logic separate from monitoring logic

### Modifying Error Logging
1. Update functions in `sentry-api-monitor.ts`
2. No changes needed to axios interceptors
3. All monitoring changes are centralized

### Debugging Issues
1. Check axios interceptor logs first
2. Check Sentry for error details
3. Use test functions to verify monitoring works

## Best Practices

1. **Keep HTTP and monitoring separate**
2. **Use descriptive function names**
3. **Provide rich context for errors**
4. **Test monitoring functions independently**
5. **Document any changes to the architecture**

## Why This Architecture?

### Before (Spaghetti Code)
- HTTP logic mixed with monitoring logic
- Hard to test and maintain
- Difficult to understand
- Violates Single Responsibility Principle

### After (Clean Architecture)
- Clear separation of concerns
- Easy to test and maintain
- Easy to understand and modify
- Follows SOLID principles

This architecture makes the code professional, maintainable, and interview-ready.