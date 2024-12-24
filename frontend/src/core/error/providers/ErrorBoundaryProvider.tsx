import { ErrorBoundary } from 'react-error-boundary';
import { ErrorFallbackPage } from '@/core/error/components/ErrorFallbackPage';
import { ErrorBoundaryProps } from '@/core/error/types/error.types';

export const ErrorBoundaryProvider = ({ children, config }: ErrorBoundaryProps) => {
  const handleError = (error: Error) => {
    // Add error reporting service integration here
    console.error('[ErrorBoundary] Caught error:', error);
  };

  return (
    <ErrorBoundary
      FallbackComponent={(props) => <ErrorFallbackPage {...props} config={config} />}
      onError={handleError}
      onReset={() => {
        // Optional: Clear error state or refresh data
        window.location.reload();
      }}
    >
      {children}
    </ErrorBoundary>
  );
};