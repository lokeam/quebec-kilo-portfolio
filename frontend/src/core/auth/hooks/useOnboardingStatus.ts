import { useAuth0 } from '@auth0/auth0-react';
import { useEffect, useState } from 'react';
import { getOnboardingDebugState } from '@/core/utils/debug/onboardingDebug';

// Extend Window interface for optimistic state
declare global {
  interface Window {
    __ONBOARDING_OPTIMISTIC_COMPLETE__?: boolean;
    __WANTS_INTRO_TOASTS__?: boolean;
  }
}

/**
 * Hook for onboarding status operations
 *
 * This hook determines onboarding status using Auth0 app metadata and custom claims.
 * It first tries to get data from custom claims in the ID token (for user name, etc.),
 * then falls back to the user.app_metadata object for basic onboarding status.
 * Custom claims are necessary for complete onboarding data.
 *
 * @example
 * ```typescript
 * // ✅ GOOD: Gets custom claims for complete onboarding data
 * function OnboardingProtectedRoute() {
 *   const { hasCompletedOnboarding } = useOnboardingStatus();
 *   // ...
 * }
 * ```
 */
export const useOnboardingStatus = () => {
  const { user, getIdTokenClaims } = useAuth0(); // Get Auth0 user object and token claims
  const debugState = getOnboardingDebugState();
  const [onboardingStatus, setOnboardingStatus] = useState<{
    hasCompletedOnboarding: boolean;
    hasCompletedNameStep: boolean;
    source: 'id_token_claims' | 'app_metadata' | 'fallback';
  } | null>(null);

  // Debug: Log Auth0 user object to see what's available
  console.log('🔍 Auth0 User Object:', {
    email: user?.email,
    name: user?.name,
    firstName: user?.firstName,
    lastName: user?.lastName,
    user_metadata: user?.user_metadata,
    app_metadata: user?.app_metadata,
  });

  useEffect(() => {
    const checkOnboardingStatus = async () => {
      if (!user) return;

      try {
        // First, try to get from ID token claims (custom claims) - NECESSARY for user name
        const claims = await getIdTokenClaims();
        const namespace = 'https://qko.app/claims';

        console.log('🔍 ID Token Claims (full object):', claims);
        console.log('🔍 Looking for claims under namespace:', namespace);

        // Check if custom claims exist
        if (claims && claims[`${namespace}/hasCompletedOnboarding`] !== undefined) {
          console.log('✅ Found onboarding status in ID token claims!');
          console.log('📋 Custom claims found:', {
            hasCompletedOnboarding: claims[`${namespace}/hasCompletedOnboarding`],
            wantsIntroToasts: claims[`${namespace}/wants_intro_toasts`],
            allClaimsWithNamespace: Object.keys(claims).filter(key => key.startsWith(namespace))
          });
          const newStatus: {
            hasCompletedOnboarding: boolean;
            hasCompletedNameStep: boolean;
            source: 'id_token_claims' | 'app_metadata' | 'fallback';
          } = {
            hasCompletedOnboarding: !!claims[`${namespace}/hasCompletedOnboarding`],
            hasCompletedNameStep: !!claims[`${namespace}/hasCompletedOnboarding`],
            source: 'id_token_claims'
          };
          setOnboardingStatus(prev => {
            if (JSON.stringify(prev) === JSON.stringify(newStatus)) return prev;
            return newStatus;
          });
          return;
        } else {
          console.log('❌ No custom claims found under namespace:', namespace);
          console.log('🔍 Available claim keys:', Object.keys(claims || {}));
        }
      } catch (error) {
        console.warn('Failed to get ID token claims:', error);
      }

      // Fallback to app_metadata from user object
      console.log('🔄 Falling back to app_metadata from user object');
      const fallbackStatus: {
        hasCompletedOnboarding: boolean;
        hasCompletedNameStep: boolean;
        source: 'id_token_claims' | 'app_metadata' | 'fallback';
      } = {
        hasCompletedOnboarding: !!user?.app_metadata?.hasCompletedOnboarding,
        hasCompletedNameStep: !!user?.app_metadata?.hasCompletedOnboarding,
        source: 'app_metadata'
      };
      setOnboardingStatus(prev => {
        if (JSON.stringify(prev) === JSON.stringify(fallbackStatus)) return prev;
        return fallbackStatus;
      });
    };

    checkOnboardingStatus();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user]);

  // Apply debug overrides and optimistic state
  const finalOnboardingStatus = debugState.forceCompletedOnboarding ? true :
  debugState.forceIncompleteOnboarding ? false :
  window.__ONBOARDING_OPTIMISTIC_COMPLETE__ || localStorage.getItem('__ONBOARDING_OPTIMISTIC_COMPLETE__') === 'true' ? true :
  onboardingStatus?.hasCompletedOnboarding ?? false;

  const finalNameStepStatus = debugState.forceCompletedOnboarding ? true :
  debugState.forceIncompleteOnboarding ? false :
  onboardingStatus?.hasCompletedNameStep ?? false;

  console.log('🔍 Onboarding Status:', {
    hasCompletedOnboarding: finalOnboardingStatus,
    hasCompletedNameStep: finalNameStepStatus,
    source: onboardingStatus?.source,
    appMetadataHasCompletedOnboarding: user?.app_metadata?.hasCompletedOnboarding,
  });

  return {
    hasCompletedOnboarding: finalOnboardingStatus,
    hasCompletedNameStep: finalNameStepStatus,
    isLoading: user ? onboardingStatus === null : false, // Loading if we have a user but haven't determined status yet
    profile: null, // No profile data since no API calls
  };
};