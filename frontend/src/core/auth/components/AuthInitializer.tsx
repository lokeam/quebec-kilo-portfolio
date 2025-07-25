import { useAuth0 } from '@auth0/auth0-react';
import { useEffect } from 'react';
import * as Sentry from '@sentry/react';
import { saveAuth0Token } from '@/core/api/utils/auth.utils';
// import { BetaAccessErrorPage } from './BetaAccessErrorPage';

// Type for Auth0 error
interface Auth0Error extends Error {
  error?: string;
  error_description?: string;
}

/**
 * AuthInitializer is a crucial component that bridges Auth0 authentication
 * with the API layer. Ensures that every API request has access to
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
 *   <AuthInitializer>     <--- We are here
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
  const { getAccessTokenSilently, user, isAuthenticated, isLoading, error } = useAuth0();
  // const [showBetaError, setShowBetaError] = useState(false);

  useEffect(() => {
    // Save Auth0's getAccessTokenSilently fn so that we may use tokens anywhere.
    // Needed bc getAccessTokenSilently only works in React components.
    // Fn wrapper used to pass different options later if needed.
    saveAuth0Token(() => getAccessTokenSilently());
  }, [getAccessTokenSilently]);

  // Handle Auth0 errors, especially access denied errors
  useEffect(() => {
    if (error) {
      console.error('Auth0 error:', error);

      const auth0Error = error as Auth0Error;

      // Handle specific error types
      if (auth0Error.error === 'unauthorized' || auth0Error.error_description?.includes('Access denied')) {
        // This is likely a beta access denial
        //setShowBetaError(true);

        Sentry.captureException(error, {
          tags: {
            error_type: 'beta_access_denied',
            component: 'AuthInitializer'
          },
          extra: {
            error_code: auth0Error.error,
            error_description: auth0Error.error_description,
            user_email: user?.email
          }
        });

        console.warn('User denied beta access:', {
          email: user?.email,
          error: auth0Error.error_description
        });
      } else {
        // Other Auth0 errors
        Sentry.captureException(error, {
          tags: {
            error_type: 'auth0_general',
            component: 'AuthInitializer'
          },
          extra: {
            error_code: auth0Error.error,
            error_description: auth0Error.error_description
          }
        });
      }
    } else {
      // Clear error state when no error
      // setShowBetaError(false);
    }
  }, [error, user]);

  // Ensure theme consistency during Auth0 callback processing
  useEffect(() => {
    // Check if we're in an Auth0 callback state
    const isAuth0Callback = window.location.search.includes('code=') && window.location.search.includes('state=');

    if (isAuth0Callback) {
      // Apply theme immediately during Auth0 callback to prevent white flash
      try {
        const storedTheme = localStorage.getItem('qko-theme-storage');
        let themeMode = 'light';

        if (storedTheme) {
          try {
            const parsed = JSON.parse(storedTheme);
            themeMode = parsed.mode || 'light';
          } catch {
            console.warn('Failed to parse stored theme during Auth0 callback');
          }
        }

        const html = document.documentElement;
        const body = document.body;
        const root = document.getElementById('root');

        if (themeMode === 'dark') {
          html.classList.add('dark');
          html.style.backgroundColor = 'hsl(222.2 47.4% 11.2%)';
          html.style.color = 'hsl(210 40% 98%)';
          if (body) {
            body.style.backgroundColor = 'hsl(222.2 47.4% 11.2%)';
            body.style.color = 'hsl(210 40% 98%)';
          }
          if (root) {
            root.style.backgroundColor = 'hsl(222.2 47.4% 11.2%)';
            root.style.color = 'hsl(210 40% 98%)';
          }
        } else {
          html.classList.remove('dark');
          html.style.backgroundColor = 'hsl(0 0% 100%)';
          html.style.color = 'hsl(222.2 47.4% 11.2%)';
          if (body) {
            body.style.backgroundColor = 'hsl(0 0% 100%)';
            body.style.color = 'hsl(222.2 47.4% 11.2%)';
          }
          if (root) {
            root.style.backgroundColor = 'hsl(0 0% 100%)';
            root.style.color = 'hsl(222.2 47.4% 11.2%)';
          }
        }
      } catch (error) {
        console.warn('Failed to apply theme during Auth0 callback:', error);
      }
    }
  }, [isLoading]);

  // Set Sentry user context when authentication state changes
  useEffect(() => {
    if (isAuthenticated && user) {
      // Set user context in Sentry for error tracking
      Sentry.setUser({
        id: user.sub,
        email: user.email,
        username: user.name || user.email,
      });

      // Add user properties for better errors
      Sentry.setContext('user', {
        email: user.email,
        name: user.name,
        email_verified: user.email_verified,
        sub: user.sub,
        beta_access: user.app_metadata?.betaAccess || false,
        beta_access_granted_at: user.app_metadata?.betaAccessGrantedAt
      });
    } else {
      // Clear user context when not authenticated
      Sentry.setUser(null);
    }
  }, [isAuthenticated, user]);

  // Show beta access error page if access was denied
  // if (showBetaError) {
  //   return <BetaAccessErrorPage />;
  // }

  return <>{children}</>;
}