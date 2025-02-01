import { useAuth0 } from '@auth0/auth0-react';
import { useCallback } from 'react';

/**
 * Hook to manage Auth0 token retrieval and caching.
 * Uses Auth0's SDK to handle token lifecycle.
 *
 * @returns Function to get the current access token
 * @throws Error if token retrieval fails or user is not authenticated
 *
 * @see https://auth0.com/docs/libraries/auth0-react
 */
export const useAuth0Token = () => {
  const { getAccessTokenSilently, isAuthenticated } = useAuth0();

  return useCallback(async () => {
    if (!isAuthenticated) {
      throw new Error('User is not authenticated');
    }

    try {
      return await getAccessTokenSilently();
    } catch (error) {
      console.error('Failed to get Auth0 token:', error);
      throw error;
    }
  }, [getAccessTokenSilently, isAuthenticated]);
};