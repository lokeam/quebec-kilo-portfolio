// Sentry error boundary
import * as Sentry from '@sentry/react';

//import { ErrorBoundary, type FallbackProps } from 'react-error-boundary';
import { DefaultErrorFallbackPage } from '../pages/DefaultErrorFallbackPage';

// Note: Removed Tanstack Query error boundary in favor of Sentry
// import { QueryErrorResetBoundary } from '@tanstack/react-query';


interface ErrorBoundaryProviderProps {
  children: React.ReactNode;
  FallbackComponent?: React.ComponentType<{
    error: Error;
    resetErrorBoundary: () => void;
  }>;
}

/**
 * Global error boundary provider that catches runtime errors and displays fallback UI
 * Based on React Error Boundary pattern: https://react.dev/reference/react/Component#catching-rendering-errors-with-an-error-boundary
 */
export const ErrorBoundaryProvider = ({
  children,
  FallbackComponent = DefaultErrorFallbackPage
}: ErrorBoundaryProviderProps) => {
  return (
    <Sentry.ErrorBoundary
      fallback={({ error, resetError }) => (
        <FallbackComponent
          error={error as Error}
          resetErrorBoundary={resetError}
        />
      )}
      beforeCapture={(scope) => {
        scope.setLevel("error");
      }}
    >
      {children}
    </Sentry.ErrorBoundary>
  );
};
