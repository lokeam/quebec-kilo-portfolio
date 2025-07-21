import { useMemo } from 'react';
import { IconCloudDataConnection } from '@/shared/components/ui/icons';
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';
import type { GamePlatformLocation } from '@/types/domain/library-types';

interface LocationIconsProps {
  /** Array of platform and location data */
  gamesByPlatformAndLocation?: GamePlatformLocation[];
  /** Optional index to select which platform/location to show */
  selectedIndex?: number;
}

/**
 * Custom hook for rendering location and sublocation icons based on location type.
 *
 * @param {LocationIconsProps} props - Platform and location data to determine which icons to display
 * @returns {{ locationIcon: ReactNode | null, subLocationIcon: ReactNode | null }} Object containing memoized location icons
 *
 * @example
 * ```tsx
 * const { locationIcon, subLocationIcon } = useLocationIcons({
 *   gamesByPlatformAndLocation: gameData.gamesByPlatformAndLocation,
 *   selectedIndex: 0
 * });
 * ```
 */
export function useLocationIcons({
  gamesByPlatformAndLocation = [],
  selectedIndex = 0
}: LocationIconsProps) {
  const { location: locationIcons, sublocation: sublocationIcons } = useDomainMaps();

  // Get the selected platform/location data
  const selectedLocation = gamesByPlatformAndLocation?.[selectedIndex];

  const locationIcon = useMemo(() => {
    if (!selectedLocation) return null;

    if (selectedLocation.Type === 'physical') {
      const IconComponent = locationIcons[selectedLocation.LocationType.toLowerCase() as keyof typeof locationIcons];
      return IconComponent != null ? <IconComponent className="h-7 w-7" /> : null;
    }

    if (selectedLocation.Type === 'digital') {
      return <IconCloudDataConnection className="h-7 w-7" />
    }
    return null;
  }, [selectedLocation, locationIcons]);

  const subLocationIcon = useMemo(() => {
    if (!selectedLocation?.SublocationType) return null;

    const key = selectedLocation.SublocationType.toLowerCase() as keyof typeof sublocationIcons;
    const IconComponent = sublocationIcons[key];
    return IconComponent ? <IconComponent className="h-7 w-7" /> : null;
  }, [selectedLocation, sublocationIcons]);

  return { locationIcon, subLocationIcon };
}
