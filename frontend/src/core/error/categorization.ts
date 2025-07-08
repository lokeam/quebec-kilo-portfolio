/**
 * Error categorization utilities for Sentry
 *
 * Prioritizes errors by type:
 * 1. UI Errors (Highest Priority) - Component rendering, user interface issues
 * 2. Authentication Errors (Medium Priority) - Login, token, permission issues
 * 3. API Errors (Lower Priority) - Network, data fetching issues
 */

export type ErrorCategory = 'ui' | 'authentication' | 'api' | 'navigation' | 'unknown';

export interface ErrorContext {
  category: ErrorCategory;
  priority: number;
  tags: Record<string, string>;
  context: Record<string, unknown>;
}

/**
 * Categorize an error based on its type and context
 */
export function categorizeError(error: Error, context?: Record<string, unknown>): ErrorContext {
  const errorMessage = error.message.toLowerCase();
  const errorName = error.name.toLowerCase();
  const stackTrace = error.stack?.toLowerCase() || '';

  // UI Errors (Highest Priority)
  if (
    errorMessage.includes('render') ||
    errorMessage.includes('component') ||
    errorMessage.includes('jsx') ||
    errorMessage.includes('react') ||
    errorName.includes('typeerror') ||
    stackTrace.includes('react') ||
    stackTrace.includes('component')
  ) {
    return {
      category: 'ui',
      priority: 1,
      tags: {
        type: 'ui_error',
        priority: 'high',
      },
      context: {
        errorType: 'UI Error',
        description: 'User interface or component rendering error',
        ...context,
      },
    };
  }

  // Authentication Errors (Medium Priority)
  if (
    errorMessage.includes('auth') ||
    errorMessage.includes('token') ||
    errorMessage.includes('login') ||
    errorMessage.includes('unauthorized') ||
    errorMessage.includes('forbidden') ||
    errorMessage.includes('permission') ||
    errorName.includes('auth') ||
    stackTrace.includes('auth')
  ) {
    return {
      category: 'authentication',
      priority: 2,
      tags: {
        type: 'auth_error',
        priority: 'medium',
      },
      context: {
        errorType: 'Authentication Error',
        description: 'User authentication or authorization error',
        ...context,
      },
    };
  }

  // API Errors (Lower Priority)
  if (
    errorMessage.includes('api') ||
    errorMessage.includes('network') ||
    errorMessage.includes('fetch') ||
    errorMessage.includes('axios') ||
    errorMessage.includes('http') ||
    errorMessage.includes('request') ||
    errorMessage.includes('response') ||
    errorName.includes('network') ||
    stackTrace.includes('api') ||
    stackTrace.includes('fetch') ||
    stackTrace.includes('axios')
  ) {
    return {
      category: 'api',
      priority: 3,
      tags: {
        type: 'api_error',
        priority: 'low',
      },
      context: {
        errorType: 'API Error',
        description: 'Network or API communication error',
        ...context,
      },
    };
  }

  // Navigation Errors
  if (
    errorMessage.includes('route') ||
    errorMessage.includes('navigation') ||
    errorMessage.includes('router') ||
    stackTrace.includes('router') ||
    stackTrace.includes('route')
  ) {
    return {
      category: 'navigation',
      priority: 2,
      tags: {
        type: 'navigation_error',
        priority: 'medium',
      },
      context: {
        errorType: 'Navigation Error',
        description: 'Routing or navigation error',
        ...context,
      },
    };
  }

  // Unknown/Other Errors
  return {
    category: 'unknown',
    priority: 4,
    tags: {
      type: 'unknown_error',
      priority: 'low',
    },
    context: {
      errorType: 'Unknown Error',
      description: 'Uncategorized error',
      ...context,
    },
  };
}

/**
 * Get error priority level for sorting and filtering
 */
export function getErrorPriority(category: ErrorCategory): number {
  const priorities: Record<ErrorCategory, number> = {
    ui: 1,
    authentication: 2,
    navigation: 2,
    api: 3,
    unknown: 4,
  };
  return priorities[category];
}

/**
 * Check if an error is high priority
 */
export function isHighPriorityError(category: ErrorCategory): boolean {
  return getErrorPriority(category) <= 2;
}