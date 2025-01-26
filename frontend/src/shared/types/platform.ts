
/**
 * High-level categories for gaming platforms
 * Used for broad classification and filtering
 *
 * @example
 * ```ts
 * const mobileGame: Platform = {
 *   category: 'Mobile',
 *   model: 'iOS'
 * };
 * ```
 */
export type PlatformCategory =
  | 'Android'
  | 'Console'
  | 'iOS'
  | 'Mobile'
  | 'PC';

/**
 * Specific platform models across different manufacturers
 * Includes historical + current gaming systems
 *
 * Note: Some platforms have multiple naming conventions
 * (e.g., 'Game Boy' and 'Gameboy' are both ok)
 *
 * @example
 * ```ts
 * const ps5Game: Platform = {
 *   category: 'Console',
 *   model: 'PlayStation 5'
 * };
 * ```
 */
  export type PlatformModel =
  | 'Atari 2600'
  | 'Atari 5200'
  | 'Atari 7800'
  | 'Atari Jaguar'
  | 'Atari Lynx'
  | 'Xbox'
  | 'Xbox 360'
  | 'Xbox One'
  | 'Xbox Series X'
  | 'Xbox Series S'
  | 'Android'
  | 'iOS'
  | 'PC Engine'
  | 'TurboGrafx 16'
  | 'PC Engine CD'
  | 'TurboGrafx CD'
  | 'TurboDuo'
  | 'PC Engine Turbo Duo'
  | 'Famicom'
  | 'Super Famicom'
  | 'Nintendo Entertainment System'
  | 'Super Nintendo Entertainment System'
  | 'NES'
  | 'SNES'
  | 'Nintendo 64'
  | 'Nintendo Wii'
  | 'Nintendo Wii U'
  | 'Nintendo Switch'
  | 'Game & Watch'
  | 'Gameboy'
  | 'Game Boy'
  | 'Gameboy Advance'
  | 'Game Boy Advance'
  | 'Nintendo DS'
  | 'Nintendo 3DS'
  | 'PC'
  | 'Mac'
  | 'MacOS'
  | 'SG-1000'
  | 'Master System'
  | 'Mega Drive'
  | 'Genesis'
  | 'Sega 32X'
  | 'Sega CD'
  | 'Sega Saturn'
  | 'Sega Dreamcast'
  | 'Sega Game Gear'
  | 'PlayStation 1'
  | 'PlayStation 2'
  | 'PlayStation 3'
  | 'PlayStation 4'
  | 'PlayStation 5'
  | 'PlayStation Portable'
  | 'PlayStation Vita'
  | 'Windows PC';

/**
 * Represents a gaming platform with its category and optional specific model
 *
 * @example
 * ```ts
 * // Console with specific model
 * const switch: Platform = {
 *   category: 'Console',
 *   model: 'Nintendo Switch'
 * };
 *
 * // PC platform without specific model
 * const pc: Platform = {
 *   category: 'PC'
 * };
 * ```
 */
export interface Platform {
  /** High-level platform category */
  category: PlatformCategory;
  /** Optional specific model for the platform */
  model?: PlatformModel;
}

