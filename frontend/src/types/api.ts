import type { Game } from "@/features/dashboard/lib/types/game.types";

export interface SearchResponse {
  games: Game[];
  total: number;
}