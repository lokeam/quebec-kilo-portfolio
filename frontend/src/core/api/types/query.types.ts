import type { UseQueryOptions, UseMutationOptions } from '@tanstack/react-query';
import type { AxiosError } from 'axios';
import type { ApiError, ApiResponse } from './api.types';

/**
 * Enhanced Query Options for TanStack Query
 * Omits queryKey + queryFn since they are required parameters within useApiQuery
 * @template ResponseData The type of data returned from API endpoint
 */
export type EnhancedQueryOptions<ResponseData> = Omit<
  UseQueryOptions<ApiResponse<ResponseData>, AxiosError<ApiError>>,
  'queryKey' | 'queryFn'
>;


/**
 * Enhanced Mutation Options for TanStack Query
 * @template ResponseData The type of data RETURNED from API endpoint
 * @template RequestData The type of data SENT to API endpoint
 * Example: EnhancedMutationOptions<User, CreateUserRequest>
 */
export type EnhancedMutationOptions<
  ResponseData,       // Example: User
  RequestData         // Example: CreateUserRequest
> = Omit<
  UseMutationOptions<
  ApiResponse<ResponseData>,   // Example: { data: User, metadata: {...} }
  AxiosError<ApiError>,        // Example: { status: 400, message: 'Bad Request', ... }
  RequestData                  // Example: { name: 'John Doe', email: 'john.doe@example.com', ... }
>,
  'mutationFn'
>;
