import type { DigitalLocation, PhysicalLocation, StorageAnalytics } from '@/types/api/storage';
import type { LocationCardData } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/MediaStoragePageAccordionCard';
import { SublocationType } from '@/features/dashboard/lib/types/media-storage/constants';

/**
 * Media Storage Adapter
 * Handles transformations for all media storage related data from API to frontend
 */

/**
 * Transforms a physical location's sublocations into card data format
 */
function adaptPhysicalLocationToCardData(location: PhysicalLocation): LocationCardData[] {
  return location.sublocations.map(sublocation => ({
    id: sublocation.id,
    name: sublocation.name,
    locationType: sublocation.location_type as SublocationType,
    bgColor: sublocation.bg_color,
    items: [], // TODO: Add items when available in API response
  }));
}

/**
 * Transforms a digital location into card data format
 */
function adaptDigitalLocationToCardData(location: DigitalLocation): LocationCardData {
  return {
    id: location.id,
    name: location.name,
    locationType: SublocationType.box, // Default to box for digital locations
    items: [], // TODO: Add items when available in API response
  };
}

/**
 * Transforms a location (physical or digital) into card data format
 */
export function adaptLocationToCardData(
  location: PhysicalLocation | DigitalLocation,
  type: 'physical' | 'digital'
): LocationCardData[] {
  if (type === 'physical') {
    return adaptPhysicalLocationToCardData(location as PhysicalLocation);
  }
  return [adaptDigitalLocationToCardData(location as DigitalLocation)];
}

/**
 * Transforms storage analytics data into a format suitable for the frontend
 */
export function adaptStorageAnalytics(analytics: StorageAnalytics) {
  return {
    digitalLocations: analytics.storage.digital_locations.map(location => ({
      ...location,
      cardData: adaptDigitalLocationToCardData(location)
    })),
    physicalLocations: analytics.storage.physical_locations.map(location => ({
      ...location,
      cardData: adaptPhysicalLocationToCardData(location)
    })),
    totalDigitalLocations: analytics.storage.total_digital_locations,
    totalPhysicalLocations: analytics.storage.total_physical_locations
  };
}