import { Home, Building2, Warehouse, Building } from 'lucide-react';
import { IconCar } from '@tabler/icons-react';
import type { PhysicalLocationType } from '@/types/domain/location-types';
import type { LocationIconBgColor } from '@/types/domain/location-types';
import { useLocationBgColor } from './getLocationBgColor';

interface PhysicalLocationIconProps {
  type: PhysicalLocationType | string | undefined;
  bgColor?: LocationIconBgColor;
}

/**
 * Renders the appropriate icon component for a physical location type with background color
 *
 * @param type - The physical location type (house, apartment, office, warehouse, vehicle)
 * @param bgColor - Optional background color for the icon
 * @returns A React component for the icon, or null if no matching icon is found
 *
 * @example
 * ```tsx
 * <PhysicalLocationIcon type="house" bgColor="blue" />
 * ```
 */
export function PhysicalLocationIcon({ type, bgColor }: PhysicalLocationIconProps) {
  const { background } = useLocationBgColor(bgColor);
  const wrapperStyle = {
    backgroundColor: background,
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 6,
    width: 42,
    height: 42,
    padding: 3,
    marginRight: 8,
  } as React.CSSProperties;

  if (!type) return null;

  const normalizedType = type.toLowerCase();
  let IconComponent;

  switch (normalizedType) {
    case 'house':
      IconComponent = <Home className="h-6 w-6" />;
      break;
    case 'apartment':
      IconComponent = <Building2 className="h-6 w-6" />;
      break;
    case 'office':
      IconComponent = <Building className="h-6 w-6" />;
      break;
    case 'warehouse':
      IconComponent = <Warehouse className="h-6 w-6" />;
      break;
    case 'vehicle':
      IconComponent = <IconCar className="h-6 w-6" />;
      break;
    default:
      return null;
  }

  return <span style={wrapperStyle}>{IconComponent}</span>;
}
