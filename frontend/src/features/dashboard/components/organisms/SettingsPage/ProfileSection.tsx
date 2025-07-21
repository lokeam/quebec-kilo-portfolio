import { useFormContext } from 'react-hook-form';
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
import { useState } from 'react';

export function ProfileSection() {
  const form = useFormContext<FormValues>();
  const updateUserProfileMutation = useUpdateUserProfile();
  const [isUpdating, setIsUpdating] = useState(false);

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

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            Profile Information
          </CardTitle>
          <CardDescription>
            Update your name displayed within QKO.
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
                      className="bg-background text-foreground w-1/2"
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
                      className="bg-background text-foreground w-1/2"
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
    </div>
  );
}