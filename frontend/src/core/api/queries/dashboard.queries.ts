// Tanstack query
import { useQueryClient } from '@tanstack/react-query';

// Base query hook
import { useAPIQuery } from '@/core/api/queries/useAPIQuery';

// Service Layer methods
import { getDashboardBFFResponse } from '@/core/api/services/dashboard.service';

// Types
import type { DashboardResponse } from '@/core/api/services/dashboard.service';

/**
 * Query keys for dashboard-related queries
 *
 * Used for caching and invalidation of dashboard data
 */
export const dashboardKeys = {
  all: ['dashboard'] as const,
  bff: () => [...dashboardKeys.all, 'bff'] as const,
};

/**
 * Hook to fetch all dashboard data for the BFF page
 *
 * This hook fetches all dashboard data in a single request, including:
 * - Basic statistics (games, subscriptions, locations)
 * - Digital and physical storage locations
 * - Platform distribution
 * - Monthly spending data
 *
 * The data is considered stale immediately (staleTime: 0) to ensure
 * it's always up to date with the latest changes in the system.
 *
 * Usage:
 * ```tsx
 * const { data, isLoading, error } = useGetDashboardBFFResponse();
 *
 * if (isLoading) return <LoadingSpinner />;
 * if (error) return <ErrorMessage error={error} />;
 *
 * return <Dashboard data={data} />;
 * ```
 */
export const useGetDashboardBFFResponse = () => {
  return useAPIQuery<DashboardResponse>({
    queryKey: dashboardKeys.bff(),
    queryFn: async () => {
      try {
        const response = await getDashboardBFFResponse();
        return response;
      } catch(error) {
        console.error('[DEBUG] useGetDashboardBFFResponse: Error fetching data:', error);
        throw error;
      }
    },
    staleTime: 0,
    refetchOnMount: true,
    refetchOnWindowFocus: true,
  });
};

/**
 * Utility function to invalidate dashboard queries
 *
 * This should be called after any mutation that affects dashboard data, such as:
 * - Adding/removing/updating games
 * - Adding/removing/updating locations
 * - Adding/removing/updating subscriptions
 * - Any changes to spending data
 *
 * Usage:
 * ```tsx
 * const queryClient = useQueryClient();
 *
 * // After mutation
 * invalidateDashboardQueries(queryClient);
 * ```
 *
 * This ensures the dashboard data is refreshed to reflect any changes
 * in the underlying data.
 */
export const invalidateDashboardQueries = (queryClient: ReturnType<typeof useQueryClient>) => {
  queryClient.invalidateQueries({ queryKey: dashboardKeys.all });
};