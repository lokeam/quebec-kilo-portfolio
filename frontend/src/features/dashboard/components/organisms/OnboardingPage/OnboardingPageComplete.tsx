
import { useNavigate } from 'react-router';
import { Button, } from '@/shared/components/ui/button';

import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';
import { Suspense } from 'react';


export default function OnboardingPageComplete() {
  const navigate = useNavigate();

  return (
    <ErrorBoundary
      FallbackComponent={OnlineServicesPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>

      <div className="mx-auto flex h-screen max-w-3xl flex-col items-center justify-center overflow-x-hidden">
        {/* Header Section */}
        <div className="text-center max-w-3xl mx-auto mb-16 md:mb-24">
          <h1 className="text-4xl md:text-5xl font-bold mb-6">You&apos;re all set</h1>
          <p className="text-lg md:text-xl text-gray-400">Here&apos;s a few quick tips to get you started:</p>
        </div>

        {/* Features Section */}
        <div className="w-full max-w-3xl mx-auto space-y-12 md:space-y-16">
          {/* Feature 1 */}
          <div className="flex items-start gap-6">
            <div className="flex-shrink-0 w-12 h-12 rounded-full border-2 border-white flex items-center justify-center">
              <span className="text-xl font-semibold">1</span>
            </div>
            <div>
              <h2 className="text-xl md:text-2xl font-semibold mb-2">First add locations or online services</h2>
              <p className="text-gray-400 text-lg">You need a place to store your stuff before you add anything</p>
            </div>
          </div>

          {/* Feature 2 */}
          <div className="flex items-start gap-6">
            <div className="flex-shrink-0 w-12 h-12 rounded-full border-2 border-white flex items-center justify-center">
              <span className="text-xl font-semibold">2</span>
            </div>
            <div>
              <h2 className="text-xl md:text-2xl font-semibold mb-2">Next add items</h2>
              <p className="text-gray-400 text-lg">
                Use top search bar to search for games to add to your library or wishlist
              </p>
            </div>
          </div>

          {/* Feature 3 */}
          <div className="flex items-start gap-6">
            <div className="flex-shrink-0 w-12 h-12 rounded-full border-2 border-white flex items-center justify-center">
              <span className="text-xl font-semibold">3</span>
            </div>
            <div>
              <h2 className="text-xl md:text-2xl font-semibold mb-2">Track spending and deals on games</h2>
              <p className="text-gray-400 text-lg">
                Get notified of whenever a wishlisted game goes on sale and how much you&apos;re spending every month
              </p>
            </div>
          </div>
        </div>

        {/* CTA Button */}
        <div className="mt-12 md:mt-16">
          <Button
            size="lg"
            className="bg-white text-black hover:bg-gray-200 text-lg px-8 py-6 h-auto"
            onClick={() => navigate('/')}
          >
            GO!
          </Button>
        </div>
      </div>
      </Suspense>
    </ErrorBoundary>
  );
}