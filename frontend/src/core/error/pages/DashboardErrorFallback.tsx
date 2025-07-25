import { type FallbackProps } from 'react-error-boundary';
import { ErrorPage } from '@/core/error/components/ErrorPage';

export function DashboardErrorFallback({ error, resetErrorBoundary }: FallbackProps) {
  console.error('DashboardErrorFallback - error:', error);
  return (
    <ErrorPage
      variant="500"
      title="Dashboard Error"
      subtext="There was an error loading this page. Please try again."
      buttonText="Retry"
      onButtonClick={resetErrorBoundary}
      role="alert"
      ariaLive="assertive"
    />
  );
}