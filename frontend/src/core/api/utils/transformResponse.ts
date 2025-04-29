import { JsonValue, snakeToCamelCase } from './caseConverter';

/**
 * Transform API response data from snake_case to camelCase
 * Use this as an opt-in utility rather than forcing all API responses
 * to be transformed automatically
 *
 * @example
 * // In a service
 * async function getLocation(id: string) {
 *   const response = await axiosInstance.get(`/locations/${id}`);
 *   return transformResponse(response);
 * }
 */
export function transformResponse<T extends JsonValue>(data: T): T {
  return snakeToCamelCase(data);
}

/**
 * Create a wrapped version of an API fetching function that transforms its result
 *
 * @example
 * // Original function
 * const getLocations = () => axiosInstance.get('/locations');
 *
 * // Wrapped function that transforms responses
 * const getLocationsWithCamelCase = withCamelCase(getLocations);
 */
export function withCamelCase<T extends JsonValue, P extends unknown[]>(
  fn: (...args: P) => Promise<T>
): (...args: P) => Promise<T> {
  return async (...args: P) => {
    const result = await fn(...args);
    return transformResponse(result);
  };
}