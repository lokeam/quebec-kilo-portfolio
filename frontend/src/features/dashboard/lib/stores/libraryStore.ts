import { create } from 'zustand';
import type { LibraryGameItemResponse } from '@/types/domain/library-types';
import {
  getStoredViewMode,
  featureViewModes
} from '@/shared/constants/viewModes';

type LibraryViewMode = typeof featureViewModes.library.allowed[number];

interface LibraryState {
  userGames: LibraryGameItemResponse[];
  setGames: (games: LibraryGameItemResponse[]) => void;
  viewMode: LibraryViewMode;
  setViewMode: (mode: LibraryViewMode) => void;
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  platformFilters: string[];
  setPlatformFilters: (filters: string[]) => void;
  locationFilters: string[];
  setLocationFilters: (filters: string[]) => void;
}

export const useLibraryStore = create<LibraryState>((set) => ({
  userGames: [],
  setGames: (games) => set({ userGames: games }),
  viewMode: getStoredViewMode(
    featureViewModes.library.storageKey,
    featureViewModes.library.default,
    featureViewModes.library.allowed, // Need to pass allowed modes here to differentiate between modes with or without table view
  ),
  setViewMode: (mode) =>{
    // Only allow modes that are valid for library
    if (featureViewModes.library.allowed.includes(mode)) {
      localStorage.setItem(featureViewModes.library.storageKey, mode);
      set({ viewMode: mode });
    }
  },
  searchQuery: '',
  setSearchQuery: (query) => set({ searchQuery: query }),
  platformFilters: [],
  setPlatformFilters: (filters) => set({ platformFilters: filters }),
  locationFilters: [],
  setLocationFilters: (filters) => set({ locationFilters: filters }),
}));

// Selector hooks for better performance
export const useLibraryGames = () => useLibraryStore((state) => state.userGames);
export const useLibraryViewMode = () => useLibraryStore((state) => state.viewMode);
export const useLibrarySetViewMode = () => useLibraryStore((state) => state.setViewMode);
export const useLibrarySearchQuery = () => useLibraryStore((state) => state.searchQuery);
export const useLibrarySetSearchQuery = () => useLibraryStore((state) => state.setSearchQuery);
export const useLibraryPlatformFilters = () => useLibraryStore((state) => state.platformFilters);
export const useLibrarySetPlatformFilters = () => useLibraryStore((state) => state.setPlatformFilters);
export const useLibraryLocationFilters = () => useLibraryStore((state) => state.locationFilters);
export const useLibrarySetLocationFilters = () => useLibraryStore((state) => state.setLocationFilters);
