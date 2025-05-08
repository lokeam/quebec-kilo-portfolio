/**
 * Digital Location Types
 *
 * Defines types for digital media storage locations.
 */

import type { GameItem } from './game-item';

/**
 * Represents a digital media storage location
 */
export interface DigitalLocation {
  /** Unique identifier for the location */
  id: string;

  /** Name of the digital location */
  name: string;

  /** Type of digital location (e.g., 'steam', 'epic', 'gog') */
  type: string;

  /** Optional description of the location */
  description?: string;

  /** Optional metadata for the location */
  metadata?: DigitalLocationMetadata;

  /** List of game items stored in this location */
  items?: GameItem[];

  /** Timestamp when the location was created */
  createdAt: string;

  /** Timestamp when the location was last updated */
  updatedAt: string;
}

/**
 * Metadata for a digital location
 */
export interface DigitalLocationMetadata {
  /** Platform-specific identifier */
  platformId?: string;

  /** Platform-specific username */
  username?: string;

  /** Platform-specific API key or token */
  apiKey?: string;

  /** Additional platform-specific metadata */
  [key: string]: unknown;
}

/**
 * Request type for creating a new digital location
 */
export interface CreateDigitalLocationRequest {
  /** Name of the digital location */
  name: string;

  /** Type of digital location */
  type: string;

  /** Optional description */
  description?: string;

  /** Optional metadata */
  metadata?: DigitalLocationMetadata;
}