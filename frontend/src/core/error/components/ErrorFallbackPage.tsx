import type { FallbackProps } from 'react-error-boundary';
import { ErrorButton } from '@/core/error/components/ErrorButton';
import type { ErrorConfig } from '@/core/error/types/error.types';
import { cn } from '@/shared/components/ui/utils';

interface ErrorFallbackPageProps extends FallbackProps {
  config?: ErrorConfig;
};

export const ErrorFallbackPage = ({
  error,
  resetErrorBoundary, // Provided by react-error-boundary
  config
}: ErrorFallbackPageProps) => {

  const handleClick = () => {
    if (config?.onAction) {
      // If custom action exists, only call that
      config.onAction();
    } else {
      // Otherwise use the default reset behavior
      resetErrorBoundary();
    }
  };

  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' || event.key === ' ') {
       // Prevent page scroll on space
      event.preventDefault();
      if (config?.onAction) {
        config.onAction();
      } else {
        resetErrorBoundary();
      }
    }
  };
  // Log error for monitoring
  console.error('[ErrorBoundary] Error:', error);

  return (
    <div className="container mx-auto max-w-2xl px-4">
      <div className={cn(
        "mt-32 p-8 rounded-lg shadow-lg bg-background",
        "flex flex-col items-center"
      )}>
        <div className="text-center">
          <h2 className={cn(
            "text-3xl font-bold mb-4",
            "text-destructive"
          )}>
            {config?.severity === 'fatal' ? 'Critical Error' : 'Oops!'}
          </h2>

          <h3 className="text-xl font-semibold mb-4">
            {config?.message || 'Something went wrong'}
          </h3>

          <p className="text-muted-foreground mb-6">
            {error?.message || 'We apologize for the inconvenience. Please try again.'}
          </p>

          <ErrorButton
            onClick={handleClick}
            onKeyDown={handleKeyDown}
            label={config?.actionLabel || 'Try Again'}
          />
        </div>
      </div>
    </div>
  );
};