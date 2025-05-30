/**
 * Digital Services Catalog Service
 *
 * Provides functions for managing digital services catalog through the backend API.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';
import type { DigitalServiceItem } from '@/types/domain/online-service';

interface CatalogResponse {
  success: boolean;
  catalog: DigitalServiceItem[];
}

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
  { id: 'nintendo', name: 'Nintendo Switch Online', logo: 'nintendo', isSubscriptionService: true, url: 'https://www.nintendo.com/' },
  { id: 'nvidia', name: 'NVIDIA', logo: 'nvidia', isSubscriptionService: true, url: 'https://www.nvidia.com/en-us/geforce-now/' },
  { id: 'primegaming', name: 'Prime Gaming', logo: 'prime', isSubscriptionService: true, url: 'https://gaming.amazon.com/home' },
  { id: 'playstation', name: 'PlayStation Network', logo: 'ps', isSubscriptionService: true, url: 'https://www.playstation.com/en-us/playstation-network/' },
  { id: 'shadow', name: 'Shadow', logo: 'shadow', isSubscriptionService: true, url: 'https://shadow.tech/' },
  { id: 'steam', name: 'Steam', logo: 'steam', isSubscriptionService: false, url: 'https://store.steampowered.com/' },
  { id: 'ubisoft', name: 'Ubisoft', logo: 'ubisoft', isSubscriptionService: false, url: 'https://www.ubisoft.com/en-us/' },
  { id: 'xboxgamepass', name: 'Xbox Game Pass', logo: 'xbox', isSubscriptionService: true, url: 'https://www.xbox.com/en-US/xbox-game-pass' }
];

const CATALOG_ENDPOINT = '/v1/locations/digital/services/catalog';

/**
 * Fetches the digital services catalog
 *
 * Uses apiRequest helper to wrap the axios call with:
 *  - async/await syntax
 *  - pre‑call debug log
 *  - post‑call success log
 *  - catch block with error log + optional Sentry/metrics
 *  - retry logic (if configured)
 *
 * Usage:
 *   return apiRequest('getServicesCatalog', () => axios.get(...));
 */
export const getServicesCatalog = (): Promise<DigitalServiceItem[]> =>
  apiRequest('getServicesCatalog', () =>
    axiosInstance
      .get<CatalogResponse>(CATALOG_ENDPOINT)
      .then(response => {
        if (!response.data.catalog) {
          return [];
        }
        return response.data.catalog;
      })
  );
