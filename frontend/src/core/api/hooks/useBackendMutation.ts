import { useMutation, type UseMutationOptions } from '@tanstack/react-query';
import { useAuth0 } from '@auth0/auth0-react';
import type { EnhancedMutationOptions } from '../types/query.types';
import type { ApiResponse } from '../types/api.types';
import { AxiosError } from 'axios';
import type { ApiError } from '../types/api.types';

/**
 * Hook for sending data to the backend (create, update, delete operations).
 *
 * Key features:
 * - Automatically handles authentication tokens
 * - Provides loading and error states
 * - Type-safe request and response data
 * - Handles cache updates after successful mutations
 *
 * Common use cases:
 * - Creating new resources (POST)
 * - Updating existing data (PUT)
 * - Deleting data (DELETE)
 *
 * @param mutationFn - Your API call function. Receives:
 *   - data: The data you want to send
 *   - token: Auth token (automatically injected)
 * @param options - Additional settings like:
 *   - onSuccess: Function to run after successful mutation
 *   - onError: Function to run if mutation fails
 *
 * @returns Mutation object containing:
 *   - mutate: Function to trigger the mutation
 *   - isLoading: Whether mutation is in progress
 *   - error: Error object if mutation failed
 *
 * @example Creating a new user
 * ```typescript
 * const createUser = useBackendMutation<User, NewUserData>(
 *   // Your API call
 *   async (userData, token) => {
 *     return axios.post('/users', userData, {
 *       headers: { Authorization: `Bearer ${token}` }
 *     });
 *   },
 *   // Additional options
 *   {
 *     onSuccess: (newUser) => {
 *       // Update UI or cache after success
 *       queryClient.invalidateQueries(['users']);
 *     },
 *     onError: (error) => {
 *       // Handle any errors
 *       toast.error('Failed to create user');
 *     }
 *   }
 * );
 *
 * // Using the mutation in a component
 * function CreateUserForm() {
 *   const handleSubmit = () => {
 *     createUser.mutate({ name: 'John', email: 'john@example.com' });
 *   };
 *
 *   if (createUser.isLoading) return <Spinner />;
 *
 *   return <form onSubmit={handleSubmit}>...</form>;
 * }
 * ```
 */
export function useBackendMutation<ResponseData, RequestData = unknown>(
  mutationFn: (data: RequestData, token: string) => Promise<ApiResponse<ResponseData>>,
  options?: EnhancedMutationOptions<ResponseData, RequestData>
) {
  const { getAccessTokenSilently } = useAuth0();

  const mutation = useMutation({
    ...options,
    mutationFn: async (data: RequestData) => {
      const token = await getAccessTokenSilently();
      return mutationFn(data, token);
    },
  });

  // Return a wrapped version with properly typed mutate function
  return {
    ...mutation,
    mutate: (
      data: RequestData,
      mutateOptions?: UseMutationOptions<
        ApiResponse<ResponseData>,
        AxiosError<ApiError>,
        RequestData,
        unknown
      >
    ) => {
      return mutation.mutate(data, mutateOptions);
    }
  };
}