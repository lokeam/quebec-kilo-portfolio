import axios, { type InternalAxiosRequestConfig } from 'axios';

import { logger } from '@/core/utils/logger/logger';
import { toCamelCase, toSnakeCase } from '@/core/api/utils/serialization';
import { getAuth0Token } from '@/core/api/utils/auth.utils';
import { logApiError, logAuthError, logApiBreadcrumb } from '@/core/monitoring/sentry-api-monitor';
import { HTTP_STATUS, HTTP_HEADERS, HTTP_STATUS_RANGES } from '@/core/api/constants/http-constants';

// Helper to check for plain objects within the response
const isPlainObject = (value: unknown): value is Record<string, unknown> =>
  Object.prototype.toString.call(value) === '[object Object]';

// Extended axios config interface for timing and retry logic
interface ExtendedAxiosRequestConfig extends InternalAxiosRequestConfig {
  _startTime?: number;
  _retry?: boolean;
}

/**
 * Axios instance for API requests.
 *
 * RESPONSIBILITIES:
 * - Make HTTP requests to the backend
 * - Handle authentication (add tokens to requests)
 * - Handle token refresh on 401 errors
 * - Log errors to monitoring service
 *
 * USAGE:
 * - Import and use: axiosInstance.get('/api/users')
 * - All errors are automatically logged to monitoring
 */
const BASE_API_URL = '/api';
const BASE_API_URL_TIMEOUT = 30_000;

const axiosInstance = axios.create({
  baseURL: BASE_API_URL,
  timeout: BASE_API_URL_TIMEOUT,
  headers: {
    [HTTP_HEADERS.CONTENT_TYPE]: HTTP_HEADERS.APPLICATION_JSON,
    [HTTP_HEADERS.ACCEPT]: HTTP_HEADERS.APPLICATION_JSON,
  },
  validateStatus: status => status >= HTTP_STATUS_RANGES.SUCCESS_MIN && status <= HTTP_STATUS_RANGES.SUCCESS_MAX,

  // Convert request data to snake_case for backend
  transformRequest: [
    (data: unknown, headers?: Record<string, string>) => {
      if (isPlainObject(data)) {
        headers![HTTP_HEADERS.CONTENT_TYPE] = HTTP_HEADERS.APPLICATION_JSON;
        try {
          const snakeCasedData = toSnakeCase(data);
          console.log("🐍 snakeCasedData data", snakeCasedData)
          return JSON.stringify(snakeCasedData);
        } catch (err) {
          logger.error('❌ Request serialization error', { error: err, data });
          throw err;
        }
      }

      return data;
    }
  ],

  // Convert response from snake_case to camelCase
  transformResponse: [
    (raw: string) => {
      try {
        console.log('🔍 TransformResponse - Raw response:', raw);

        const parsed = JSON.parse(raw) as { success: boolean; error?: string; data: unknown };

        // Only transform successful responses
        if (!isPlainObject(parsed)) {
          console.log('🔍 TransformResponse - Not a plain object, returning raw');
          // If it's not a plain object, return as-is (could be an error response)
          return raw;
        }

        console.log('🔍 TransformResponse - Parsed response:', parsed);

        // Check if this is a successful response with the expected structure
        if (typeof parsed.success === 'boolean' && parsed.success === true) {
          if (!isPlainObject(parsed.data)) {
            throw new Error('Missing API data field');
          }
          const camelCasedData = toCamelCase(parsed.data);
          console.log("🐫 camelCasedData data", camelCasedData);
          return camelCasedData;
        }

        // For unsuccessful responses or unexpected structures, return as-is
        // This allows error responses to be handled by the response interceptor
        console.log('🔍 TransformResponse - Unsuccessful response, returning raw');
        return raw;
      } catch (err) {
        logger.error('❌ Response parsing error', { error: err, raw });
        // Return raw response on parsing error to allow error handling
        return raw;
      }
    }
  ]
});

