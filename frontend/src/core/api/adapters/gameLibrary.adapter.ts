import type { AddToLibraryFormPayload } from '@/features/dashboard/components/organisms/GameSearchAndSelectDialog/AddGameToLibraryForm/AddGameToLibraryForm';
import type {
  CreateLibraryGameRequest,
  LibraryGameItemResponse,
  LibraryItemsBFFResponse,
  LibraryGameItemRefactoredResponse,
  LibraryItemsBFFRefactoredResponse,
} from '@/types/domain/library-types';

export const adaptAddToLibraryFromToRequest = (
  formPayload: AddToLibraryFormPayload
): CreateLibraryGameRequest => {

  // Since types are identical, this is a passthrough
  return formPayload;
};

// Type guard focusing on critical properties
const isValidLibraryItem = (item: unknown): item is LibraryGameItemResponse => {
  if (!item || typeof item !== 'object' || item === null) {
    return false;
  }

  const obj = item as Record<string, unknown>;

  return (
    'id' in obj &&
    'name' in obj &&
    'gamesByPlatformAndLocation' in obj &&
    typeof obj.id === 'number' &&
    typeof obj.name === 'string' &&
    Array.isArray(obj.gamesByPlatformAndLocation)
  );
};

// Simple adapter with basic error handling
const adaptLibraryBFFResponse = (response: LibraryItemsBFFResponse | undefined) => {
  return {
    libraryItems: response?.libraryItems?.filter(isValidLibraryItem) ?? [],
    recentlyAdded: response?.recentlyAdded?.filter(isValidLibraryItem) ?? []
  };
};

// REFACTORED RESPONSE - updated adapters
const adaptLibraryBFFRefactoredResponse = (response: LibraryItemsBFFRefactoredResponse | undefined) => {
  if (!response) {
    return {
      libraryItems: [],
      recentlyAdded: []
    };
  }

  return {
    libraryItems: response.libraryItems?.filter(isValidLibraryItemRefactored) ?? [],
    recentlyAdded: response.recentlyAdded?.filter(isValidLibraryItemRefactored) ?? []
  };
};

// Add validation function for refactored items
const isValidLibraryItemRefactored = (item: LibraryGameItemRefactoredResponse): boolean => {
  return !!(
    item &&
    typeof item.id === 'number' &&
    typeof item.name === 'string' &&
    typeof item.coverUrl === 'string' &&
    Array.isArray(item.physicalLocations) &&
    Array.isArray(item.digitalLocations)
  );
};


// Export both functions
export {
  isValidLibraryItem,
  adaptLibraryBFFResponse,
  adaptLibraryBFFRefactoredResponse
};
