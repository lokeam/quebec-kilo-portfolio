import { Suspense } from 'react';

// Components
import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton'
import { MediaPageSublocationForm, FormSchema } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaPageSublocationForm/MediaPageSublocationForm';

// Hooks
import { useNavigate } from 'react-router-dom';
import { useOnboardingStore } from '@/features/dashboard/lib/stores/onboarding/onboardingStore';

// Types
import type { z } from 'zod';
import { SubLocationType } from '@/features/dashboard/lib/types/media-storage/constants';

// Constants
import { NAVIGATION_ROUTES } from '@/features/dashboard/lib/types/onboarding/constants';

export default function OnboardingPagePhysicalSublocations() {
  const navigate = useNavigate();
  const physicalDraft = useOnboardingStore((state) => state.physicalDraft);
  const updatePhysicalDraft = useOnboardingStore((state) => state.updatePhysicalDraft);

  /* Build partial default values */
  const defaultValues = {
    locationName: physicalDraft?.name ?? '',
    locationType: physicalDraft?.locationType ?? '',
    bgColor: physicalDraft?.bgColor ?? '',
  };

  const handleFormSuccess = (data: z.infer<typeof FormSchema>) => {
    console.log('sublocation form data received', data);

    if (!data) {
      console.error('No sublocation form data received');
      return;
    }

    updatePhysicalDraft({
      name: data.locationName,
      locationType: data.locationType as SubLocationType,
      bgColor: data.bgColor,
    });

    navigate(NAVIGATION_ROUTES.ONBOARDING_SELECT_DIGITAL);
  }

  return (
    <ErrorBoundary
      FallbackComponent={OnlineServicesPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>

        <div className="mx-auto flex h-screen max-w-3xl flex-col items-center justify-center overflow-x-hidden">
          <h1 className="text-3xl font-bold mb-6">Where do we find your games in your home?</h1>
          <p className="text-lg mb-8">We'll start with one area, and we'll add more as you go.</p>

          <MediaPageSublocationForm
            onSuccess={handleFormSuccess}
            defaultValues={defaultValues}
            buttonText="Continue"
          />
        </div>
      </Suspense>
    </ErrorBoundary>
  );
}