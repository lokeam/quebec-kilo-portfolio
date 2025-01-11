import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';
import { WishListPageContent } from './WishListPageContent';

function WishListPage() {
  return (
    <ErrorBoundary FallbackComponent={OnlineServicesPageErrorFallback} resetKeys={[location.pathname]}>
      <Suspense fallback={<HomePageSkeleton />}>
        <WishListPageContent />
      </Suspense>
    </ErrorBoundary>
  );
}

export default WishListPage;
