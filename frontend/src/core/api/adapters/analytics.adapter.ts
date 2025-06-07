import type { AnalyticsResponseWrapper, LocationSummary } from '@/core/api/services/analytics.service';
import type { MediaStorageMetadata } from '@/types/api/storage';
import type { PhysicalLocation } from '@/types/domain/physical-location';
import type { DigitalLocation } from '@/types/domain/digital-location';
import type { GamePlatform } from '@/types/domain/game-platform';
import type { GameItem } from '@/types/domain/game-item';
import type { Sublocation } from '@/types/domain/sublocation';
import { GamePlatform as GamePlatformEnum } from '@/types/domain/game-platform';
import { PhysicalLocationType, SublocationType } from '@/types/domain/location-types';
import type { LocationIconBgColor } from '@/types/domain/location-types';

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

  console.log('[DEBUG] Analytics Data:', analyticsData);
  console.log('[DEBUG] Storage Data:', storage);
  console.log('[DEBUG] Physical Locations:', storage.physicalLocations);

  return storage.physicalLocations.map((location: LocationSummary) => {
    // Get bg_color from the raw data
    const rawData = location as unknown as { bgColor?: string };
    const bgColor = rawData.bgColor as LocationIconBgColor | undefined;

    const transformed = {
      id: location.id,
      name: location.name,
      locationType: location.locationType as PhysicalLocationType,  // This is the PhysicalLocationType
      bgColor,
      mapCoordinates: location.mapCoordinates ? {
        coords: location.mapCoordinates.coords,
        googleMapsLink: location.mapCoordinates.googleMapsLink
      } : undefined,
      sublocations: location.sublocations?.map(subloc => ({
        id: subloc.id,
        name: subloc.name,
        type: subloc.locationType as SublocationType,
        parentLocationId: location.id,
        metadata: {
          bgColor: location.bgColor,
          notes: subloc.storedItems?.toString()
        },
        items: subloc.items || [],
        createdAt: new Date(subloc.createdAt),
        updatedAt: new Date(subloc.updatedAt)
      })) || [],
      items: [],
      createdAt: new Date(location.createdAt),
      updatedAt: new Date(location.updatedAt)
    };

    console.log('[DEBUG] Transformed Location:', transformed);
    return transformed;
  });
}

/**
 * Maps backend payment method IDs to PaymentIcon types
 */
const PAYMENT_METHOD_MAP: Record<string, string> = {
  'alipay': 'Alipay',
  'amex': 'Amex',
  'diners': 'Diners',
  'discover': 'Discover',
  'elo': 'Elo',
  'generic': 'Generic',
  'hiper': 'Hiper',
  'hipercard': 'Hipercard',
  'jcb': 'Jcb',
  'maestro': 'Maestro',
  'mastercard': 'Mastercard',
  'mir': 'Mir',
  'paypal': 'Paypal',
  'unionpay': 'Unionpay',
  'visa': 'Visa'
};

export function adaptAnalyticsToDigitalLocations(
  analyticsData: AnalyticsResponseWrapper,
): DigitalLocation[] {
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

    // Transform payment method to the format expected by PaymentIcon
    const paymentMethod = location.paymentMethod
      ? PAYMENT_METHOD_MAP[location.paymentMethod.toLowerCase()] || 'Generic'
      : 'Generic';

    return {
      id: location.id,
      name: location.name,
      type: platform,
      isSubscription: location.isSubscription || false,
      monthlyCost: location.monthlyCost || 0,
      items: (location.items || []).map(item => ({
        ...item,
        id: item.id.toString(),
        label: item.name.toLowerCase().replace(/\s+/g, '-'),
        acquiredDate: new Date(item.acquiredDate)  // Convert string to Date
      })) as GameItem[],
      createdAt: location.createdAt,
      updatedAt: location.updatedAt,
      isActive: Boolean(location.isActive),
      url: location.url || '',
      paymentMethod,
      paymentDate: location.paymentDate?.toString() || '',
      billingCycle: location.billingCycle || '',
      costPerCycle: location.costPerCycle || 0,
      nextPaymentDate: location.nextPaymentDate?.toString() || ''
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
    bgColor: location.bgColor,
    items: [],
    sublocations: [],
    mapCoordinates: location.mapCoordinates?.googleMapsLink || '',
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

// Types for flattened sublocation data
export type SublocationItemData = {
  id?: string;
  name?: string;
  locationType?: string;
  sublocationId: string;
  sublocationName: string;
  sublocationType: string;
  storedItems: number;
  parentLocationId?: string;
  parentLocationName: string;
  parentLocationType: string;
  parentLocationBgColor: LocationIconBgColor;
  mapCoordinates: {
    coords: string;
    googleMapsLink: string;
  };
  createdAt: string;
  updatedAt: string;
};

/**
 * Transforms physical locations into a flattened array of SublocationItemData
 * This is the single source of truth for sublocation data transformation
 */
export function adaptPhysicalLocationsToSublocationRows(physicalLocations: PhysicalLocation[]): SublocationItemData[] {
  return physicalLocations.flatMap(location =>
    (location.sublocations || []).map(sublocation => ({
      sublocationId: sublocation.id,
      sublocationName: sublocation.name,
      sublocationType: sublocation.type,
      storedItems: sublocation.metadata?.notes ? parseInt(sublocation.metadata.notes) : 0,
      parentLocationId: location.id,
      parentLocationName: location.name,
      parentLocationType: location.locationType,
      parentLocationBgColor: location.bgColor as LocationIconBgColor,
      mapCoordinates: {
        coords: location.mapCoordinates?.coords || '',
        googleMapsLink: location.mapCoordinates?.googleMapsLink || ''
      },
      createdAt: sublocation.createdAt.toISOString(),
      updatedAt: sublocation.updatedAt.toISOString()
    }))
  );
}