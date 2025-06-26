import { Package } from 'lucide-react';
import { BookshelfIcon } from '@/shared/components/ui/CustomIcons/BookshelfIcon';
import { MediaConsoleIcon } from '@/shared/components/ui/CustomIcons/MediaConsoleIcon';
import { DrawerIcon } from '@/shared/components/ui/CustomIcons/DrawerIcon';
import { CabinetIcon } from '@/shared/components/ui/CustomIcons/CabinetIcon';
import { ClosetIcon } from '@/shared/components/ui/CustomIcons/ClosetIcon';
//import type { LocationIconBgColor } from '@/types/domain/location-types';
import { useLocationBgColor } from './getLocationBgColor';

interface SublocationIconProps {
  type: string;
  bgColor?: string;
}

/**
 * Types of storage subdivisions within physical locations.
 * Represents specific storage units or furniture.
 *
 * @constant SublocationType
 */
export const SublocationType = {
  shelf: 'shelf',
  console: 'console',
  cabinet: 'cabinet',
  closet: 'closet',
  drawer: 'drawer',
  box: 'box'
} as const;

export type SublocationType = typeof SublocationType[keyof typeof SublocationType];


/**
 * Renders the appropriate icon component for a sublocation type with background color
 *
 * @param type - The sublocation type (shelf, drawer, box, cabinet, closet, console)
 * @param bgColor - Optional background color for the icon
 * @returns A React component for the icon, or a default icon if no matching icon is found
 *
 * @example
 * ```tsx
 * <SublocationIcon type="box" bgColor="blue" />
 * ```
 */
export function SublocationIcon({ type, bgColor }: SublocationIconProps) {
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

  let IconComponent;
  switch (type) {
    case SublocationType.shelf:
      IconComponent = <BookshelfIcon size={30} color='#fff' />;
      break;
    case SublocationType.drawer:
      IconComponent = <DrawerIcon size={30} color='#fff' />;
      break;
    case SublocationType.box:
      IconComponent = <Package size={30} color='#fff' />;
      break;
    case SublocationType.cabinet:
      IconComponent = <CabinetIcon size={30} color='#fff' />;
      break;
    case SublocationType.closet:
      IconComponent = <ClosetIcon size={30} color='#fff' />;
      break;
    case SublocationType.console:
      IconComponent = <MediaConsoleIcon size={30} color='#fff' />;
      break;
    default:
      IconComponent = <Package size={30} color='#fff' />;
  }

  return <span style={wrapperStyle}>{IconComponent}</span>;
}