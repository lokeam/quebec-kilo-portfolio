import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiDebug } from '@/core/utils/debug/debug-utils';

export interface DigitalServiceItem {
  id: string;
  name: string;
  logo: string;
}

// Fallback data to use when API fails
const FALLBACK_SERVICES: DigitalServiceItem[] = [
  { id: 'amazonluna', name: 'Amazon Luna', logo: 'amazon' },
  { id: 'applearcade', name: 'Apple Arcade', logo: 'apple' },
  { id: 'blizzard', name: 'Blizzard Battle.net', logo: 'blizzard' },
  { id: 'ea', name: 'EA Play', logo: 'ea' },
  { id: 'epicgames', name: 'Epic Games', logo: 'epicgames' },
  { id: 'fanatical', name: 'Fanatical', logo: 'fanatical' },
  { id: 'gog', name: 'GOG', logo: 'gog' },
  { id: 'googleplaypass', name: 'Google Play Pass', logo: 'google' },
  { id: 'greenmangaming', name: 'Green Man Gaming', logo: 'greenmangaming' },
  { id: 'humblebundle', name: 'Humble Bundle', logo: 'humblebundle' },
  { id: 'itchio', name: 'itch.io', logo: 'itchio' },
  { id: 'meta', name: 'Meta', logo: 'meta' },
  { id: 'netflix', name: 'Netflix', logo: 'netflix' },
  { id: 'nintendo', name: 'Nintendo', logo: 'nintendo' },
  { id: 'nvidia', name: 'NVIDIA', logo: 'nvidia' },
  { id: 'primegaming', name: 'Prime Gaming', logo: 'prime' },
  { id: 'playstation', name: 'PlayStation Network', logo: 'ps' },
  { id: 'shadow', name: 'Shadow', logo: 'shadow' },
  { id: 'ubisoft', name: 'Ubisoft', logo: 'ubisoft' },
  { id: 'xboxlive', name: 'Xbox Live', logo: 'xbox' },
  { id: 'xboxgamepass', name: 'Xbox Game Pass', logo: 'xbox' }
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

      const data = await axiosInstance.get<DigitalServiceItem[]>('/v1/locations/digital/services/catalog', config);

      apiDebug.logResponse({ status: 200, data });
      return data;
    } catch (error) {
      apiDebug.logError(error);
      console.log('Using fallback data due to API error');
      return FALLBACK_SERVICES;
    }
  }
};