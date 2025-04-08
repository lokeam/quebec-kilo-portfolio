import { Box, Columns, BookOpen, BookmarkIcon } from 'lucide-react';
import { SubLocationType } from '@/features/dashboard/lib/types/media-storage/constants';

/**
 * Gets an icon component for the specified sublocation type
 */
export function getSublocationTypeIcon(type: string): JSX.Element {
  // Different icon mapping based on sublocation type
  switch (type) {
    case SubLocationType.SHELF:
      return <Columns className="h-4 w-4 mr-1" />;
    case SubLocationType.DRAWER:
      return <Box className="h-4 w-4 mr-1" />;
    case SubLocationType.BOX:
      return <BookmarkIcon className="h-4 w-4 mr-1" />;
    case SubLocationType.CABINET:
      return <BookOpen className="h-4 w-4 mr-1" />;
    default:
      return <Box className="h-4 w-4 mr-1" />;
  }
}