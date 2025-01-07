import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export const ViewModes = {
  GRID: 'grid',
  LIST: 'list',
  TABLE: 'table'
} as const;

export type ViewMode = typeof ViewModes[keyof typeof ViewModes];

interface OnlineServicesState {
  viewMode: ViewMode;
  searchQuery: string;
  setViewMode: (mode: ViewMode) => void;
  setSearchQuery: (query: string) => void;
}

export const useOnlineServicesStore = create<OnlineServicesState>((set) => ({
  viewMode: ViewModes.GRID,
  searchQuery: '',
  setViewMode: (mode) => set({ viewMode: mode }),
  setSearchQuery: (query) => set({ searchQuery: query }),
}));

// Add a selector hook for better performance
export const useOnlineServicesSearch = () => useOnlineServicesStore((state) => state.searchQuery);
