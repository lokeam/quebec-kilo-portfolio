import type { AnalyticsResponseWrapper, LocationSummary } from '@/core/api/services/analytics.service';
import type { MediaStorageMetadata } from '@/types/api/storage';
import type { PhysicalLocation } from '@/types/domain/physical-location';
import type { DigitalLocation } from '@/types/domain/digital-location';
import type { GamePlatform } from '@/types/domain/game-platform';
import type { GameItem } from '@/types/domain/game-item';
import type { Sublocation } from '@/types/domain/sublocation';
import { GamePlatform as GamePlatformEnum } from '@/types/domain/game-platform';
import { PhysicalLocationType, SublocationType } from '@/types/domain/location-types';

export function adaptAnalyticsToStorageMetadata(analyticsData: AnalyticsResponseWrapper): MediaStorageMetadata {
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

  const physicalItemCount = storage.physicalLocations?.reduce((sum: number, loc: LocationSummary) => sum + loc.itemCount, 0) || 0;
  const digitalItemCount = storage.digitalLocations?.reduce((sum: number, loc: LocationSummary) => sum + loc.itemCount, 0) || 0;

  // Create byLocation map
  const byLocation: Record<string, { total: number; inSublocations: number }> = {};

  // Add physical locations
  storage.physicalLocations?.forEach((loc: LocationSummary) => {
    byLocation[loc.id] = {
      total: loc.itemCount,
      inSublocations: 0 // NOTE: We don't have this info in analytics data
    };
  });

  // Add digital locations
  storage.digitalLocations?.forEach((loc: LocationSummary) => {
    byLocation[loc.id] = {
      total: loc.itemCount,
      inSublocations: 0 // NOTE: Digital locations don't have sublocations
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

export function adaptAnalyticsToPhysicalLocations(analyticsData: AnalyticsResponseWrapper): PhysicalLocation[] {
  const storage = analyticsData.storage;
  if (!storage?.physicalLocations) {
    return [];
  }

  return storage.physicalLocations.map((location: LocationSummary) => ({
    id: location.id,
    name: location.name,
    type: PhysicalLocationType.HOUSE,
    sublocations: location.sublocations?.map(subloc => ({
      id: subloc.id,
      name: subloc.name,
      type: subloc.locationType as SublocationType,
      parentLocationId: location.id,
      metadata: {
        bgColor: subloc.bgColor,
        notes: subloc.storedItems?.toString()
      },
      items: subloc.items || [],
      createdAt: new Date(subloc.createdAt),
      updatedAt: new Date(subloc.updatedAt)
    })) || [],
    items: [],
    createdAt: new Date(),
    updatedAt: new Date()
  }));
}

export function adaptAnalyticsToDigitalLocations(analyticsData: AnalyticsResponseWrapper): DigitalLocation[] {
  const storage = analyticsData.storage;
  if (!storage?.digitalLocations) {
    return [];
  }

  return storage.digitalLocations.map((location: LocationSummary) => {
    // Map the location type to a valid GamePlatform value
    let platform: GamePlatform;
    switch (location.locationType.toLowerCase()) {
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
      id: location.id,
      name: location.name,
      type: platform,
      isSubscription: location.isSubscription || false,
      monthlyCost: location.monthlyCost || 0,
      items: location.items || [],
      createdAt: location.createdAt,
      updatedAt: location.updatedAt,
      isActive: location.isActive,
      url: location.url,
      paymentMethod: location.paymentMethod,
      paymentDate: location.paymentDate,
      billingCycle: location.billingCycle,
      costPerCycle: location.costPerCycle,
      nextPaymentDate: location.nextPaymentDate
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
export function adaptAnalyticsToUIFormat(analyticsData: AnalyticsResponseWrapper): {
  physicalCards: LocationCardData[];
  digitalCards: LocationCardData[];
  metadata: MediaStorageMetadata;
} {
  const physicalLocations = adaptAnalyticsToPhysicalLocations(analyticsData);
  const digitalLocations = adaptAnalyticsToDigitalLocations(analyticsData);

  const physicalCards = physicalLocations.flatMap(adaptPhysicalLocationToCardData);
  const digitalCards = digitalLocations.map(adaptDigitalLocationToCardData);

  return {
    physicalCards,
    digitalCards,
    metadata: adaptAnalyticsToStorageMetadata(analyticsData)
  };
}