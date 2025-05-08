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
  createdAt: Date;
  updatedAt: Date;
}