import axios, { type AxiosRequestConfig, AxiosError } from 'axios';
//import type { AxiosError } from 'axios';
import type { ApiError } from '@/core/api/types/api.types';
import { logger } from '@/core/utils/logger/logger';

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
  baseURL: '/api', // This will proxy through Vite's dev server
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
  (response) => {
    logger.debug('🔄 Response received:', {
      url: response.config.url,
      method: response.config.method,
      status: response.status,
      data: response.data,
      headers: response.headers,
    });

    return response.data;
  },

  (error) => {
    // Enhance error logging
    if (error.code === 'ECONNABORTED') {
      logger.error('⏱️ Request timeout - Timed out waiting for response', {
        url: error.config?.url,
        baseURL: error.config?.baseURL,
        method: error.config?.method,
        timeout: error.config?.timeout,
        data: error.config?.data,
      });
    } else if (error.response) {
      // Server responded with a status code outside 2xx range
      logger.error('❌ Server error response:', {
        url: error.config?.url,
        status: error.response?.status,
        data: error.response?.data,
        method: error.config?.method,
      });
    } else if (error.request) {
      // Request made but no response received (network error)
      logger.error('❌ Network error - No response received:', {
        url: error.config?.url,
        method: error.config?.method,
      });
    }

    return Promise.reject(error);
  }
);

// We need to trick TypeScript into accepting this modified behavior
const typedAxiosInstance: CustomAxiosInstance = axiosInstance as unknown as CustomAxiosInstance;
export { typedAxiosInstance as axiosInstance };
