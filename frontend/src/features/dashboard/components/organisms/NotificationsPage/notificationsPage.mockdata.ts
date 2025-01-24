import type { Notification } from '@/features/dashboard/lib/types/service.types';

export const notificationsMockData: Notification[] = [
  {
    id: "1",
    type: "appUpdate",
    timestamp: "2025-01-22T06:28:00Z",
    isRead: false,
    icon: "check",
    title: "New Features Available",
    message: "We added new features to the Wishlist tracking and Monthly Spending pages!",
    update: {
      version: "1.2",
      infoUrl: "http://locahttp://localhost:5173/blog/updates/",
      changes: [
        "Added Wishlist tracking",
        "Enhanced Monthly Spending page",
      ]
    },
  },
  {
    id: "2",
    type: "report",
    timestamp: "2025-01-19T08:01:00Z",
    isRead: false,
    icon: "iconchart",
    title: "Monthly Spend Report Ready",
    report: {
      type: "monthly",
      period: {
        month: "December",
        year: "2024",
      },
      downloadUrl: "/reports/dec-2024.pdf",
      fileSize: "2.4MB"
    }
  },
  {
    id: "3",
    type: "wishlist",
    timestamp: "2025-01-05T08:30:00Z",
    isRead: true,
    icon: "tag",
    title: "Wishlist Item on Sale",
    message: "Deal Alert!",
    item: {
      name: "Shardpunk",
      salePrice: "8.39",
      originalPrice: "29.99",
      discountPercentage: 72,
      coverUrl: "https://images.igdb.com/igdb/image/upload/t_cover_big/co7eqr.webp",
    }
  }
];
