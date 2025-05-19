import type { AddToLibraryFormPayload } from '@/features/dashboard/components/organisms/GameSearchAndSelectDialog/AddGameToLibraryForm/AddGameToLibraryForm';
import type { CreateLibraryGameRequest } from '@/types/domain/library-types';

export const adaptAddToLibraryFromToRequest = (
  formPayload: AddToLibraryFormPayload
): CreateLibraryGameRequest => {

  // Since types are identical, this is a passthrough
  return formPayload;
};