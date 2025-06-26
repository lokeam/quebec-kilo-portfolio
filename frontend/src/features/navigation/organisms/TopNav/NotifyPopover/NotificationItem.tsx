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
import type { Notification } from '@/features/dashboard/lib/types/notifications/event-variants';
import type { NotificationIconType } from '@/features/navigation/types/navigation.types';

function transformNotificationIcon(icon: NotificationIcon): string {
  const iconMap: Record<NotificationIcon, string> = {
    [NOTIFICATION_ICONS.CHECK]: 'check',
    [NOTIFICATION_ICONS.TAG]: 'tag',
    [NOTIFICATION_ICONS.BAR_CHART]: 'barChart',
    [NOTIFICATION_ICONS.ALERT_TRIANGLE]: 'alertTriangle',
  };

  return iconMap[icon] || 'QuestionMarkIcon'; // Fallback icon for unexpected values
}

// Guards
//import { transformNotificationIcon } from '@/features/dashboard/lib/types/library/guards';
import { NOTIFICATION_ICONS, type NotificationIcon } from '@/features/dashboard/lib/types/notifications/constants';

interface NotificationItemProps {
  notification: Notification;
  onRemove: (id: string) => void;
}

export const MemoizedNotificationItem = memo(function NotifcationItem({
  notification,
  onRemove
}: NotificationItemProps) {
  console.log('notification', notification);

  const {
    timestamp,
    icon,
    title,
    id,
    message,
    isRead,
  } = notification;

  const { notifications } = useDomainMaps();
  const iconName = transformNotificationIcon(icon);
  const NoticiationIcon = notifications[iconName as NotificationIconType]
    || notifications.default;

  // Memoize click handler to prevent unnecessary re-renders
  const handleRemove = () => {
    onRemove(id);
  };

  return (
    <div
      className="flex items-start gap-4 p-4 cursor-default relative group hover:bg-white/5 transition-colors"
      // Add proper accessibility attributes
      role="listitem"
      aria-label={title}
    >
      {/* Unread indicator */}
      {!isRead && (
        <Badge
          className="absolute left-2 top-2 h-2 w-2 rounded-full bg-primary p-0"
          aria-label="Unread notification"
        />
      )}

      {/* Notification icon */}
      <Avatar className="h-9 w-9 bg-green-500/20">
        <AvatarFallback className="bg-transparent">
          <NoticiationIcon className="h-5 w-5 text-green-500"/>
        </AvatarFallback>
      </Avatar>

      {/* Notification content */}
      <div className="flex-1 space-y-3">
        <p className="text-sm font-medium leading-none">
          {title}
        </p>
        <p className="text-sm text-muted-foreground">
          {message && ` ${message}`}
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
