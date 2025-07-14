import { type FallbackProps } from 'react-error-boundary';
import { ErrorPage } from '@/core/error/components/ErrorPage';

export function FormErrorFallback({ error, resetErrorBoundary }: FallbackProps) {
  return (
    <ErrorPage
      variant="500"
      title="Something went wrong"
      subtext={error?.message || "An unexpected error occurred."}
      buttonText="Try Again"
      onButtonClick={resetErrorBoundary}
      role="alert"
      ariaLive="assertive"
    />
  );
}