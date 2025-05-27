import { Box, Columns, BookOpen, BookmarkIcon } from 'lucide-react';
import { SublocationType } from '@/features/dashboard/lib/types/media-storage/constants';

/**
 * Gets an icon component for the specified sublocation type
 */
export function getSublocationTypeIcon(type: string): JSX.Element {
  // Different icon mapping based on sublocation type
  switch (type) {
    case SublocationType.shelf:
      return <Columns className="h-4 w-4 mr-1" />;
    case SublocationType.drawer:
      return <Box className="h-4 w-4 mr-1" />;
    case SublocationType.box:
      return <BookmarkIcon className="h-4 w-4 mr-1" />;
    case SublocationType.cabinet:
      return <BookOpen className="h-4 w-4 mr-1" />;
    default:
      return <Box className="h-4 w-4 mr-1" />;
  }
}