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

// Extend Window interface for optimistic state
declare global {
  interface Window {
    __ONBOARDING_OPTIMISTIC_COMPLETE__?: boolean;
    __WANTS_INTRO_TOASTS__?: boolean;
  }
}

export default function OnboardingToastSetup() {
  const navigate = useNavigate();
  const setWantsSetup = useOnboardingStore((state) => state.setWantsSetup);
  const { user } = useAuthContext();
  const updateUserMetadataMutation = useUpdateUserMetadata();

  // Set global optimistic state to bypass Auth0 checks
  const setOptimisticOnboardingComplete = () => {
    // Set a global flag that the onboarding hooks can check
    window.__ONBOARDING_OPTIMISTIC_COMPLETE__ = true;
    // Also store in localStorage to persist across refreshes
    localStorage.setItem('__ONBOARDING_OPTIMISTIC_COMPLETE__', 'true');
  };

  // Set optimistic intro toasts preference
  const setOptimisticIntroToasts = (wantsToasts: boolean) => {
    window.__WANTS_INTRO_TOASTS__ = wantsToasts;
    // Also store in localStorage to persist across refreshes
    localStorage.setItem('__WANTS_INTRO_TOASTS__', wantsToasts.toString());
  };

  const handleStartOnboardFlow = async () => {
    console.log('🚀 Starting onboarding flow with intro toasts enabled');

    // OPTIMISTIC UPDATE: Mark onboarding as complete and set intro toasts preference
    setOptimisticOnboardingComplete();
    setOptimisticIntroToasts(true);
    setWantsSetup(true);

    navigate(NAVIGATION_ROUTES.ONBOARDING_COMPLETE);

    // Then update backend asynchronously (fire and forget)
    if (user) {
      console.log('📤 Sending wantsIntroToasts: true to backend');
      updateUserMetadataMutation.mutateAsync({
        wantsIntroToasts: true
      }).then(() => {
        console.log('✅ Successfully updated user metadata with wantsIntroToasts: true');
      }).catch(error => {
        console.error('❌ Failed to update user metadata (background):', error);
        // User is already on dashboard, just log the error
      });
    }
  }

  const handleSkipOnboardFlow = async () => {
    console.log('🚀 Starting onboarding flow with intro toasts disabled');

    // OPTIMISTIC UPDATE: Mark onboarding as complete and set intro toasts preference
    setOptimisticOnboardingComplete();
    setOptimisticIntroToasts(false);
    setWantsSetup(false);

    navigate(NAVIGATION_ROUTES.ONBOARDING_COMPLETE);

    // Then update backend asynchronously (fire and forget)
    if (user) {
      console.log('📤 Sending wantsIntroToasts: false to backend');
      updateUserMetadataMutation.mutateAsync({
        wantsIntroToasts: false
      }).then(() => {
        console.log('✅ Successfully updated user metadata with wantsIntroToasts: false');
      }).catch(error => {
        console.error('❌ Failed to update user metadata (background):', error);
        // User is already on dashboard, just log the error
      });
    }
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