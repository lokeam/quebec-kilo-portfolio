import { ErrorBoundary } from 'react-error-boundary';
import type { ErrorBoundaryProps } from '../types/error.types';
import { ErrorFallbackPage } from '../components/ErrorFallbackPage';

export const ErrorBoundaryProvider = ({
  children,
  config
}: ErrorBoundaryProps) => {
  const handleError = (error: Error, errorInfo: React.ErrorInfo) => {
    console.error('[ErrorBoundary]:', error);
    config?.onError?.(error, errorInfo);
  };

  const handleAction = () => {
    // Custom action should trigger reset
    if (config?.onAction) {
      config.onAction();
      config.onReset?.();
    }
  };

  return (
    <ErrorBoundary
      FallbackComponent={(props) => (
        <ErrorFallbackPage
          {...props}
          config={{
            ...config,
            onAction: config?.onAction ? handleAction : undefined
          }}
        />
      )}
      onError={handleError}
      onReset={config?.onReset}
    >
      {children}
    </ErrorBoundary>
  );
};
