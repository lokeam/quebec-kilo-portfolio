// Tanstack query
import { useMutation, useQueryClient } from '@tanstack/react-query';

// Base query hook
import { useAPIQuery } from '@/core/api/queries/useAPIQuery';

// Service Layer methods
import {
  //getAllLibraryGames,
  getLibraryGameById,
  createLibraryGame,
  updateLibraryGame,
  deleteLibraryGame,
  deleteLibraryGameVersions,
  getLibraryPageBFFResponse,
} from '@/core/api/services/gameLibrary.service';

// Utils
import { showToast } from '@/shared/components/ui/TanstackMutationToast/showToast';

// Adapters
import { adaptAddToLibraryFromToRequest, adaptLibraryBFFRefactoredResponse } from '@/core/api/adapters/gameLibrary.adapter';

// Type
import type { AddToLibraryFormPayload } from '@/features/dashboard/components/organisms/GameSearchAndSelectDialog/AddGameToLibraryForm/AddGameToLibraryForm';
import type { LibraryGameItemResponse } from '@/types/domain/library-types';
import type {
   CreateLibraryGameRequest,
   LibraryItemsBFFRefactoredResponse,
} from '@/types/domain/library-types';

// Constants
import { TOAST_SUCCESS_MESSAGES } from '@/shared/constants/toast.success.messages';
import { TOAST_ERROR_MESSAGES } from '@/shared/constants/toast.error.messages';
import { TOAST_DURATIONS } from '@/shared/constants/toast.config';

export const gameLibraryKeys = {
  all: ['library-games'] as const,
  lists: () => [...gameLibraryKeys.all, 'list'] as const,
  list: (filters: string) => [...gameLibraryKeys.lists(), { filters }] as const,
  details: () => [...gameLibraryKeys.all, 'detail'] as const,
  detail: (id: string) => [...gameLibraryKeys.details(), id] as const,
};

/*
 * Hook to fetch all library items for the BFF page
*/
export const useGetLibraryPageBFFResponse = () => {
  return useAPIQuery<LibraryItemsBFFRefactoredResponse>({
    queryKey: gameLibraryKeys.lists(),
    queryFn: async () => {
      try {
        const response = await getLibraryPageBFFResponse();

        // Use our adapter to safely extract and validate data
        const adaptedResponse = adaptLibraryBFFRefactoredResponse(response);

        return adaptedResponse;
      } catch(error) {
        console.error('[DEBUG] useGetLibraryPageBFFResponse: Error fetching data:', error);
        throw error;
      }
    },
    staleTime: 0, // Consider data stale immediately
    refetchOnMount: true, // Refetch when component mounts
    refetchOnWindowFocus: true, // Refetch when window regains focus
  });
}

/**
 * Hook to fetch all games in the library
 */
// export const useGetAllLibraryGames = () => {
//   return useAPIQuery<LibraryGameItem[]>({
//     queryKey: gameLibraryKeys.lists(),
//     queryFn: async () => {
//       const games = await getAllLibraryGames();
//       return games;
//     },
//   });
// };

/*
 * Hook to fetch a single game from the library
*/
export const useGetSingleGame = (id: string) => {
  return useAPIQuery<LibraryGameItemResponse>({
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
    onSuccess: (data) => {
      console.log('🔍 DEBUG: useCreateLibraryGame onSuccess:', {
        data,
        queryKey: gameLibraryKeys.lists(),
        currentCache: queryClient.getQueryData(gameLibraryKeys.lists())
      });
      queryClient.invalidateQueries({ queryKey: gameLibraryKeys.lists() });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.GAME.ADD_TO_LIBRARY,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onMutate: async (newGame) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({ queryKey: gameLibraryKeys.lists() });

      // Snapshot the previous value
      const previousGames = queryClient.getQueryData(gameLibraryKeys.lists());

      // Optimistically update the cache
      queryClient.setQueryData(gameLibraryKeys.lists(), (old: LibraryItemsBFFRefactoredResponse | undefined) => ({
        ...old,
        libraryItems: [...(old?.libraryItems || []), newGame]
      }));

      return { previousGames };
    },

    onError: (error) => {
      console.log('❔🔎 useCreateLibraryGame query onError, error - ', error);
      showToast({
        message: TOAST_ERROR_MESSAGES.GAME.ADD_TO_LIBRARY.DEFAULT,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onSettled: () => {
      // Always refetch after error or success
      queryClient.invalidateQueries({ queryKey: gameLibraryKeys.lists() });
    }
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
 * Hook to delete a game from the library (legacy - deletes entire game)
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

/*
 * Hook to delete specific platform versions of a game from the library
*/
export const useDeleteLibraryGameVersions = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteLibraryGameVersions,
    onSuccess: (data, variables) => {
      // Invalidate the library list to refresh the data
      queryClient.invalidateQueries({ queryKey: gameLibraryKeys.lists() });

      // Invalidate the specific game detail if it exists
      queryClient.invalidateQueries({ queryKey: gameLibraryKeys.detail(String(variables.gameId)) });

      // Show success toast
      const versionCount = variables.deleteAll ? 'all versions' : `${variables.versions.length} version${variables.versions.length === 1 ? '' : 's'}`;
      showToast({
        message: `Successfully removed ${versionCount} of the game from your library`,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error) => {
      console.error('Error deleting game versions:', error);
      showToast({
        message: TOAST_ERROR_MESSAGES.GAME.REMOVE_FROM_LIBRARY?.DEFAULT || 'Failed to remove game versions from library',
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
  });
};
