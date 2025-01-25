import { useMemo } from 'react';
import { formatDate } from '@/features/dashboard/lib/utils/formatDate';
import  { type PurchasedMediaCategory, PURCHASED_MEDIA_CATEGORIES } from "@/features/dashboard/lib/types/service.types";
import type { ISO8601Date } from "@/shared/types/types";

export function useFormattedDate(
  spendTransactionType: PurchasedMediaCategory,
  nextBillingDate: ISO8601Date,
  purchaseDate: ISO8601Date,
) {
  const dateString = useMemo(
    () => spendTransactionType === PURCHASED_MEDIA_CATEGORIES.SUBSCRIPTION
      ? nextBillingDate
      : purchaseDate,
    [spendTransactionType, nextBillingDate, purchaseDate]
  );

  const formattedDate = useMemo(() => {
    return formatDate(dateString);
  }, [dateString]);

  const dateDisplay = useMemo(
    () => `${formattedDate.dayStr} ${formattedDate.monthStr}`.trim(),
    [formattedDate]
  );

  return dateDisplay;
}
