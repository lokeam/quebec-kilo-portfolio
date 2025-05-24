import type { AddToLibraryFormPayload } from '@/features/dashboard/components/organisms/GameSearchAndSelectDialog/AddGameToLibraryForm/AddGameToLibraryForm';
import type { CreateLibraryGameRequest } from '@/types/domain/library-types';

export const adaptAddToLibraryFromToRequest = (
  formPayload: AddToLibraryFormPayload
): CreateLibraryGameRequest => {
  // Log the incoming payload
  console.log('🔍 DEBUG: gameLibrary.adapter - Incoming form payload:', {
    gameId: formPayload.gameId,
    gameName: formPayload.gameName,
    gameCoverUrl: formPayload.gameCoverUrl,
    gameFirstReleaseDate: formPayload.gameFirstReleaseDate,
    gameType: formPayload.gameType,
    gameThemeNames: formPayload.gameThemeNames,
    gamesByPlatformAndLocation: formPayload.gamesByPlatformAndLocation,
    gameRating: formPayload.gameRating,
  });

  // Since types are identical, this is a passthrough
  const result = formPayload;

  // Log the outgoing result
  console.log('🔍 DEBUG: gameLibrary.adapter - Outgoing request:', {
    gameId: result.gameId,
    gameName: result.gameName,
    gameCoverUrl: result.gameCoverUrl,
    gameFirstReleaseDate: result.gameFirstReleaseDate,
    gameType: result.gameType,
    gameThemeNames: result.gameThemeNames,
    gamesByPlatformAndLocation: result.gamesByPlatformAndLocation,
    gameRating: result.gameRating,
  });

  return result;
};
