import type { AnalyticsResponse } from '../services/analytics.service';
import type { MediaStorageMetadata } from '@/types/api/storage';
import type { PhysicalLocation } from '@/types/domain/physical-location';
import type { DigitalLocation } from '@/types/domain/digital-location';
import type { GamePlatform } from '@/types/domain/game-platform';
import type { GameItem } from '@/types/domain/game-item';
import type { Sublocation } from '@/types/domain/sublocation';
import { GamePlatform as GamePlatformEnum } from '@/types/domain/game-platform';
import { PhysicalLocationType, SublocationType } from '@/types/domain/location-types';

export function adaptAnalyticsToStorageMetadata(analyticsData: AnalyticsResponse): MediaStorageMetadata {
  const storage = analyticsData.storage;
  if (!storage) {
    return {
      counts: {
        locations: {
          total: 0,
          physical: 0,
          digital: 0
        },
        items: {
          total: 0,
          physical: 0,
          digital: 0,
          byLocation: {}
        }
      },
      lastUpdated: new Date(),
      version: '1.0'
    };
  }

  const physicalItemCount = storage.physicalLocations?.reduce((sum, loc) => sum + loc.itemCount, 0) || 0;
  const digitalItemCount = storage.digitalLocations?.reduce((sum, loc) => sum + loc.itemCount, 0) || 0;

  // Create byLocation map
  const byLocation: Record<string, { total: number; inSublocations: number }> = {};

  // Add physical locations
  storage.physicalLocations?.forEach(loc => {
    byLocation[loc.id] = {
      total: loc.itemCount,
      inSublocations: 0 // We don't have this info in analytics data
    };
  });

  // Add digital locations
  storage.digitalLocations?.forEach(loc => {
    byLocation[loc.id] = {
      total: loc.itemCount,
      inSublocations: 0 // Digital locations don't have sublocations
    };
  });

  return {
    counts: {
      locations: {
        total: (storage.totalPhysicalLocations || 0) + (storage.totalDigitalLocations || 0),
        physical: storage.totalPhysicalLocations || 0,
        digital: storage.totalDigitalLocations || 0
      },
      items: {
        total: physicalItemCount + digitalItemCount,
        physical: physicalItemCount,
        digital: digitalItemCount,
        byLocation
      }
    },
    lastUpdated: new Date(),
    version: '1.0'
  };
}

export const adaptAnalyticsToPhysicalLocations = (
  analyticsData: AnalyticsResponse
): PhysicalLocation[] => {
  if (!analyticsData.storage?.physicalLocations) {
    return [];
  }

  return analyticsData.storage.physicalLocations.map(location => ({
    id: location.id,
    name: location.name,
    type: location.locationType as PhysicalLocationType,
    description: undefined,
    metadata: undefined,
    mapCoordinates: location.mapCoordinates,
    sublocations: location.sublocations?.map(subloc => ({
      id: subloc.id,
      name: subloc.name,
      type: subloc.locationType as SublocationType,
      parentLocationId: location.id,
      description: undefined,
      metadata: {
        bgColor: subloc.bgColor,
        notes: undefined
      },
      items: [],
      createdAt: new Date(subloc.createdAt),
      updatedAt: new Date(subloc.updatedAt)
    })) || [],
    createdAt: new Date(location.created_at),
    updatedAt: new Date(location.updated_at)
  }));
};

export function adaptAnalyticsToDigitalLocations(analyticsData: AnalyticsResponse): DigitalLocation[] {
  if (!analyticsData.storage?.digitalLocations) {
    return [];
  }

  return analyticsData.storage.digitalLocations.map(loc => {
    // Map the location type to a valid GamePlatform value
    let platform: GamePlatform;
    switch (loc.locationType.toLowerCase()) {
      case 'steam':
        platform = GamePlatformEnum.STEAM;
        break;
      case 'epic':
        platform = GamePlatformEnum.EPIC;
        break;
      case 'gog':
        platform = GamePlatformEnum.GOG;
        break;
      case 'playstation':
      case 'sony':
        platform = GamePlatformEnum.PLAYSTATION;
        break;
      case 'xbox':
        platform = GamePlatformEnum.XBOX;
        break;
      case 'nintendo':
        platform = GamePlatformEnum.NINTENDO;
        break;
      default:
        // Default to Steam if unknown
        platform = GamePlatformEnum.STEAM;
    }

    return {
      id: loc.id,
      name: loc.name,
      type: platform,
      description: undefined,
      metadata: undefined,
      items: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
  });
}

// New types for UI components
export interface LocationCardData {
  id: string;
  name: string;
  description?: string;
  locationType: string;
  bgColor?: string;
  items: GameItem[];
  sublocations: Sublocation[];
  mapCoordinates?: string;
  createdAt: Date;
  updatedAt: Date;
  platform?: string; // For digital locations
}

// Transform physical location to card data
export function adaptPhysicalLocationToCardData(location: PhysicalLocation): LocationCardData[] {
  return (location.sublocations || []).map(subloc => ({
    id: subloc.id,
    name: subloc.name,
    description: subloc.description,
    locationType: subloc.type,
    bgColor: subloc.metadata?.bgColor,
    items: (subloc.items || []) as GameItem[],
    sublocations: [],
    mapCoordinates: subloc.metadata?.notes,
    createdAt: subloc.createdAt,
    updatedAt: subloc.updatedAt
  }));
}

// Transform digital location to card data
export function adaptDigitalLocationToCardData(location: DigitalLocation): LocationCardData {
  return {
    id: location.id,
    name: location.name,
    description: location.description,
    locationType: location.type,
    items: (location.items || []) as GameItem[],
    sublocations: [],
    createdAt: new Date(location.createdAt),
    updatedAt: new Date(location.updatedAt),
    platform: location.type // For digital locations, type is the platform
  };
}

// Transform analytics data to UI-ready format
export function adaptAnalyticsToUIFormat(analyticsData: AnalyticsResponse): {
  physicalCards: LocationCardData[];
  digitalCards: LocationCardData[];
  metadata: MediaStorageMetadata;
} {
  const physicalLocations = adaptAnalyticsToPhysicalLocations(analyticsData);
  const digitalLocations = adaptAnalyticsToDigitalLocations(analyticsData);

  // Transform physical locations to card data
  const physicalCards = physicalLocations.flatMap(loc =>
    adaptPhysicalLocationToCardData(loc)
  );

  // Transform digital locations to card data
  const digitalCards = digitalLocations.map(loc =>
    adaptDigitalLocationToCardData(loc)
  );

  // Get metadata
  const metadata = adaptAnalyticsToStorageMetadata(analyticsData);

  return {
    physicalCards,
    digitalCards,
    metadata
  };
}