import axios, { type AxiosRequestConfig, AxiosError, type AxiosResponse } from 'axios';
//import type { AxiosError } from 'axios';
import type { ApiError } from '@/core/api/types/api.types';
import { logger } from '@/core/utils/logger/logger';
import { toCamelCase, toSnakeCase } from '@/core/api/utils/serialization';

/**
 * Axios Instance Configuration
 *
 * For API standards and best practices, see:
 * @see {@link ../../docs/api-standards.md}
 */

/**
 * Even though Axios methods return just (T) ie: data, not the full response.
 *
 * This is complete bullshit that I need to write this just to get Axios to work
 */
interface CustomAxiosInstance {
  post<T = unknown, R = T, D = unknown>(
    url: string,
    data?: D,
    config?: AxiosRequestConfig<D>
  ): Promise<R>;

  get<T = unknown, R = T, D = unknown>(
    url: string,
    config?: AxiosRequestConfig<D>
  ): Promise<R>;

  put<T = unknown, R = T, D = unknown>(
    url: string,
    data?: D,
    config?: AxiosRequestConfig<D>
  ): Promise<R>;

  delete<T = unknown, R = T, D = unknown>(
    url: string,
    config?: AxiosRequestConfig<D>
  ): Promise<R>;
}

/**
 * Axios instance configured for backend API requests.
 * Base configuration only - auth tokens are handled by useBackendQuery.
 *
 * @see https://axios-http.com/docs/config_defaults
 */
const axiosInstance = axios.create({
  baseURL: '/api', // Use Vite's proxy to avoid CORS issues
  timeout: 30000, // Increase to 30 seconds to match backend timeout
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
  validateStatus: (status) => status >= 200 && status < 300,
});

/**
 * Global error handler for consistency across requests
 * Combines production error handling with development logging
 */
const handleAxiosError = (error: AxiosError<ApiError>) => {
  // Always log in development
  logger.error('❌ API Error:', {
    url: error.config?.url,
    status: error.response?.status,
    message: error.message,
    details: error.response?.data
  });

  // Production error handling
  if (error.response?.status === 429) {
    // Rate limiting
    logger.warn('Rate limit exceeded');
  } else if (error.code === 'ECONNABORTED') {
    // Timeout
    logger.warn('Request timeout');
  }

  return Promise.reject(error);
};

// Request interceptor
axiosInstance.interceptors.request.use(
  (config) => {
    // Convert request data to snake_case if it exists
    if (config.data) {
      config.data = toSnakeCase(config.data);
      console.log('🐍 Converting request data to snake_case');
    }

    logger.debug('🚀 Outgoing request:', {
      url: config.url,
      method: config.method,
      data: config.data,
      params: config.params,
      headers: config.headers,
      baseURL: config.baseURL,
    });
    return config;
  },
  handleAxiosError // Use same error handler
);

// Response interceptor
axiosInstance.interceptors.response.use(
  (response: AxiosResponse) => {
    logger.debug('🔄 Response received:', {
      url: response.config.url,
      method: response.config.method,
      status: response.status,
      data: response.data,
      headers: response.headers,
    });

    // Convert ONLY the response data to camelCase
    const camelCaseData = toCamelCase(response.data);
    console.log('🐫 camelCase converted Response data:', camelCaseData);

    return { ...response, data: camelCaseData };
  },

  (error) => {
    // If there's an error response with data, convert that to camelCase too
    if (error.response?.data) {
      error.response.data = toCamelCase(error.response.data);
    }
    return Promise.reject(error);
  }
);

// We need to trick TypeScript into accepting this modified behavior
const typedAxiosInstance: CustomAxiosInstance = axiosInstance as unknown as CustomAxiosInstance;
export { typedAxiosInstance as axiosInstance };
