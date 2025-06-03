import { ErrorBoundary } from 'react-error-boundary';
import { PhysicalLocationsPageErrorFallback } from './PhysicalLocationsPageErrorFallback';
import { PhysicalLocationsPageContent } from '@/features/dashboard/pages/PhysicalLocations/PhysicalLocationsPageContent';
import { Suspense } from 'react';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';

function PhysicalLocationsPage() {
  return (
    <ErrorBoundary
      FallbackComponent={PhysicalLocationsPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <PhysicalLocationsPageContent />
      </Suspense>
    </ErrorBoundary>
  );
}
export default PhysicalLocationsPage;

