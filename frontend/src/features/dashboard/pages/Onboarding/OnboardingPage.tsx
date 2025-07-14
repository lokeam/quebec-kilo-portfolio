import { ErrorBoundary } from 'react-error-boundary';
import { DashboardErrorFallback } from '@/core/error/pages/DashboardErrorFallback';
import { OnboardingPageContent } from '@/features/dashboard/pages/Onboarding/OnboardingPageContent';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';

function OnboardingPage() {
  return (
    <ErrorBoundary
      FallbackComponent={DashboardErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <OnboardingPageContent />
      </Suspense>
    </ErrorBoundary>
  );
}
export default OnboardingPage;

