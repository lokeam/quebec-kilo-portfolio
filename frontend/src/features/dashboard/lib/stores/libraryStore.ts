import { create } from 'zustand';
import { type LibraryItem } from '@/features/dashboard/lib/types/library/items';
import {
  featureViewModes,
  getStoredViewMode,
} from '@/shared/constants/viewModes';


export type LibraryViewMode = typeof featureViewModes.library.allowed[number];

interface LibraryState {
  platformFilter: string;
  setPlatformFilter: (filter: string) => void;
  userGames: LibraryItem[];
  setGames: (games: LibraryItem[]) => void;
  viewMode: LibraryViewMode;
  setViewMode: (mode: LibraryViewMode) => void;
  searchQuery: string;
  setSearchQuery: (query: string) => void;
}

export const useLibraryStore = create<LibraryState>((set) => ({
  platformFilter: '',
  setPlatformFilter: (platform) => set({ platformFilter: platform}),
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
}));

// Selector hooks
export const useLibraryGames = () => useLibraryStore((state) => state.userGames);
export const useLibraryViewMode = () => useLibraryStore((state) => state.viewMode);
export const useLibrarySetGames = () => useLibraryStore((state) => state.setGames);
export const useLibrarySetViewMode = () => useLibraryStore((state) => state.setViewMode);
export const useLibrarySearchQuery = () => useLibraryStore((state) => state.searchQuery);
export const useLibrarySetSearchQuery = () => useLibraryStore((state) => state.setSearchQuery);
export const useLibraryPlatformFilter = () => useLibraryStore((state) => state.platformFilter);
export const useLibrarySetPlatformFilter = () => useLibraryStore((state) => state.setPlatformFilter);
