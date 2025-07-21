import type { UseQueryOptions, UseMutationOptions } from '@tanstack/react-query';
import type { AxiosError } from 'axios';
import type { ApiError, ApiResponse } from '@/types/api/response';

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
 * @template TData The type of data RETURNED from API endpoint
 * @template TVariables The type of data SENT to API endpoint
 */
export type EnhancedMutationOptions<TData, TVariables> = UseMutationOptions<
  ApiResponse<TData>,
  AxiosError<ApiError>,
  TVariables
>;
