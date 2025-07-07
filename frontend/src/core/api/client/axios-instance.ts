import axios, { AxiosError } from 'axios';
import type { ApiError } from '@/types/api/response';

import { logger } from '@/core/utils/logger/logger';
import { toCamelCase, toSnakeCase } from '@/core/api/utils/serialization';
import { getAuth0Token } from '@/core/api/utils/auth.utils';

// Helper to check for plain objects
const isPlainObject = (value: unknown): value is Record<string, unknown> =>
  Object.prototype.toString.call(value) === '[object Object]';

/**
 * Axios instance configured for backend API requests.
 * Serialization and deserialization are handled via transforms;
 * interceptors are used for logging, auth, and global error handling.
 */
const axiosInstance = axios.create({
  baseURL: '/api',
  timeout: 30_000,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
  validateStatus: status => status >= 200 && status < 300,

  // 1) transformRequest: serialize JS object to snake_case JSON
  transformRequest: [
    (data: unknown, headers?: Record<string, string>) => {
      if (isPlainObject(data)) {
        headers!['Content-Type'] = 'application/json';
        try {
          const snakeCasedData = toSnakeCase(data);
          console.log("🐍 snakeCasedData data", snakeCasedData)
          return JSON.stringify(snakeCasedData);
        } catch (err) {
          logger.error('❌ Request serialization error', { error: err, data });
          throw err;
        }
      }
      // Leave other data types (string, FormData, etc.) unchanged
      return data;
    }
  ],

  // 2) transformResponse: parse JSON, validate envelope, camelCase
  transformResponse: [
    (raw: string) => {
      try {
        const parsed = JSON.parse(raw) as { success: boolean; error?: string; data: unknown };
        if (!isPlainObject(parsed) || typeof parsed.success !== 'boolean') {
          throw new Error('Invalid API response structure');
        }
        if (!parsed.success) {
          throw new Error(parsed.error ?? 'API returned unsuccessful status');
        }
        if (!isPlainObject(parsed.data)) {
          throw new Error('Missing API data field');
        }
        const camelCasedData = toCamelCase(parsed.data);
        console.log("🐫 camelCasedData data", camelCasedData)
        return camelCasedData;
      } catch (err) {
        logger.error('❌ Response parsing error', { error: err, raw });
        // Propagate parsing errors to be caught by the response interceptor
        throw err;
      }
    }
  ]
});

// Global error handler for interceptors
const handleAxiosError = (error: AxiosError<ApiError>): Promise<never> => {
  logger.error('❌ API Error', {
    url: error.config?.url,
    status: error.response?.status,
    message: error.message,
    details: error.response?.data,
  });
  if (error.response?.status === 429) {
    logger.warn('Rate limit exceeded');
  } else if (error.code === 'ECONNABORTED') {
    logger.warn('Request timeout');
  }
  return Promise.reject(error);
};

// Request interceptor: attach auth token, log requests
axiosInstance.interceptors.request.use(
  async config => {
    // Add Auth0 token to all requests
    try {
      const token = await getAuth0Token();
      config.headers = config.headers || {};
      config.headers['Authorization'] = `Bearer ${token}`;
    } catch (error) {
      // If we can't get a token, continue without it
      // The backend will handle unauthorized requests
      logger.debug('No auth token available for request', {
        url: config.url,
        error: error instanceof Error ? error.message : 'Unknown error'
      });
    }

    logger.debug('❓ INTERCEPTOR Request 📢', {
      method: config.method,
      url: config.url,
      data: config.data,
      params: config.params,
    });
    return config;
  },
  handleAxiosError
);

// Response interceptor: handle token refresh and retry
axiosInstance.interceptors.response.use(
  (response) => {
    // Successful response, just return it
    return response;
  },
  async (error) => {
    const originalRequest = error.config;

    // If we get a 401 and haven't already retried this request
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        // Try to get a fresh token
        logger.debug('Attempting to refresh token after 401 error');
        const newToken = await getAuth0Token();

        // Update the original request with the new token
        originalRequest.headers['Authorization'] = `Bearer ${newToken}`;

        // Retry the original request
        logger.debug('Retrying request with fresh token');
        return axiosInstance(originalRequest);
      } catch (refreshError) {
        logger.error('Failed to refresh token', {
          error: refreshError instanceof Error ? refreshError.message : 'Unknown error'
        });
        // If token refresh fails, let the error propagate
        return Promise.reject(error);
      }
    }

    // For other errors, just propagate them
    return Promise.reject(error);
  }
);

export { axiosInstance };
