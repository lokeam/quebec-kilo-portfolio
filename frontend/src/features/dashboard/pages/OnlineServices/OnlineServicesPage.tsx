import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from './OnlineServicesPageErrorFallback';
import { OnlineServicesPageContent } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageContent';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';

function OnlineServicesPage() {
  return (
    <ErrorBoundary
      FallbackComponent={OnlineServicesPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <OnlineServicesPageContent />
      </Suspense>
    </ErrorBoundary>
  );
}
export default OnlineServicesPage;

