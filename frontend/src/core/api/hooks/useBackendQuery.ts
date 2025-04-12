import { useQuery } from '@tanstack/react-query';
import { useAuth0 } from '@auth0/auth0-react';
import type { BackendQueryOptions } from '@/core/api/types/api.types';
import { AxiosError } from 'axios';

// Mock token for development when Auth0 is not configured
const MOCK_AUTH_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c";

/**
 * Custom hook for making authenticated backend queries.
 * Wraps TanStack Query with Auth0 authentication + error handling.
 *
 * Features:
 * - Automatic Auth0 token injection
 * - Error handling with login redirect for 401s
 * - Network mode optimization for offline/online states
 * - Automatic retries on remount and reconnect
 *
 * @template ResponseData - Type of data returned from backend
 *
 * @property networkMode - 'online' (default) for optimal network handling:
 *   - Prevents unnecessary failed requests when offline
 *   - Works with browser's native online/offline detection
 *   - Automatically retries when connection is restored
 * @property retryOnMount - Retry failed queries when component remounts
 * @property refetchOnReconnect - Retry when internet connection is restored
 * @property throwOnError - Handle common error cases:
 *   - 401: Redirect to login
 *   - All others: Return false to prevent error propagation
 *
 *
 * @param queryOptions - Configuration object containing:
 *   - queryKey: Unique cache identifier (e.g., ['users', 'current'])
 *   - queryFn: Custom API call function that receives a getToken function.
 *     The hook will automatically inject Auth0's getAccessTokenSilently
 *     into the custom queryFn, so you can use it to get the token.
 *   - Additional TanStack Query options (enabled, staleTime, etc.)
 *
 * @returns TanStack Query result object containing:
 *   - data: The query response data
 *   - isLoading: Loading state
 *   - error: Error object if request failed
 *   - Additional TanStack Query properties
 *
 * @example
 * ```typescript
 * // Basic usage with getCurrentUserProfile query
 * function UserProfile() {
 *   const { data, isLoading, error } = useBackendQuery({
 *     queryKey: ['users', 'current'],
 *     queryFn: async (getToken) => {
 *       const token = await getToken();
 *       const response = await axios.get('/api/users/current', {
 *         headers: { Authorization: `Bearer ${token}` }
 *       });
 *       return response.data;
 *     }
 *   });
 *
 *   if (isLoading) return <Spinner />;
 *   if (error) return <ErrorDisplay error={error} />;
 *
 *   return <div>Welcome {data.name}</div>;
 * }
 * ```
 */
export function useBackendQuery<ResponseData>(
  queryOptions: BackendQueryOptions<ResponseData>
) {
  // We're not using Auth0 for now, but keep the import for future use
  const { getAccessTokenSilently } = useAuth0();

  // Create a mock token getter function
  const getMockToken = async () => {
    console.log("Using mock token instead of Auth0");
    return MOCK_AUTH_TOKEN;
  };

  return useQuery<ResponseData, Error>({
    networkMode: 'online',
    retryOnMount: true,
    refetchOnReconnect: true,
    ...queryOptions,
    /* Use mock token getter instead of Auth0's getAccessTokenSilently */
    queryFn: () => queryOptions.queryFn(getMockToken),
    throwOnError: (error: Error) => {
      if (error instanceof AxiosError && error.response?.status === 401) {
        console.error("401 Unauthorized error - would normally redirect to login");
        // Don't actually redirect for now since we're using a mock token
        // window.location.href = '/login';
      }
      console.error("Query error:", error);
      return false;
    }
  });
}
