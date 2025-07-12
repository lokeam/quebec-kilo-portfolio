import { useAuth0 } from '@auth0/auth0-react';
import { getOnboardingDebugState } from '@/core/utils/debug/onboardingDebug';

/**
 * Hook for onboarding status operations
 *
 * This hook determines onboarding status using Auth0 user metadata.
 * No API calls are made - it uses Auth0's built-in user data.
 *
 * @example
 * ```typescript
 * // ✅ GOOD: No API calls, uses Auth0 user metadata
 * function OnboardingProtectedRoute() {
 *   const { hasCompletedOnboarding } = useOnboardingStatus();
 *   // ...
 * }
 * ```
 */
export const useOnboardingStatus = () => {
  const { user } = useAuth0(); // Get Auth0 user object
  const debugState = getOnboardingDebugState();

  // Debug: Log Auth0 user object to see what's available
  console.log('🔍 Auth0 User Object:', {
    email: user?.email,
    name: user?.name,
    firstName: user?.firstName,
    lastName: user?.lastName,
    user_metadata: user?.user_metadata,
    app_metadata: user?.app_metadata,
  });

  // Check if user has completed onboarding based on Auth0 user metadata
  // Auth0 user object contains user_metadata with custom fields
  const hasCompletedOnboarding = debugState.forceCompletedOnboarding ||
    debugState.forceIncompleteOnboarding ? false :
    !!(
      user?.user_metadata?.firstName &&
      user?.user_metadata?.lastName
    );

  // Check if user has completed the name step specifically
  const hasCompletedNameStep = debugState.forceCompletedOnboarding ||
    debugState.forceIncompleteOnboarding ? false :
    !!(
      user?.user_metadata?.firstName &&
      user?.user_metadata?.lastName
    );

  console.log('🔍 Onboarding Status:', {
    hasCompletedOnboarding,
    hasCompletedNameStep,
    firstName: user?.user_metadata?.firstName,
    lastName: user?.user_metadata?.lastName,
  });

  return {
    hasCompletedOnboarding,
    hasCompletedNameStep,
    isLoading: false, // No loading since no API calls
    profile: null, // No profile data since no API calls
  };
};