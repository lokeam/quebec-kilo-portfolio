import { memo } from 'react';
import {
  IconBoxMultiple2,
  IconBoxMultiple3,
  IconBoxMultiple4,
  IconBoxMultiple5,
  IconBoxMultiple6,
  IconBoxMultiple7,
  IconBoxMultiple8,
  IconBoxMultiple9,
} from '@/shared/components/ui/icons';

interface LibraryCountIconProps {
  /** The count of the library item */
  count: number;
  /** Optional className for styling the icon */
  className?: string;
}

export const LibraryCountIcon = memo(function LibraryCountIcon({
  count,
  className = "h-4 w-4"
}: LibraryCountIconProps) {
  // If count is 1 or less, no need to show multiple boxes
  if (count <= 1) {
    return null;
  }

  // Map count to appropriate icon
  const IconComponent = (() => {
    switch (count) {
      case 2:
        return IconBoxMultiple2;
      case 3:
        return IconBoxMultiple3;
      case 4:
        return IconBoxMultiple4;
      case 5:
        return IconBoxMultiple5;
      case 6:
        return IconBoxMultiple6;
      case 7:
        return IconBoxMultiple7;
      case 8:
        return IconBoxMultiple8;
      case 9:
        return IconBoxMultiple9;
      default:
        // For counts greater than 9, use IconBoxMultiple9
        return IconBoxMultiple9;
    }
  })();

  return <IconComponent className={className} />;
});
