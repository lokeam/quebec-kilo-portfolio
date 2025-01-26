
import { useState, useRef, memo } from 'react';
import { Link } from 'react-router-dom';

// ShadCN UI Components
import { Badge } from '@/shared/components/ui/badge';
import { Button } from '@/shared/components/ui/button';
import { Popover, PopoverContent, PopoverTrigger } from '@/shared/components/ui/popover';
import { ScrollArea } from '@/shared/components/ui/scroll-area';

// Components
import { MemoizedNotificationItem } from '@/features/navigation/organisms/TopNav/NotifyPopover/NotificationItem';
import { NoNotificationsMessage } from '@/features/dashboard/components/molecules/NoNotificationsMessage';

// Hooks + Utils
import { useNotifications } from '@/features/dashboard/lib/hooks/useNotifications';

// Types
import type { Notification } from '@/features/dashboard/lib/types/notifications/event-variants';

// Icons
import { Bell } from 'lucide-react';
import { IconMailSpark, IconMailOpened } from '@tabler/icons-react';

// Constants
import { POPOVER_DIMENSIONS } from '@/features/navigation/constants/navigation.constants';

// Mock Data
import { notificationsMockData } from '@/features/dashboard/components/organisms/NotificationsPage/notificationsPage.mockdata';

export interface NotificationState {
  notifications: Notification[];
};


export const NotifyPopover = memo(function NotifyPopover() {
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [hasOpenedPopover, setHasOpenedPopover] = useState<boolean>(false);
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

  const handleRemoveNotification = (id: string) => {
    const isEmpty = removeNotification(id);
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
            <div className="flex items-center gap-2">
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
              <span className="text-sm">Mark all {allRead ? 'unread' : 'read'}</span>
            </div>
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
