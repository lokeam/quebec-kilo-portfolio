// Sentry initialization - Recommended to be imported first
import '@/instrument';

import { StrictMode, lazy, Suspense } from 'react';
import { createRoot } from 'react-dom/client';
import { Auth0Provider } from '@auth0/auth0-react';
import { AuthInitializer } from '@/core/api/components/AuthInitializer';
import { AuthProvider } from '@/core/auth/context-provider/AuthProvider';

import { QueryClientProvider } from '@tanstack/react-query';
import { logger } from '@/core/utils/logger/logger';
import { createSentryQueryClient } from '@/core/api/sentryQueryClient';

import App from './App.tsx'
import './index.css'

// Sentry init code
const container = document.getElementById("root");
const root = createRoot(container!);

// Debug
logger.configure({
  enabled: process.env.NODE_ENV === 'development',
  level: process.env.NODE_ENV === 'development' ? 'debug' : 'error'
});

const queryClient = createSentryQueryClient();

// Lazy load React Query DevTools only in development
const ReactQueryDevtools = process.env.NODE_ENV === 'development'
  ? lazy(() => import('@tanstack/react-query-devtools').then(module => ({ default: module.ReactQueryDevtools })))
  : null;

// Note: We removed skipRedirectCallback logic to prevent white flash during Auth0 redirects
root.render(
  <StrictMode>
    <Auth0Provider
      domain={import.meta.env.VITE_AUTH0_DOMAIN}
      clientId={import.meta.env.VITE_AUTH0_CLIENT_ID}
      authorizationParams={{
        redirect_uri: window.location.origin,
        audience: import.meta.env.VITE_AUTH0_AUDIENCE,
        scope: 'openid profile email read:user_metadata update:user_metadata app_metadata offline_access'
      }}
      useRefreshTokens={true}
      cacheLocation="localstorage"
      sessionCheckExpiryDays={30}
      skipRedirectCallback={false} // Always process redirects to prevent white flash
    >

      <AuthProvider>
        <AuthInitializer>
          <QueryClientProvider client={queryClient}>
            <App />
            {process.env.NODE_ENV === 'development' && ReactQueryDevtools && (
              <Suspense fallback={null}>
                <ReactQueryDevtools />
              </Suspense>
            )}
          </QueryClientProvider>
        </AuthInitializer>
      </AuthProvider>
    </Auth0Provider>
  </StrictMode>,
);
