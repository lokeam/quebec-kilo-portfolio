import type { GamePlatform } from '@/features/dashboard/lib/types/media-storage/constants';

/**
 * Base interface for all game items in media storage system.
 * Contains common properties shared between physical + digital games.
 *
 * @interface GameItem
 */
export interface GameItem {
  /** Unique identifier for the game item */
  id: string;

  /** Full game title */
  name: string;

  /** URL-friendly game identifier */
  label: string;

  /** Discriminator to distinguish between physical + digital items */
  type: 'physical' | 'digital';

  /** Platform the game is associated to */
  platform: GamePlatform;

  /** Specific version or generation of the platform */
  platformVersion: string;

  /** Date when the game was acquired */
  acquiredDate: Date;

  /** Last time game was played (optional) */
  lastPlayed?: Date;

  /** Total time spent playing in minutes (optional) */
  playtime?: number;
}

/**
 * Interface for physical game items.
 * Extends GameItem with properties specific to physical media.
 *
 * @interface PhysicalGameItem
 * @extends {GameItem}
 */
export interface PhysicalGameItem extends GameItem {
  /** Physical condition of the item */
  condition: 'new' | 'excellent' | 'good' | 'fair' | 'poor';

  /** Whether the original case is present */
  hasOriginalCase: boolean;

  /** Whether the manual is included */
  hasManual?: boolean;

  /** Product serial number (optional) */
  serialNumber?: string;
}

/**
 * Interface for digital game items.
 * Extends GameItem with properties specific to digital distribution.
 *
 * @interface DigitalGameItem
 * @extends {GameItem}
 */
export interface DigitalGameItem extends GameItem {
  /** Size of the game installation in megabytes */
  installationSize: number;

  /** Whether the game is currently installed */
  isInstalled: boolean;

  /** Timestamp of the last game update */
  lastUpdated?: Date;

  /** URL for downloading/reinstalling the game */
  downloadUrl?: string;
}

/**
 * Input type for creating new game items.
 * Omits server-generated fields from GameItem.
 *
 * @type CreateGameItemInput
 */
export type CreateGameItemInput = Omit<GameItem, 'id'>;

/**
 * Input type for updating existing game items.
 * Makes all fields optional except id.
 *
 * @type UpdateGameItemInput
 */
export type UpdateGameItemInput = Partial<Omit<GameItem, 'id'>>;
