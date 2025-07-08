/**
 * HTTP Constants for API Communication
 *
 * Contains all HTTP-related constants used throughout the API layer.
 * Improves maintainability and reduces magic strings.
 *
 * USAGE:
 * import { HTTP_STATUS, HTTP_HEADERS, HTTP_METHODS } from '@/core/api/constants/http-constants';
 */

/**
 * HTTP Status Codes
 * Common status codes used in API responses and error handling
 */
export const HTTP_STATUS = {
  // Success responses
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,

  // Client error responses
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  UNPROCESSABLE_ENTITY: 422,

  // Server error responses
  INTERNAL_SERVER_ERROR: 500,
  BAD_GATEWAY: 502,
  SERVICE_UNAVAILABLE: 503,
} as const;

/**
 * HTTP Headers
 * Standard HTTP header names and values used in API requests/responses
 */
export const HTTP_HEADERS = {
  // Header names
  AUTHORIZATION: 'Authorization',
  CONTENT_TYPE: 'Content-Type',
  ACCEPT: 'Accept',
  USER_AGENT: 'User-Agent',
  CACHE_CONTROL: 'Cache-Control',

  // Header values
  BEARER_PREFIX: 'Bearer ',
  APPLICATION_JSON: 'application/json',
  APPLICATION_XML: 'application/xml',
  TEXT_PLAIN: 'text/plain',
  MULTIPART_FORM_DATA: 'multipart/form-data',
} as const;

/**
 * HTTP Methods
 * Standard HTTP methods used in API requests
 */
export const HTTP_METHODS = {
  GET: 'GET',
  POST: 'POST',
  PUT: 'PUT',
  PATCH: 'PATCH',
  DELETE: 'DELETE',
  HEAD: 'HEAD',
  OPTIONS: 'OPTIONS',
} as const;

/**
 * HTTP Status Code Ranges
 * Used for status validation and categorization
 */
export const HTTP_STATUS_RANGES = {
  SUCCESS_MIN: 200,
  SUCCESS_MAX: 299,
  CLIENT_ERROR_MIN: 400,
  CLIENT_ERROR_MAX: 499,
  SERVER_ERROR_MIN: 500,
  SERVER_ERROR_MAX: 599,
} as const;

/**
 * Type exports for TypeScript usage
 */
export type HttpStatusCode = typeof HTTP_STATUS[keyof typeof HTTP_STATUS];
export type HttpHeaderName = typeof HTTP_HEADERS[keyof typeof HTTP_HEADERS];
export type HttpMethod = typeof HTTP_METHODS[keyof typeof HTTP_METHODS];

/**
 * Utility functions for HTTP status code validation
 */
export const HTTP_UTILS = {
  /**
   * Check if a status code indicates success (2xx range)
   */
  isSuccess: (status: number): boolean =>
    status >= HTTP_STATUS_RANGES.SUCCESS_MIN && status <= HTTP_STATUS_RANGES.SUCCESS_MAX,

  /**
   * Check if a status code indicates a client error (4xx range)
   */
  isClientError: (status: number): boolean =>
    status >= HTTP_STATUS_RANGES.CLIENT_ERROR_MIN && status <= HTTP_STATUS_RANGES.CLIENT_ERROR_MAX,

  /**
   * Check if a status code indicates a server error (5xx range)
   */
  isServerError: (status: number): boolean =>
    status >= HTTP_STATUS_RANGES.SERVER_ERROR_MIN && status <= HTTP_STATUS_RANGES.SERVER_ERROR_MAX,

  /**
   * Check if a status code indicates an error (4xx or 5xx range)
   */
  isError: (status: number): boolean =>
    HTTP_UTILS.isClientError(status) || HTTP_UTILS.isServerError(status),
} as const;