import { useForm, FormProvider } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

import { AppearanceSection } from '@/features/dashboard/components/organisms/SettingsPage/AppearanceSection';
import { DangerZoneSection } from '@/features/dashboard/components/organisms/SettingsPage/DangerZoneSection';

import { useThemeStore } from '@/core/theme/stores/useThemeStore';


const formSchema = z.object({
  notifications: z.boolean(),
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
      notifications: false,
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