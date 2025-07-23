import type { SpendTrackingBFFResponse, MediaCategory } from "@/types/domain/spend-tracking";

const camelCaseMediaTypeMap: Record<string, string> = {
  'in_game_purchase': 'inGamePurchase',
  'physical_game': 'physicalGame',
  'digital_game': 'digitalGame',
  'hardware': 'hardware',
  'dlc': 'dlc',
  'misc': 'misc',
  'subscription': 'subscription'
};

const camelCaseMediaType = (snakeCaseMediaType: string): string => {
  const normalizedStr = camelCaseMediaTypeMap[snakeCaseMediaType];

  if (!normalizedStr) {
    console.warn(`Unknown media type: ${snakeCaseMediaType}, defaulting to misc`);
    return 'misc';
  }

  return normalizedStr;
}

export const spendTrackingAdapter = {
  transformSpendTrackingResponse: (response: SpendTrackingBFFResponse): SpendTrackingBFFResponse => {
    return {
      ...response,
      currentTotalThisMonth: response.currentTotalThisMonth?.map(item => ({
        ...item,
        mediaType: camelCaseMediaType(item.mediaType) as MediaCategory
      })) || null,
      oneTimeThisMonth: response.oneTimeThisMonth?.map(item => ({
        ...item,
        mediaType: camelCaseMediaType(item.mediaType) as MediaCategory
      })) || null,
      recurringNextMonth: response.recurringNextMonth?.map(item => ({
        ...item,
        mediaType: camelCaseMediaType(item.mediaType) as MediaCategory
      })) || null,
      totalMonthlySpending: {
        ...response.totalMonthlySpending,
        percentageChange: Math.round(response.totalMonthlySpending.percentageChange * 10) / 10
      }
    };
  }
}