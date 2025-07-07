import { useForm, FormProvider } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

import { AppearanceSection } from '@/features/dashboard/components/organisms/SettingsPage/AppearanceSection';
//import { NotificationSection } from '@/features/dashboard/components/organisms/SettingsPage/NotificationSection';
import { DangerZoneSection } from '@/features/dashboard/components/organisms/SettingsPage/DangerZoneSection';
import { ExperimentalSection } from '@/features/dashboard/components/organisms/SettingsPage/Experimental';

import { useThemeStore } from '@/core/theme/stores/useThemeStore';


const formSchema = z.object({
  notificationLevel: z.enum(["everything", "available", "ignoring"]),
  theme: z.enum(["light", "dark", "system"]),
  ui: z.enum(["web", "console"]),
});

export type FormValues = z.infer<typeof formSchema>;

export function SettingsForm() {
  const { mode, uiMode, isSystemPreference } = useThemeStore();
  const { actions } = useThemeStore();

  const formMethods = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      notificationLevel: 'available',
      theme: isSystemPreference ? 'system' : mode,
      ui: uiMode,
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

        <AppearanceSection />
        {/* <NotificationSection /> */}
        <ExperimentalSection />
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