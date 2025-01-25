import { memo } from 'react';

// Components
import { MemoizedDashboardBadge } from '@/features/dashboard/components/molecules/DashboardBadge/DashboardBadge';

// Hooks + Utils
import { LogoOrIcon } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/LogoOrIcon';
import { useFormattedDate } from '@/features/dashboard/lib/hooks/useFormattedDate';

// Icons
// import { Check } from "lucide-react"


// Types
import type { SubscriptionSpend } from '@/features/dashboard/lib/types/spend-tracking/subscription';
import type { OneTimeSpend } from '@/features/dashboard/lib/types/spend-tracking/purchases';
import type { ISO8601Date } from '@/shared/types/types';

// Guards
import { isSubscriptionSpend } from '@/features/dashboard/lib/types/spend-tracking/guards';

// Constants
import { BADGE_STYLES } from '@/features/dashboard/lib/constants/style.constants';
import { type PurchasedMediaCategory } from '@/features/dashboard/lib/types/spend-tracking/media';


export interface MonthlySpendingAccordionItemProps {
  item: SubscriptionSpend | OneTimeSpend;
  onClick?: () => void;
}

export const MemoizedMonthlySpendingAccordionItem = memo(function MonthlySpendingAccordionItem({
  item,
  onClick,
}: MonthlySpendingAccordionItemProps) {
  const { spendTransactionType } = item;
  const nextBillingDate = isSubscriptionSpend(item) ? item.nextBillingDate : undefined;
  const purchaseDate = !isSubscriptionSpend(item) ? item.purchaseDate : undefined;
  const dateDisplay = useFormattedDate(
    spendTransactionType as PurchasedMediaCategory,
    nextBillingDate as ISO8601Date,
    purchaseDate as ISO8601Date
  )

  const handleClick = () => {
    console.log(`Clicked on ${item.title} payment`)
    onClick?.()
  };


  return (
    <div
      className="flex items-center justify-between p-4 hover:bg-slate-800/50 cursor-pointer transition-colors rounded-lg"
      onClick={handleClick}
    >
      <div className="flex items-center gap-1">
        <span className="text-slate-400 w-16">
          {dateDisplay}
        </span>
        <div className="flex items-center gap-3">
          <div className="h-9 w-9 flex items-center justify-center">
            <LogoOrIcon
              name={item.provider?.id ?? ''} // Need to provide a nullish check bc OneTimeSpends may or may not be purchased from an online service provider
              mediaType={item.mediaType}
            />
          </div>
          <div className="flex flex-col">
            <span
              className="text-slate-200 truncate text-wrap max-w-[120px] lg:max-w-full"
            >{item.title}</span>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-4">
        {isSubscriptionSpend(item) && (
            <MemoizedDashboardBadge className={BADGE_STYLES.mediaType.subscription}>
              {(item as SubscriptionSpend).billingCycle}
            </MemoizedDashboardBadge>
        )}

        <MemoizedDashboardBadge
          className={
            `hidden md:inline-flex
              ${BADGE_STYLES.spendTransactionType[item.spendTransactionType as keyof typeof BADGE_STYLES.spendTransactionType]}
            `}>
          {item.spendTransactionType}
        </MemoizedDashboardBadge>

        {isSubscriptionSpend(item) && (
          <MemoizedDashboardBadge className={BADGE_STYLES.mediaType[item.mediaType]}>
            {item.mediaType}
          </MemoizedDashboardBadge>
        )}
        <span className="text-slate-200 w-24 text-right">${item.amount}</span>
        {/* {isPaid && <Check className="h-4 w-4 text-green-500" />} */}
      </div>
    </div>
  )
});
