import { memo } from 'react';
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';
import { IconCloudDataConnection, IconDeviceGamepad2, IconPackage, IconDisc, IconCpu2, IconSparkles } from '@tabler/icons-react';
import { normalizeDigitalLocationName } from '../constants/digital-location-logos';
import { MediaCategory } from '@/types/domain/spend-tracking';

interface MediaIconProps {
  /** The name of the service or media (e.g., 'steam', 'epic', 'gog') */
  name?: string;
  /** The type of media (e.g., 'subscription', 'dlc', 'hardware') */
  mediaType?: MediaCategory;
  /** Optional className for styling the icon */
  className?: string;
}

/**
 * A unified component for rendering media icons based on service name or media type.
 * Falls back to appropriate icons based on the context.
 *
 * @example
 * ```tsx
 * // Service icon
 * <MediaIcon name="steam" className="h-4 w-4" />
 *
 * // Media type icon
 * <MediaIcon mediaType={MediaCategory.HARDWARE} className="h-4 w-4" />
 * ```
 */
export const MediaIcon = memo(function MediaIcon({
  name,
  mediaType,
  className = "h-4 w-4"
}: MediaIconProps) {
  const { games } = useDomainMaps();

  // Try to get icon from service name first
  const normalizedName = name ? normalizeDigitalLocationName(name) : undefined;
  const ServiceIcon = normalizedName ? games[normalizedName] : null;

  if (ServiceIcon) {
    return <ServiceIcon className={className} />;
  }

  // If no service icon, try to get icon from media type
  if (mediaType) {
    switch (mediaType) {
      case MediaCategory.SUBSCRIPTION:
        return <IconCloudDataConnection className={className} />;
      case MediaCategory.DLC:
        return <IconDeviceGamepad2 className={className} />;
      case MediaCategory.IN_GAME_PURCHASE:
        return <IconSparkles className={className} />;
      case MediaCategory.MISC:
        return <IconPackage className={className} />;
      case MediaCategory.PHYSICAL:
        return <IconDisc className={className} />;
      case MediaCategory.HARDWARE:
        return <IconCpu2 className={className} />;
      default:
        return <IconPackage className={className} />;
    }
  }

  // Default fallback
  return <IconPackage className={className} />;
});