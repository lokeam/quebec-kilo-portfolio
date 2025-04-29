import type { JsonObject, JsonValue } from './caseConverter';

/**
 * A base interface for all snake_case API models
 * This ensures they're compatible with our transformation utilities
 */
export interface SnakeCaseModel extends JsonObject {
  [key: string]: JsonValue; // Use JsonValue instead of any
}

/**
 * Type utility that helps represent snake_case versions of camelCase properties
 *
 * Usage example:
 * interface UserSnakeCase extends SnakeCaseModel {
 *   user_id: string;
 *   first_name: string;
 *   last_name: string;
 * }
 */
export type SnakeCase<S extends string> = string extends S ? string :
  S extends `${infer T}${infer U}` ?
    T extends Uppercase<T> ?
      `_${Lowercase<T>}${SnakeCase<U>}` :
      `${T}${SnakeCase<U>}` :
    S;

/**
 * Type utility to convert entire interface from camelCase to snake_case
 *
 * Usage example:
 * type UserPayload = {
 *   userId: string;
 *   firstName: string;
 *   lastName: string;
 * };
 *
 * type UserPayloadSnakeCase = ToSnakeCase<UserPayload>;
 * // Equivalent to: { user_id: string; first_name: string; last_name: string; }
 */
export type ToSnakeCase<T> = {
  [K in keyof T as SnakeCase<string & K>]: T[K] extends object
    ? ToSnakeCase<T[K]>
    : T[K]
};

/**
 * Type utility for camelCase to snake_case conversion
 *
 * This makes it easier to create models for both representations
 * without having to manually type everything twice
 *
 * Usage example:
 * export interface User {
 *   userId: string;
 *   firstName: string;
 *   lastName: string;
 * }
 *
 * export type UserSnakeCase = SnakeCaseVersion<User>;
 */
export type SnakeCaseVersion<T> = ToSnakeCase<T> & SnakeCaseModel;