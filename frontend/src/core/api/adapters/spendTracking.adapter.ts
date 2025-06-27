import type { SpendTrackingBFFResponse, MediaCategory } from "@/types/domain/spend-tracking";

const mediaTypeMap: Record<string, string> = {
  'in_game_purchase': 'inGamePurchase',
  'physical_game': 'physicalGame',
  'digital_game': 'digitalGame',
  'hardware': 'hardware',
  'dlc': 'dlc',
  'misc': 'misc',
  'subscription': 'subscription'
};

const normalizeMediaType = (snakeCaseMediaType: string): string => {
  console.log("-----")
  console.log(`Normalizing media type: ${snakeCaseMediaType}`);
  const normalizedStr = mediaTypeMap[snakeCaseMediaType];

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
      currentTotalThisMonth: response.currentTotalThisMonth.map(item => ({
        ...item,
        mediaType: normalizeMediaType(item.mediaType) as MediaCategory
      }))
    };
  }
}