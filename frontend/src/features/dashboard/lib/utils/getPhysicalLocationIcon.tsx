import { Home, Building2, Warehouse, Building } from 'lucide-react';
import type { PhysicalLocationType } from '@/types/domain/location-types';

/**
 * Gets the appropriate icon component for a physical location type
 *
 * @param type - The physical location type (house, apartment, office, warehouse)
 * @returns A React component for the icon, or null if no matching icon is found
 *
 * @example
 * ```tsx
 * const icon = getPhysicalLocationIcon('house');
 * return <div>{icon}</div>;
 * ```
 */
export function getPhysicalLocationIcon(type: PhysicalLocationType | string | undefined) {
  if (!type) return null;

  const normalizedType = type.toLowerCase();

  switch (normalizedType) {
    case 'house':
      return <Home className="h-4 w-4" />;
    case 'apartment':
      return <Building2 className="h-4 w-4" />;
    case 'office':
      return <Building className="h-4 w-4" />;
    case 'warehouse':
      return <Warehouse className="h-4 w-4" />;
    default:
      return null;
  }
}
