import type { PhysicalLocation as DomainPhysicalLocation } from '@/types/domain/physical-location';
import type { PhysicalLocation as UIPhysicalLocation } from '@/types/api/storage';
import type { Sublocation as DomainSublocation } from '@/types/domain/sublocation';
import type { LocationCardData } from '@/types/domain/location-card';
import type { GameItem as DomainGameItem } from '@/types/domain/game-item';
import type { UIGameItem } from '@/types/api/storage';
import { SublocationType } from '@/types/domain/location-types';
import type { DigitalLocation as DomainDigitalLocation } from '@/types/domain/digital-location';
import type { UIDigitalLocation } from '@/types/api/storage';
import type { GamePlatform } from '@/types/domain/game-platform';

/**
 * Media Storage Adapter
 * Handles transformations for physical storage locations and their sublocations
 * from API response format to frontend component format
 */

interface SublocationMetadata {
  bgColor?: string;
  shelf?: string;
  box?: string;
  notes?: string;
}

interface DigitalLocationMetadata {
  url?: string;
  isActive?: boolean;
  isSubscriptionService?: boolean;
  billing?: {
    monthlyFee?: string;
    nextRenewalDate?: string;
  };
  features?: string[];
  storage?: {
    used: number;
    total: number;
    unit: 'GB' | 'TB';
  };
}

/**
 * Transforms a sublocation into card data format
 * @param sublocation - The sublocation to transform
 * @returns LocationCardData object for the sublocation
 */
function adaptSublocationToCardData(sublocation: DomainSublocation): LocationCardData {
  const metadata = sublocation.metadata as SublocationMetadata;
  return {
    id: sublocation.id,
    name: sublocation.name,
    description: sublocation.description as string | undefined,
    locationType: sublocation.type as SublocationType,
    bgColor: metadata?.bgColor,
    items: (sublocation.items || []) as DomainGameItem[],
  };
}

/**
 * Transforms a physical location's sublocations into card data format
 * @param location - The physical location containing sublocations
 * @returns Array of LocationCardData objects for each sublocation
 */
export function adaptPhysicalLocationToCardData(location: DomainPhysicalLocation): LocationCardData[] {
  if (!location.sublocations) {
    return [];
  }

  return location.sublocations.map(sublocation => adaptSublocationToCardData(sublocation));
}

/**
 * Transforms a physical location into a format suitable for the frontend
 * @param location - The physical location to transform
 * @returns Transformed physical location with card data
 */
export function adaptPhysicalLocation(location: DomainPhysicalLocation) {
  return {
    ...location,
    cardData: adaptPhysicalLocationToCardData(location),
  };
}

/**
 * Transforms an array of physical locations into a format suitable for the frontend
 * @param locations - Array of physical locations to transform
 * @returns Array of transformed physical locations with card data
 */
export function adaptPhysicalLocations(locations: DomainPhysicalLocation[]) {
  return locations.map(location => adaptPhysicalLocation(location));
}

/**
 * Transforms a sublocation into a format suitable for the frontend
 * @param sublocation - The sublocation to transform
 * @returns Transformed sublocation with card data
 */
export function adaptSublocation(sublocation: DomainSublocation) {
  return {
    ...sublocation,
    cardData: adaptSublocationToCardData(sublocation),
  };
}

/**
 * Transforms an array of sublocations into a format suitable for the frontend
 * @param sublocations - Array of sublocations to transform
 * @returns Array of transformed sublocations with card data
 */
export function adaptSublocations(sublocations: DomainSublocation[]) {
  return sublocations.map(sublocation => adaptSublocation(sublocation));
}

/**
 * Converts a domain GameItem to a UI GameItem
 */
function adaptDomainToUIGameItem(item: DomainGameItem): UIGameItem {
  return {
    ...item,
    type: 'digital' as const
  };
}

/**
 * Converts a domain DigitalLocation to a UI DigitalLocation
 */
