import { useState, useEffect } from 'react';
import { onlineServicesPageMockData } from '@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata';
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import type { ServiceTierName } from '../types/online-services/tiers';


export interface UseAvailableServicesResult {
  availableServices: OnlineService[];
  isLoading: boolean;
  error: Error | null;
}

export function useAvailableServices(searchQuery: string): UseAvailableServicesResult {
  const [availableServices, setAvailableServices] = useState<OnlineService[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!searchQuery.trim()) {
      setAvailableServices([]);
      return;
    }

    setIsLoading(true);
    setError(null);

    // Simulate network delay - replace with Tanstack Query
    const timeoutId = setTimeout(() => {
      try {
        // Filter using label instead of name
        // NOTE: Replace with proper data from Tanstack Query
        const services = onlineServicesPageMockData.services.map(service => ({
          ...service,
          tier: {
            currentTier: (service.tier.name || 'free') as ServiceTierName,
            availableTiers: [{
              id: '1',
              name: service.tier.name || 'free',
              features: service.tier.features,
              isDefault: true
            }],
          }
        })) as OnlineService[];
        const filtered = services.filter((service) =>
          service.label.toLowerCase().includes(searchQuery.toLowerCase())
        );
        setAvailableServices(filtered as OnlineService[]);
        setIsLoading(false);
      } catch (err) {
        setError(err instanceof Error ? err : new Error('An error occurred'));
        setIsLoading(false);
      }
    }, 500);

    return () => clearTimeout(timeoutId);
  }, [searchQuery]);

  return { availableServices, isLoading, error };
}