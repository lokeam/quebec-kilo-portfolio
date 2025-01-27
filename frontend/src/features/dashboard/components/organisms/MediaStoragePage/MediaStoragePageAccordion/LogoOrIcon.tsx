import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';
import type { BaseMediaCategory } from '@/features/dashboard/lib/types/spend-tracking/constants';

const ICON_CLASS = 'w-full h-full';

interface LogoOrIconProps {
  name: string;
  mediaType: BaseMediaCategory;
};

export function LogoOrIcon({ name, mediaType }: LogoOrIconProps) {
  const { games, physicalMedia, digitalMedia } = useDomainMaps();

  if (!mediaType) return null;
  let IconComponent = null;

  switch (mediaType) {
    case 'subscription': {
      const LogoComponent = games[name];
      return LogoComponent ? <LogoComponent className={ICON_CLASS} /> : null;
    }
    case 'dlc':
    case 'inGamePurchase':
      IconComponent = digitalMedia[mediaType];
      return IconComponent ? <IconComponent className={ICON_CLASS} /> : null;
    case 'disc':
    case 'hardware':
      IconComponent = physicalMedia[mediaType];
      return IconComponent ? <IconComponent className={ICON_CLASS} /> : null;
    default:
      return null;
  }
}