export function adaptDomainToUIDigitalLocation(location: DomainDigitalLocation): UIDigitalLocation {
  const metadata = location.metadata as DigitalLocationMetadata;
  return {
    id: location.id,
    name: location.name,
    label: location.name.toLowerCase().replace(/\s+/g, '-'),
    type: 'digital' as const,
    description: location.description,
    platform: location.type as GamePlatform,
    url: metadata?.url || '',
    subscription: {
      isActive: metadata?.isActive || false,
      isFree: !metadata?.isSubscriptionService,
      monthlyFee: metadata?.billing?.monthlyFee,
      renewalDate: metadata?.billing?.nextRenewalDate ? new Date(metadata.billing.nextRenewalDate) : undefined,
      benefits: metadata?.features || []
    },
    lastSync: location.updatedAt ? new Date(location.updatedAt) : undefined,
    totalStorage: metadata?.storage,
    createdAt: new Date(location.createdAt),
    updatedAt: new Date(location.updatedAt),
    items: (location.items || []).map(adaptDomainToUIGameItem),
    itemsInSublocations: 0
  };
}

/**
 * Converts a UI DigitalLocation to a domain DigitalLocation
 */
export function adaptUIToDomainDigitalLocation(location: UIDigitalLocation): DomainDigitalLocation {
  return {
    id: location.id,
    name: location.name,
    description: location.description,
    type: location.platform,
    metadata: {
      url: location.url,
      isActive: location.subscription?.isActive,
      isSubscriptionService: !location.subscription?.isFree,
      billing: {
        monthlyFee: location.subscription?.monthlyFee,
        nextRenewalDate: location.subscription?.renewalDate?.toISOString()
      },
      features: location.subscription?.benefits,
      storage: location.totalStorage
    },
    createdAt: location.createdAt.toISOString(),
    updatedAt: location.updatedAt.toISOString(),
    items: location.items || []
  };
}

/**
 * Converts an array of domain DigitalLocations to UI DigitalLocations
 */
export function adaptDomainToUIDigitalLocations(locations: DomainDigitalLocation[]): UIDigitalLocation[] {
  return locations.map(adaptDomainToUIDigitalLocation);
}

/**
 * Converts an array of UI DigitalLocations to domain DigitalLocations
 */
export function adaptUIToDomainDigitalLocations(locations: UIDigitalLocation[]): DomainDigitalLocation[] {
  return locations.map(adaptUIToDomainDigitalLocation);
}

/**
 * Transforms a domain PhysicalLocation to a UI PhysicalLocation
 */
export function adaptDomainToUIPhysicalLocation(location: DomainPhysicalLocation): UIPhysicalLocation {
  return {
    id: location.id,
    name: location.name,
    created_at: location.createdAt.toISOString(),
    updated_at: location.updatedAt.toISOString(),
    item_count: location.sublocations?.reduce((total, sublocation) => total + (sublocation.items?.length || 0), 0) || 0,
    location_type: location.type as 'house' | 'apartment' | 'office' | 'warehouse' | 'vehicle',
    sublocations: location.sublocations?.map((sublocation: DomainSublocation) => ({
      id: sublocation.id,
      name: sublocation.name,
      bg_color: (sublocation.metadata?.bgColor || 'gray') as 'red' | 'green' | 'blue' | 'orange' | 'gold' | 'purple' | 'brown' | 'gray',
      created_at: sublocation.createdAt.toISOString(),
      updated_at: sublocation.updatedAt.toISOString(),
      location_type: sublocation.type as 'shelf' | 'console' | 'cabinet' | 'closet' | 'drawer' | 'box' | 'device',
      stored_items: sublocation.items?.length || 0
    })) || []
  };
}

/**
 * Transforms an array of domain PhysicalLocations to UI PhysicalLocations
 */
export function adaptDomainToUIPhysicalLocations(locations: DomainPhysicalLocation[]): UIPhysicalLocation[] {
  return locations.map(adaptDomainToUIPhysicalLocation);
}