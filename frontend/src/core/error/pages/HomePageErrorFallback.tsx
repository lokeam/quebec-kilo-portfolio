import { type FallbackProps } from 'react-error-boundary';
import { ErrorPage } from '@/core/error/components/ErrorPage';

export function HomePageErrorFallback({ error, resetErrorBoundary }: FallbackProps) {
  console.error('HomePageErrorFallback - error:', error);
  return (
    <ErrorPage
      variant="500"
      title="Dashboard Error"
      subtext="There was an error loading the dashboard. Please try again."
      buttonText="Retry Dashboard"
      onButtonClick={resetErrorBoundary}
      role="alert"
      ariaLive="assertive"
    />
  );
}
