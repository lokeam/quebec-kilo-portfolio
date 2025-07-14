import { ErrorBoundary } from 'react-error-boundary';
import { DashboardErrorFallback } from '@/core/error/pages/DashboardErrorFallback';
import { NotificationsPageContent } from '@/features/dashboard/pages/NotificationsPage/NotificationsPageContent';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';

function NotificationsPage() {
  return (
    <ErrorBoundary
      FallbackComponent={DashboardErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <NotificationsPageContent />
      </Suspense>
    </ErrorBoundary>
  );
}

export default NotificationsPage;
