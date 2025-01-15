import { useState, useEffect } from 'react';
import { onlineServicesPageMockData } from '@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata';
import type { OnlineService } from '@/features/dashboard/lib/types/service.types';


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

    // Simulate network delay
    const timeoutId = setTimeout(() => {
      try {
        // Filter using label instead of name
        const filtered = onlineServicesPageMockData.services.filter((service: OnlineService) =>
          service.label.toLowerCase().includes(searchQuery.toLowerCase())
        );
        setAvailableServices(filtered);
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