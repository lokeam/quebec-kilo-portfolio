// Sentry initialization - Recommended to be imported first
import '@/instrument';

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { Auth0Provider } from '@auth0/auth0-react';
import { AuthInitializer } from '@/core/api/components/AuthInitializer';
import { QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
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

root.render(
  <StrictMode>
    <Auth0Provider
      domain={import.meta.env.VITE_AUTH0_DOMAIN}
      clientId={import.meta.env.VITE_AUTH0_CLIENT_ID}
      authorizationParams={{
        redirect_uri: window.location.origin,
        audience: import.meta.env.VITE_AUTH0_AUDIENCE,
        scope: 'openid profile email'
      }}
      useRefreshTokens={true}
      cacheLocation="localstorage"
      sessionCheckExpiryDays={30}
      skipRedirectCallback={window.location.pathname === '/login'}
    >
      <AuthInitializer>
        <QueryClientProvider client={queryClient}>
          <App />
          {process.env.NODE_ENV === 'development' && <ReactQueryDevtools />}
        </QueryClientProvider>
      </AuthInitializer>
    </Auth0Provider>
  </StrictMode>,
);
