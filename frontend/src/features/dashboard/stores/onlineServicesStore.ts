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
  billingCycleFilters: string[];
  paymentMethodFilters: string[];
  setViewMode: (mode: ViewMode) => void;
  setSearchQuery: (query: string) => void;
  setBillingCycleFilters: (filters: string[]) => void;
  setPaymentMethodFilters: (filters: string[]) => void;
}

export const useOnlineServicesStore = create<OnlineServicesState>((set) => ({
  viewMode: ViewModes.GRID,
  searchQuery: '',
  billingCycleFilters: [],
  paymentMethodFilters: [],
  setViewMode: (mode) => set({ viewMode: mode }),
  setSearchQuery: (query) => set({ searchQuery: query }),
  setBillingCycleFilters: (filters) => set({ billingCycleFilters: filters }),
  setPaymentMethodFilters: (filters) => set({ paymentMethodFilters: filters }),
}));

// Add a selector hook for better performance
export const useOnlineServicesSearch = () => useOnlineServicesStore((state) => state.searchQuery);
export const useOnlineServicesBillingFilters = () => useOnlineServicesStore((state) => state.billingCycleFilters);
export const useOnlineServicesPaymentFilters = () => useOnlineServicesStore((state) => state.paymentMethodFilters);
