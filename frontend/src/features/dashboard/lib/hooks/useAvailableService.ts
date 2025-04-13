import { useDigitalServicesCatalog } from '@/core/api/queries/useDigitalServicesCatalog';
import type { OnlineService } from '../types/online-services/services';
import { useMemo } from 'react';

export interface UseAvailableServicesResult {
  availableServices: OnlineService[];
  isLoading: boolean;
  error: Error | null;
}

export function useAvailableServices(searchQuery: string): UseAvailableServicesResult {
  // Get the full list from cache (single API call)
  const { data = [], isLoading, error } = useDigitalServicesCatalog();

  // Filter the data locally using useMemo for performance
  const availableServices = useMemo(() => {
    const services = data.map(service => ({
      id: service.id,
      name: service.name,
      label: service.name,
      logo: service.logo,
      url: '#',
      type: service.isSubscriptionService ? 'subscription' : 'online',
      isSubscriptionService: service.isSubscriptionService,
      status: 'active',
      features: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      tier: {
        currentTier: 'free',
        availableTiers: [{
          id: '1',
          name: 'free',
          features: [],
          isDefault: true
        }]
      },
      billing: {
        cycle: 'NA',
        fees: { monthly: '0' },
        paymentMethod: 'Generic'
      }
    })) as OnlineService[];

    // Local filtering
    if (!searchQuery.trim()) return services;

    const lowercaseQuery = searchQuery.toLowerCase();
    return services.filter(service =>
      service.name.toLowerCase().includes(lowercaseQuery) ||
      service.label.toLowerCase().includes(lowercaseQuery)
    );
  }, [data, searchQuery]);

  return {
    availableServices,
    isLoading,
    error: error as Error | null
  };
}