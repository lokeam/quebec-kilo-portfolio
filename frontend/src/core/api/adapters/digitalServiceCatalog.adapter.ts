/**
 * Digital Service Catalog Adapters
 *
 * Provides transformation functions for digital service catalog data between API and domain types.
 */

import type { DigitalServiceItem } from '@/core/api/services/digitalServicesCatalog.service';
import type { DigitalLocation } from '@/types/domain/online-service';

/**
 * Transforms a DigitalServiceItem from the API to a DigitalLocation domain type
 *
 * @param service - The service item from the API
 * @returns A DigitalLocation object with default values for required fields
 */
export function adaptDigitalServiceToLocation(service: DigitalServiceItem): DigitalLocation {
  return {
    id: service.id,
    name: service.name,
    locationType: 'subscription',
    itemCount: 0,
    isSubscription: service.isSubscriptionService,
    monthlyCost: 0,
    isActive: true,
    url: service.url || '#',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    items: [],
    paymentMethod: 'Generic',
    paymentDate: new Date().toISOString(),
    billingCycle: 'NA',
    costPerCycle: 0,
    nextPaymentDate: new Date().toISOString()
  };
}

/**
 * Transforms an array of DigitalServiceItems to DigitalLocations
 *
 * @param services - Array of service items from the API
 * @returns Array of DigitalLocation objects
 */
export function adaptDigitalServicesToLocations(services: DigitalServiceItem[]): DigitalLocation[] {
  return services.map(adaptDigitalServiceToLocation);
}
