import { useFormContext } from 'react-hook-form';
import { User } from 'lucide-react';
import { Button } from '@/shared/components/ui/button';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/shared/components/ui/card';
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/components/ui/form';
import { Input } from '@/shared/components/ui/input';

import type { FormValues } from '@/features/dashboard/pages/SettingsPage/SettingsForm';
import { useUpdateUserProfile } from '@/core/api/queries/user.queries';
import { useIntroToasts } from '@/core/hooks/useIntroToasts';
import { useState } from 'react';

export function ProfileSection() {
  const form = useFormContext<FormValues>();
  const updateUserProfileMutation = useUpdateUserProfile();
  const { wantsIntroToasts, updatePreference } = useIntroToasts();
  const [isUpdating, setIsUpdating] = useState(false);
  const [isUpdatingToasts, setIsUpdatingToasts] = useState(false);

  const handleUpdateProfile = async () => {
    setIsUpdating(true);

    try {
      const values = form.getValues();
      await updateUserProfileMutation.mutateAsync({
        firstName: values.firstName,
        lastName: values.lastName || '',
      });

      console.log('✅ User profile updated successfully');
      // You could show a success toast here
    } catch (error) {
      console.error('❌ Failed to update user profile:', error);
      // Handle error (show toast, etc.)
    } finally {
      setIsUpdating(false);
    }
  };

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
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <User className="h-5 w-5" />
            Profile Information
          </CardTitle>
          <CardDescription>
            Update your personal information. Changes will be saved to both your account and Auth0.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-6">
          <div className="grid gap-4">
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

            <div className="flex justify-end">
              <Button
                type="button"
                onClick={handleUpdateProfile}
                disabled={isUpdating}
                size="sm"
              >
                {isUpdating ? 'Updating...' : 'Update Profile'}
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Intro Toasts</CardTitle>
          <CardDescription>
            Show helpful tips and guidance when using new features
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium">Enable Intro Toasts</p>
              <p className="text-xs text-muted-foreground">
                Show contextual help when using features for the first time
              </p>
            </div>
            <Button
              variant={wantsIntroToasts ? "default" : "outline"}
              onClick={handleToggleIntroToasts}
              disabled={isUpdatingToasts}
            >
              {isUpdatingToasts ? "Updating..." : wantsIntroToasts ? "Enabled" : "Disabled"}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}