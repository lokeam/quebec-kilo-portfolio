import { useState, useEffect } from 'react';
import { onlineServicesPageMockData } from '@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata'; // Import your mock data
import type{ AvailableService } from '@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata';


export interface UseAvailableServicesResult {
  availableServices: AvailableService[];
  isLoading: boolean;
  error: Error | null;
}

export function useAvailableServices(searchQuery: string): UseAvailableServicesResult {
  const [availableServices, setAvailableServices] = useState<AvailableService[]>([]);
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
        const filtered = onlineServicesPageMockData.filter(service =>
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