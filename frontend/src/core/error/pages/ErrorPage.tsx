import { ErrorPage } from '@/core/error/components/ErrorPage';

export default function ErrorPageComponent() {
  return (
    <ErrorPage
      variant="500"
      title="Something went wrong"
      subtext="An unexpected error occurred. Please try again."
      buttonText="Go Home"
      onButtonClick={() => window.location.href = '/'}
      role="alert"
      ariaLive="assertive"
    />
  );
}