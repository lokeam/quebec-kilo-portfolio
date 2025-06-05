import { useMemo } from 'react';
import { useOnlineServicesSearch } from '@/features/dashboard/lib/stores/onlineServicesStore';
import type { PhysicalLocation } from '@/types/domain/physical-location';

export function useFilteredPhysicalLocations(locations: PhysicalLocation[]) {
  const searchQuery = useOnlineServicesSearch();

  return useMemo(() => {
    return locations.filter((location) => {
      // Search filter - match against BOTH location name and sublocation names
      const matchesSearch = !searchQuery ||
        location.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        location.sublocations?.some(sublocation =>
          sublocation.name.toLowerCase().includes(searchQuery.toLowerCase())
        );

      return matchesSearch;
    });
  }, [locations, searchQuery]);
}