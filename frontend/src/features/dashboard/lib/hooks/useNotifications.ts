import { useState } from 'react';
import type { Notification } from '@/features/dashboard/lib/types/service.types';

/**
 * Interface for the return value of useNotifications hook
 */
interface UseNotificationsReturn {
  /** Array of current notifications */
  notifications: Notification[];
  /** Boolean indicating if all notifications are marked as read */
  allRead: boolean;
  /** Count of unread notifications */
  unreadCount: number;
  /** Function to toggle read status of all notifications */
  markAllAsRead: () => void;
  /** Function to remove a specific notification */
  removeNotification: (timestamp: string) => boolean;
}

/**
 * Custom hook for managing notification state and operations
 *
 * @param initialNotifications - Initial array of notifications to populate the state
 * @returns Object containing notification state and management functions
 *
 * @example
 * ```tsx
 * const {
 *   notifications,
 *   allRead,
 *   unreadCount,
 *   markAllAsRead,
 *   removeNotification
 * } = useNotifications(initialNotifications);
 * ```
 */
export function useNotifications(
  initialNotifications: Notification[]
): UseNotificationsReturn {
  const [notifications, setNotifications] = useState<Notification[]>(initialNotifications);
  const [allRead, setAllRead] = useState(() => notifications.every((n) => n.isRead));

  // Memoized count of unread notifications
  const unreadCount = notifications.filter((n) => !n.isRead).length;

  /**
   * Toggles the read status of all notifications
   */
  const markAllAsRead = () => {
    const newAllRead = !allRead;
    setAllRead(newAllRead);
    setNotifications(notifications.map((n) => ({ ...n, isRead: newAllRead })));
  };

    /**
   * Removes a specific notification by its timestamp
   *
   * @param timestamp - Unique timestamp identifier of the notification to remove
   * @returns boolean indicating if the notifications list is now empty
   */
  const removeNotification = (timestamp: string) => {
    const newNotifications = notifications.filter((n) => n.timestamp !== timestamp);
    setNotifications(newNotifications);
    setAllRead(newNotifications.every((n) => n.isRead));
    return newNotifications.length === 0;
  };

  return {
    notifications,
    allRead,
    unreadCount,
    markAllAsRead,
    removeNotification,
  };
}
