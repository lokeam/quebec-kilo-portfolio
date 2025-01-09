import { useState, useEffect } from 'react';
import { searchPageMockData } from '@/features/navigation/molecules/SearchButton/searchPageMockData';
import type { Game } from '@/features/navigation/molecules/SearchButton/searchPageMockData';

export interface UseSearchGamesResult {
  games: Game[];
  isLoading: boolean;
  error: Error | null;
}

export function useSearchGames(searchQuery: string): UseSearchGamesResult {
  const [games, setGames] = useState<Game[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!searchQuery.trim()) {
      setGames([]);
      return;
    }

    setIsLoading(true);
    setError(null);

    // Simulate network delay
    const timeoutId = setTimeout(() => {
      try {
        const filteredGames = searchPageMockData.filter(game =>
          game.name.toLowerCase().includes(searchQuery.toLowerCase())
        );
        setGames(filteredGames);
        setIsLoading(false);
      } catch (err) {
        setError(err instanceof Error ? err : new Error(
          'An error occured fetching games. Please try again later.'
        ));
        setIsLoading(false);
      }
    }, 400);

    return () => clearTimeout(timeoutId);
  }, [searchQuery]);

  return { games, isLoading, error };
}