/* Library Page */
export type LibraryItem = {
  id: string;
  image: string;
  title: string;
  favorite: boolean;
  dateAdded: string;
  diskSize?: string;
  platformVersion?: string;
  physicalLocation?: string;
  physicalLocationType?: string;
  physicalSublocation?: string;
  physicalSublocationType?: string;
  digitalLocation?: string;
};

/* Wishlist/Deals Page */
export type Platform = 'pc' | 'console' | 'mobile';

export type WishlistItem = {
  id: string;
  title: string;
  thumbnailUrl: string;
  tags: string[];
  releaseDate: string;
  rating: {
    positive: number;
    negative: number;
    totalReviews: number;
  };
  price: {
    original: number;
    discounted?: number | undefined | null;
    discountPercentage?: number | undefined | null;
    vendor: string;
  }
  platform: Platform;
  hasMacOSVersion?: boolean | undefined | null;
  hasAndroidVersion?: boolean | undefined | null;
  hasIOSVersion?: boolean | undefined | null;
};

export type WishListPageData = {
  pc: WishlistItem[];
  console: WishlistItem[];
  mobile: WishlistItem[];
};
