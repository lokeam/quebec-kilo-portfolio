import { create } from 'zustand';
import { type LibraryItem } from '@/features/dashboard/lib/types/page.types';

export const ViewModes = {
  GRID: 'grid',
  LIST: 'list',
} as const;

export type ViewMode = typeof ViewModes[keyof typeof ViewModes];

interface LibraryState {
  platformFilter: string;
  setPlatformFilter: (filter: string) => void;
  userGames: LibraryItem[];
  setGames: (games: LibraryItem[]) => void;
  viewMode: ViewMode;
  setViewMode: (mode: ViewMode) => void;
  searchQuery: string;
  setSearchQuery: (query: string) => void;
}

export const useLibraryStore = create<LibraryState>((set) => ({
  platformFilter: '',
  setPlatformFilter: (platform) => set({ platformFilter: platform}),
  userGames: [],
  setGames: (games) => set({ userGames: games }),
  viewMode: ViewModes.GRID,
  setViewMode: (mode) =>{ set({ viewMode: mode })},
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
