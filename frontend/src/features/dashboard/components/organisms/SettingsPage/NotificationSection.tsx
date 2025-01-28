
import { useFormContext } from 'react-hook-form';

// ShadCN UI Components
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/shared/components/ui/card';

// Icons
import { Bell, EyeOff } from 'lucide-react';
import { IconHeart } from '@tabler/icons-react';

/* Types from main form */
import type { FormValues } from '@/features/dashboard/pages/SettingsPage/SettingsForm';

export function NotificationSection() {
  /* Access form methods via useFormContext */
  const { watch, setValue } = useFormContext<FormValues>();

  /* Destructure the current notification level */
  const notificationLevel = watch("notificationLevel");

  /* Helper to set notification level */
  const handleNotificationChange = (value: "everything" | "available" | "ignoring") => {
    setValue("notificationLevel", value)
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Notifications</CardTitle>
        <CardDescription>Choose what you want to be notified about.</CardDescription>
      </CardHeader>
      <CardContent className="grid gap-1">
        {/* EVERYTHING */}
        <div
          className={`-mx-2 flex items-start space-x-4 rounded-md p-2 transition-all hover:bg-gray-100 dark:hover:bg-gray-800 cursor-pointer ${
            notificationLevel === "everything" ? "bg-gray-100 dark:bg-gray-800" : ""
          }`}
          onClick={() => handleNotificationChange("everything")}
        >
          <Bell className="mt-px h-5 w-5" />
          <div className="space-y-1">
            <p className="text-sm font-medium leading-none">Everything</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              App critical notifications, email digest, and wishlist deals.
            </p>
          </div>
        </div>

        {/* AVAILABLE */}
        <div
          className={`-mx-2 flex items-start space-x-4 rounded-md p-2 transition-all hover:bg-gray-100 dark:hover:bg-gray-800 cursor-pointer ${
            notificationLevel === "available" ? "bg-gray-100 dark:bg-gray-800" : ""
          }`}
          onClick={() => handleNotificationChange("available")}
        >
          <IconHeart className="mt-px h-5 w-5" />
          <div className="space-y-1">
            <p className="text-sm font-medium leading-none">Available</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Only when wishlist items are available.
            </p>
          </div>
        </div>

        {/* IGNORING */}
        <div
          className={`-mx-2 flex items-start space-x-4 rounded-md p-2 transition-all hover:bg-gray-100 dark:hover:bg-gray-800 cursor-pointer ${
            notificationLevel === "ignoring" ? "bg-gray-100 dark:bg-gray-800" : ""
          }`}
          onClick={() => handleNotificationChange("ignoring")}
        >
          <EyeOff className="mt-px h-5 w-5" />
          <div className="space-y-1">
            <p className="text-sm font-medium leading-none">Ignoring</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">Turn off all notifications.</p>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
