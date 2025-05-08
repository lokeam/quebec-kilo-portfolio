import type { GamePlatform } from './game-platform';

/**
 * Base interface for all game items in media storage system.
 */
export interface GameItem {
  /** Unique identifier for the game item */
  id: string;

  /** Full game title */
  name: string;

  /** URL-friendly game identifier */
  label: string;

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

  /** Physical condition of the item (for physical games) */
  condition?: 'new' | 'excellent' | 'good' | 'fair' | 'poor';

  /** Whether the original case is present (for physical games) */
  hasOriginalCase?: boolean;

  /** Whether the manual is included (for physical games) */
  hasManual?: boolean;

  /** Product serial number (for physical games) */
  serialNumber?: string;

  /** Size of the game installation in megabytes (for digital games) */
  installationSize?: number;

  /** Whether the game is currently installed (for digital games) */
  isInstalled?: boolean;

  /** Timestamp of the last game update (for digital games) */
  lastUpdated?: Date;

  /** URL for downloading/reinstalling the game (for digital games) */
  downloadUrl?: string;
}

/**
 * Input type for creating new game items.
 * Omits server-generated fields from GameItem.
 */
export type CreateGameItemInput = Omit<GameItem, 'id'>;

/**
 * Input type for updating existing game items.
 * Makes all fields optional except id.
 */
export type UpdateGameItemInput = Partial<Omit<GameItem, 'id'>>;