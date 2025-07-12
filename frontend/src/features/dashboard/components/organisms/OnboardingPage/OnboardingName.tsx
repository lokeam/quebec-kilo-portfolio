import { Suspense } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

// Components
import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton';

// ShadCN UI Components
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader } from '@/shared/components/ui/card';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/shared/components/ui/form';
import { Input } from '@/shared/components/ui/input';

// Hooks
import { useAuth } from '@/core/auth/hooks/useAuth';
import { useCreateUserProfile } from '@/core/api/queries/user.queries';

// Constants
import { NAVIGATION_ROUTES, STAGGER_CHILD_VARIANTS } from '@/features/dashboard/lib/types/onboarding/constants';

import { BrandLogo } from '@/features/navigation/organisms/SideNav/BrandLogo';

// Motion
import { motion } from 'framer-motion';

// Form Schema
const formSchema = z.object({
  firstName: z.string().min(1, 'First name is required'),
  lastName: z.string().optional(),
});

type FormValues = z.infer<typeof formSchema>;

export default function OnboardingName() {
  const navigate = useNavigate();
  const { user } = useAuth(); // Only get auth data
  const createUserProfileMutation = useCreateUserProfile(); // Get the mutation

  // Check if user already has a complete profile (has firstName and lastName)
  useEffect(() => {
    if (user?.user_metadata?.firstName && user?.user_metadata?.lastName) {
      console.log('🚫 User already has complete profile, redirecting to intro');
      navigate(NAVIGATION_ROUTES.ONBOARDING_INTRO, { replace: true });
    }
  }, [user, navigate]);

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      firstName: '',
      lastName: '',
    },
  });

  const onSubmit = async (data: FormValues) => {
    console.log('📝 Form submitted:', data);

    try {
      // Create user profile ONLY when user explicitly submits the form
      await createUserProfileMutation.mutateAsync({
        email: user?.email || '',
        auth0UserID: user?.sub || '',
        firstName: data.firstName,
        lastName: data.lastName || '',
      });

      console.log('✅ User profile created successfully');

      // Navigate to next onboarding step
      navigate(NAVIGATION_ROUTES.ONBOARDING_INTRO);
    } catch (error) {
      console.error('❌ Failed to create user profile:', error);
      // Handle error (show toast, etc.)
    }
  };

  return (
    <ErrorBoundary
      FallbackComponent={OnlineServicesPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <div className="mx-auto flex h-screen max-w-3xl flex-col items-center justify-center overflow-x-hidden">
          <motion.div
            className="z-10"
            exit={{ opacity: 0, scale: 0.95 }}
            transition={{ duration: 0.3, type: "spring" }}
          >
            <motion.div
              variants={{
                show: {
                  transition: {
                    staggerChildren: 0.2,
                  },
                },
              }}
              initial="hidden"
              animate="show"
              className="mx-5 flex flex-col items-center space-y-10 text-center sm:mx-auto"
            >
              <motion.div
                className="font-display text-4xl font-bold text-foreground transition-colors sm:text-5xl"
                variants={STAGGER_CHILD_VARIANTS}
              >
                <BrandLogo />
              </motion.div>

              <motion.div
                className="w-full max-w-md"
                variants={STAGGER_CHILD_VARIANTS}
              >
                <Card>
                  <CardHeader>
                    <h2 className="text-2xl font-semibold text-center">What's your name?</h2>
                    <p className="text-sm text-muted-foreground text-center">
                      You can change this anytime in your profile settings after completing onboarding.
                    </p>
                  </CardHeader>
                  <CardContent className="space-y-6">
                    <Form {...form}>
                      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <FormField
                          control={form.control}
                          name="firstName"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>First Name *</FormLabel>
                              <FormControl>
                                <Input
                                  placeholder="Enter your first name"
                                  {...field}
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <FormField
                          control={form.control}
                          name="lastName"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Last Name (Optional)</FormLabel>
                              <FormControl>
                                <Input
                                  placeholder="Enter your last name"
                                  {...field}
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <motion.div
                          variants={STAGGER_CHILD_VARIANTS}
                          className="pt-4"
                        >
                          <Button
                            type="submit"
                            size="lg"
                            className="w-full"
                            disabled={createUserProfileMutation.isPending}
                          >
                            {createUserProfileMutation.isPending ? 'Creating Profile...' : 'Nice to meet you!'}
                          </Button>
                        </motion.div>
                      </form>
                    </Form>
                  </CardContent>
                </Card>
              </motion.div>
            </motion.div>
          </motion.div>
        </div>
      </Suspense>
    </ErrorBoundary>
  );
}