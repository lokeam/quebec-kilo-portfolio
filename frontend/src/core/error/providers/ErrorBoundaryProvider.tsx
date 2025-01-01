import { ErrorBoundary, type FallbackProps } from 'react-error-boundary';
import { DefaultErrorFallbackPage } from '../pages/DefaultErrorFallbackPage';
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
   * Pattern from: https://github.com/bvaughn/react-error-boundary#error-boundary
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
    <ErrorBoundary
      FallbackComponent={FallbackComponent}
      // onError={handleError}
      // Don't retry on same error to prevent infinite loops
      // https://github.com/bvaughn/react-error-boundary#reset-keys
      resetKeys={[]}
    >
      {children}
    </ErrorBoundary>
  );
};
