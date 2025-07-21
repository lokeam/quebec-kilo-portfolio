import { useCallback, useState } from 'react';

/**
 * Hook for preloading navigation-related components
 *
 * Preloads components when user hovers over navigation items
 * to reduce loading time when they actually click
 */
export const usePreloadNavigation = () => {
  const [preloadedRoutes, setPreloadedRoutes] = useState<Set<string>>(new Set());

  const preloadRoute = useCallback(async (route: string) => {
    if (preloadedRoutes.has(route)) return;

    try {
      switch (route) {
        case 'library':
          await import('@/features/dashboard/pages/LibraryPage/LibraryPage');
          break;
        case 'online-services':
          await import('@/features/dashboard/pages/OnlineServices/OnlineServicesPage');
          break;
        case 'physical-locations':
          await import('@/features/dashboard/pages/PhysicalLocations/PhysicalLocationsPage');
          break;
        case 'spend-tracking':
          await import('@/features/dashboard/pages/SpendTrackingPage/SpendTrackingPage');
          break;
        case 'settings':
          await import('@/features/dashboard/pages/SettingsPage/SettingsPage');
          break;
        default:
          return;
      }

      setPreloadedRoutes(prev => new Set([...prev, route]));
      console.log(`📦 Preloaded route: ${route}`);
    } catch (error) {
      console.error(`Failed to preload route ${route}:`, error);
    }
  }, [preloadedRoutes]);

  const preloadOnHover = useCallback((route: string) => {
    // Use a small delay to avoid preloading on accidental hovers
    const timer = setTimeout(() => {
      preloadRoute(route);
    }, 100);

    return () => clearTimeout(timer);
  }, [preloadRoute]);

  return { preloadOnHover, preloadedRoutes };
};