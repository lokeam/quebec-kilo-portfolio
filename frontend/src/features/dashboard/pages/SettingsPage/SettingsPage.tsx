import { ErrorBoundary } from 'react-error-boundary';
import { DashboardErrorFallback } from '@/core/error/pages/DashboardErrorFallback';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';
import { SettingsPageContent } from '@/features/dashboard/pages/SettingsPage/SettingsPageContent';

function SettingsPage() {
  return (
    <ErrorBoundary FallbackComponent={DashboardErrorFallback} resetKeys={[location.pathname]}>
      <Suspense fallback={<HomePageSkeleton />}>
        <SettingsPageContent />
      </Suspense>
    </ErrorBoundary>
  )
}

export default SettingsPage;
