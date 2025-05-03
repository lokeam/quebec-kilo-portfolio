import type { WishlistItem } from '@/features/dashboard/lib/types/wishlist/base';

export interface SearchResponse {
  games: WishlistItem[];
  total: number;
}