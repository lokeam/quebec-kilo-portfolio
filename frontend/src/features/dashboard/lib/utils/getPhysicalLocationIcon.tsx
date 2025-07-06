import { Home, Building2, Warehouse, Building } from 'lucide-react';
import { IconCar } from '@tabler/icons-react';
import type { PhysicalLocationType } from '@/types/domain/location-types';
import type { LocationIconBgColor } from '@/types/domain/location-types';
import { LOCATION_ICON_COLORS, DEFAULT_COLORS } from '../constants/location-icon-colors';

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
  // Use dark mode colors consistently for better visibility in both themes
  const background = bgColor
    ? LOCATION_ICON_COLORS[bgColor].dark.background
    : DEFAULT_COLORS.dark.background;

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
      IconComponent = <Home className="h-6 w-6 text-white" />;
      break;
    case 'apartment':
      IconComponent = <Building2 className="h-6 w-6 text-white" />;
      break;
    case 'office':
      IconComponent = <Building className="h-6 w-6 text-white" />;
      break;
    case 'warehouse':
      IconComponent = <Warehouse className="h-6 w-6 text-white" />;
      break;
    case 'vehicle':
      IconComponent = <IconCar className="h-6 w-6 text-white" />;
      break;
    default:
      return null;
  }

  return <span style={wrapperStyle}>{IconComponent}</span>;
}
