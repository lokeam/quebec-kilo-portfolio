// Tanstack query
import { useMutation, useQueryClient } from '@tanstack/react-query';

// base query hook
import { useAPIQuery } from '@/core/api/queries/useAPIQuery';

// query hooks
import {
  getAllLibraryGames,
  getLibraryGameById,
  createLibraryGame,
  updateLibraryGame,
  deleteLibraryGame,
} from '@/core/api/services/gameLibrary.service';

// adapters
import { adaptAddToLibraryFromToRequest } from '@/core/api/adapters/gameLibrary.adapter.ts';

// type
import type { AddToLibraryFormPayload } from '@/features/dashboard/components/organisms/GameSearchAndSelectDialog/AddGameToLibraryForm/AddGameToLibraryForm';
import type { Game } from '@/types/game';
import type { CreateLibraryGameRequest, LibraryGameItem } from '@/types/domain/library-types';


export const gameLibraryKeys = {
  all: ['library-games'] as const,
  lists: () => [...gameLibraryKeys.all, 'list'] as const,
  list: (filters: string) => [...gameLibraryKeys.lists(), { filters }] as const,
  details: () => [...gameLibraryKeys.all, 'detail'] as const,
  detail: (id: string) => [...gameLibraryKeys.details(), id] as const,
};


/**
 * Hook to fetch all games in the library
 */
export const useGetAllLibraryGames = () => {
  return useAPIQuery<LibraryGameItem[]>({
    queryKey: gameLibraryKeys.lists(),
    queryFn: async () => {
      const games = await getAllLibraryGames();
      return games;
    },
  });
};

/*
 * Hook to fetch a single game from the library
*/
export const useGetSingleGame = (id: string) => {
  return useAPIQuery<Game>({
    queryKey: gameLibraryKeys.detail(id),
    queryFn: async () => {
      const game = await getLibraryGameById(id);
      return game;
    }
  })
}

/*
 * Hook to create a game in the library
*/
export const useCreateLibraryGame = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: AddToLibraryFormPayload) => {
      // Transform form payload to API request
      const apiRequest = adaptAddToLibraryFromToRequest(data);
      return createLibraryGame(apiRequest);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: gameLibraryKeys.lists() });
    },
  });
};

/*
 * Hook to update a game in the library
*/
export const useUpdateLibraryGame = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateLibraryGameRequest> }) =>
      updateLibraryGame(id, data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: gameLibraryKeys.lists() });
      queryClient.invalidateQueries({ queryKey: gameLibraryKeys.detail(String(data.id)) });
    },
  })
};

/*
 * Hook to delete a game from the library
*/
export const useDeleteLibraryGame = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteLibraryGame,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: gameLibraryKeys.lists() });
      queryClient.invalidateQueries({ queryKey: gameLibraryKeys.detail(id) });
    },
  });
};
