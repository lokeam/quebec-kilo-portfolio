import type { PlatformCategory, PlatformModel } from '@/shared/types/platform';


export interface PlatformOption {
  /** Platform category identifier */
  readonly key: PlatformCategory;
  /** Display name for the platform */
  readonly label: string;
  /** Optional array of search terms for improved filtering */
  readonly searchTerms?: string[];
  /** Optional array of specific models for this platform */
  readonly models?: PlatformModel[];
}

/**
 * Mapping of platform manufacturers to respective platform options
 * Organized by manufacturer (e.g., 'nintendo', 'sony') for grouping in UI
 *
 * @example
 * ```tsx
 * // Access PlayStation platforms
 * const playstationPlatforms = CONSOLE_PLATFORMS.sony;
 *
 * // Loop over all manufacturers and platforms
 * Object.entries(CONSOLE_PLATFORMS).forEach(([manufacturer, platforms]) => {
 *   console.log(`${manufacturer}: ${platforms.map(p => p.label).join(', ')}`);
 * });
 * ```
 */
export const CONSOLE_PLATFORMS: Record<string, ReadonlyArray<PlatformOption>> = {
  atari: [
    {
      key: 'Console',
      label: 'Atari',
      searchTerms: [
        '2600',
        '5200',
        '7800',
        'atari',
        'console',
        'jaguar',
        'lynx',
        'atari 2600',
        'atari 5200',
        'atari 7800',
        'atari jaguar',
        'atari lynx',
      ],
      models: [
        'Atari 2600',
        'Atari 5200',
        'Atari 7800',
        'Atari Jaguar',
        'Atari Lynx',
      ],
    },
  ],
  microsoft: [
    {
      key: 'Console',
      label: 'Xbox',
      searchTerms: [
        '360',
        'console',
        'original xbox',
        'og xbox',
        'one',
        'series',
        'xbox',
        'microsoft',
        'xbox classic',
        'xbox 360',
        'xbox one',
        'xbox series',
        'xbox series x',
        'xbox series s',
      ],
      models: [
        'Xbox',
        'Xbox 360',
        'Xbox One',
        'Xbox Series X',
        'Xbox Series S',
      ],
    },
  ],
  mobile: [
    {
      key: 'Mobile',
      label: 'Mobile',
      searchTerms: [
        'android',
        'apple',
        'droid',
        'ios',
        'iphone',
        'ipad',
        'mobile',
        'phone',
        'tablet',
      ],
      models: ['Android', 'iOS'],
    },
  ],
  nec: [
    {
      key: 'Console',
      label: 'NEC',
      searchTerms: [
        '16',
        'cd',
        'nec',
        'pc engine',
        'pc engine cd',
        'pc engine turbo duo',
        'turbo duo',
        'turbo',
        'turbo duo',
        'turboduo',
        'turbo grafx',
        'turbografx',
        'turbografx 16',
        'turbografx cd',
      ],
      models: [
        'PC Engine',
        'TurboGrafx 16',
        'PC Engine CD',
        'TurboGrafx CD',
        'TurboDuo',
        'PC Engine Turbo Duo',
      ],
    },
  ],
  nintendo: [
    {
      key: 'Console',
      label: 'Nintendo',
      searchTerms: [
        '3ds',
        '64',
        'boy',
        'ds',
        'family computer',
        'famicom',
        'famicom classic',
        'gba',
        'game',
        'game & watch',
        'gameboy',
        'gameboy advance',
        'gameboy classic',
        'n64',
        'nes',
        'nintendo',
        'nintendo classic',
        'nintendo entertainment system',
        'nintendo 3ds',
        'nintendo 64',
        'nintendo ds',
        'nintendo ds',
        'nintendo 3ds',
        'nintendo wii',
        'nintendo wii u',
        'nintendo switch',
        'snes',
        'super',
        'super famicom',
        'super famicom classic',
        'super nes',
        'super nintendo',
        'switch',
        'switch classic',
        'wii',
        'wii classic',
        'wii u',
        'wii u classic',
      ],
      models: [
        'Famicom',
        'Super Famicom',
        'Nintendo Entertainment System',
        'Super Nintendo Entertainment System',
        'Nintendo 64',
        'Game & Watch',
        'Game Boy',
        'Game Boy Advance',
        'Nintendo DS',
        'Nintendo 3DS',
        'Nintendo Wii',
        'Nintendo Wii U',
        'Nintendo Switch',
      ],
    },
  ],
  pc: [
    {
      key: 'PC',
      label: 'PC',
      searchTerms: [
        'computer',
        'pc',
        'windows',
       ],
      models: ['Windows PC', 'Mac'],
    },
  ],
  sega: [
    {
      key: 'Console',
      label: 'Sega',
      searchTerms: [
        '32',
        '32x',
        '1000',
        'cd',
        'console',
        'dreamcast',
        'game gear',
        'gear',
        'genesis',
        'master system',
        'mega',
        'mega drive',
        'saturn',
        'sega',
        'sega 32x',
        'sega cd',
        'sega game gear',
        'sega genesis',
        'sega saturn',
        'sega dreamcast',
        'sg',
        'sg-1000',
        'system',
      ],
      models: [
        'SG-1000',
        'Master System',
        'Mega Drive',
        'Genesis',
        'Sega 32X',
        'Sega CD',
        'Sega Saturn',
        'Sega Dreamcast',
        'Sega Game Gear',
      ],
    },
  ],
  sony: [
    {
      key: 'Console',
      label: 'PlayStation',
      searchTerms: [
        '2',
        '3',
        '4',
        '5',
        'playstation',
        'playstation 2',
        'playstation 3',
        'playstation 4',
        'playstation 5',
        'playstation portable',
        'playstation vita',
        'ps',
        'ps1',
        'ps2',
        'ps3',
        'ps4',
        'ps5',
        'psp',
        'vita',
      ],
      models: [
        'PlayStation 1',
        'PlayStation 2',
        'PlayStation 3',
        'PlayStation 4',
        'PlayStation 5',
        'PlayStation Portable',
        'PlayStation Vita',
      ],
    },
  ],
} as const;
