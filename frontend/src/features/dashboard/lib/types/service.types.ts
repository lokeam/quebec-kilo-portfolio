
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
export interface MediaStorageService extends BaseService {
  name: string;
  label: string;
  physicalLocations: PhysicalMediaStorageLocation[];
  digitalLocations: DigitalMediaStorageLocation[];
  totalItems: number;
  totalPhysicalItems: number;
  totalDigitalItems: number;
}

export interface PhysicalMediaStorageLocation {
  name: string;
  label: string;
  notes: string;
  storedItems: number;
  items: MediaStorageItem[];
};

export interface DigitalMediaStorageLocation extends PhysicalMediaStorageLocation {
  url: string;
  isActive: boolean;
  isFree: boolean;
  monthlyFee?: string;
}

export interface MediaStorageItem {
  name: string;
  label: string;
  createdAt?: string;
  lastUpdatedAt?: string;
}

export interface GameItem extends MediaStorageItem {
  platform: string;
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
  isDigital: boolean;
  isOwned: boolean;
  isWishlisted: boolean;
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
export type Service = OnlineService | MediaStorageService;