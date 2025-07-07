import { useFormContext } from 'react-hook-form';
import { IconDeviceDesktopUp, IconDeviceGamepad3 } from '@tabler/icons-react';
import { Button } from '@/shared/components/ui/button';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/shared/components/ui/card';

import { Switch } from '@/shared/components/ui/switch';

import type { FormValues } from '@/features/dashboard/pages/SettingsPage/SettingsForm';
import { useThemeStore } from '@/core/theme/stores/useThemeStore';

export function ExperimentalSection() {
  const form = useFormContext<FormValues>();
  const { watch, setValue } = form;
  const { actions } = useThemeStore();

  const handleUIModeChange = (value: string) => {
    setValue('ui', value as FormValues['ui']);
    actions.changeUIMode(value as 'web' | 'console');
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Experimental</CardTitle>
        <CardDescription>
          Features not yet fully ready for production. Toggle to be notified when new features are added.
        </CardDescription>
      </CardHeader>
      <CardContent className="grid gap-6">
        {/* Steam batch import */}
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <p className="text-sm font-medium leading-none">Import Games from Steam Library</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Sync your Q-Ko account with your Steam library.
            </p>
          </div>
          <Switch />
        </div>
        {/* Wishlist */}
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <p className="text-sm font-medium leading-none">Wishlist</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Save games to your wishlist and get notified when they are on sale.
            </p>
          </div>
          <Switch />
        </div>
      </CardContent>
    </Card>
  )
}
