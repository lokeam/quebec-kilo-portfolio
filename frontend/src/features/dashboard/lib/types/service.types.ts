
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
 * Spend Tracking Service
*/
export type SpendTrackingMediaType =
  | 'hardware'
  | 'dlc'
  | 'inGamePurchase'
  | 'subscription'
  | 'physical'
  | 'disc';

export interface SpendTrackingService extends BaseService {
  day?: string;
  month?: string;
  year?: string;
  onlineService?: string;
  title?: string;
  amount?: string;
  billingCycle?: string;
  spendType?: string; // subscription, one-time
  paymentMethod?: string;
  isActive?: boolean;
  isDigital?: boolean;
  isRecurring?: boolean;
  billingDate?: string;
  annualTotalSpend?: string;
  isWishlisted?: boolean;
  isPaid?: boolean;
  mediaType?: SpendTrackingMediaType;
}

export interface SpendTrackingData {
  recurringThisMonth: SpendTrackingService[];
  recurringNextMonth: SpendTrackingService[];
  oneTimeThisMonth: SpendTrackingService[];
  totalSpendsThisMonth: string;
  totalSpendsThisYear: string;
  top5SpendsAll: SpendTrackingService[];
  top5SpendsDigital: SpendTrackingService[];
  top5SpendsPhysical: SpendTrackingService[];
}

/**
 * Online Service
 */
export type PaymentMethodType = "Alipay" | "Amex" | "Code" | "CodeFront" | "Diners" | "Discover" |
"Elo" | "Generic" | "Hiper" | "Hipercard" | "Jcb" | "Maestro" | "Mastercard" |
"Mir" | "Paypal" | "Unionpay" | "Visa";

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
  paymentMethod?: PaymentMethodType;
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