import { useMemo } from 'react';
import type { OnlineService } from '@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata';
import { useOnlineServicesSearch } from '../stores/onlineServicesStore';

export function useFilteredServices(services: OnlineService[]) {
  const searchQuery = useOnlineServicesSearch();

  return useMemo(() => {
    if (!searchQuery) return services;

    return services.filter((service) =>
      service.label.toLowerCase().includes(searchQuery)
    );
  }, [services, searchQuery]);
}