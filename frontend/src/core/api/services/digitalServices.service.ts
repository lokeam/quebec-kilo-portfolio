import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiDebug } from '@/core/utils/debug/debug-utils';
import type { AxiosResponse } from 'axios';

interface ApiDigitalServiceItem {
  id: string;
  name: string;
  logo: string;
  is_subscription_service: boolean;
}

export interface DigitalServiceItem {
  id: string;
  name: string;
  logo: string;
  isSubscriptionService: boolean;
}

// Transform API response to match our interface
const transformApiService = (service: ApiDigitalServiceItem): DigitalServiceItem => ({
  id: service.id,
  name: service.name,
  logo: service.logo,
  isSubscriptionService: service.is_subscription_service
});

// Fallback data to use when API fails
export const FALLBACK_SERVICES: DigitalServiceItem[] = [
  { id: 'amazonluna', name: 'Amazon Luna', logo: 'amazon', isSubscriptionService: true },
  { id: 'applearcade', name: 'Apple Arcade', logo: 'apple', isSubscriptionService: true },
  { id: 'blizzard', name: 'Blizzard Battle.net', logo: 'blizzard', isSubscriptionService: false },
  { id: 'ea', name: 'EA Play', logo: 'ea', isSubscriptionService: true },
  { id: 'epicgames', name: 'Epic Games', logo: 'epicgames', isSubscriptionService: false },
  { id: 'fanatical', name: 'Fanatical', logo: 'fanatical', isSubscriptionService: false },
  { id: 'gog', name: 'GOG', logo: 'gog', isSubscriptionService: false },
  { id: 'googleplaypass', name: 'Google Play Pass', logo: 'google', isSubscriptionService: true },
  { id: 'greenmangaming', name: 'Green Man Gaming', logo: 'greenmangaming', isSubscriptionService: false },
  { id: 'humblebundle', name: 'Humble Bundle', logo: 'humblebundle', isSubscriptionService: false },
  { id: 'itchio', name: 'itch.io', logo: 'itchio', isSubscriptionService: false },
  { id: 'meta', name: 'Meta', logo: 'meta', isSubscriptionService: false },
  { id: 'netflix', name: 'Netflix', logo: 'netflix', isSubscriptionService: true },
  { id: 'nintendo', name: 'Nintendo', logo: 'nintendo', isSubscriptionService: true },
  { id: 'nvidia', name: 'NVIDIA', logo: 'nvidia', isSubscriptionService: true },
  { id: 'primegaming', name: 'Prime Gaming', logo: 'prime', isSubscriptionService: true },
  { id: 'playstation', name: 'PlayStation Network', logo: 'ps', isSubscriptionService: true },
  { id: 'shadow', name: 'Shadow', logo: 'shadow', isSubscriptionService: true },
  { id: 'steam', name: 'Steam', logo: 'steam', isSubscriptionService: false },
  { id: 'ubisoft', name: 'Ubisoft', logo: 'ubisoft', isSubscriptionService: false },
  { id: 'xboxlive', name: 'Xbox Live', logo: 'xbox', isSubscriptionService: true },
  { id: 'xboxgamepass', name: 'Xbox Game Pass', logo: 'xbox', isSubscriptionService: true }
];

export const digitalServicesService = {
  getServicesCatalog: async (token?: string): Promise<DigitalServiceItem[]> => {
    try {
      const config = {
        headers: token ? {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        } : undefined
      };

      apiDebug.logRequest({ url: '/v1/locations/digital/services/catalog', method: 'GET', ...config });

      const response = await axiosInstance.get<ApiDigitalServiceItem[], AxiosResponse<ApiDigitalServiceItem[]>>('/v1/locations/digital/services/catalog', config);

      apiDebug.logResponse({ status: 200, data: response.data });
      return response.data.map(transformApiService);
    } catch (error) {
      apiDebug.logError(error);
      console.log('Using fallback data due to API error');
      return FALLBACK_SERVICES;
    }
  }
};