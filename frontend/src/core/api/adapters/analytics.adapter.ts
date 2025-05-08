import type { AnalyticsResponse } from '../services/analytics.service';
import type { MediaStorageMetadata } from '@/types/api/storage';
import type { PhysicalLocation } from '@/types/domain/physical-location';
import type { DigitalLocation } from '@/types/domain/digital-location';
import type { GamePlatform } from '@/types/domain/game-platform';
import { GamePlatform as GamePlatformEnum } from '@/types/domain/game-platform';
import { PhysicalLocationType } from '@/types/domain/location-types';

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

export function adaptAnalyticsToPhysicalLocations(analyticsData: AnalyticsResponse): PhysicalLocation[] {
  if (!analyticsData.storage?.physicalLocations) {
    return [];
  }

  return analyticsData.storage.physicalLocations.map(loc => ({
    id: loc.id,
    name: loc.name,
    type: PhysicalLocationType.HOUSE, // Default to HOUSE since we don't have this info
    description: undefined,
    metadata: undefined,
    sublocations: [], // Analytics data doesn't include sublocations
    createdAt: new Date(),
    updatedAt: new Date()
  }));
}

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