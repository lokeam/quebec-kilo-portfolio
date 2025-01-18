import { create } from 'zustand';
import { type LibraryItem } from '@/features/dashboard/lib/types/page.types';

export const ViewModes = {
  GRID: 'grid',
  LIST: 'list',
} as const;

export type ViewMode = typeof ViewModes[keyof typeof ViewModes];

interface LibraryState {
  userGames: LibraryItem[];
  viewMode: ViewMode;
  setGames: (games: LibraryItem[]) => void;
  setViewMode: (mode: ViewMode) => void;
  searchQuery: string;
  setSearchQuery: (query: string) => void;
}

export const useLibraryStore = create<LibraryState>((set) => ({
  userGames: [],
  viewMode: ViewModes.GRID,
  searchQuery: '',
  setGames: (games) => set({ userGames: games }),
  setViewMode: (mode) =>{ set({ viewMode: mode })},
  setSearchQuery: (query) => set({ searchQuery: query }),
}));

// Selector hooks
export const useLibraryGames = () => useLibraryStore((state) => state.userGames);
export const useLibraryViewMode = () => useLibraryStore((state) => state.viewMode);
export const useLibrarySetGames = () => useLibraryStore((state) => state.setGames);
export const useLibrarySetViewMode = () => useLibraryStore((state) => state.setViewMode);
export const useLibrarySearchQuery = () => useLibraryStore((state) => state.setViewMode);
