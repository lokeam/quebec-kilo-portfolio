// types/domain/location-with-cards.ts
import type { LocationCardData } from './location-card';

export interface LocationWithCardData {
  id: string;
  name: string;
  type: 'physical' | 'digital';
  cards: LocationCardData[];
  metadata?: {
    bgColor?: string;
    notes?: string;
  };
}