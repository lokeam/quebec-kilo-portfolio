import { useMemo } from 'react';
//import { formatDate } from '@/features/dashboard/lib/utils/formatDate';
import { TransactionType } from '@/types/domain/spend-tracking';
// import type { ISO8601Date } from "@/shared/types/types";
import { format } from 'date-fns';

export function useFormattedDate(
  spendTransactionType: 'subscription' | 'one-time',
  nextBillingDate: number | undefined,
  purchaseDate: number | undefined,
) {
  const dateString = useMemo(
    () => spendTransactionType === TransactionType.SUBSCRIPTION
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

export function formatDate(timestamp: number | undefined): { dayStr: string; monthStr: string } {
  if (!timestamp) {
    return { dayStr: '--', monthStr: '---' };
  }

  try {
    const date = new Date(timestamp);
    if (isNaN(date.getTime())) {
      return { dayStr: '--', monthStr: '---' };
    }
    return {
      dayStr: format(date, 'd'),
      monthStr: format(date, 'MMM')
    };
  } catch {
    return { dayStr: '--', monthStr: '---' };
  }
}
