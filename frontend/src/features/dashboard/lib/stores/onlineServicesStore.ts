import { create } from 'zustand';
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import {
  getStoredViewMode,
  featureViewModes
} from '@/shared/constants/viewModes';
import { SERVICE_STATUS_CODES } from '@/shared/constants/service.constants';

// export const ViewModes = {
//   GRID: 'grid',
//   LIST: 'list',
//   TABLE: 'table'
// } as const;

// export type ViewMode = typeof ViewModes[keyof typeof ViewModes];

type OnlineServicesViewMode = typeof featureViewModes.onlineServices.allowed[number];

interface OnlineServicesState {
  viewMode: OnlineServicesViewMode;
  searchQuery: string;
  billingCycleFilters: string[];
  paymentMethodFilters: string[];
  setViewMode: (mode: OnlineServicesViewMode) => void;
  setSearchQuery: (query: string) => void;
  setBillingCycleFilters: (filters: string[]) => void;
  setPaymentMethodFilters: (filters: string[]) => void;
  services: OnlineService[];
  setServices: (services: OnlineService[]) => void;
  toggleActiveOnlineService: (serviceName: string, isActive: boolean) => void;
}

export const useOnlineServicesStore = create<OnlineServicesState>((set) => ({
  viewMode: getStoredViewMode(
    featureViewModes.onlineServices.storageKey,
    featureViewModes.onlineServices.default,
    featureViewModes.onlineServices.allowed, // Need to pass allowed modes here to differentiate between modes with or without table view
  ),
  searchQuery: '',
  billingCycleFilters: [],
  paymentMethodFilters: [],
  setViewMode: (mode) => {
    // Only allow modes that are valid for online services
    if (featureViewModes.onlineServices.allowed.includes(mode)) {
      localStorage.setItem(featureViewModes.onlineServices.storageKey, mode);
      set({ viewMode: mode });
    }
  },
  setSearchQuery: (query) => set({ searchQuery: query }),
  setBillingCycleFilters: (filters) => set({ billingCycleFilters: filters }),
  setPaymentMethodFilters: (filters) => set({ paymentMethodFilters: filters }),
  services: [],
  setServices: (services) => set({ services }),
  toggleActiveOnlineService: (serviceName, isActive) =>
    set((state) => ({
    services: state.services.map((service) =>
      service.name === serviceName
        ? { ...service, status: isActive ? SERVICE_STATUS_CODES.ACTIVE : SERVICE_STATUS_CODES.INACTIVE }
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
  useOnlineServicesStore((state) =>
    state.services.find((service) => service.name === serviceName)?.status === SERVICE_STATUS_CODES.ACTIVE || false
  );