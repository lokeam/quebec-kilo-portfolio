import type { DigitalLocation } from '@/features/dashboard/lib/types/media-storage/digital-location.types';
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import { SERVICE_STATUS_CODES, SERVICE_TYPES } from '@/shared/constants/service.constants';

/**
 * Transforms a DigitalLocation from the API into the format expected by SingleOnlineServiceCard
 */
export const adaptDigitalLocationToService = (location: DigitalLocation): OnlineService => ({
  id: location.id,
  name: location.name,
  label: location.label,
  logo: location.logo,
  url: location.url,
  status: location.isActive ? SERVICE_STATUS_CODES.ACTIVE : SERVICE_STATUS_CODES.INACTIVE,
  createdAt: new Date(location.createdAt),
  updatedAt: new Date(location.updatedAt),
  isSubscriptionService: location.isSubscriptionService,
  serviceType: location.serviceType === 'subscription' ? SERVICE_TYPES.SUBSCRIPTION : SERVICE_TYPES.ONLINE,
  isActive: location.isActive,
  type: location.serviceType === 'subscription' ? SERVICE_TYPES.SUBSCRIPTION : SERVICE_TYPES.ONLINE,
  billing: location.billing || {
    cycle: 'monthly',
    fees: {
      monthly: '0',
      quarterly: '0',
      annual: '0'
    },
    paymentMethod: 'None'
  },
  tier: {
    currentTier: 'standard',
    availableTiers: [{
      id: 'tier-standard',
      name: 'standard',
      features: [],
      isDefault: true
    }]
  },
  features: []
});