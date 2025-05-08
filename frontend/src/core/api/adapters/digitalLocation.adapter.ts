import type { DigitalLocation } from '@/types/domain/digital-location';
import type { OnlineService as DomainOnlineService } from '@/types/domain/online-service';
import { SERVICE_STATUS_CODES, SERVICE_TYPES } from '@/shared/constants/service.constants';

/**
 * Transforms a DigitalLocation from the API into the domain OnlineService format
 */
export const adaptDigitalLocationToService = (location: DigitalLocation): DomainOnlineService => ({
  id: location.id,
  name: location.name,
  label: location.name,
  logo: '',
  url: '',
  status: SERVICE_STATUS_CODES.ACTIVE,
  createdAt: new Date(location.createdAt),
  updatedAt: new Date(location.updatedAt),
  isSubscriptionService: false,
  serviceType: SERVICE_TYPES.ONLINE,
  isActive: true,
  type: SERVICE_TYPES.ONLINE,
  billing: {
    cycle: 'monthly',
    fees: {
      monthly: '0',
      quarterly: '0',
      annual: '0'
    },
    paymentMethod: 'None'
  },
  tier: {
    currentTier: 'standard' as const,
    availableTiers: [{
      id: 'tier-standard',
      name: 'standard',
      features: [],
      isDefault: true
    }]
  },
  features: [],
  games: location.items
});