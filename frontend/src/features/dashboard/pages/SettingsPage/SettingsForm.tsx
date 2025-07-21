import { useForm, FormProvider } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

import { AppearanceSection } from '@/features/dashboard/components/organisms/SettingsPage/AppearanceSection';
import { DangerZoneSection } from '@/features/dashboard/components/organisms/SettingsPage/DangerZoneSection';
// import { ProfileSection } from '@/features/dashboard/components/organisms/SettingsPage/ProfileSection';
// import { NotificationSection } from '@/features/dashboard/components/organisms/SettingsPage/NotificationSection';

import { useThemeStore } from '@/core/theme/stores/useThemeStore';
import { useAuthContext } from '@/core/auth/context-provider/AuthContext';


const formSchema = z.object({
  notifications: z.boolean(),
  theme: z.enum(["light", "dark", "system"]),
  ui: z.enum(["web", "console"]),
  firstName: z.string().min(1, 'First name is required'),
  lastName: z.string().optional(),
});

export type FormValues = z.infer<typeof formSchema>;

export function SettingsForm() {
  const { mode, uiMode, isSystemPreference } = useThemeStore();
  const { actions } = useThemeStore();
  const { user } = useAuthContext();

  const formMethods = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      notifications: false,
      theme: isSystemPreference ? 'system' : mode,
      ui: uiMode,
      firstName: user?.user_metadata?.firstName || '',
      lastName: user?.user_metadata?.lastName || '',
    },
  })

  function onSubmit(data: FormValues) {
    console.log("Submitting form data:", data);

    // Apply theme changes based on form data
    if (data.theme === 'system') {
      actions.enableSystemPreference();
    } else {
      actions.disableSystemPreference();
      actions.changeTheme(data.theme as 'light' | 'dark');
    }

    // Apply UI mode changes
    actions.changeUIMode(data.ui);
  }

  return (
    <FormProvider {...formMethods}>
      <form
        onSubmit={formMethods.handleSubmit(onSubmit)}
        className="col-span-full space-y-6 p-6"
      >
        {/* <ProfileSection /> */}
        {/* <NotificationSection /> */}
        <AppearanceSection />
        <DangerZoneSection />

        <button
          type="submit"
          className="mt-4 px-6 py-2 bg-blue-600 text-white rounded"
        >
          Save Settings
        </button>
      </form>
    </FormProvider>
  )
}