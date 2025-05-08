/**
 * API Types
 *
 * Defines common types used across the API layer.
 */

/**
 * Standard API error response
 */
export interface ApiError {
  /** Error code */
  code: string;

  /** Human-readable error message */
  message: string;

  /** Optional detailed error information */
  details?: Record<string, unknown>;

  /** HTTP status code */
  status: number;
}

/**
 * Standard API response wrapper
 */
export interface ApiResponse<T> {
  /** Whether the request was successful */
  success: boolean;

  /** Response data */
  data: T;

  /** Response metadata */
  metadata: {
    /** Timestamp of the response */
    timestamp: string;

    /** Unique request identifier */
    request_id: string;
  };
}