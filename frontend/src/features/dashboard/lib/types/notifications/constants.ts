/**
 * Main categories of notifications in the application
 * @readonly
 */
export const NOTIFICATION_CATEGORIES = {
  /** System-level updates and changes */
  APP_UPDATE: 'appUpdate',
  /** Generated reports and analytics */
  REPORT: 'report',
  /** Updates related to wishlist items */
  WISHLIST: 'wishlist',
  /** Important system alerts and warnings */
  ALERT: 'alert',
} as const;

/** Union type of all possible notification categories */
export type NotificationCategory = typeof NOTIFICATION_CATEGORIES[keyof typeof NOTIFICATION_CATEGORIES];

/**
 * Types of reports available in the system
 * @readonly
 */
export const REPORT_CATEGORIES = {
  MONTHLY: 'monthly',
  ANNUAL: 'annual',
} as const;

export type ReportCategory = typeof REPORT_CATEGORIES[keyof typeof REPORT_CATEGORIES];

/**
 * Available icons for notifications
 * @readonly
 */
export const NOTIFICATION_ICONS = {
  CHECK: 'check',
  TAG: 'tag',
  BAR_CHART: 'barChart',
  ALERT_TRIANGLE: 'alertTriangle',
} as const;

/** Union type of all possible notification icons */
export type NotificationIcon = typeof NOTIFICATION_ICONS[keyof typeof NOTIFICATION_ICONS];
