import type { LocationsBFFResponse } from '@/types/domain/physical-location';
import type { SublocationRowData } from './analytics.adapter';

export function adaptBFFToSublocationRows(data: LocationsBFFResponse | undefined): SublocationRowData[] {
  if (!data?.physicalLocations || !data?.sublocations) {
    return [];
  }

  // Create a map of physical locations for quick lookup
  const physicalLocationMap = new Map(
    data.physicalLocations.map(loc => [loc.physicalLocationID, loc])
  );

  return data.sublocations.map(subloc => {
    const parentLocation = physicalLocationMap.get(subloc.parentLocationID);
    if (!parentLocation) {
      return null;
    }

    return {
      sublocationId: subloc.sublocationID,
      sublocationName: subloc.sublocationName,
      sublocationType: subloc.sublocationType,
      parentLocationId: parentLocation.physicalLocationID,
      parentLocationName: parentLocation.name,
      parentLocationType: parentLocation.physicalLocationType,
      mapCoordinates: {
        coords: parentLocation.mapCoordinates.coords,
        googleMapsLink: parentLocation.mapCoordinates.googleMapsLink
      },
      storedItems: subloc.storedItems,
      parentLocationBgColor: parentLocation.bgColor
    };
  }).filter((row): row is SublocationRowData => row !== null);
}