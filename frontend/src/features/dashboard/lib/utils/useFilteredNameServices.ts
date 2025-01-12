import { useMemo } from 'react';
import type { OnlineService } from '@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata';

export const useFilteredNameServices = (services: OnlineService[], filter: string) => {
  return useMemo(() => {
    const normalizedFilter = filter.toLowerCase();
    return services.filter((service) =>
      service.label.toLowerCase().includes(normalizedFilter) ||
      service.name.toLowerCase().includes(normalizedFilter)
    );
  }, [services, filter]);
};