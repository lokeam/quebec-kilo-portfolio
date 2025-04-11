import { Suspense } from 'react';

// Components
import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton'
import { MediaPageLocationForm, FormSchema } from '@/features/dashboard/components/organisms/MediaStoragePage/PhysicalLocationFormSingle/PhysicalLocationFormSingle';

// Hooks
import { useNavigate } from 'react-router-dom';
import { useOnboardingStore } from '@/features/dashboard/lib/stores/onboarding/onboardingStore';


// Types
import { PhysicalLocationType } from '@/features/dashboard/lib/types/media-storage/constants';
import type { z } from 'zod';


// Constants
import { NAVIGATION_ROUTES } from '@/features/dashboard/lib/types/onboarding/constants';

export default function OnboardingPagePhysical() {
  const navigate = useNavigate();
  const physicalDraft = useOnboardingStore((state) => state.physicalDraft);
  const updatePhysicalDraft = useOnboardingStore((state) => state.updatePhysicalDraft);

  /* Build partial default values */
  const defaultValues = {
    locationName: physicalDraft?.name ?? '',
    locationType: physicalDraft?.locationType ?? '',
    coordinates: {
      enabled: Boolean(physicalDraft?.mapCoordinates),
      value: physicalDraft?.mapCoordinates ?? '',
    },
  };

  const handleFormSuccess = (data: z.infer<typeof FormSchema>) => {
    console.log('Form data received:', data);

    if (!data) {
      console.error('No form data received');
      return;
    }

    updatePhysicalDraft({
      name: data.locationName,
      locationType: data.locationType as PhysicalLocationType,
      mapCoordinates: data.coordinates.enabled ? data.coordinates.value : undefined,
    });

    navigate(NAVIGATION_ROUTES.ONBOARDING_SELECT_PHYSICAL_SUB);
  }

  return (
    <ErrorBoundary
      FallbackComponent={OnlineServicesPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>

        <div className="mx-auto flex h-screen max-w-3xl flex-col items-center justify-center overflow-x-hidden">
          <h1 className="text-3xl font-bold mb-6">Let's learn where you are storing your games</h1>
          <p className="text-lg mb-8">Where do we go to find them?</p>

          <MediaPageLocationForm
            onSuccess={handleFormSuccess}
            defaultValues={defaultValues}
            buttonText="Continue"
          />
        </div>
      </Suspense>
    </ErrorBoundary>
  );
}