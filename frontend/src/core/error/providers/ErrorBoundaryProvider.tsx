import { ErrorBoundary, type FallbackProps } from 'react-error-boundary';
import { DefaultErrorFallbackPage } from '../pages/DefaultErrorFallbackPage';
import { QueryErrorResetBoundary } from '@tanstack/react-query';
// import type { ErrorInfo } from 'react';
// import { sendToSentry } from '@/services/monitoring'; // Assuming you have monitoring

interface ErrorBoundaryProviderProps {
  children: React.ReactNode;
  FallbackComponent?: React.ComponentType<FallbackProps>;
}

/**
 * Global error boundary provider that catches runtime errors and displays fallback UI
 * Based on React Error Boundary pattern: https://react.dev/reference/react/Component#catching-rendering-errors-with-an-error-boundary
 */
export const ErrorBoundaryProvider = ({
  children,
  FallbackComponent = DefaultErrorFallbackPage
}: ErrorBoundaryProviderProps) => {
  /**
   * Handles logging of caught errors
   * @param error - The error that was thrown
   * @param errorInfo - React's error info object containing component stack
   * Patterns from:
   * https://github.com/bvaughn/react-error-boundary#error-boundary
   * https://github.com/bvaughn/react-error-boundary#reset-keys
   */
  //const handleError = (error: Error, errorInfo: ErrorInfo) => {
    // TODO: Wire into monitoring service
    // sendToSentry(error, {
    //   extra: {
    //     componentStack: info.componentStack,
    //   },
    // });
  //};

  return (
    <QueryErrorResetBoundary>
      {({ reset }) => (
        <ErrorBoundary
          FallbackComponent={FallbackComponent}
          onReset={reset}
          resetKeys={[]}
        >
          {children}
        </ErrorBoundary>
      )}
    </QueryErrorResetBoundary>
  );
};
