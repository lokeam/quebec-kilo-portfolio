import { ErrorBoundary } from 'react-error-boundary';
import { MediaStoragePageErrorFallback } from './MediaStoragePageErrorFallback';
import { MediaStoragePageContent } from '@/features/dashboard/pages/MediaStoragePage/MediaStoragePageContent';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';

function MediaStoragePage() {
  return (
    <ErrorBoundary
      FallbackComponent={MediaStoragePageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <MediaStoragePageContent />
      </Suspense>
    </ErrorBoundary>
  );
}

export default MediaStoragePage;
