/**
 * Simple utility for converting between camelCase and snake_case
 * Used to transform data between frontend (camelCase) and backend (snake_case)
 *
 * This follows the principle of keeping transformation logic in one place
 * while maintaining type safety and simplicity.
 */

/**
 * Converts a snake_case object to camelCase
 * Handles nested objects and arrays
 */
export function toCamelCase<T>(obj: unknown): T {
  if (Array.isArray(obj)) {
    return obj.map(toCamelCase) as T;
  }

  if (obj !== null && typeof obj === 'object') {
    return Object.keys(obj as Record<string, unknown>).reduce((acc: Record<string, unknown>, key) => {
      const camelKey = key.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
      acc[camelKey] = toCamelCase((obj as Record<string, unknown>)[key]);
      return acc;
    }, {}) as T;
  }

  return obj as T;
}

/**
 * Converts a camelCase object to snake_case
 * Handles nested objects and arrays
 */
export function toSnakeCase<T>(obj: unknown): T {
  if (Array.isArray(obj)) {
    return obj.map(toSnakeCase) as T;
  }

  if (obj !== null && typeof obj === 'object') {
    return Object.keys(obj as Record<string, unknown>).reduce((acc: Record<string, unknown>, key) => {
      const snakeKey = key.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`);
      acc[snakeKey] = toSnakeCase((obj as Record<string, unknown>)[key]);
      return acc;
    }, {}) as T;
  }

  return obj as T;
}