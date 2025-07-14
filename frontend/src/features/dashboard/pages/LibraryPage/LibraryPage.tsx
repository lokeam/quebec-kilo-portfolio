import { Suspense } from 'react';
import { ErrorBoundary } from 'react-error-boundary';

// Custom Components
import { LibraryPageContent } from '@/features/dashboard/pages/LibraryPage/LibraryPageContent';

// Error Boundary
import { DashboardErrorFallback } from '@/core/error/pages/DashboardErrorFallback';

// Skeleton
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';


function LibraryPage() {
  return (
    <ErrorBoundary
      FallbackComponent={DashboardErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <LibraryPageContent />
      </Suspense>
    </ErrorBoundary>
  );
}

export default LibraryPage;
