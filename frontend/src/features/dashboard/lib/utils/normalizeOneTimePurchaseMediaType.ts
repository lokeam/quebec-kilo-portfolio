
/**
 * Helper function for normalizing the camelCased media type strings to a more readable format.
 * Used in MemoizedDashboardBadge within the MonthlySpendintItemDetails component.
*/
const normalizeMediaTypeMap: Record<string, string> = {
  'inGamePurchase': 'In-Game Purchase',
  'physicalGame': 'Physical Game',
  'digitalGame': 'Digital Game',
  'hardware': 'Hardware',
  'dlc': 'DLC',
  'misc': 'Misc',
  'subscription': 'Subscription'
}

export const normalizeOneTimePurchaseMediaType = (camelCaseMediaType: string): string => {
  //console.log("-----")
  //console.log(`Normalizing one time purchase media type: ${camelCaseMediaType}`);
  const normalizedStr = normalizeMediaTypeMap[camelCaseMediaType];

  if (!normalizedStr) {
    //console.warn(`Unknown one time purchase media type: ${camelCaseMediaType}, defaulting to misc`);
    return 'misc';
  }

  return normalizedStr;
}