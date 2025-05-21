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
  // Handle arrays - preserve array structure
  if (Array.isArray(obj)) {
    return obj.map(item => toCamelCase(item)) as T;
  }

  // Handle objects - preserve object structure
  if (obj !== null && typeof obj === 'object') {
    const result = {} as Record<string, unknown>;
    for (const [key, value] of Object.entries(obj)) {
      // Convert the key to camelCase
      const camelKey = key.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
      // Recursively convert nested values
      result[camelKey] = toCamelCase(value);
    }
    return result as T;
  }

  // Return primitives as-is
  return obj as T;
}

/**
 * Converts a camelCase object to snake_case
 * Handles nested objects and arrays
 */
export function toSnakeCase<T>(obj: unknown): T {
  // Handle arrays - preserve array structure
  if (Array.isArray(obj)) {
    return obj.map(item => toSnakeCase(item)) as T;
  }

  // Handle objects - preserve object structure
  if (obj !== null && typeof obj === 'object') {
    const result = {} as Record<string, unknown>;
    for (const [key, value] of Object.entries(obj)) {
      // Convert the key to snake_case
      const snakeKey = key.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`);
      // Recursively convert nested values
      result[snakeKey] = toSnakeCase(value);
    }
    return result as T;
  }

  // Return primitives as-is
  return obj as T;
}