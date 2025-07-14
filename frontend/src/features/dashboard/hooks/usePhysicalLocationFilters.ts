import { useMemo } from 'react';
import type { LocationsBFFResponse } from '@/types/domain/physical-location';
import { PhysicalLocationType } from '@/types/domain/location-types';
import { SublocationType } from '@/types/domain/location-types';

interface FilterOption {
  key: string;
  label: string;
}

interface FilterOptions {
  sublocationTypes: FilterOption[];
  parentTypes: FilterOption[];
}

export function usePhysicalLocationFilters(data: LocationsBFFResponse | undefined): FilterOptions {
  return useMemo(() => {
    if (!data?.sublocations) {
      return { sublocationTypes: [], parentTypes: [] };
    }

    // Get unique sublocation types and format them for display
    const uniqueSublocationTypes = Array.from(new Set(
      data.sublocations.map(sublocation => sublocation.sublocationType)
    ))
    .filter((type): type is SublocationType =>
      Object.values(SublocationType).includes(type as SublocationType)
    )
    .map(type => ({
      key: type,
      label: type.charAt(0).toUpperCase() + type.slice(1) // Capitalize first letter
    }));

    // Get unique physical location types and format them for display
    const uniqueParentTypes = Array.from(new Set(
      data.sublocations.map(sublocation => sublocation.parentLocationType)
    ))
    .filter((type): type is PhysicalLocationType =>
      Object.values(PhysicalLocationType).includes(type as PhysicalLocationType)
    )
    .map(type => ({
      key: type,
      label: type.charAt(0).toUpperCase() + type.slice(1) // Capitalize first letter
    }));

    return {
      sublocationTypes: uniqueSublocationTypes,
      parentTypes: uniqueParentTypes
    };
  }, [data?.sublocations]);
}