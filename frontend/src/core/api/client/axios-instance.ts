import axios from 'axios';
import type { AxiosError } from 'axios';
import type { ApiError } from '@/core/api/types/api.types';
import { logger } from '@/core/utils/logger/logger';

/**
 * Axios instance configured for backend API requests.
 * Base configuration only - auth tokens are handled by useBackendQuery.
 *
 * @see https://axios-http.com/docs/config_defaults
 */
const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
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
      params: config.params,
      headers: config.headers,
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
      status: response.status,
      data: response.data,
      headers: response.headers,
    });
    return response.data;
  },
  handleAxiosError // Use same error handler
);

export { axiosInstance };
