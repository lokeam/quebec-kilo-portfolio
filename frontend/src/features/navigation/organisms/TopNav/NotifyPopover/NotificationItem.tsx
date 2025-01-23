import { memo } from 'react';

// ShadCN UI Components
import { Avatar, AvatarFallback } from '@/shared/components/ui/avatar';
import { Badge } from '@/shared/components/ui/badge';
import { Button } from '@/shared/components/ui/button';

// HOoks + Utils
import { formatTimestamp } from '@/features/navigation/utils/formatTimestamp';
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';

// Icons
import { X } from 'lucide-react'

// Types
import type { Notification } from '@/features/dashboard/lib/types/service.types';
import type { NotificationIconType } from '@/features/navigation/types/navigation.types';

interface NotificationItemProps {
  notification: Notification;
  onRemove: (timestamp: string) => void;
}

export const MemoizedNotificationItem = memo(function NotifcationItem({
  notification,
  onRemove
}: NotificationItemProps) {
  const {
    timestamp,
    notificationIcon,
    notificationTitle,
    notificationHd,
    notificationMsg,
    isRead,
  } = notification;

  const { notifications } = useDomainMaps();
  const NoticiationIcon = notifications[notificationIcon as NotificationIconType]
    || notifications.default;

  // Memoize click handler to prevent unnecessary re-renders
  const handleRemove = () => {
    onRemove(timestamp);
  };



  return (
    <div
      className="flex items-start gap-4 p-4 cursor-default relative group hover:bg-accent transition-colors"
      // Add proper accessibility attributes
      role="listitem"
      aria-label={notificationTitle}
    >
      {/* Unread indicator */}
      {!isRead && (
        <Badge
          className="absolute left-2 top-2 h-2 w-2 rounded-full bg-primary p-0"
          aria-label="Unread notification"
        />
      )}

      {/* Notification icon */}
      <Avatar className="h-9 w-9 bg-secondary">
        <AvatarFallback>
          <NoticiationIcon />
        </AvatarFallback>
      </Avatar>

      {/* Notification content */}
      <div className="flex-1 space-y-1">
        <p className="text-sm font-medium leading-none">
          {notificationTitle}
        </p>
        <p className="text-sm text-muted-foreground">
          {notificationHd}
          {notificationMsg && ` ${notificationMsg}`}
        </p>
        <p className="text-xs text-muted-foreground">
          {formatTimestamp(timestamp)}
        </p>
      </div>

      {/* Remove button */}
      <Button
        variant="ghost"
        size="icon"
        className="absolute right-2 top-2 opacity-0 group-hover:opacity-100 transition-opacity h-6 w-6"
        onClick={handleRemove}
        aria-label="Remove notification"
      >
        <X className="h-3 w-3" />
      </Button>
    </div>
  );
});

// Note: Experimenting adding with display name for debugging purposes
MemoizedNotificationItem.displayName = 'NotificationItem';
