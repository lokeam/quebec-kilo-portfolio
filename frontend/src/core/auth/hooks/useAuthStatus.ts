import { useAuthContext } from '@/core/auth/context-provider/AuthContext';
import { useOnboardingStatus } from '@/core/auth/hooks/useOnboardingStatus';
import { useMemo } from 'react';

/**
 * Consolidated auth status hook that combines authentication and onboarding status
 * into a single, predictable loading state.
 *
 * This eliminates race conditions between Auth0 loading and onboarding status loading,
 * providing a better user experience with a single loading state.
 *
 * @returns {Object} Consolidated auth status
 * @returns {boolean} isAuthenticated - Whether user is authenticated
 * @returns {boolean} isLoading - Single loading state for all auth operations
 * @returns {boolean} hasCompletedOnboarding - Whether user completed onboarding
 * @returns {boolean} needsOnboarding - Whether user needs to complete onboarding
 * @returns {boolean} canAccessProtectedRoutes - Whether user can access main app
 */
export const useAuthStatus = () => {
  const { isAuthenticated, isLoading: authLoading } = useAuthContext();
  const { hasCompletedOnboarding, isLoading: onboardingLoading } = useOnboardingStatus();

  const status = useMemo(() => {
    // If Auth0 is still loading, we're loading
    if (authLoading) {
      return {
        isAuthenticated: false,
        isLoading: true,
        hasCompletedOnboarding: false,
        needsOnboarding: false,
        canAccessProtectedRoutes: false,
      };
    }

    // If not authenticated, we're not loading and user needs to login
    if (!isAuthenticated) {
      return {
        isAuthenticated: false,
        isLoading: false,
        hasCompletedOnboarding: false,
        needsOnboarding: false,
        canAccessProtectedRoutes: false,
      };
    }

    // If onboarding status is still loading, we're loading
    if (onboardingLoading) {
      return {
        isAuthenticated: true,
        isLoading: true,
        hasCompletedOnboarding: false,
        needsOnboarding: false,
        canAccessProtectedRoutes: false,
      };
    }

    // User is authenticated, check onboarding status
    const needsOnboarding = !hasCompletedOnboarding;
    const canAccessProtectedRoutes = hasCompletedOnboarding;

    return {
      isAuthenticated: true,
      isLoading: false,
      hasCompletedOnboarding,
      needsOnboarding,
      canAccessProtectedRoutes,
    };
  }, [isAuthenticated, authLoading, hasCompletedOnboarding, onboardingLoading]);

  return status;
};