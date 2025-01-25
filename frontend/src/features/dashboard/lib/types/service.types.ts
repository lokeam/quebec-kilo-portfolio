import type { PaymentMethod, SpendTransaction } from "./spend-tracking/constants";

export interface BaseService {
  name: string;
  status?: 'active' | 'inactive' | 'error';
  createdAt?: string;
  updatedAt?: string;
};

/**
 * Home Service
*/

/**
 * Wishlist Service
* /

/**
 * Notifications Service
*/
interface BaseNotification {
  id: string;
  timestamp: string;
  isRead: boolean;
  icon: string;
title: string;
  message?: string;
}

interface AppUpdateNotification extends BaseNotification {
  type: 'appUpdate';
  update: {
    version?: string;
    infoUrl: string;
    changes: string[];
  }
}

interface ReportNotification extends BaseNotification {
  type: 'report';
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

interface WishlistNotification extends BaseNotification {
  type: 'wishlist';
  item: {
    name: string;
    salePrice: string;
    originalPrice: string;
    discountPercentage: number;
    coverUrl?: string;
  };
}

export type Notification =
  | AppUpdateNotification
  | ReportNotification
  | WishlistNotification;

export type NotificationReportType = 'monthly' | 'annual';

export type NotificationType =
  | 'wishlist'
  | 'appUpdate'
  | 'report'
  | 'alert';

/**
 * Spend Tracking Service
*/

/**
 * Core date and currency types used throughout the spending system
 */
type ISO8601Date = string; // YYYY-MM-DD format
type Currency = string; // TODO: Use decimal library?

/**
 * Defines all possible billing cycles for subscription services
 * Used for calculating recurring charges and forecasting expenses
 */
type BillingCycle =
  | '1 month'
  | '3 month'
  | '6 month'
  | '1 year';


  /**
 * Media types represent the category of purchased content or item
 * This helps with reporting and analytics
 */
export const PURCHASED_MEDIA_CATEGORIES = {
  HARDWARE: 'hardware',
  DLC: 'dlc',
  IN_GAME_PURCHASE: 'inGamePurchase',
  SUBSCRIPTION: 'subscription',
  PHYSICAL: 'physical',
  DISC: 'disc'
} as const;

export type PurchasedMediaCategory = typeof PURCHASED_MEDIA_CATEGORIES[keyof typeof PURCHASED_MEDIA_CATEGORIES];

/**
 * Types of in-game purchases available in the system
 * Used for categorizing microtransactions and virtual goods
 */
export type InGamePurchaseType =
  | 'currency'
  | 'item'
  | 'cosmetic'
  | 'feature';

/**
 * Physical item conditions for inventory tracking
 * Applies to hardware and physical media
 */
export type ItemCondition =
  | 'new'
  | 'used'
  | 'refurbished';


/**
 * Base interface for all spending transactions
 * Contains common fields required for any type of purchase
 */
interface BaseSpendTracking {
  id: string;
  amount: Currency;
  title: string;
  spendTransactionType: SpendTransaction;
  paymentMethod: PaymentMethod;
  mediaType: PurchasedMediaCategory;
  createdAt: ISO8601Date;
  updatedAt: ISO8601Date;
};

// Specific Spend Tracking interfaces extending base
export interface SubscriptionSpend extends BaseSpendTracking {
  spendTransactionType: 'subscription';
  billingCycle: BillingCycle;
  nextBillingDate: ISO8601Date;
  isActive: boolean;
}

export interface OneTimeSpend extends BaseSpendTracking {
  spendTransactionType: 'one-time';
  isDigital: boolean;
  isWishlisted: boolean;
}

// Union type for spend tracking
export type SpendTrackingService = SubscriptionSpend | OneTimeSpend;


export type SpendTrackingMediaType =
  | 'hardware'
  | 'dlc'
  | 'inGamePurchase'
  | 'subscription'
  | 'physical'
  | 'disc';

// export interface SpendTrackingService extends BaseService {
//   day?: string;
//   month?: string;
//   year?: string;
//   onlineService?: string;
//   title?: string;
//   amount?: string;
//   billingCycle?: string;
//   spendTransactionType?: string; // subscription, one-time
//   paymentMethod?: string;
//   isActive?: boolean;
//   isDigital?: boolean;
//   isRecurring?: boolean;
//   billingDate?: string;
//   annualTotalSpend?: string;
//   isWishlisted?: boolean;
//   isPaid?: boolean;
//   mediaType?: SpendTrackingMediaType;
// }

export interface SpendTrackingData {
  currentTotalThisMonth: SpendTrackingService[];
  recurringNextMonth: SpendTrackingService[];
  oneTimeThisMonth: SpendTrackingService[];
  totalSpendsThisMonth: string;
  totalSpendsThisYear: string;
}

/**
 * Online Service
 */

export interface OnlineService extends BaseService {
  label: string;
  logo: string;
  tierName: string;
  billingCycle: string;
  url: string;
  monthlyFee: string;
  quarterlyFee: string;
  annualFee: string;
  renewalDay: string;
  renewalMonth: string;
  isActive: boolean;
  paymentMethod?: PaymentMethod;
  plan?: string;
};

export interface OnlineServicesData {
  services: OnlineService[];
  totalServices: number;
};


/**
 * Library Service
*/

/**
 * Media Storage
*/
export interface MediaStorageData extends BaseService {
  name: string;
  label: string;
  physicalLocations: PhysicalMediaStorageLocation[];
  digitalLocations: DigitalMediaStorageLocation[];
  totalPhysicalLocations: number;
  totalDigitalLocations: number;
  totalItems: number;
  totalPhysicalItems: number;
  totalDigitalItems: number;
}

export interface PhysicalMediaStorageLocation {
  name: string;
  label: string;
  locationType?: string;
  locationImage: string;
  mapCoordinates: string;
  subLocations: SubLocation[];
  subLocationCount: number;
  totalSublocationItems: number;
  totalStoredItems: number;
  items: MediaItem[];
};

export interface SubLocation {
  name: string;
  label: string;
  totalItems: number;
  items: MediaItem[];
}

export interface DigitalMediaStorageLocation extends PhysicalMediaStorageLocation {
  url: string;
  isActive: boolean;
  isFree: boolean;
  monthlyFee?: string;
}

export interface MediaStorageItem {
  itemName: string;
  itemLabel: string;
  createdAt?: string;
  lastUpdatedAt?: string;
}

export interface MediaStorageServicesData {
  services: MediaStorageData[];
  totalServices: number;
};

// Domain Specific Media Storage Items
export interface GameItem extends MediaStorageItem {
  itemPlatform: string;
  itemPlatformVersion: string;
  coverArtUrl?: string;
  backdropArtUrl?: string;
  artworks?: string[];
  alternativeTitles?: string[];
  story?: string;
  platformVersion?: string;
  genres?: string[];
  themes?: string[];
  rating?: string;
  playerModes?: string[];
  playerPerspectives?: string[];
  developers?: string[];
  releaseDate?: string;
  isDigital?: boolean;
  isOwned?: boolean;
  isWishlisted?: boolean;
}

export interface MovieItem extends MediaStorageItem {
  title: string;
  localizedTitles?: string[];
  posterArtUrl?: string;
  backdropArtUrl?: string;
  releaseDate?: string;
  genres?: string[];
  themes?: string[];
  rating?: string;
  isDigital: boolean;
  isOwned: boolean;
  isWishlisted: boolean;
  streamingServices?: string[];
}

/**
 * Union Type for all services
*/
export type Service = OnlineService | MediaStorageData;
export type MediaItem = GameItem | MovieItem;