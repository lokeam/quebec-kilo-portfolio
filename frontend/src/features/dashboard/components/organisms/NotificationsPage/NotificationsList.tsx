import { memo } from 'react';

// Components
import { NotificationTimelineItem } from '@/features/dashboard/components/organisms/NotificationsPage/NotificationTimelineItem';
import { NoNotificationsMessage } from '@/features/dashboard/components/molecules/NoNotificationsMessage';

// Hooks + Utils
import { useNotifications } from '@/features/dashboard/lib/hooks/useNotifications';

// Types
import type { Notification } from '@/features/dashboard/lib/types/service.types';

interface NotificationListProps {
  initialNotifications: Notification[]
}

export const NotificationList = memo( function NotificationList({ initialNotifications }: NotificationListProps) {
  const { notifications, removeNotification } = useNotifications(initialNotifications);

  if (notifications.length === 0) {
    return <NoNotificationsMessage />
  }

  return (
    <div className="w-full rounded-lg bg-[#1e1b2e] p-4 shadow-lg">
      <div className="space-y-3" role="list">
        {notifications.map((notification) => (
          <NotificationTimelineItem
            key={notification.timestamp}
            notification={notification}
            onRemove={removeNotification}
          />
        ))}
      </div>
    </div>
  );
});

NotificationList.displayName = 'NotificationList';
