import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';
import { IconCloudDataConnection } from '@tabler/icons-react';
import { memo } from 'react';
import { normalizeDigitalLocationName } from '../constants/digital-location-logos';

interface DigitalLocationIconProps {
  /** The name of the digital location (e.g., 'steam', 'epic', 'gog') */
  name?: string;
  /** Optional className for styling the icon */
  className?: string;
}

/**
 * Renders a digital location logo based on the service name.
 * Falls back to a cloud icon if no matching logo is found.
 *
 * @example
 * ```tsx
 * <DigitalLocationIcon name="steam" className="h-4 w-4" />
 * ```
 */
export const DigitalLocationIcon = memo(function DigitalLocationIcon({
  name,
  className = "h-4 w-4"
}: DigitalLocationIconProps) {
  const { games } = useDomainMaps();

  console.log('Input name:', name);

  const normalizedName = name ? normalizeDigitalLocationName(name) : undefined;

  console.log('Normalized name:', normalizedName);
  console.log('Available games:', games);

  const LogoComponent = normalizedName ? games[normalizedName] : null;

  console.log('Found Logo component:', LogoComponent);

  return LogoComponent ? (
    <LogoComponent className={className} />
  ) : (
    <IconCloudDataConnection className={className} />
  );
});
