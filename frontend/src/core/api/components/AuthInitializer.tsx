import { useAuth0 } from '@auth0/auth0-react';
import { useEffect } from 'react';
import { setAuth0TokenFn } from '@/core/api/utils/auth.utils';

/**
 * AuthInitializer is a crucial component that bridges Auth0 authentication
 * with the API layer. It ensures that every API request has access to
 * the latest auth token.
 *
 * Why this exists:
 * 1. Our API layer needs auth tokens for requests
 * 2. Auth0 manages these tokens
 * 3. We need a way to get fresh tokens when making requests
 * 4. This component provides that connection
 *
 * Component Hierarchy:
 * ```
 * <Auth0Provider>
 *   <AuthInitializer>     <--- You are here
 *     <QueryClientProvider>
 *       <App />
 *     </QueryClientProvider>
 *   </AuthInitializer>
 * </Auth0Provider>
 * ```
 *
 * @example
 * ```typescript
 * // main.tsx
 * root.render(
 *   <Auth0Provider {...auth0Config}>
 *     <AuthInitializer>
 *       <App />
 *     </AuthInitializer>
 *   </Auth0Provider>
 * );
 *
 * // Later used in API calls:
 * const response = await axiosInstance.get('/api/data'); // Token automatically included
 * ```
 *
 * @param props.children - Child components that need access to authenticated API calls
 */
export function AuthInitializer({ children }: { children: React.ReactNode }) {
  const { getAccessTokenSilently } = useAuth0();

  useEffect(() => {
    // Register token getter function with API layer
    // This function will be called whenever we need a fresh token
    setAuth0TokenFn(() => getAccessTokenSilently());
  }, [getAccessTokenSilently]);

  return <>{children}</>;
}

/**
 * Auth Initializer Component
 *
 * For API standards and best practices, see:
 * @see {@link ../../docs/api-standards.md}
 */
