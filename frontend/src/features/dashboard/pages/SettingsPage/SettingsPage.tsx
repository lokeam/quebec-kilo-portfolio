import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';
import { SettingsPageContent } from '@/features/dashboard/pages/SettingsPage/SettingsPageContent';

function SettingsPage() {
  return (
    <ErrorBoundary FallbackComponent={OnlineServicesPageErrorFallback} resetKeys={[location.pathname]}>
      <Suspense fallback={<HomePageSkeleton />}>
        <SettingsPageContent />
      </Suspense>
    </ErrorBoundary>
  )
}

export default SettingsPage;
