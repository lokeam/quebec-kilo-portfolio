import { useForm, FormProvider } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

import { AppearanceSection } from '@/features/dashboard/components/organisms/SettingsPage/AppearanceSection';
import { NotificationSection } from '@/features/dashboard/components/organisms/SettingsPage/NotificationSection';
import { DangerZoneSection } from '@/features/dashboard/components/organisms/SettingsPage/DangerZoneSection';


const formSchema = z.object({
  notificationLevel: z.enum(["everything", "available", "ignoring"]),
  theme: z.enum(["light", "dark", "system"]),
});

export type FormValues = z.infer<typeof formSchema>;

export function SettingsForm() {
  const formMethods = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      notificationLevel: "available",
      theme: "light",
    },
  })

  function onSubmit(data: FormValues) {
    console.log("Submitting form data:", data)
  }

  return (
    <FormProvider {...formMethods}>
      <form
        onSubmit={formMethods.handleSubmit(onSubmit)}
        className="col-span-full space-y-6 p-6"
      >

        <AppearanceSection />
        <NotificationSection />
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