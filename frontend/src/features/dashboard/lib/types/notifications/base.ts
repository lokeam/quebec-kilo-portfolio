import type { NotificationIcon } from '@/features/dashboard/lib/types/notifications/constants';

/**
 * Base interface for all notification types in the application.
 * This provides the common properties that every notification must have.
 */
export interface BaseNotification {
  id: string;
  timestamp: string;
  isRead: boolean;
  icon: NotificationIcon;
  title: string;
  message?: string;
};
