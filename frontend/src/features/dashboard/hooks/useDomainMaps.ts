import { useMemo } from 'react';
import type { ComponentType } from 'react';

// Game related logos
import AmazonLunaLogo from '@/shared/components/ui/LogoMap/domains/games/amazon_luna_logo';
import AppleLogo from '@/shared/components/ui/LogoMap/domains/games/apple_logo';
import BlizzardLogo from '@/shared/components/ui/LogoMap/domains/games/blizzard_logo';
import EALogo from '@/shared/components/ui/LogoMap/domains/games/ea_logo';
import EpicGamesLogo from '@/shared/components/ui/LogoMap/domains/games/epicgames_logo';
import FanaticalLogo from '@/shared/components/ui/LogoMap/domains/games/fanatical_logo';
import GreenManGamingLogo from '@/shared/components/ui/LogoMap/domains/games/green_man_gaming_logo';
import GoGLogo from '@/shared/components/ui/LogoMap/domains/games/gog_logo';
import HumbleBundleLogo from '@/shared/components/ui/LogoMap/domains/games/humble_bundle_logo';
import ItchIoLogo from '@/shared/components/ui/LogoMap/domains/games/itchio_logo';
import LunaLogo from '@/shared/components/ui/LogoMap/domains/games/amazon_luna_logo';
import MetaLogo from '@/shared/components/ui/LogoMap/domains/games/meta_logo';
import NetflixLogo from '@/shared/components/ui/LogoMap/domains/games/netflix_logo';
import NintendoLogo from '@/shared/components/ui/LogoMap/domains/games/nintendo_logo';
import NvidiaLogo from '@/shared/components/ui/LogoMap/domains/games/nvidia_logo';
import PlayPassLogo from '@/shared/components/ui/LogoMap/domains/games/playpass_logo';
import PlaystationLogo from '@/shared/components/ui/LogoMap/domains/games/playstation_logo';
import PrimeGamingLogo from '@/shared/components/ui/LogoMap/domains/games/prime_gaming_logo';
import ShadowLogo from '@/shared/components/ui/LogoMap/domains/games/shadow_logo';
import SteamLogo from '@/shared/components/ui/LogoMap/domains/games/steam_logo';
import UbisoftLogo from '@/shared/components/ui/LogoMap/domains/games/ubisoft_logo';
import XboxLogo from '@/shared/components/ui/LogoMap/domains/games/xbox_logo';

// TODO: Refactor out other domain maps into their own files

// Platform related logos
import {
  IconBrandAndroid,
  IconBrandApple,
  IconBrandWindowsFilled,
  IconDeviceGamepad,
  IconDeviceMobile,
 } from '@tabler/icons-react';

/**
 * @file useDomainMaps.ts
 * @description Custom hook for managing domain-specific SVG logo mappings.
 * Provides memoized access to logo components to prevent unnecessary re-renders.
 */


/**
 * Type definition for logo components that accept an optional className prop
 */
type LogoComponent = ComponentType<{ className?: string }>;

/**
 * Central mapping of all logos organized by domain
 * Each domain contains a record of logo name to component mappings
 */
export const LOGO_MAP: Record<string, Record<string, LogoComponent>> = {
  books: {},
  games: {
    amazon_luna: AmazonLunaLogo,
    apple: AppleLogo,
    blizzard: BlizzardLogo,
    ea: EALogo,
    epic: EpicGamesLogo,
    fanatical: FanaticalLogo,
    gog: GoGLogo,
    greenman: GreenManGamingLogo,
    humble: HumbleBundleLogo,
    itchio: ItchIoLogo,
    luna: LunaLogo,
    meta: MetaLogo,
    netflix: NetflixLogo,
    nintendo: NintendoLogo,
    nvidia: NvidiaLogo,
    prime: PrimeGamingLogo,
    playpass: PlayPassLogo,
    playstation: PlaystationLogo,
    shadow: ShadowLogo,
    steam: SteamLogo,
    sony: PlaystationLogo,
    ubisoft: UbisoftLogo,
    xbox: XboxLogo,
  },
  movies: {},
  music: {},
  platforms: {
    android: IconBrandAndroid,
    ios: IconBrandApple,
    pc: IconBrandWindowsFilled,
    console: IconDeviceGamepad,
    mobile: IconDeviceMobile,
  },
};

/**
 * Custom hook that provides memoized access to domain-specific logo mappings
 * @returns Memoized object containing domain-specific logo mappings
 */
export function useDomainMaps() {
  return {
    games: useMemo(() => LOGO_MAP.games, []),
    movies: useMemo(() => LOGO_MAP.movies, []),
    music: useMemo(() => LOGO_MAP.music, []),
    platforms: useMemo(() => LOGO_MAP.platforms, []),
  };
}
