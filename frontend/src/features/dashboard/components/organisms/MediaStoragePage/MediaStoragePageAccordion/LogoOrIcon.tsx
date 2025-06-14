import { memo } from 'react';
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';

const ICON_CLASS = 'w-full h-full';

type PhysicalMediaType = 'disc' | 'hardware';

interface LogoOrIconProps {
  name: string;
  mediaType: 'subscription' | 'dlc' | 'inGamePurchase' | 'disc' | 'hardware';
}

export const LogoOrIcon = memo(function LogoOrIcon({ name, mediaType }: LogoOrIconProps) {
  const { games, physicalMedia, digitalMedia } = useDomainMaps();

  if (!mediaType) return null;
  let IconComponent = null;

  switch (mediaType) {
    case 'subscription': {
      const LogoComponent = games[name];
      return LogoComponent ? <LogoComponent className={ICON_CLASS} /> : null;
    }
    case 'dlc':
    case 'inGamePurchase': {
      const mediaTypeKey = mediaType === 'inGamePurchase' ? 'inGamePurchase' : 'dlc';
      IconComponent = digitalMedia[mediaTypeKey];
      return IconComponent ? <IconComponent className={ICON_CLASS} /> : null;
    }
    case 'disc':
    case 'hardware': {
      const mediaTypeKey = mediaType.toLowerCase() as PhysicalMediaType;
      IconComponent = physicalMedia[mediaTypeKey];
      return IconComponent ? <IconComponent className={ICON_CLASS} /> : null;
    }
    default:
      return null;
  }
});
