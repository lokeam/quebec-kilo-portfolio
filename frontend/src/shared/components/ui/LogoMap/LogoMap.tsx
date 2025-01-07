import type { ReactElement, ComponentType } from 'react';
import { IconCloudDataConnection } from '@tabler/icons-react';

// Games
import AmazonLunaLogo from './domains/games/amazon_luna_logo';
import AppleLogo from './domains/games/apple_logo';
import BlizzardLogo from './domains/games/blizzard_logo';
import EALogo from './domains/games/ea_logo';
import EpicGamesLogo from './domains/games/epicgames_logo';
import FanaticalLogo from './domains/games/fanatical_logo';
import GreenManGamingLogo from './domains/games/green_man_gaming_logo';
import GoGLogo from './domains/games/gog_logo';
import HumbleBundleLogo from './domains/games/humble_bundle_logo';
import ItchIoLogo from './domains/games/itchio_logo';
import LunaLogo from './domains/games/amazon_luna_logo';
import MetaLogo from './domains/games/meta_logo';
import NetflixLogo from './domains/games/netflix_logo';
import NintendoLogo from './domains/games/nintendo_logo';
import NvidiaLogo from './domains/games/nvidia_logo';
import PlayPassLogo from './domains/games/playpass_logo';
import PlaystationLogo from './domains/games/playstation_logo';
import PrimeGamingLogo from './domains/games/prime_gaming_logo';
import ShadowLogo from './domains/games/shadow_logo';
import SteamLogo from './domains/games/steam_logo';
import UbisoftLogo from './domains/games/ubisoft_logo';
import XboxLogo from './domains/games/xbox_logo';


export const LOGO_MAP: Record<string, Record<string, ComponentType<{ className?: string}>>> = {
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
};

// Derived domain and logo types
type DomainName = keyof typeof LOGO_MAP;
export type LogoName<D extends DomainName> = keyof typeof LOGO_MAP[D];

interface SVGLogoProps<D extends DomainName> {
  // Constrain domain to keys in LOGO_MAP
  domain: D;

  // Constrain logo to keys in domain map
  name: LogoName<D>;

  className?: string;
}

// Looks up correct SVG based on domain and name

export default function SVGLogo<D extends DomainName>({
  domain,
  name,
  className,
}: SVGLogoProps<D>): ReactElement | null {
  const domainMap = LOGO_MAP[domain];

  console.log('Received logo props:', { domain, name });
  console.log('Available logos for domain:', Object.keys(LOGO_MAP[domain] || {}));

  // Return fallback if domain or logo doesn't exist
  if (!domainMap || !domainMap[name]) {
    console.log('Logo not found, falling back to default');
    return <IconCloudDataConnection className={className} />;
  }

  const SelectedLogo = domainMap[name];
  return <SelectedLogo className={className} />;
}