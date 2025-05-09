import type { WishlistItem } from '@/features/dashboard/lib/types/wishlist/base';

export interface SearchResponse {
  success: boolean;
  data: {
    games: WishlistItem[];
    total: number;
  };
  metadata: {
    timestamp: string;
    request_id: string;
  };
}