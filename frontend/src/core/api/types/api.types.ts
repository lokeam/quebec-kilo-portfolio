import type { UseQueryOptions } from '@tanstack/react-query';

/**
 * Base API Response wrapper
 * @template ResponseData The type of data returned from API endpoint
 */
export interface ApiResponse<ResponseData> {
  data: ResponseData;
  metadata?: {
    timestamp: string;  // ISO-8601 format
    requestId: string;  // UUID v4
  };
}

/**
 * Base API Error shape returned from backend
 */
export interface ApiError {
  code: string;      // Error code (e.g., 'AUTH_001', 'VALIDATION_ERROR')
  message: string;   // Human-readable error message
  details?: Record<string, unknown>;  // Additional error context
}

/**
 * Parameters for paginated API requests
 */
export interface PaginationParams {
  page: number;   // 0-based page number
  limit: number;  // Items per page
}

/**
 * Response shape for paginated data
 * @template ListData The type of items being paginated
 */
export type PaginatedResponse<ListData> = ApiResponse<ListData> & {
  pagination: {
    total: number;      // Total number of items
    currentPage: number; // Current page (0-based)
    totalPages: number;  // Total number of pages
  };
}

/**
 * Backend Query Options for API calls
 * Extends TanStack Query options with Auth0 token handling
 *
 * @template ResponseData The type of data returned from API
 *
 * Type Parameters for UseQueryOptions:
 * - ResponseData: The expected data shape from backend
 * - Error: The error type (defaults to Error)
 *
 * Omitted Properties (handled by current implementation):
 * - queryKey: Managed separately for type safety
 * - queryFn: Redefined to include token handling
 *
 * Inherited Properties from UseQueryOptions:
 * - enabled: Boolean to control query execution
 * - staleTime: Duration before data is considered stale
 * - cacheTime: Duration to keep inactive data in cache
 * - retry: Number of retry attempts
 * - retryDelay: Delay between retry attempts
 * - networkMode: 'online' | 'always' | 'offlineFirst'
 *
 * @example
 * ```typescript
 * const userQuery: BackendQueryOptions<User> = {
 *   queryKey: ['users', userId],
 *   queryFn: async (getToken) => {
 *     const token = await getToken();
 *     return api.get(`/users/${userId}`, {
 *       headers: { Authorization: `Bearer ${token}` }
 *     });
 *   },
 *   enabled: !!userId,
 *   staleTime: 5 * 60 * 1000, // 5 minutes
 *   retry: 3
 * }
 * ```
 */
export interface BackendQueryOptions<ResponseData> extends
  Omit<UseQueryOptions<ResponseData, Error>, 'queryKey' | 'queryFn'> {
  queryKey: readonly unknown[];
  queryFn: (getToken: () => Promise<string>) => Promise<ResponseData>;
}

/**
 * Valid service types for digital locations
 */
export const SERVICE_TYPES = {
  SUBSCRIPTION: 'subscription',
  BASIC: 'basic'
} as const;

export type ServiceTypeValue = typeof SERVICE_TYPES[keyof typeof SERVICE_TYPES];

/**
 * Creates a properly typed service type value based on subscription status
 */
export function createServiceType(isSubscription: boolean): ServiceTypeValue {
  return isSubscription ? 'subscription' : 'basic';
}

/**
 * Maps frontend billing cycle display values to backend billing cycle values
 */
export const BILLING_CYCLE_MAP: Record<string, string> = {
  '1 month': 'monthly',
  '3 months': 'quarterly',
  '6 months': 'semi_annual',
  '1 year': 'annual'
};

/**
 * Default values for digital services
 */
export const DIGITAL_SERVICE_DEFAULTS = {
  URL: 'https://example.com',
  PAYMENT_METHOD: 'Visa',
  BILLING_CYCLE: '1 month'
};
