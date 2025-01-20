import { useMemo } from 'react';
import { IconCloudDataConnection } from '@tabler/icons-react';
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';

interface LocationIconsProps {
  /** Physical location string (e.g., "Warehouse A") */
  physicalLocation?: string;
  /** Type of physical location (e.g., "warehouse", "store") */
  physicalLocationType?: string;
  /** Specific sublocation within the physical location */
  physicalSublocation?: string;
  /** Type of sublocation (e.g., "shelf", "bin") */
  physicalSublocationType?: string;
  /** Digital location or service URL */
  digitalLocation?: string;
}

/**
 * Custom hook for rendering location and sublocation icons based on location type.
 *
 * @param {LocationIconsProps} props - Location properties to determine which icons to display
 * @returns {{ locationIcon: ReactNode | null, subLocationIcon: ReactNode | null }} Object containing memoized location icons
 *
 * @example
 * ```tsx
 * const { locationIcon, subLocationIcon } = useLocationIcons({
 *   physicalLocation: "Warehouse A",
 *   physicalLocationType: "warehouse",
 *   physicalSublocationType: "shelf"
 * });
 * ```
 */
export function useLocationIcons({
  physicalLocation,
  physicalLocationType,
  digitalLocation,
  physicalSublocationType
}: LocationIconsProps) {
  const { location: locationIcons, sublocation: sublocationIcons } = useDomainMaps();

  const locationIcon = useMemo(() => {
    if (physicalLocation && physicalLocationType) {
      const IconComponent = locationIcons[physicalLocationType.toLowerCase() as keyof typeof locationIcons];
      return IconComponent != null ? <IconComponent className="h-7 w-7" /> : null;
    }

    if (digitalLocation) {
      return <IconCloudDataConnection className="h-7 w-7" />
    }
    return null;
  }, [physicalLocation, physicalLocationType, digitalLocation, locationIcons]);

  const subLocationIcon = useMemo(() => {
    if (physicalSublocationType) {
      const IconComponent = sublocationIcons[physicalSublocationType.toLowerCase() as keyof typeof sublocationIcons];
      return IconComponent ? <IconComponent className="h-6 w-6 mt-1" /> : null;
    }

    return null;
  }, [physicalSublocationType, sublocationIcons]);

  return { locationIcon, subLocationIcon };
}
