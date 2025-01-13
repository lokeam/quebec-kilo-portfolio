import { create } from 'zustand';
import type { OnlineService } from '@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata';

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
  services: OnlineService[];
  setServices: (services: OnlineService[]) => void;
  toggleActiveOnlineService: (serviceName: string, isActive: boolean) => void;
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
  services: [],
  setServices: (services) => set({ services }),
  toggleActiveOnlineService: (serviceName, isActive) =>
    set((state) => ({
      services: state.services.map((service) =>
        service.name === serviceName
          ? { ...service, isActive }
          : service
      ),
    })),
}));

// Add a selector hook for better performance
export const useOnlineServicesSearch = () => useOnlineServicesStore((state) => state.searchQuery);
export const useOnlineServicesBillingFilters = () => useOnlineServicesStore((state) => state.billingCycleFilters);
export const useOnlineServicesPaymentFilters = () => useOnlineServicesStore((state) => state.paymentMethodFilters);
export const useOnlineServicesToggleActive = () => useOnlineServicesStore((state) => state.toggleActiveOnlineService);
export const useOnlineServices = () => useOnlineServicesStore((state) => state.services);
export const useSetOnlineServices = () => useOnlineServicesStore((state) => state.setServices);
export const useOnlineServicesIsActive = (serviceName: string) =>
  useOnlineServicesStore((state) => state.services.find((service) => service.name === serviceName)?.isActive ?? false);
