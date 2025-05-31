/**
 * Digital services catalog.
 * Represents all available digital game distribution and subscription services.
 */
export const DigitalLocationName = {
  AMAZON_LUNA: 'amazonluna',
  APPLE_ARCADE: 'applearcade',
  BLIZZARD: 'blizzard',
  EA_PLAY: 'ea',
  EPIC_GAMES: 'epicgames',
  FANATICAL: 'fanatical',
  GOG: 'gog',
  GOOGLE_PLAY_PASS: 'googleplaypass',
  GREEN_MAN_GAMING: 'greenmangaming',
  HUMBLE_BUNDLE: 'humblebundle',
  ITCH_IO: 'itchio',
  META: 'meta',
  NINTENDO_SWITCH_ONLINE: 'nintendo',
  NVIDIA: 'nvidia',
  PRIME_GAMING: 'primegaming',
  PLAYSTATION_NETWORK: 'playstation',
  SHADOW: 'shadow',
  STEAM: 'steam',
  UBISOFT: 'ubisoft',
  XBOX_GAME_PASS: 'xboxgamepass'
} as const;

export type DigitalLocationName = typeof DigitalLocationName[keyof typeof DigitalLocationName];

/**
 * Human-readable display names for digital services.
 */
export const DIGITAL_LOCATION_DISPLAY_NAMES: Record<DigitalLocationName, string> = {
  [DigitalLocationName.AMAZON_LUNA]: 'Amazon Luna',
  [DigitalLocationName.APPLE_ARCADE]: 'Apple Arcade',
  [DigitalLocationName.BLIZZARD]: 'Blizzard Battle.net',
  [DigitalLocationName.EA_PLAY]: 'EA Play',
  [DigitalLocationName.EPIC_GAMES]: 'Epic Games',
  [DigitalLocationName.FANATICAL]: 'Fanatical',
  [DigitalLocationName.GOG]: 'GOG',
  [DigitalLocationName.GOOGLE_PLAY_PASS]: 'Google Play Pass',
  [DigitalLocationName.GREEN_MAN_GAMING]: 'Green Man Gaming',
  [DigitalLocationName.HUMBLE_BUNDLE]: 'Humble Bundle',
  [DigitalLocationName.ITCH_IO]: 'itch.io',
  [DigitalLocationName.META]: 'Meta',
  [DigitalLocationName.NINTENDO_SWITCH_ONLINE]: 'Nintendo Switch Online',
  [DigitalLocationName.NVIDIA]: 'NVIDIA',
  [DigitalLocationName.PRIME_GAMING]: 'Prime Gaming',
  [DigitalLocationName.PLAYSTATION_NETWORK]: 'PlayStation Plus',
  [DigitalLocationName.SHADOW]: 'Shadow',
  [DigitalLocationName.STEAM]: 'Steam',
  [DigitalLocationName.UBISOFT]: 'Ubisoft',
  [DigitalLocationName.XBOX_GAME_PASS]: 'Xbox Game Pass'
};

/**
 * Store URLs for digital services.
 */
export const DIGITAL_LOCATION_URLS: Record<DigitalLocationName, string> = {
  [DigitalLocationName.AMAZON_LUNA]: 'https://luna.amazon.com/',
  [DigitalLocationName.APPLE_ARCADE]: 'https://www.apple.com/apple-arcade/',
  [DigitalLocationName.BLIZZARD]: 'https://www.blizzard.com/en-us/',
  [DigitalLocationName.EA_PLAY]: 'https://www.ea.com/ea-play',
  [DigitalLocationName.EPIC_GAMES]: 'https://store.epicgames.com/en-US/',
  [DigitalLocationName.FANATICAL]: 'https://www.fanatical.com/en/',
  [DigitalLocationName.GOG]: 'https://www.gog.com/en/',
  [DigitalLocationName.GOOGLE_PLAY_PASS]: 'https://play.google.com/store/pass/getstarted/',
  [DigitalLocationName.GREEN_MAN_GAMING]: 'https://www.greenmangaming.com/',
  [DigitalLocationName.HUMBLE_BUNDLE]: 'https://www.humblebundle.com/',
  [DigitalLocationName.ITCH_IO]: 'https://itch.io/',
  [DigitalLocationName.META]: 'https://www.meta.com/nz/meta-quest-plus/',
  [DigitalLocationName.NINTENDO_SWITCH_ONLINE]: 'https://www.nintendo.com/',
  [DigitalLocationName.NVIDIA]: 'https://www.nvidia.com/en-us/geforce-now/',
  [DigitalLocationName.PRIME_GAMING]: 'https://gaming.amazon.com/home',
  [DigitalLocationName.PLAYSTATION_NETWORK]: 'https://www.playstation.com/en-us/playstation-network/',
  [DigitalLocationName.SHADOW]: 'https://shadow.tech/',
  [DigitalLocationName.STEAM]: 'https://store.steampowered.com/',
  [DigitalLocationName.UBISOFT]: 'https://www.ubisoft.com/en-us/',
  [DigitalLocationName.XBOX_GAME_PASS]: 'https://www.xbox.com/en-US/xbox-game-pass'
};

/**
 * Subscription status for digital services.
 */
export const SUBSCRIPTION_SERVICES: DigitalLocationName[] = [
  DigitalLocationName.AMAZON_LUNA,
  DigitalLocationName.APPLE_ARCADE,
  DigitalLocationName.EA_PLAY,
  DigitalLocationName.GOOGLE_PLAY_PASS,
  DigitalLocationName.META,
  DigitalLocationName.NINTENDO_SWITCH_ONLINE,
  DigitalLocationName.NVIDIA,
  DigitalLocationName.PRIME_GAMING,
  DigitalLocationName.PLAYSTATION_NETWORK,
  DigitalLocationName.SHADOW,
  DigitalLocationName.XBOX_GAME_PASS
];

/**
 * Service logos for display.
 */
export const DIGITAL_LOCATION_LOGOS: Record<DigitalLocationName, string> = {
  [DigitalLocationName.AMAZON_LUNA]: 'amazon',
  [DigitalLocationName.APPLE_ARCADE]: 'apple',
  [DigitalLocationName.BLIZZARD]: 'blizzard',
  [DigitalLocationName.EA_PLAY]: 'ea',
  [DigitalLocationName.EPIC_GAMES]: 'epicgames',
  [DigitalLocationName.FANATICAL]: 'fanatical',
  [DigitalLocationName.GOG]: 'gog',
  [DigitalLocationName.GOOGLE_PLAY_PASS]: 'google',
  [DigitalLocationName.GREEN_MAN_GAMING]: 'greenmangaming',
  [DigitalLocationName.HUMBLE_BUNDLE]: 'humblebundle',
  [DigitalLocationName.ITCH_IO]: 'itchio',
  [DigitalLocationName.META]: 'meta',
  [DigitalLocationName.NINTENDO_SWITCH_ONLINE]: 'nintendo',
  [DigitalLocationName.NVIDIA]: 'nvidia',
  [DigitalLocationName.PRIME_GAMING]: 'prime',
  [DigitalLocationName.PLAYSTATION_NETWORK]: 'ps',
  [DigitalLocationName.SHADOW]: 'shadow',
  [DigitalLocationName.STEAM]: 'steam',
  [DigitalLocationName.UBISOFT]: 'ubisoft',
  [DigitalLocationName.XBOX_GAME_PASS]: 'xbox'
};