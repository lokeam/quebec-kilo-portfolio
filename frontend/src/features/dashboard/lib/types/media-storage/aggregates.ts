import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';
import type { DigitalLocation } from '@/features/dashboard/lib/types/media-storage/digital';
import type { PhysicalGameItem } from '@/features/dashboard/lib/types/media-storage/items';
import type { DigitalGameItem } from '@/features/dashboard/lib/types/media-storage/items';

// Additional exports for commonly used type combinations
export type StorageLocation = PhysicalLocation | DigitalLocation;
export type GameItemType = PhysicalGameItem | DigitalGameItem;
