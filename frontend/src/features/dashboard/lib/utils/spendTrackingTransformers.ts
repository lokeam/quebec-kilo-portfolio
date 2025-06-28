import type { SpendingItemBFFResponse } from '@/types/domain/spend-tracking';
import type { SpendTrackingData } from '@/features/dashboard/components/organisms/SpendTrackingPage/SpendTrackingForm/SpendTrackingForm';

/**
 * Helper function to transform Spend Tracking BFF response items to form data
 * Converts SpendingItemBFFResponse types to FormValues for the edit form
 * Used in the MonthlySpendingAccordion component to populate the edit form
*/
export const transformSpendingItemResponseToFormData = (item: SpendingItemBFFResponse): SpendTrackingData => {
  return {
    id: item.id,
    title: item.title,
    spending_category_id: getSpendingCategoryId(item.mediaType),
    amount: item.amount,
    payment_method: item.paymentMethod,
    purchase_date: convertTimestampToDate(item.purchaseDate ?? 0),
    digital_location_id: item.serviceName?.id,
  }
}

/**
 * Helper function to convert the mediaType string to a spending category ID
 * Maps MediaCategory values to database category IDs
*/
const getSpendingCategoryId = (mediaType: string): number => {
  const categoryMap: Record<string, number> = {
    'hardware': 1,
    'dlc': 2,
    'in_game_purchase': 3,
    'physical_game': 4,
    'digital_game': 5,
    'misc': 6,
    'subscription': 5, // Default to digital for subscriptions
  };

  return categoryMap[mediaType] || 5; // Default to digital_game value for unknown media types
}

/**
 * Helper function to convert a timestamp to a date object
 * Handles both seconds and millisecond timestamps
*/
const convertTimestampToDate = (timestamp: number): Date => {
  if (!timestamp) return new Date();

  // If the timestamp is less than 10billion, its probably in seconds, if greater its likely in milliseconds
  const milliseconds = timestamp < 10000000000 ? timestamp * 1000 : timestamp;

  return new Date(milliseconds);
}

/**
 * Checks if a spend tracking item may be edited.
 * Only one-time purchases (not subscriptions) may be edited
 */
export const isEditableSpendItem = (item: SpendingItemBFFResponse): boolean => {
  return item.spendTransactionType !== 'subscription';
}