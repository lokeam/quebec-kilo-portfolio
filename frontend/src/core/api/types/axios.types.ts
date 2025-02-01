import type { AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';

/**
 * Extended Axios request configuration
 * Adds custom behavior flags to the standard Axios config
 */
export interface CustomRequestConfig extends AxiosRequestConfig {
  /** Skip authentication header injection */
  skipAuth?: boolean;

  /** Skip global error handling middleware */
  skipErrorHandling?: boolean;
}


/**
 * Interceptor configuration for Axios
 * Defines success and error handlers for requests and responses
 */
export interface InterceptorConfig {
  request?: {
    /** Transform or augment request before sending */
    onSuccess?: (config: CustomRequestConfig) => Promise<CustomRequestConfig>;

    /** Handle request preparation errors */
    onError?: (error: AxiosError) => Promise<never>;
  };
  response?: {
    /** Transform or process response before returning */
    onSuccess?: (response: AxiosResponse) => AxiosResponse | Promise<AxiosResponse>;

    /** Handle response errors (4xx, 5xx, network) */
    onError?: (error: AxiosError) => Promise<never>;
  };
}
