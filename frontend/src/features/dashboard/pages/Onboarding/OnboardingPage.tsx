import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { OnboardingPageContent } from '@/features/dashboard/pages/Onboarding/OnboardingPageContent';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';

function OnlineServicesPage() {
  return (
    <ErrorBoundary
      FallbackComponent={OnlineServicesPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <OnboardingPageContent />
      </Suspense>
    </ErrorBoundary>
  );
}
export default OnlineServicesPage;

