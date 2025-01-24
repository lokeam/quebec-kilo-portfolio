import { memo } from 'react';

// Shadcn UI Components
import { Avatar, AvatarFallback } from '@/shared/components/ui/avatar';

// Icons
import { Trash2 } from 'lucide-react';
import { PdfIcon } from '@/shared/components/ui/LogoMap/misc/pdfFile';

// Hooks + Utils
import { formatTimestamp } from '@/features/navigation/utils/formatTimestamp';
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';

// Types
import type { Notification } from '@/features/dashboard/lib/types/service.types';
import type { NotificationIconType } from '@/features/navigation/types/navigation.types';



interface ReportNotificationProps {
  report: {
    type: 'monthly' | 'annual';
    period: {
      month?: string;
      year: string;
    };
    downloadUrl: string;
    fileSize?: string;
  };
}

interface AppUpdateNotificationProps {
  update: {
    version?: string;
    infoUrl: string;
    changes: string[];
  };
};

interface WishlistNotificationProps {
  item: {
    name: string;
    salePrice: string;
    originalPrice: string;
    discountPercentage: number;
    coverUrl?: string;
  };
}

// Specialized Content Components
const ReportNotification = memo(function ReportNotification({
  report
}: ReportNotificationProps) {
  return (
    <div className="mt-2 flex flex-col text-left gap-2">
      <span className="inline-flex items-center rounded bg-[#474649d9] px-2 py-1 text-xs text-slate-300">
        <PdfIcon size={18} className="mr-1" />
        {report.downloadUrl.split('/').pop()}
      </span>
      <p className="mt-1 text-sm text-gray-400">Size: {report.fileSize}</p>

    </div>
  );
});

const AppUpdateNotification = memo(function AppUpdateNotification({
  update
}: AppUpdateNotificationProps) {
  return (
    <div className="mt-2 flex flex-col gap-2">
      <div className="text-sm text-gray-400">
        {update.changes.map((change, index) => (
          <p key={index}>• {change}</p>
        ))}
      </div>
      <a
        href={update.infoUrl}
        className="text-left mt-2 text-sm text-indigo-400 hover:text-indigo-300"
      >
        Learn more →
      </a>
    </div>
  );
});

const WishlistNotification = memo(function WishlistNotification({
  item
}: WishlistNotificationProps) {
  return (
    <div className="mt-2 flex items-center gap-4">
      {item.coverUrl && (
        <img
          src={item.coverUrl}
          alt={item.name}
          className="h-16 w-12 rounded object-cover"
        />
      )}
      <div className="flex flex-col gap-1">
        <p className="text-sm text-gray-400">
          ${item.salePrice} <span className="text-xs line-through">${item.originalPrice}</span>
        </p>
        <p className="text-xs text-green-500">Save {item.discountPercentage}%</p>
      </div>
    </div>
  );
});

interface NotificationTimelineItemProps {
  notification: Notification;
  onRemove: (id: string) => void;
}

const renderNotificationContent = (notification: Notification) => {
  switch (notification.type) {
    case 'report':
      return <ReportNotification report={notification.report}/>;
    case 'appUpdate':
      return <AppUpdateNotification update={notification.update} />;
    case 'wishlist':
      return <WishlistNotification item={notification.item} />;
    default:
      return null;
  }
};

/**
 * Timeline-styled notification item for the notifications page
 * Shares core notification logic but with specialized timeline UI
 */
export const NotificationTimelineItem = memo(function NotificationTimelineItem({
  notification,
  onRemove
}: NotificationTimelineItemProps) {
  const { notifications } = useDomainMaps();
  const NotificationIcon = notifications[notification.icon as NotificationIconType]
    || notifications.default;

  return (
    <div
      data-testid={`notification-${notification.type}`}
      data-notification-card
      className="group relative rounded-md p-3 transition-colors hover:bg-white/5 min-h-[70px]"
      role="listitem"
      aria-label={`${notification.title} notification`}
    >
      <div className="flex items-start gap-3">
        {/* Timeline Separator */}
        <div className="flex flex-col items-center relative h-full min-h-[116px]">
          {/* Timeline Dot */}
          <Avatar className="h-8 w-8 bg-green-500/20">
            <AvatarFallback className="bg-transparent">
              <NotificationIcon className="h-4 w-4 text-green-500" />
            </AvatarFallback>
          </Avatar>
          {/* Timeline Vertical Line */}
          <div className="w-[1px] absolute top-12 bottom-0 bg-white/30" />
        </div>

        {/* Content */}
        <div className="flex-1">
          <h3 className="text-lg font-medium text-white">{notification.title}</h3>
          {notification.message && (
            <p className="mt-1 text-sm text-gray-400">{notification.message}</p>
          )}
          {renderNotificationContent(notification)}
        </div>

        {/* Timestamp */}
        <span className="text-xs text-gray-500">
          {formatTimestamp(notification.timestamp)}
        </span>

        {/* Remove Button */}
        <button
          onClick={() => onRemove(notification.id)}
          className="absolute right-2 bottom-2 opacity-0 transition-opacity group-hover:opacity-100"
          aria-label="Remove notification"
        >
          <Trash2 className="h-5 w-5 text-gray-500 hover:text-gray-400" />
        </button>
      </div>
    </div>
  );
});

NotificationTimelineItem.displayName = 'NotificationTimelineItem';
