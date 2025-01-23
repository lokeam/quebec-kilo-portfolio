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

interface NotificationTimelineItemProps {
  notification: Notification;
  onRemove: (timestamp: string) => void;
}

/**
 * Timeline-styled notification item for the notifications page
 * Shares core notification logic but with specialized timeline UI
 */
export const NotificationTimelineItem = memo(function NotificationTimelineItem({
  notification,
  onRemove
}: NotificationTimelineItemProps) {
  const {
    timestamp,
    notificationTitle,
    notificationIcon,
    notificationHd,
    notificationMsg,
  } = notification;

  const { notifications } = useDomainMaps();
  const NotificationIcon = notifications[notificationIcon as NotificationIconType]
    || notifications.default;

  return (
    <div
      data-notification-card
      className="group relative rounded-md p-3 transition-colors hover:bg-white/5 min-h-[70px]"
      role="listitem"
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
          <h3 className="text-lg font-medium text-white">{notificationHd}</h3>
          {notificationMsg && (
            <p className="mt-1 text-sm text-gray-400">{notificationMsg}</p>
          )}

          {/* Specialized Content Renderers */}
          {notificationTitle === "Monthly Spend report generated" && (
            <div className="mt-2 flex flex-col text-left gap-2">
              <p className="mt-1 text-sm text-gray-400">{notificationHd}</p>
              <span className="inline-flex items-center rounded bg-[#474649d9] px-2 py-1 text-xs text-slate-300">
                <PdfIcon size={18} className="mr-1" /> monthlySpendDec2024.pdf
              </span>
            </div>
          )}

          {notificationTitle === "App updates" && (
            <div className="mt-2 flex flex-col gap-2">
              <p className="text-sm text-gray-400">{notificationHd}</p>
              <button className="text-left mt-2 text-sm text-indigo-400 hover:text-indigo-300">
                Learn more →
              </button>
            </div>
          )}
        </div>

        {/* Timestamp */}
        <span className="text-xs text-gray-500">
          {formatTimestamp(timestamp)}
        </span>

        {/* Remove Button */}
        <button
          onClick={() => onRemove(timestamp)}
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
