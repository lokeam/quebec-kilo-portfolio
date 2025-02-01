import { useAuth0 } from '@auth0/auth0-react';
import { useCallback } from 'react';
import type { User } from '@auth0/auth0-react';

/**
 * Interface for the return value of useAuth hook
 */
interface UseAuthReturn {
  /** Whether the user is currently authenticated */
  isAuthenticated: boolean;

  /** Whether the authentication state is still being loaded */
  isLoading: boolean;

  /** The authenticated user's information, undefined if not authenticated */
  user: User | undefined;

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
 * @returns {UseAuthReturn} Object containing authentication state and methods
 *
 * @example
 * ```typescript
 * function ProtectedComponent() {
 *   const { isAuthenticated, isLoading, user, logout } = useAuth();
 *
 *   if (isLoading) return <Loading />;
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
    isLoading,
    user,
    loginWithRedirect,
    logout: auth0Logout,
    getAccessTokenSilently,
  } = useAuth0();

  /**
   * Initiates the login process by redirecting to Auth0's login page
   * User will be redirected back to the application after login
   */
  const login = useCallback(async () => {
    await loginWithRedirect();
  }, [loginWithRedirect]);

  /**
   * Logs user out + redirects them to home page
   * Clears all authentication state and local storage
   */
  const logout = useCallback(async () => {
    await auth0Logout({
      logoutParams: {
        returnTo: window.location.origin,
      }
    });
  }, [auth0Logout]);

  /**
   * Gets the current access token silently (without prompting user)
   * This token is used for authenticated API requests
   *
   * @returns {Promise<string>} The access token
   */
  const getAccessToken = useCallback(async () => {
    return await getAccessTokenSilently();
  }, [getAccessTokenSilently]);

  return {
    isAuthenticated,
    isLoading,
    user,
    login,
    logout,
    getAccessToken,
  };
};
