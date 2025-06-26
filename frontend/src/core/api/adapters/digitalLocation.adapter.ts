import type { DigitalLocation } from '@/types/domain/digital-location';

/**
 * Transforms a DigitalLocation from the API into the domain OnlineService format
 */
export const adaptDigitalLocationToService = (location: DigitalLocation): DigitalLocation => ({
  id: location.id,
  name: location.name,
  locationType: 'subscription',
  itemCount: location.itemCount,
  isSubscription: location.isSubscription,
  monthlyCost: location.monthlyCost,
  isActive: location.isActive,
  url: location.url,
  createdAt: location.createdAt,
  updatedAt: location.updatedAt,
  items: location.items,
  paymentMethod: location.paymentMethod,
  paymentDate: location.paymentDate,
  billingCycle: location.billingCycle,
  costPerCycle: location.costPerCycle,
  nextPaymentDate: location.nextPaymentDate
});