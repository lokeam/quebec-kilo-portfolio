import { ErrorBoundary } from 'react-error-boundary';
import { DashboardErrorFallback } from '@/core/error/pages/DashboardErrorFallback';
import { Suspense } from 'react';
import { useLocation } from 'react-router-dom';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';
import { SpendTrackingPageContent } from '@/features/dashboard/pages/SpendTrackingPage/SpendTrackingPageContent';

function SpendTrackingPage() {
  const location = useLocation();

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
