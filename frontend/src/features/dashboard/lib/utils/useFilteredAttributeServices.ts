import { useMemo } from 'react';
import type { OnlineService } from '@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata';

type FilterAttribute = 'billingCycle' | 'paymentMethod';

export const useFilteredServicesByAttribute = (
  services: OnlineService[],
  attribute: FilterAttribute,
  filterValues: string[]
) => {
  return useMemo(() => {
    if (!filterValues.length) return services;

    return services.filter((service) => {
      const serviceValue = service[attribute];

      if (serviceValue === undefined) return false;
      return filterValues.includes(serviceValue);
    })
  }, [services, attribute, filterValues]);
};
