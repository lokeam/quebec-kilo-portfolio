import { useQuery } from '@tanstack/react-query';
import { digitalServicesService } from '../services/digitalServices.service';
import { hookDebug } from '@/core/utils/debug/debug-utils';

// Create a query key constant for consistency
export const digitalServicesCatalogQueryKey = ['digitalServices', 'catalog'];

/**
 * Hook for fetching the digital services catalog with authentication
 */
export function useDigitalServicesCatalog() {
  // Use TanStack Query to fetch and cache the data
  return useQuery({
    queryKey: digitalServicesCatalogQueryKey,
    queryFn: async () => {
      hookDebug.logHookCall('useDigitalServicesCatalog', {});
      const data = await digitalServicesService.getServicesCatalog();
      hookDebug.logHookResult('useDigitalServicesCatalog', { dataLength: data.length });
      return data;
    },
    staleTime: 24 * 60 * 60 * 1000, // 24 hours
    gcTime: 24 * 60 * 60 * 1000, // 24 hours
    retry: 1,
  });
}

/**
 * Hook to ensure digital services catalog is prefetched during app initialization
 */
export function usePrefetchDigitalServicesCatalog() {
  // This hook is now a no-op since the main hook handles fetching
  // We keep it for backward compatibility
}