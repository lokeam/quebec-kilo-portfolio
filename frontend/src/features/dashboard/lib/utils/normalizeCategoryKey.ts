/**
 * Converts camelCase or PascalCase strings to snake_case
 * Used specifically for normalizing category keys for consistent color/display name mapping
 */
export function normalizeToSnakeCase(str: string): string {
  return str
    .replace(/([A-Z])/g, '_$1')
    .replace(/^_/, '')
    .toLowerCase();
}