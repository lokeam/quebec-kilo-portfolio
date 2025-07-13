import { useAuth0 } from '@auth0/auth0-react';
import { useCallback, useEffect } from 'react';
import type { User } from '@auth0/auth0-react';
import { getOnboardingDebugState, logDebugInfo } from '@/core/utils/debug/onboardingDebug';

const detectNewUserSignup = (user: User): boolean => {
  if (!user) return false;

  // If created_at field in response is missing, its an existing user
  if (!user.created_at) {
    console.log('🔍 User signup detection: Existing user (no created_at field)');
    return false;
  }

  // Check if created_at and updated_at are the same (within 1 minute)
  const createdAt = new Date(user.created_at);
  const updatedAt = new Date(user.updated_at || '');

  // Check if dates are valid
  if (isNaN(createdAt.getTime()) || isNaN(updatedAt.getTime())) {
    console.log('🔍 User signup detection: Invalid date fields');
    return false;
  }

  const timeDiff = Math.abs(updatedAt.getTime() - createdAt.getTime());

  // If they're within 1 minute of each other, it's a new user signup
  const isNewUserSignup = timeDiff < 60000;

  console.log('🔍 User signup detection:', {
    createdAt,
    updatedAt,
    timeDiff,
    isNewUserSignup,
  });
  return isNewUserSignup;
}

/**
 * Interface for the return value of useAuth hook
 */
export interface UseAuthReturn {
  /** Whether the user is currently authenticated */
  isAuthenticated: boolean;

  /** Whether the authentication state is still being loaded */
  isLoading: boolean;

  /** The authenticated user's information, undefined if not authenticated */
  user: User | undefined;

  /** Whether the user is a new user signup */
  isNewUserSignup: boolean;

  /** Function to initiate the login process */
  login: () => Promise<void>;

  /** Function to log the user out */
  logout: () => Promise<void>;

  /** Function to get the current access token */
  getAccessToken: () => Promise<string>;
}

/**
 * Custom hook that wraps Auth0's authentication functionality
 *
 * This hook provides a simplified interface for handling authentication
 * in the application. It wraps Auth0's useAuth0 hook and provides
 * only the necessary functionality for our use case.
 *
 * IMPORTANT: This hook ONLY handles authentication state.
 * For onboarding status, use useOnboardingStatus() hook instead.
 *
 * IMPORTANT: DO NOT USE THIS HOOK DIRECTLY.
 * Instead, use useAuthContext() hook from AuthContext.tsx
 *
 * @returns {UseAuthReturn} Object containing authentication state and methods
 *
 * @example
 * ```typescript
 * function ProtectedComponent() {
 *   const { isAuthenticated, isLoading, user, logout } = useAuthContext();
 *
 *   if (isLoading) return <LoadingPage />;
 *   if (!isAuthenticated) return <LoginPrompt />;
 *
 *   return (
 *     <div>
 *       <p>Welcome, {user?.name}!</p>
 *       <button onClick={logout}>Logout</button>
 *     </div>
 *   );
 * }
 * ```
 */
export const useAuth = (): UseAuthReturn => {
  const {
    isAuthenticated,
    isLoading: authLoading,
    user,
    loginWithRedirect,
    logout: auth0Logout,
    getAccessTokenSilently,
  } = useAuth0();

  console.log('-- useAuth rendered --');

  // Get debug state
  const debugState = getOnboardingDebugState();

  // Check if this is a new user (with debug override)
  const isNewUserSignup = debugState.forceNewUser || (user ? detectNewUserSignup(user) : false);

  // Debug authentication state changes
  logDebugInfo('Auth State', {
    isAuthenticated,
    isLoading: authLoading,
    hasUser: !!user,
    userEmail: user?.email,
    isNewUserSignup,
    debugMode: debugState.bypassOnboarding || debugState.forceNewUser || debugState.forceIncompleteOnboarding || debugState.forceCompletedOnboarding,
  });

  // Listen for authentication state changes to backup theme data on automatic logout
  useEffect(() => {
    // If user was authenticated but is now not authenticated, backup theme data
    if (!isAuthenticated && !authLoading) {
      // Backup theme data before Auth0 clears localStorage
      const themeData = localStorage.getItem('qko-theme-storage');

      if (themeData) {
        // sessionStorage.setItem('qko-theme-backup', themeData);
        // console.log('🎨 Theme data backed up due to automatic logout');
      }
    }
  }, [isAuthenticated, authLoading]);

  /**
   * Initiates the login process by redirecting to Auth0's login page
   * User will be redirected back to the application after login
   */
  const login = useCallback(async () => {
    console.log('🔑 Initiating login...');
    await loginWithRedirect();
  }, [loginWithRedirect]);

  /**
   * Logs user out + redirects them to home page
   * Preserves theme preferences before Auth0 clears localStorage
   */
  const logout = useCallback(async () => {
    console.log('🚪 Initiating logout...');

    await auth0Logout({
      logoutParams: {
        returnTo: `${window.location.origin}/login`,
      }
    });
  }, [auth0Logout]);

  /**
   * Gets the current access token without prompting user
   * This token is used for authenticated API requests
   *
   * @returns {Promise<string>} The access token
   */
  const getAccessToken = useCallback(async () => {
    try {
      const token = await getAccessTokenSilently();
      console.log('✅ Token retrieved successfully');
      return token;
    } catch (error) {
      console.error('❌ Failed to get access token:', error);
      throw error;
    }
  }, [getAccessTokenSilently]);

  return {
    isAuthenticated,
    isLoading: authLoading,
    user,
    isNewUserSignup,
    login,
    logout,
    getAccessToken,
  };
};
