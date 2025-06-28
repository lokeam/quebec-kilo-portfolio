import { useMemo } from 'react';

import {
  useOnlineServicesSearch,
  useOnlineServicesBillingFilters,
  useOnlineServicesPaymentFilters,
} from '@/features/dashboard/lib/stores/onlineServicesStore';

//import type { DigitalLocation } from '@/types/domain/digital-location';
import type { DigitalLocationBFFResponseItem } from '@/types/domain/digital-location';


export function useFilteredServices(services: DigitalLocationBFFResponseItem[]) {
  const searchQuery = useOnlineServicesSearch();
  const billingCycleFilters = useOnlineServicesBillingFilters();
  const paymentMethodFilters = useOnlineServicesPaymentFilters();

  console.log('[DEBUG] useFilteredServices: Services:', services);

  return useMemo(() => {
    return services.filter((service) => {
      // Search filter - match against both label and name fields
      const matchesSearch = !searchQuery ||
        service.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        (service.name && service.name.toLowerCase().includes(searchQuery.toLowerCase()));

      // Billing cycle filter
      const matchesBillingCycle =
        billingCycleFilters.length === 0 || // if no filters, show all
        (service.billingCycle != null &&
         billingCycleFilters.some(filter => filter === service.billingCycle));

      // Payment method filter - case insensitive comparison
      const matchesPaymentMethod =
        paymentMethodFilters.length === 0 || // if no filters, show all
        (service.paymentMethod != null &&
         paymentMethodFilters.some(filter =>
           filter.toLowerCase() === service.paymentMethod.toLowerCase()
         ));
      return matchesSearch && matchesBillingCycle && matchesPaymentMethod;
    });
  }, [services, searchQuery, billingCycleFilters, paymentMethodFilters]);
}