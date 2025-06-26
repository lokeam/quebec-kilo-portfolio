import type {
  DigitalLocationBFFResponseItem,
  DigitalLocation
} from '@/types/domain/digital-location';

/**
 * Transforms BFF response items to DigitalLocation format
 */
export const adaptBFFResponseItemToDigitalLocation = (
  item: DigitalLocationBFFResponseItem
):DigitalLocation => ({
  id: item.id,
  name: item.name,
  items: item.items,
  createdAt: item.createdAt,
  updatedAt: item.updatedAt,
  isSubscription: item.isSubscription,
  monthlyCost: item.monthlyCost,
  isActive: item.isActive,
  url: item.url,
  paymentMethod: item.paymentMethod,
  paymentDate: item.paymentDate,
  billingCycle: item.billingCycle,
  costPerCycle: item.costPerCycle,
  nextPaymentDate: item.nextPaymentDate,
});

/**
 * Transforms BFF response to DigitalLocation array
 */
export const adaptBFFResponseToDigitalLocations = (
  response: { digitalLocations: DigitalLocationBFFResponseItem[] | null }
): DigitalLocation[] => {
  // Defensive check: ensure we always have an array, even if backend returns null
  const digitalLocations = response.digitalLocations || [];
  return digitalLocations.map(adaptBFFResponseItemToDigitalLocation);
};