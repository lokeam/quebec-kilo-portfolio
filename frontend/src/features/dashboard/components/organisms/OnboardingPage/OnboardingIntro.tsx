import { Suspense } from 'react';

// Components
import { ErrorBoundary } from 'react-error-boundary';
import { DashboardErrorFallback } from '@/core/error/pages/DashboardErrorFallback';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton'

// ShadCN UI Components
import { Button, } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader } from '@/shared/components/ui/card';

// Hooks
import { useNavigate } from 'react-router-dom';
import { useOnboardingStore } from '@/features/dashboard/lib/stores/onboarding/onboardingStore';
import { useAuthContext } from '@/core/auth/context-provider/AuthContext';

// API
import { useUpdateUserMetadata } from '@/core/api/queries/user.queries';

// Constants
import { NAVIGATION_ROUTES } from '@/types/domain/onboarding';

export default function OnboardingIntro() {
  const navigate = useNavigate();
  const setWantsSetup = useOnboardingStore((state) => state.setWantsSetup);
  const { user } = useAuthContext();
  const updateUserMetadataMutation = useUpdateUserMetadata();

  const handleStartOnboardFlow = async () => {
    if (user) {
      try {
        // Call our backend endpoint to update Auth0 metadata
        await updateUserMetadataMutation.mutateAsync({
          wantsIntroToasts: true
        });
      } catch (error) {
        console.error('Failed to update user metadata:', error);
        // Continue with onboarding even if metadata update fails
      }
    }
    setWantsSetup(true);
    navigate(NAVIGATION_ROUTES.ONBOARDING_SELECT_STORAGE);
  }

  const handleSkipOnboardFlow = async () => {
    if (user) {
      try {
        // Call our backend endpoint to update Auth0 metadata
        await updateUserMetadataMutation.mutateAsync({
          wantsIntroToasts: false
        });
      } catch (error) {
        console.error('Failed to update user metadata:', error);
        // Continue with onboarding even if metadata update fails
      }
    }
    setWantsSetup(false);
    navigate(NAVIGATION_ROUTES.ONBOARDING_COMPLETE);
  }


  return (
    <ErrorBoundary
      FallbackComponent={DashboardErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>

        <div className="mx-auto flex h-screen max-w-3xl flex-col items-center justify-center overflow-x-hidden">
          <Card className="w-full max-w-3xl">
            <CardHeader>
              <h1 className="text-3xl font-bold text-center">Is this your first time here?</h1>
            </CardHeader>
            <CardContent className="space-y-8">

              {/* Display component copy */}
              {/* <OnboardingIntroContent /> */}
              <div className="mt-3 text-lg text-center">
                <p className="mb-3">Shall we show helper messages to guide you after logging in?</p>
                <p className="italic text-muted-foreground">You'll only see them once.</p>
              </div>

              <div className="flex flex-col items-center gap-4 pt-4">
                <Button
                  size="lg"
                  className="w-full max-w-md"
                  onClick={handleStartOnboardFlow}
                >
                  Yes, please!
                </Button>
                <Button
                  variant="outline"
                  size="lg"
                  className="w-full max-w-md text-muted-foreground"
                  onClick={handleSkipOnboardFlow}
                >
                  No thanks, I&apos;ll figure QKO out myself
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </Suspense>
    </ErrorBoundary>
  );
}