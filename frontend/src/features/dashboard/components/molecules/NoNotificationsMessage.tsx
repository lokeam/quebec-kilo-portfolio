// Icons
import { IconMailOff } from '@/shared/components/ui/icons';

export const NoNotificationsMessage = () => (
  <div className="flex flex-col items-center justify-center h-full text-center p-4">
    <IconMailOff className="h-12 w-12 text-muted-foreground mb-4" />
    <h3 className="text-lg font-semibold mb-2">You're all caught up!</h3>
    <p className="text-sm text-muted-foreground mb-4">You have no new notifications at this time.</p>
    <p className="text-sm text-muted-foreground">
      This is where you'll see notifications about wishlist deals, generated reports and more.
    </p>
  </div>
);

NoNotificationsMessage.displayName = 'NoNotificationsMessage';
