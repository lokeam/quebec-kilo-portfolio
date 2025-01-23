import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { NotificationsPageContent } from '@/features/dashboard/pages/NotificationsPage/NotificationsPageContent';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';

function NotificationsPage() {
  return (
    <ErrorBoundary
      FallbackComponent={OnlineServicesPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <NotificationsPageContent />
      </Suspense>
    </ErrorBoundary>
  );
}

export default NotificationsPage;
