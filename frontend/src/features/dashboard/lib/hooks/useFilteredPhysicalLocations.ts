import { useMemo } from 'react';
import { useOnlineServicesSearch, useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';
import type { SublocationItemData } from '@/core/api/adapters/analytics.adapter';

export function useFilteredPhysicalLocations(locations: SublocationItemData[]) {
  const searchQuery = useOnlineServicesSearch();
  const { sublocationTypeFilters, parentLocationTypeFilters } = useOnlineServicesStore();

  return useMemo(() => {
    if (!searchQuery && sublocationTypeFilters.length === 0 && parentLocationTypeFilters.length === 0) {
      return locations;
    }

    const query = searchQuery.toLowerCase();
    return locations.filter((location) => {
      const matchesSearch = !searchQuery ||
        location.sublocationName.toLowerCase().includes(query) ||
        location.parentLocationName.toLowerCase().includes(query);

      const matchesSublocationType = sublocationTypeFilters.length === 0 ||
        sublocationTypeFilters.includes(location.sublocationType);

      const matchesParentType = parentLocationTypeFilters.length === 0 ||
        parentLocationTypeFilters.includes(location.parentLocationType);

      return matchesSearch && matchesSublocationType && matchesParentType;
    });
  }, [locations, searchQuery, sublocationTypeFilters, parentLocationTypeFilters]);
}