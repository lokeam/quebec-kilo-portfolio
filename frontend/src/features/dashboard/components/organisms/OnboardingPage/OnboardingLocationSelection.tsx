
import { Suspense } from 'react';
import { useForm, Controller } from "react-hook-form"

// Components
import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton'

// ShadCN UI Components
import { Button, } from '@/shared/components/ui/button';
import { Form } from '@/shared/components/ui/form';
import { RadioGroup, RadioGroupItem } from '@radix-ui/react-radio-group';
import { cn } from '@/shared/components/ui/utils';

// Hooks
import { useNavigate } from 'react-router-dom';
import { useSetOnboardingStorageType } from '@/features/dashboard/lib/stores/onboarding/onboardingStore';
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"

// Types
import type { OnboardingStorageType } from '@/features/dashboard/lib/types/onboarding/base';

// Icons
import { Disc, Cloud, HardDrive } from 'lucide-react';

// Consts
import { NAVIGATION_ROUTES } from '@/features/dashboard/lib/types/onboarding/constants';

// Form Schema
const formSchema = z.object({
  libraryType: z.enum(['physical', 'digital', 'both'], {
    required_error: 'Please select a library type.',
  }),
});
type FormValues = z.infer<typeof formSchema>


// Constants
const libraryOptions = [
  {
    value: 'physical',
    label: 'I only have physical games (discs)',
    icon: Disc,
  },
  {
    value: 'digital',
    label: 'All my games live on an online service (Steam, Playstation Network, etc.)',
    icon: Cloud,
  },
  {
    value: 'both',
    label: 'I have both physical and digital games',
    icon: HardDrive,
  },
];

export default function OnboardingLocationSelection() {
  const navigate = useNavigate();
  const setStorageType = useSetOnboardingStorageType();

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: { libraryType: 'physical' },
  });

  function onSubmit(data: FormValues) {

    let storageSelection: OnboardingStorageType;
    switch (data.libraryType) {
      case 'physical':
        storageSelection = 'PHYSICAL';
        break;
      case 'digital':
        storageSelection = 'DIGITAL';
        break;
      case 'both':
        storageSelection = 'PHYSICAL_AND_DIGITAL';
        break;
    }

    /* Update Zustand store */
    setStorageType(storageSelection);

    if (storageSelection === 'DIGITAL') {
      navigate(NAVIGATION_ROUTES.ONBOARDING_SELECT_DIGITAL);
    } else {
      navigate(NAVIGATION_ROUTES.ONBOARDING_SELECT_PHYSICAL_MAIN);
    }
  }

  return (
    <ErrorBoundary
      FallbackComponent={OnlineServicesPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
      <div className="mx-auto flex h-screen max-w-screen-2xl flex-col items-center justify-center overflow-x-hidden p-4">
        <p className="text-lg mb-8">(1 of 3)</p>
        <h1 className="text-3xl font-bold mb-6">Tell us a bit about your game library</h1>
        <p className="text-lg mb-8">This helps us set up your account correctly</p>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="space-y-8"
          >
            <Controller
              name="libraryType"
              control={form.control}
              render={({ field }) => (
                <div className="space-y-4">
                  <p className="text-xl text-center font-semibold">(Select one)</p>
                  <RadioGroup
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                    className="flex flex-col sm:flex-row gap-4"
                  >
                    {libraryOptions.map((option) => (
                      <RadioGroupItem
                        key={option.value}
                        value={option.value}
                        className={cn(
                          "relative flex-1 group ring-1 ring-border rounded-lg p-4 cursor-pointer",
                          "data-[state=checked]:ring-2 data-[state=checked]:ring-primary",
                        )}
                      >
                        <div className="flex flex-col items-center text-center">
                          <option.icon className="w-12 h-12 mb-4 text-primary" />
                          <span className="font-medium">{option.label}</span>
                        </div>
                      </RadioGroupItem>
                    ))}
                  </RadioGroup>
                </div>
              )}
            />

            <p className="text-sm text-center text-muted-foreground">We'll learn more about one of these options in a moment.</p>
            <p className="text-sm text-center text-muted-foreground">This choice just helps us organize your library</p>

            <div className="flex justify-center">
              <Button
                className="w-80"
                type="submit"
              >
                Continue
              </Button>
            </div>
          </form>
        </Form>
      </div>
      </Suspense>
    </ErrorBoundary>
  );
}
