import { Suspense } from 'react';

// Components
import { OnboardingIntroContent } from '@/features/dashboard/components/organisms/OnboardingPage/OnboardingIntroContent';
import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton'

// ShadCN UI Components
import { Button, } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader } from '@/shared/components/ui/card';

// Hooks
import { useNavigate } from 'react-router-dom';
import { useOnboardingStore } from '@/features/dashboard/lib/stores/onboarding/onboardingStore';

// Consts
import { NAVIGATION_ROUTES } from '@/features/dashboard/lib/types/onboarding/constants';
export default function OnboardingIntro() {
  const navigate = useNavigate();
  const setWantsSetup = useOnboardingStore((state) => state.setWantsSetup);

  const handleStartOnboardFlow = () => {
    setWantsSetup(true);
    navigate(NAVIGATION_ROUTES.ONBOARDING_SELECT_STORAGE);
  }

  const handleSkipOnboardFlow = () => {
    setWantsSetup(false);
    navigate(NAVIGATION_ROUTES.ONBOARDING_COMPLETE);
  }


  return (
    <ErrorBoundary
      FallbackComponent={OnlineServicesPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>

        <div className="mx-auto flex h-screen max-w-3xl flex-col items-center justify-center overflow-x-hidden">
          <Card className="w-full max-w-3xl">
            <CardHeader>
              <h1 className="text-3xl font-bold text-center">Welcome to Q-Ko</h1>
            </CardHeader>
            <CardContent className="space-y-8">
              <div className="text-center text-muted-foreground">Here&apos;s what to expect:</div>

              {/* Display component copy */}
              <OnboardingIntroContent />

              <div className="flex flex-col items-center gap-4 pt-4">
                <Button
                  size="lg"
                  className="w-full max-w-md"
                  onClick={handleStartOnboardFlow}
                >
                  Let&apos;s go!
                </Button>
                <Button
                  variant="outline"
                  size="lg"
                  className="w-full max-w-md text-muted-foreground"
                  onClick={handleSkipOnboardFlow}
                >
                  No thanks, I&apos;ll figure out Q-ko by myself
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </Suspense>
    </ErrorBoundary>
  );
}