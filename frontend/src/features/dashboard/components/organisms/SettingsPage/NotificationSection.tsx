import { useFormContext } from 'react-hook-form';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/shared/components/ui/card';

import { Switch } from '@/shared/components/ui/switch';

import type { FormValues } from '@/features/dashboard/pages/SettingsPage/SettingsForm';

export function NotificationSection() {
  const form = useFormContext<FormValues>();
  const { watch, setValue } = form;

  const handleNotificationChange = (value: boolean) => {
    setValue('notifications', value);
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Notifications</CardTitle>
        <CardDescription>
          Toggle to be notified when new features are added.
        </CardDescription>
      </CardHeader>
      <CardContent className="grid gap-6">
        {/* Steam batch import */}
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <p className="text-sm font-medium leading-none">Announcements and update emails</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">
            Receive occasional emails about product launches and new features from Q-KO
            </p>
          </div>
          <Switch
            checked={watch('notifications')}
            onCheckedChange={handleNotificationChange}
          />
        </div>
      </CardContent>
    </Card>
  )
}