// Request interceptor: Add authentication token
axiosInstance.interceptors.request.use(
  async (config: ExtendedAxiosRequestConfig) => {
    // Record start time for performance tracking
    config._startTime = Date.now();

    // Add auth token to all requests
    try {
      const token = await getAuth0Token();
      config.headers = config.headers || {};
      config.headers[HTTP_HEADERS.AUTHORIZATION] = `${HTTP_HEADERS.BEARER_PREFIX}${token}`;
    } catch (error) {
      // If we can't get a token, continue without it
      // The backend will handle unauthorized requests
      logger.debug('No auth token available for request', {
        url: config.url,
        error: error instanceof Error ? error.message : 'Unknown error'
      });
    }

    // Sentry track breadcrumb: Track what the user is trying to do
    logApiBreadcrumb({
      url: config.url,
      method: config.method,
      message: `API request started: ${config.method} ${config.url}`,
      data: {
        hasAuthToken: !!config.headers?.[HTTP_HEADERS.AUTHORIZATION],
        requestSize: config.data ? JSON.stringify(config.data).length : 0
      }
    });

    logger.debug('❓ API Request', {
      method: config.method,
      url: config.url,
    });

    return config;
  },
  (error) => {
    logger.error('❌ Request interceptor error', error);
    return Promise.reject(error);
  }
);

// Response interceptor: Handle errors and token refresh
axiosInstance.interceptors.response.use(
  (response) => {
    const config = response.config as ExtendedAxiosRequestConfig;

    // Sentry track breadcrumb: Track successful requests
    if (config._startTime) {
      logApiBreadcrumb({
        url: response.config.url,
        method: response.config.method,
        message: `API request succeeded: ${response.status}`,
        data: {
          responseTime: Date.now() - config._startTime,
          responseSize: JSON.stringify(response.data).length
        }
      });
    }

    // Success - just return the response
    return response;
  },
  async (error) => {
    const originalRequest = error.config as ExtendedAxiosRequestConfig;

    // Try to extract error message from response
    let errorMessage = error.message;
    if (error.response?.data) {
      try {
        const errorData = typeof error.response.data === 'string'
          ? JSON.parse(error.response.data)
          : error.response.data;

        // Extract error message from various possible formats
        if (errorData.error) {
          errorMessage = errorData.error;
        } else if (errorData.message) {
          errorMessage = errorData.message;
        } else if (typeof errorData === 'string') {
          errorMessage = errorData;
        }
      } catch (parseError) {
        // If we can't parse the error response, use the original message
        logger.debug('Could not parse error response', { error: parseError, data: error.response.data });
      }
    }

    // Update the error message
    error.message = errorMessage;

    // Sentry track error: Log the API error to Sentry's monitoring service
    logApiError(error, {
      url: originalRequest?.url,
      method: originalRequest?.method,
      status: error.response?.status,
      errorMessage: errorMessage,
      responseData: error.response?.data,
    });

    // Handle 401 errors by refreshing the token and retrying
    if (error.response?.status === HTTP_STATUS.UNAUTHORIZED && !originalRequest._retry) {
      originalRequest._retry = true;

      // Sentry track breadcrumb: Track token refresh attempts
      logApiBreadcrumb({
        url: originalRequest?.url,
        method: originalRequest?.method,
        message: 'Token expired, attempting refresh',
        data: {
          originalStatus: error.response?.status,
          retryCount: 1
        }
      });

      try {
        logger.debug('Token expired, trying to refresh...');
        const newToken = await getAuth0Token();
        originalRequest.headers[HTTP_HEADERS.AUTHORIZATION] = `${HTTP_HEADERS.BEARER_PREFIX}${newToken}`;

        logger.debug('Retrying request with new token');
        return axiosInstance(originalRequest);
      } catch (refreshError) {
        logger.error('Failed to refresh token', refreshError);

        // Sentry track error: Log the auth error to monitoring service
        logAuthError(refreshError as Error, {
          originalUrl: originalRequest?.url,
          originalMethod: originalRequest?.method,
          originalStatus: error.response?.status,
          errorMessage: refreshError instanceof Error ? refreshError.message : 'Unknown error',
        });

        return Promise.reject(error);
      }
    }

    // For all other errors, just pass them through
    return Promise.reject(error);
  }
);

export { axiosInstance };
