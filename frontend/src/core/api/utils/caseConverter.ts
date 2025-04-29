// Types for the converter
export type JsonValue =
  | string
  | number
  | boolean
  | null
  | undefined
  | JsonObject
  | JsonArray;

export interface JsonObject {
  [key: string]: JsonValue;
}

export type JsonArray = JsonValue[];

/**
 * Converts snake_case keys to camelCase recursively through objects and arrays
 * Example: { user_id: 1 } → { userId: 1 }
 */
export function snakeToCamelCase<T extends JsonValue>(data: T): T {
  if (data === null || data === undefined || typeof data !== 'object') {
    return data;
  }

  if (Array.isArray(data)) {
    return data.map(item => snakeToCamelCase(item)) as T;
  }

  return Object.keys(data as JsonObject).reduce((result, key) => {
    // Convert key from snake_case to camelCase
    const camelKey = key.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());

    // Handle nested objects and arrays
    const value = (data as JsonObject)[key];
    result[camelKey] = snakeToCamelCase(value);

    return result;
  }, {} as Record<string, JsonValue>) as T;
}