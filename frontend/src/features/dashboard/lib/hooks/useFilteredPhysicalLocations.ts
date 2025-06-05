import { useMemo } from 'react';
import { useOnlineServicesSearch } from '@/features/dashboard/lib/stores/onlineServicesStore';
import type { SublocationRowData } from '@/core/api/adapters/analytics.adapter';

export function useFilteredPhysicalLocations(locations: SublocationRowData[]) {
  const searchQuery = useOnlineServicesSearch();

  return useMemo(() => {
    if (!searchQuery) {
      return locations;
    }

    const query = searchQuery.toLowerCase();
    return locations.filter((location) => {
      // Search in both sublocation name and parent location name
      return location.sublocationName.toLowerCase().includes(query) ||
             location.parentLocationName.toLowerCase().includes(query);
    });
  }, [locations, searchQuery]);
}