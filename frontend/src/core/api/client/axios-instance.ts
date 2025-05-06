import axios, { type AxiosResponse, AxiosError } from 'axios';
import type { ApiError } from '@/core/api/types/api.types';
import { logger } from '@/core/utils/logger/logger';
import { toCamelCase, toSnakeCase } from '@/core/api/utils/serialization';

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
  config => {
    // Example: attach auth token
    // const token = getAuthToken();
    // if (token) config.headers['Authorization'] = `Bearer ${token}`;
    logger.debug('→ Request', {
      method: config.method,
      url: config.url,
      data: config.data,
      params: config.params,
    });
    return config;
  },
  handleAxiosError
);

// Response interceptor: log responses and handle errors
axiosInstance.interceptors.response.use(
  (response: AxiosResponse) => {
    logger.debug('← Response', {
      status: response.status,
      url: response.config.url,
    });
    return response;
  },
  handleAxiosError
);

export { axiosInstance };
