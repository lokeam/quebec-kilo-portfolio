import { useFormContext } from 'react-hook-form';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/shared/components/ui/card';
import { Switch } from '@/shared/components/ui/switch';

import type { FormValues } from '@/features/dashboard/pages/SettingsPage/SettingsForm';

export function PrivacySection() {
  const { watch, setValue } = useFormContext<FormValues>()

  return (
    <Card>

      <CardHeader>
        <CardTitle>Privacy</CardTitle>
        <CardDescription>Manage your privacy settings.</CardDescription>
      </CardHeader>

      <CardContent className="grid gap-4">
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <p className="text-sm font-medium leading-none">Share Usage Data</p>
            <p className="text-sm text-gray-400">
              Help us improve the product by sharing anonymous usage data.
            </p>
          </div>
          <Switch
            checked={watch("shareUsageData")}
            onCheckedChange={(checked) => setValue("shareUsageData", checked)}
          />
        </div>

        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <p className="text-sm font-medium leading-none">Allow Third-Party Cookies</p>
            <p className="text-sm text-gray-400">
              Enable third-party cookies for personalized content.
            </p>
          </div>
          <Switch
            checked={watch("allowThirdPartyCookies")}
            onCheckedChange={(checked) => setValue("allowThirdPartyCookies", checked)}
          />
        </div>
      </CardContent>
    </Card>
  );
}
