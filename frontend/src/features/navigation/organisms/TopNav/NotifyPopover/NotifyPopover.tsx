
import { useState, useRef, memo } from 'react';
import { Link } from 'react-router-dom';

// ShadCN UI Components
import { Badge } from '@/shared/components/ui/badge';
import { Button } from '@/shared/components/ui/button';
import { Popover, PopoverContent, PopoverTrigger } from '@/shared/components/ui/popover';
import { ScrollArea } from '@/shared/components/ui/scroll-area';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/shared/components/ui/tooltip';

// Components
import { MemoizedNotificationItem } from '@/features/navigation/organisms/TopNav/NotifyPopover/NotificationItem';

// Hooks + Utils
import { useNotifications } from '@/features/dashboard/lib/hooks/useNotifications';

// Types
import type { Notification } from '@/features/dashboard/lib/types/service.types';

// Icons
import { Bell } from 'lucide-react';
import { IconMailSpark, IconMailOpened, IconMailOff } from '@tabler/icons-react';

// Constants
import { POPOVER_DIMENSIONS } from '@/features/navigation/constants/navigation.constants';

// Mock Data
import { notificationsMockData } from '@/features/dashboard/components/organisms/NotificationsPage/notificationsPage.mockdata';

export interface NotificationState {
  notifications: Notification[];
};

const NoNotificationsMessage = () => (
  <div className="flex flex-col items-center justify-center h-full text-center p-4">
    <IconMailOff className="h-12 w-12 text-muted-foreground mb-4" />
    <h3 className="text-lg font-semibold mb-2">You're all caught up!</h3>
    <p className="text-sm text-muted-foreground mb-4">You have no new notifications at this time.</p>
    <p className="text-sm text-muted-foreground">
      This is where you'll see notifications about wishlist deals, generated reports and more.
    </p>
  </div>
);


export const NotifyPopover = memo(function NotifyPopover() {
  const [isOpen, setIsOpen] = useState(false);
  const [hasOpenedPopover, setHasOpenedPopover] = useState(false);
  const popoverRef = useRef<HTMLDivElement>(null);

  const {
    notifications,
    allRead,
    unreadCount,
    markAllAsRead,
    removeNotification
  } = useNotifications(notificationsMockData);

  const handleOpenChange = (open: boolean) => {
    setIsOpen(open);
    if (open) {
      setHasOpenedPopover(true);
    }
  };

  const handleRemoveNotification = (timestamp: string) => {
    const isEmpty = removeNotification(timestamp);
    if (isEmpty) {
      setIsOpen(false);
      setHasOpenedPopover(false);
    }
  };

  return (
    <Popover open={isOpen} onOpenChange={handleOpenChange}>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="relative"
          aria-label="Open notifications"
        >
          <Bell className="h-5 w-5" />
          {unreadCount > 0 && (
            <Badge
              className="absolute right-3 top-2 h-2 w-2 rounded-full bg-destructive p-0"
              aria-label={`${unreadCount} unread notifications`}
            />
          )}
        </Button>
      </PopoverTrigger>

      <PopoverContent
        align="end"
        className={`${POPOVER_DIMENSIONS.WIDTH} p-0`}
        ref={popoverRef}
      >
        <div className="flex flex-col h-[430px]">
          {/* Header */}
          <div className="flex items-center justify-between px-4 py-2 border-b">
            <h2 className="text-sm font-semibold">Notifications</h2>
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8"
                    onClick={markAllAsRead}
                    aria-label={`Mark all as ${allRead ? 'unread' : 'read'}`}
                  >
                    {allRead ? (
                      <IconMailOpened className="h-4 w-4" />
                    ) : (
                      <IconMailSpark className="h-4 w-4" />
                    )}
                  </Button>
                </TooltipTrigger>
                <TooltipContent side="left">
                  <p>Mark all as {allRead ? 'unread' : 'read'}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </div>

          {/* Notification List */}
          <ScrollArea className="flex-grow">
            {notifications.length > 0 || !hasOpenedPopover ? (
              notifications.map((notification) => (
                <MemoizedNotificationItem
                  key={notification.timestamp}
                  notification={notification}
                  onRemove={handleRemoveNotification}
                />
              ))
            ) : (
              <NoNotificationsMessage />
            )}
          </ScrollArea>

          {/* Footer */}
          <div className="p-2 border-t mt-auto">
            {notifications.length > 0 && hasOpenedPopover && (
              <Link to="/notifications">
                <Button
                  variant="secondary"
                  className="w-full"
                  onClick={() => setIsOpen(false)}
                >
                  View All Notifications
                </Button>
              </Link>
            )}
          </div>
        </div>
      </PopoverContent>
    </Popover>
  );
});

// Note: Experimenting adding with display name for debugging purposes
NotifyPopover.displayName = 'NotifyPopover';
