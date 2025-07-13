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
import { useIntroToasts } from '@/core/hooks/useIntroToasts';
import { useState } from 'react';

export function NotificationSection() {
  const form = useFormContext<FormValues>();
  const { watch, setValue } = form;
  const { wantsIntroToasts, updatePreference } = useIntroToasts();
  const [isUpdatingToasts, setIsUpdatingToasts] = useState(false);

  const handleToggleIntroToasts = async () => {
    setIsUpdatingToasts(true);
    try {
      await updatePreference(!wantsIntroToasts);
    } catch (error) {
      console.error('Failed to update intro toasts preference:', error);
    } finally {
      setIsUpdatingToasts(false);
    }
  };

  return (
    <Card>
    <CardHeader>
      <CardTitle>Tutorial Notifications</CardTitle>
      <CardDescription>
        Show helpful tips and guidance when using new features.
      </CardDescription>
    </CardHeader>
    <CardContent>
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium">Enable Intro Toasts</p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Display contextual help when using features for the first time.
          </p>
        </div>
        <Switch
          checked={wantsIntroToasts}
          onCheckedChange={handleToggleIntroToasts}
          disabled={isUpdatingToasts}
        />
      </div>
    </CardContent>
  </Card>
  )
}
