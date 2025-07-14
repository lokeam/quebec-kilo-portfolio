import { ErrorBoundary } from 'react-error-boundary';
import { DashboardErrorFallback } from '@/core/error/pages/DashboardErrorFallback';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';
import { SpendTrackingPageContent } from '@/features/dashboard/pages/SpendTrackingPage/SpendTrackingPageContent';

function SpendTrackingPage() {
  return (
    <ErrorBoundary
      FallbackComponent={DashboardErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <SpendTrackingPageContent />
      </Suspense>
    </ErrorBoundary>
  )
}

export default SpendTrackingPage;
