import { ErrorBoundary } from 'react-error-boundary';
import { HomePageErrorFallback } from './HomePageErrorFallback';
import { HomePageContent } from '@/features/dashboard/pages/HomePage/HomePageContent';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';

function HomePage() {
  return (
    <ErrorBoundary
      FallbackComponent={HomePageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <HomePageContent />
      </Suspense>
    </ErrorBoundary>
  );
}
export default HomePage;

