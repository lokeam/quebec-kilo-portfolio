/**
 * Base API Response Types
 * These types exactly match the shape of data returned by our backend API
 */

export interface ApiResponse<T> {
  success: boolean;
  user_id: string;
  data: T;
}

export interface ApiError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}

export interface ApiErrorResponse {
  success: false;
  error: ApiError;
}