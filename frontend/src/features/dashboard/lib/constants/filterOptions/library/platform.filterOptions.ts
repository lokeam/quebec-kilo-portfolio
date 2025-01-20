export interface PlatformOption {
  readonly key: string;
  readonly label: string;
  readonly generation?: string;
  readonly searchTerms: string[];
}

export const CONSOLE_PLATFORMS: Record<string, ReadonlyArray<PlatformOption>> = {
  atari:[
    {
      key: 'atari2600',
      label: 'Atari 2600',
      searchTerms: ['atari', '2600', 'atari 2600']
    },
    {
      key: 'atari5200',
      label: 'Atari 5200',
      searchTerms: ['atari', '5200', 'atari 5200']
    },
    {
      key: 'atari7800',
      label: 'Atari 7800',
      searchTerms: ['atari', '7800', 'atari 7800']
    },
    {
      key: 'atarijaguar',
      label: 'Atari Jaguar',
      searchTerms: ['atari', 'jaguar', 'atari jaguar']
    },
    {
      key: 'atarilynx',
      label: 'Atari Lynx',
      searchTerms: ['atari', 'lynx', 'atari lynx']
    },
  ],
  microsoft: [
    {
      key: 'xbox',
      label: 'Xbox',
      searchTerms: ['microsoft', 'xbox', 'xbox classic', 'original xbox', 'og xbox']
    },
    {
      key: 'xbox360',
      label: 'Xbox 360',
      searchTerms: ['microsoft', 'xbox', '360', 'xbox 360']
    },
    {
      key: 'xboxone',
      label: 'Xbox One',
      searchTerms: ['microsoft', 'xbox', 'one', 'xbox one']
    },
    {
      key: 'xboxseriesx',
      label: 'Xbox Series X',
      searchTerms: ['microsoft', 'xbox', 'xbox series', 'series x', 'xbox series x']
    },
    {
      key: 'xboxseriesxs',
      label: 'Xbox Series S',
      searchTerms: ['microsoft', 'xbox', 'xbox series', 'series s', 'xbox series s']
    },
  ],
  mobile: [
    {
      key: 'android',
      label: 'Android',
      searchTerms: ['android', 'droid']
    },
    {
      key: 'ios',
      label: 'iOS',
      searchTerms: ['ios', 'iphone', 'ipad', 'apple']
    },
  ],
  nec: [
    {
      key: 'pcengine',
      label: 'PC Engine',
      searchTerms: ['nec', 'pc engine', 'pc engine cd', 'pc engine turbo duo']
    },
    {
      key: 'turbografx16',
      label: 'TurboGrafx 16',
      searchTerms: ['nec', 'turbo', 'turbo grafx', 'turbografx', '16', 'turbografx 16']
    },
    {
      key: 'pcenginecd',
      label: 'PC Engine CD',
      searchTerms: ['nec', 'pc', 'pc engine', 'cd', 'pc engine cd']
    },
    {
      key: 'turbografxcd',
      label: 'TurboGrafx CD',
      searchTerms: ['nec', 'turbo', 'turbo grafx', 'turbografx', 'cd', 'turbografx cd']
    },
    {
      key: 'turboduo',
      label: 'TurboDuo',
      searchTerms: ['nec', 'turbo', 'turbo duo', 'turboduo']
    },
    {
      key: 'pcengineturbo',
      label: 'PC Engine Turbo Duo',
      searchTerms: ['nec', 'pc', 'pc engine', 'turbo', 'turbo duo', 'pc engine turbo duo']
    },
  ],
  nintendo: [
    {
      key: 'famicon',
      label: 'Famicom',
      searchTerms: ['nintendo', 'famicom', 'famicon', 'famicom classic', 'original famicom']
    },
    {
      key: 'superfamicom',
      label: 'Super Famicom',
      searchTerms: ['nintendo', 'super', 'super famicom', 'super famicon', 'super famicom classic']
    },
    {
      key: 'nes',
      label: 'NES',
      searchTerms: ['nes', 'nintendo', 'nintendo entertainment system', 'nintendo classic']
    },
    {
      key: 'snes',
      label: 'Super Nintendo',
      searchTerms: ['snes','nintendo', 'super', 'super nintendo', 'super nintendo classic']
    },
    {
      key: 'n64',
      label: 'Nintendo 64',
      searchTerms: ['64', 'n64', 'nintendo', 'nintendo 64']
    },
    {
      key: 'wii',
      label: 'Nintendo Wii',
      searchTerms: ['wii', 'nintendo', 'wii classic', 'original wii']
    },
    {
      key: 'wii_u',
      label: 'Nintendo Wii U',
      searchTerms: ['wii', 'nintendo', 'wii u', 'wii u classic']
    },
    {
      key: 'switch',
      label: 'Nintendo Switch',
      searchTerms: ['switch', 'nintendo', 'switch classic', 'original switch']
    },
    {
      key: 'gameandwatch',
      label: 'Game & Watch',
      searchTerms: ['nintendo', 'game', 'watch', 'game & watch']
    },
    {
      key: 'gameboy',
      label: 'Gameboy',
      searchTerms: ['nintendo', 'gameboy', 'gameboy classic', 'original gameboy']
    },
    {
      key: 'gameboy_advance',
      label: 'Gameboy Advance',
      searchTerms: ['gba', 'nintendo', 'gameboy', 'advance', 'gameboy advance']
    },
    {
      key: 'nintendods',
      label: 'Nintendo DS',
      searchTerms: ['nintendo', 'ds', 'nintendo ds']
    },
    {
      key: 'nintendo3ds',
      label: 'Nintendo 3DS',
      searchTerms: ['nintendo', '3ds', 'nintendo 3ds']
    },
  ],
  pc: [
    {
      key: 'pc',
      label: 'PC',
      searchTerms: ['pc', 'windows']
    },
    {
      key: 'mac',
      label: 'Mac',
      searchTerms: ['mac', 'macos']
    },
  ],
  sega: [
    {
      key: 'sg1000',
      label: 'SG-1000',
      searchTerms: ['sg', 'sega', '1000', 'sg-1000']
    },
    {
      key: 'mastersystem',
      label: 'Master System',
      searchTerms: ['sega', 'master', 'system', 'master system']
    },
    {
      key: 'megadrive',
      label: 'Mega Drive',
      searchTerms: ['sega', 'mega', 'drive', 'mega drive']
    },
    {
      key: 'genesis',
      label: 'Genesis',
      searchTerms: ['sega', 'genesis', 'sega genesis']
    },
    {
      key: 'sega_32x',
      label: 'Sega 32X',
      searchTerms: ['sega', '32', 'sega 32x']
    },
    {
      key: 'sega_cd',
      label: 'Sega CD',
      searchTerms: ['sega', 'cd', 'sega cd']
    },
    {
      key: 'saturn',
      label: 'Sega Saturn',
      searchTerms: ['sega', 'saturn', 'sega saturn']
    },
    {
      key: 'dreamcast',
      label: 'Sega Dreamcast',
      searchTerms: ['sega', 'dreamcast', 'sega dreamcast']
    },
    {
      key: 'gamegear',
      label: 'Game Gear',
      searchTerms: ['sega', 'game', 'gear', 'sega game gear']
    },
  ],
  sony: [
    {
      key: 'ps1',
      label: 'PlayStation 1',
      searchTerms: ['1', 'ps1', 'ps', 'sony', 'playstation', 'playstation 1']
    },
    {
      key: 'ps2',
      label: 'PlayStation 2',
      searchTerms: ['2', 'ps2', 'ps', 'sony', 'playstation', 'playstation 2']
    },
    {
      key: 'ps3',
      label: 'PlayStation 3',
      searchTerms: ['3', 'ps3', 'ps', 'sony', 'playstation', 'playstation 3']
    },
    {
      key: 'ps4',
      label: 'PlayStation 4',
      searchTerms: ['4', 'ps4', 'ps', 'sony', 'playstation', 'playstation 4']
    },
    {
      key: 'ps5',
      label: 'PlayStation 5',
      searchTerms: ['5', 'ps5', 'ps', 'sony', 'playstation', 'playstation 5']
    },
    {
      key: 'psp',
      label: 'PlayStation Portable',
      searchTerms: ['psp', 'sony', 'playstation', 'playstation portable']
    },
    {
      key: 'psvita',
      label: 'PlayStation Vita',
      searchTerms: ['vita', 'sony', 'playstation', 'playstation vita']
    },
  ],
} as const;
