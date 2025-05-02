/**
 * Digital Services API
 *
 * For API standards and best practices, see:
 * @see {@link ../../../docs/api-standards.md}
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiDebug } from '@/core/utils/debug/debug-utils';

interface ApiDigitalServiceItem {
  id: string;
  name: string;
  logo: string;
  is_subscription_service: boolean;
  url: string;
}

export interface DigitalServiceItem {
  id: string;
  name: string;
  logo: string;
  isSubscriptionService: boolean;
  url: string;
}

// Transform API response to match our interface
const transformApiService = (service: ApiDigitalServiceItem): DigitalServiceItem => ({
  id: service.id,
  name: service.name,
  logo: service.logo,
  isSubscriptionService: service.is_subscription_service,
  url: service.url
});

// Fallback data to use when API fails
export const FALLBACK_SERVICES: DigitalServiceItem[] = [
  { id: 'amazonluna', name: 'Amazon Luna', logo: 'amazon', isSubscriptionService: true, url: 'https://luna.amazon.com/' },
  { id: 'applearcade', name: 'Apple Arcade', logo: 'apple', isSubscriptionService: true, url: 'https://www.apple.com/apple-arcade/' },
  { id: 'blizzard', name: 'Blizzard Battle.net', logo: 'blizzard', isSubscriptionService: false, url: 'https://www.blizzard.com/en-us/' },
  { id: 'ea', name: 'EA Play', logo: 'ea', isSubscriptionService: true, url: 'https://www.ea.com/ea-play' },
  { id: 'epicgames', name: 'Epic Games', logo: 'epicgames', isSubscriptionService: false, url: 'https://store.epicgames.com/en-US/' },
  { id: 'fanatical', name: 'Fanatical', logo: 'fanatical', isSubscriptionService: false, url: 'https://www.fanatical.com/en/' },
  { id: 'gog', name: 'GOG', logo: 'gog', isSubscriptionService: false, url: 'https://www.gog.com/en/' },
  { id: 'googleplaypass', name: 'Google Play Pass', logo: 'google', isSubscriptionService: true, url: 'https://play.google.com/store/pass/getstarted/' },
  { id: 'greenmangaming', name: 'Green Man Gaming', logo: 'greenmangaming', isSubscriptionService: false, url: 'https://www.greenmangaming.com/' },
  { id: 'humblebundle', name: 'Humble Bundle', logo: 'humblebundle', isSubscriptionService: false, url: 'https://www.humblebundle.com/' },
  { id: 'itchio', name: 'itch.io', logo: 'itchio', isSubscriptionService: false, url: 'https://itch.io/' },
  { id: 'meta', name: 'Meta', logo: 'meta', isSubscriptionService: false, url: 'https://www.meta.com/nz/meta-quest-plus/' },
  { id: 'netflix', name: 'Netflix', logo: 'netflix', isSubscriptionService: true, url: 'https://www.netflix.com/' },
  { id: 'nintendo', name: 'Nintendo Switch Online', logo: 'nintendo', isSubscriptionService: true, url: 'https://www.nintendo.com/' },
  { id: 'nvidia', name: 'NVIDIA', logo: 'nvidia', isSubscriptionService: true, url: 'https://www.nvidia.com/en-us/geforce-now/' },
  { id: 'primegaming', name: 'Prime Gaming', logo: 'prime', isSubscriptionService: true, url: 'https://gaming.amazon.com/home' },
  { id: 'playstation', name: 'PlayStation Network', logo: 'ps', isSubscriptionService: true, url: 'https://www.playstation.com/en-us/playstation-network/' },
  { id: 'shadow', name: 'Shadow', logo: 'shadow', isSubscriptionService: true, url: 'https://shadow.tech/' },
  { id: 'steam', name: 'Steam', logo: 'steam', isSubscriptionService: false, url: 'https://store.steampowered.com/' },
  { id: 'ubisoft', name: 'Ubisoft', logo: 'ubisoft', isSubscriptionService: false, url: 'https://www.ubisoft.com/en-us/' },
  { id: 'xboxlive', name: 'Xbox Live', logo: 'xbox', isSubscriptionService: true, url: 'https://www.xbox.com/en-US/live' },
  { id: 'xboxgamepass', name: 'Xbox Game Pass', logo: 'xbox', isSubscriptionService: true, url: 'https://www.xbox.com/en-US/xbox-game-pass' }
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

      // The API returns the array directly, not wrapped in an Axios response
      const data = await axiosInstance.get<ApiDigitalServiceItem[]>('/v1/locations/digital/services/catalog', config);

      // Debug log the data
      console.log('Raw data:', data);
      console.log('Data type:', typeof data);
      console.log('Is array?', Array.isArray(data));
      if (Array.isArray(data)) {
        console.log('First item:', data[0]);
      }

      // Check if data exists and is an array
      if (!data) {
        apiDebug.logError(new Error('Response data is undefined'));
        return FALLBACK_SERVICES;
      }

      if (!Array.isArray(data)) {
        apiDebug.logError(new Error(`Response data is not an array, got ${typeof data}`));
        return FALLBACK_SERVICES;
      }

      // Type guard function to validate the data structure
      const isApiDigitalServiceItem = (item: unknown): item is ApiDigitalServiceItem => {
        const isValid = typeof item === 'object' &&
          item !== null &&
          'id' in item &&
          'name' in item &&
          'logo' in item &&
          'is_subscription_service' in item &&
          'url' in item;

        if (!isValid) {
          console.log('Invalid item structure:', item);
        }
        return isValid;
      };

      // Validate the data structure
      const isValidData = data.every(isApiDigitalServiceItem);

      if (!isValidData) {
        apiDebug.logError(new Error('Response data items do not match expected structure'));
        console.log('Invalid data structure:', data[0]);
        return FALLBACK_SERVICES;
      }

      apiDebug.logResponse({ status: 200, data });
      return data.map(transformApiService);
    } catch (error) {
      apiDebug.logError(error);
      console.log('Using fallback data due to API error');
      return FALLBACK_SERVICES;
    }
  }
};