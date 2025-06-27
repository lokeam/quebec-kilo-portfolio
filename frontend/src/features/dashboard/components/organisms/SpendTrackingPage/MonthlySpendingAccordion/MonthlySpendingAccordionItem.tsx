import { memo } from 'react';

// Components
import { MemoizedDashboardBadge } from '@/features/dashboard/components/molecules/DashboardBadge/DashboardBadge';

// Hooks + Utils
import { useFormattedDate } from '@/features/dashboard/lib/hooks/useFormattedDate';
import { MediaIcon } from '@/features/dashboard/lib/utils/getMediaIcon';
import { DigitalLocationIcon } from '@/features/dashboard/lib/utils/getDigitalLocationIcon';

// Types
import type { SpendItem } from '@/types/domain/spend-tracking';
import { MediaCategory } from '@/types/domain/spend-tracking';

const getMediaTypeStyle = (type: string) => {
  switch(type) {
    case 'hardware': return "bg-green-700/50 text-slate-200";
    case 'dlc': return "bg-orange-700/50 text-slate-200";
    case 'inGamePurchase': return "bg-blue-600/50 text-slate-200";
    case 'disc': return "bg-blue-400/50 text-slate-200";
    case 'physical': return "bg-yellow-400/50 text-slate-200";
    case 'subscription': return "bg-red-800/50 text-slate-200";
    default: return "bg-slate-700/50 text-slate-200";
  }
}

const getTransactionTypeStyle = (type: string) => {
  switch(type) {
    case 'subscription': return "bg-purple-900/50 text-purple-200";
    case 'one-time': return "bg-slate-700/50 text-slate-200";
    default: return "bg-slate-700/50 text-slate-200";
  }
}

interface MonthlySpendingAccordionItemProps {
  item: SpendItem;
  onClick?: () => void;
}

export const MemoizedMonthlySpendingAccordionItem = memo(function MonthlySpendingAccordionItem({
  item,
  onClick,
}: MonthlySpendingAccordionItemProps) {
  const { spendTransactionType } = item;
  const dateDisplay = useFormattedDate(
    spendTransactionType,
    item.nextBillingDate,
    item.purchaseDate
  );

  const handleClick = () => {
    console.log(`Clicked on ${item.title} payment`);
    onClick?.();
  };

  const renderIcon = () => {
    // For subscriptions, use the digital location icon
    console.log('renderIcon, item: ', item);

    if (item.mediaType === MediaCategory.SUBSCRIPTION) {
      return (
        <DigitalLocationIcon
          name={item?.provider ?? ''}
          className="h-6 w-6"
        />
      );
    }

    // For other media types, use the media icon
    return (
      <MediaIcon
        mediaType={item.mediaType}
        className="h-6 w-6"
      />
    );
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
            {renderIcon()}
          </div>
          <div className="flex flex-col">
            <span
              className="text-slate-200 truncate text-wrap max-w-[120px] lg:max-w-full"
            >{item.title}</span>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-4">
        {spendTransactionType === 'subscription' && item.billingCycle && (
            <MemoizedDashboardBadge className={getMediaTypeStyle('subscription')}>
              {item.billingCycle}
            </MemoizedDashboardBadge>
        )}

        <MemoizedDashboardBadge
          className={`hidden md:inline-flex ${getTransactionTypeStyle(spendTransactionType)}`}>
          {spendTransactionType}
        </MemoizedDashboardBadge>

        {spendTransactionType === 'subscription' && (
          <MemoizedDashboardBadge className={getMediaTypeStyle(item.mediaType)}>
            {item.mediaType}
          </MemoizedDashboardBadge>
        )}
        <span className="text-slate-200 w-24 text-right">${item.amount}</span>
      </div>
    </div>
  );
});
