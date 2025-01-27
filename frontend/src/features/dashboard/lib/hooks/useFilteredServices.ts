import { useMemo } from 'react';
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
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
      const matchesSearch = !searchQuery ||
        service.label.toLowerCase().includes(searchQuery.toLowerCase());

      // Billing cycle filter
      const matchesBillingCycle =
        billingCycleFilters.length === 0 || // if no filters, show all
        (service.billing?.cycle &&
         billingCycleFilters.some(filter => filter === service.billing.cycle));

      // Payment method filter
      const matchesPaymentMethod =
        paymentMethodFilters.length === 0 || // if no filters, show all
        (service.billing?.paymentMethod &&
         paymentMethodFilters.some(filter => filter === service.billing.paymentMethod));
      return matchesSearch && matchesBillingCycle && matchesPaymentMethod;
    });
  }, [services, searchQuery, billingCycleFilters, paymentMethodFilters]);
}