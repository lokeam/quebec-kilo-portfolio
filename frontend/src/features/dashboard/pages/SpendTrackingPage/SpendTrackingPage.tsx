import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';
import { SpendTrackingPageContent } from '@/features/dashboard/pages/SpendTrackingPage/SpendTrackingPageContent';

function SpendTrackingPage() {
  return (
    <ErrorBoundary
      FallbackComponent={OnlineServicesPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <SpendTrackingPageContent />
      </Suspense>
    </ErrorBoundary>
  )
}

export default SpendTrackingPage;
