import type { GameItem } from './game-item';
import type { Sublocation } from './sublocation';
import type { SublocationType } from './location-types';

/**
 * Data structure for location cards in the media storage UI
 */
export interface LocationCardData {
  /** Unique identifier for the card */
  id: string;

  /** Display name of the location */
  name: string;

  /** Optional description of the location */
  description?: string;

  /** Type of sublocation (shelf, box, etc.) */
  locationType: SublocationType;

  /** Background color for the card */
  bgColor?: string;

  /** List of game items in this location */
  items?: GameItem[];

  /** List of sublocations within this location */
  sublocations?: Sublocation[];

  /** Platform identifier for digital locations */
  platform?: string;

  /** Map coordinates for the location */
  mapCoordinates?: string;

  /** Timestamp when the location was created */
  createdAt?: Date;

  /** Timestamp when the location was last updated */
  updatedAt?: Date;
}