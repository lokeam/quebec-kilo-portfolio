import { SERVICE_TYPES, type ServiceStatusCode, type ServiceType } from '@/shared/constants/service.constants';
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import type { ServiceTierName } from '../types/online-services/tiers';
import type { DigitalLocation } from '../types/media-storage/digital-location.types';

/**
 * Determines if a service is a subscription service based on multiple properties
 * This handles different ways a service might indicate it's a subscription
 */
export function isSubscriptionService(service: OnlineService): boolean {
  // Check multiple indicators
  return (
    service.type === SERVICE_TYPES.SUBSCRIPTION ||
    service.isSubscriptionService === true ||
    (service.billing?.cycle !== undefined && service.billing?.cycle !== 'NA')
  );
}

/**
 * Derives a logo identifier from a service name
 * @param serviceName The name of the service to get a logo for
 * @returns A standardized logo identifier or undefined if not found
 */
export function getLogo(serviceName: string): string | undefined {
  if (!serviceName) return undefined;

  // Convert to lowercase and trim whitespace
  const normalizedName = serviceName.toLowerCase().trim();

  // Common mappings for service names to logo keys
  const logoMappings: Record<string, string> = {
    'playstation': 'playstation',
    'playstation network': 'playstation',
    'playstation plus': 'playstation',
    'psn': 'playstation',
    'xbox': 'xbox',
    'xbox network': 'xbox',
    'xbox game pass': 'xbox',
    'steam': 'steam',
    'epic games': 'epic',
    'epic games store': 'epic',
    'nintendo': 'nintendo',
    'nintendo switch online': 'nintendo',
    'ea play': 'ea',
    'electronic arts': 'ea',
    'ubisoft': 'ubisoft',
    'ubisoft+': 'ubisoft',
    'gog': 'gog',
    'gog.com': 'gog',
    'humble bundle': 'humble',
    'humble': 'humble',
    'green man gaming': 'greenman',
    'fanatical': 'fanatical',
    'apple arcade': 'apple',
    'netflix games': 'netflix',
    'geforce now': 'nvidia',
    'nvidia': 'nvidia',
    'prime gaming': 'prime',
    'amazon luna': 'luna',
    'luna': 'luna',
    'meta quest': 'meta',
    'meta': 'meta',
    'google': 'playpass',
    'google play': 'playpass',
    'google play pass': 'playpass',
    'playpass': 'playpass',
    'play pass': 'playpass',
  };

  if (logoMappings[normalizedName]) {
    return logoMappings[normalizedName];
  }

  // Simplify the name by removing spaces and special characters
  return normalizedName
    .replace(/\s+/g, '')
    .replace(/[^\w]/g, '');
}

/**
 * Transforms a backend DigitalLocation object to the frontend OnlineService format
 * @param location The DigitalLocation from the backend
 * @returns A properly formatted OnlineService object for the frontend
 */
export function transformDigitalLocationToService(location: DigitalLocation): OnlineService {
  // Determine if this is a subscription service
  const isSubscriptionService = location.isSubscriptionService || false;
  const isActive = location.isActive || false;
  const serviceType = location.serviceType as ServiceType;

  return {
    id: location.id,
    name: location.name,
    logo: location.logo || 'default-logo',
    url: location.url || '#',
    type: serviceType,
    status: isActive ? 'active' as ServiceStatusCode : 'inactive' as ServiceStatusCode,
    features: [],
    label: location.label || location.name,
    createdAt: location.createdAt,
    updatedAt: location.updatedAt,
    isSubscriptionService,
    serviceType,
    isActive,
    tier: {
      currentTier: 'Basic' as ServiceTierName,
      availableTiers: [{
        name: 'Basic' as ServiceTierName,
        features: [],
        id: `tier-basic`,
        isDefault: true
      }]
    },
    billing: location.billing ? {
      cycle: location.billing.cycle || 'NA',
      fees: {
        monthly: location.billing.fees?.monthly || '0',
        quarterly: location.billing.fees?.quarterly || '0',
        annual: location.billing.fees?.annual || '0'
      },
      paymentMethod: location.billing.paymentMethod || 'Generic',
      renewalDate: location.billing.renewalDate ? {
        month: location.billing.renewalDate.month,
        day: location.billing.renewalDate.day
      } : undefined
    } : undefined
  };
}
