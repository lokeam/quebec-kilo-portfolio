import axios from 'axios';
import type { AxiosError } from 'axios';
import type { ApiError } from '@/core/api/types/api.types';

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
 * Auth errors (401) are handled by useBackendQuery
 */
const handleAxiosError = (error: AxiosError<ApiError>) => {
  console.error('API Error:', error);
  return Promise.reject(error);
};

// Response interceptor
axiosInstance.interceptors.response.use(
  (response) => response.data,
  handleAxiosError
);

export { axiosInstance };
