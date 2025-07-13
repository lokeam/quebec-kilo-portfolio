import { Package } from 'lucide-react';
import { BookshelfIcon } from '@/shared/components/ui/CustomIcons/BookshelfIcon';
import { MediaConsoleIcon } from '@/shared/components/ui/CustomIcons/MediaConsoleIcon';
import { DrawerIcon } from '@/shared/components/ui/CustomIcons/DrawerIcon';
import { CabinetIcon } from '@/shared/components/ui/CustomIcons/CabinetIcon';
import { ClosetIcon } from '@/shared/components/ui/CustomIcons/ClosetIcon';

interface SublocationIconProps {
  type: string;
  size?: number;
  className?: string;
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
 * @param size - Optional size for the icon
 * @param className - Optional className for the wrapper span
 * @returns A React component for the icon, or a default icon if no matching icon is found
 *
 * @example
 * ```tsx
 * <SublocationIcon type="box" size={80} className="custom-class" />
 * ```
 */
export function SublocationIconAlt({ type, size = 80, className }: SublocationIconProps) {
  // Transparent background, large icon
  const wrapperStyle = {
    backgroundColor: 'transparent',
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    width: size,
    height: size,
    marginRight: 16,
  } as React.CSSProperties;

  let IconComponent;
  switch (type) {
    case SublocationType.shelf:
      IconComponent = <BookshelfIcon size={size} color='#fff' />;
      break;
    case SublocationType.drawer:
      IconComponent = <DrawerIcon size={size} color='#fff' />;
      break;
    case SublocationType.box:
      IconComponent = <Package size={size} color='#fff' />;
      break;
    case SublocationType.cabinet:
      IconComponent = <CabinetIcon size={size} color='#fff' />;
      break;
    case SublocationType.closet:
      IconComponent = <ClosetIcon size={size} color='#fff' />;
      break;
    case SublocationType.console:
      IconComponent = <MediaConsoleIcon size={size} color='#fff' />;
      break;
    default:
      IconComponent = <Package size={size} color='#fff' />;
  }

  return <span style={wrapperStyle} className={className}>{IconComponent}</span>;
}