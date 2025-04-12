import type { Game } from '@/features/dashboard/lib/types/games/base';

export interface SearchResponse {
  games: Game[];
  total: number;
}