/**
 * Mapping of digital location names to their logo keys.
 * Used to normalize various service names to their canonical logo identifiers.
 */
export const DIGITAL_LOCATION_LOGO_KEYS = {
  // Gaming Platforms
  'steam': 'steam',
  'epic games': 'epic',
  'epic games store': 'epic',
  'gog': 'gog',
  'gog.com': 'gog',
  'ubisoft': 'ubisoft',
  'ubisoft+': 'ubisoft',
  'ea play': 'ea',
  'electronic arts': 'ea',
  'origin': 'ea',
  'battlenet': 'blizzard',
  'blizzard': 'blizzard',
  'rockstar': 'rockstar',
  'rockstar games': 'rockstar',

  // Console Platforms
  'playstation': 'playstation',
  'playstation network': 'playstation',
  'psn': 'playstation',
  'xbox': 'xbox',
  'xbox network': 'xbox',
  'xbox game pass': 'xbox',
  'nintendo': 'nintendo',
  'nintendo switch online': 'nintendo',

  // Cloud Gaming
  'geforce now': 'nvidia',
  'nvidia': 'nvidia',
  'amazon luna': 'luna',
  'luna': 'luna',
  'xcloud': 'xbox',
  'stadia': 'stadia',

  // Mobile/Arcade
  'apple arcade': 'apple',
  'google play pass': 'playpass',
  'play pass': 'playpass',
  'meta quest': 'meta',
  'meta': 'meta',

  // Game Stores
  'humble bundle': 'humble',
  'humble': 'humble',
  'green man gaming': 'greenman',
  'fanatical': 'fanatical',
  'itch.io': 'itchio',
  'itchio': 'itchio',

  // Streaming Services
  'netflix': 'netflix',
  'netflix games': 'netflix',
  'prime gaming': 'prime',
  'amazon prime gaming': 'prime',
} as const;

/**
 * Type for valid logo keys
 */
export type DigitalLocationLogoKey = typeof DIGITAL_LOCATION_LOGO_KEYS[keyof typeof DIGITAL_LOCATION_LOGO_KEYS];

/**
 * Normalizes a digital location name to its corresponding logo key.
 * Returns undefined if no matching logo is found.
 */
export function normalizeDigitalLocationName(name: string): DigitalLocationLogoKey | undefined {
  if (!name) return undefined;
  return DIGITAL_LOCATION_LOGO_KEYS[name.toLowerCase().trim() as keyof typeof DIGITAL_LOCATION_LOGO_KEYS];
}