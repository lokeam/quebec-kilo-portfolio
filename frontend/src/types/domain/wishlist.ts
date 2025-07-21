
export type WishlistItemPrice = {
  original: number;
  discounted: number;
  discountPercentage: number;
  vendor: string;
}

export type WishlistItemPlatformCompatibility = {
  hasMacOSVersion?: boolean;
  hasAndroidIOSVersion?: boolean;
  hasIOSVersion?: boolean;
}

export type WishlistItem = {
  id: string;
  name: string;
  coverUrl: string;
  themeNames: string[];
  releaseDate: string;
  currentPrice: WishlistItemPrice;
  platformCompatibility?: WishlistItemPlatformCompatibility;

}

export type WishlistItemBFFResponse = {
  pc: WishlistItem[];
  console: WishlistItem[];
  mobile: WishlistItem[];
}