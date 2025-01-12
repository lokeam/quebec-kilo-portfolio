import { useMemo } from 'react';
import type { OnlineService } from '@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata';
import {
  useOnlineServicesSearch,
  useOnlineServicesBillingFilters,
  useOnlineServicesPaymentFilters,
} from '@/features/dashboard/lib/stores/onlineServicesStore';

export function useFilteredServices(services: OnlineService[]) {
  const searchQuery = useOnlineServicesSearch();
  const billingCycleFilters = useOnlineServicesBillingFilters();
  const paymentMethodFilters = useOnlineServicesPaymentFilters();

  return useMemo(() => {

    return services.filter((service) => {
      // Search filter
      const matchesSearch = !searchQuery || service.label.toLowerCase().includes(searchQuery.toLowerCase());

      // Billing cycle filter
      const matchesBillingCycle = billingCycleFilters.length === 0 ||
        billingCycleFilters.includes(service.billingCycle);

      // Payment method filter
      const matchesPaymentMethod = paymentMethodFilters.length === 0 ||
        (service.paymentMethod && paymentMethodFilters.includes(service.paymentMethod))

      console.log('Service filtering:', {
        service: service.label,
        billingCycle: service.billingCycle,
        paymentMethod: service.paymentMethod,
        matchesBillingCycle,
        matchesPaymentMethod
      });

      return matchesSearch && matchesBillingCycle && matchesPaymentMethod;
    });
  }, [services, searchQuery, billingCycleFilters, paymentMethodFilters]);
}