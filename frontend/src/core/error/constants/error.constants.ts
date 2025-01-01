export const ERROR_MESSAGES = {
  // Generic messages that don't leak implementation details
  UNAUTHORIZED: 'Please sign in to continue',
  FORBIDDEN: 'You don\'t have permission to access this resource',
  NOT_FOUND: 'The requested resource was not found',
  SERVER_ERROR: 'We\'re experiencing technical difficulties',
  UNKNOWN: 'An unexpected error occurred'
} as const;

export const ERROR_ROUTES = {
  HOME: '/',
  LOGIN: '/login',
} as const;
