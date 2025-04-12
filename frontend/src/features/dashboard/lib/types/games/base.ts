import type { PlatformCategory } from '@/shared/types/platform';
import type { Price } from '@/shared/types/pricing';
import type { Rating } from '@/features/dashboard/lib/types/wishlist/ratings';

export interface Game {
  id: string;
  name: string;
  cover_url: string;
  theme_names: string[];
  is_in_library: boolean;
  is_in_wishlist: boolean;
  releaseDate: string;
  platform: PlatformCategory;
  price: Price;
  platformSupport?: {
    hasMacOSVersion?: boolean;
    hasAndroidVersion?: boolean;
    hasIOSVersion?: boolean;
  };
  rating?: Rating;
}