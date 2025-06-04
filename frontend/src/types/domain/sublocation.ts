import type { SublocationType } from './location-types';

export interface SublocationMetadata {
  bgColor?: string;
  shelf?: string;
  box?: string;
  notes?: string;
}

export interface Sublocation {
  id: string;
  name: string;
  type: SublocationType;
  parentLocationId: string;
  description?: string;
  metadata?: SublocationMetadata;
  items?: unknown[];
  storedItems?: number;
  createdAt: Date;
  updatedAt: Date;
}

/**
 * Request type for creating a new sublocation
 */
export interface CreateSublocationRequest {
  name: string;
  type: SublocationType;
  parentLocationId: string;
  description?: string;
  metadata?: SublocationMetadata;
}