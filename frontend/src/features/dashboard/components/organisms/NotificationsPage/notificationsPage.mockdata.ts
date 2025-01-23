import type { Notification } from '@/features/dashboard/lib/types/service.types';

export const notificationsMockData: Notification[] = [
  {
    notificationTitle: "App updates",
    notificationIcon: "check",
    timestamp: "2025-01-22T06:28:00Z",
    notificationHd: "New features added in Quebec Kilo version 1.2!",
    notificationMsg: "We added new features to the Wishlist tracking and Monthly Spending pages! Learn more",
    isRead: false,
  },
  {
    notificationTitle: "Monthly Spend report generated",
    notificationIcon: "iconchart",
    timestamp: "2025-01-19T08:01:00Z",
    notificationHd: "December 2024 Spend Tracking report successfully generated",
    isRead: false,
  },
  {
    notificationTitle: "Wishlist Sale",
    notificationIcon: "tag",
    timestamp: "2025-01-05T08:30:00Z",
    notificationHd: "Shardpunk",
    notificationMsg: "Is on sale for $8.39",
    isRead: true,
  },
]