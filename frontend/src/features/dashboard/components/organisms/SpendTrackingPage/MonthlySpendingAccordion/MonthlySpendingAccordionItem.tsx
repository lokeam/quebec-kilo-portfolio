import { memo } from 'react';

// Components
import { MemoizedDashboardBadge } from '@/features/dashboard/components/molecules/DashboardBadge/DashboardBadge';

// Hooks + Utils
import { useFormattedDate } from '@/features/dashboard/lib/hooks/useFormattedDate';
import { MediaIcon } from '@/features/dashboard/lib/utils/getMediaIcon';
import { DigitalLocationIcon } from '@/features/dashboard/lib/utils/getDigitalLocationIcon';
import { formatCurrency } from '@/features/dashboard/lib/utils/formatCurrency';

// Types
import type { SpendItem } from '@/types/domain/spend-tracking';
import { MediaCategory } from '@/types/domain/spend-tracking';

const getMediaTypeStyle = (type: string) => {
  switch(type) {
    case 'hardware': return "bg-green-700/50 text-foreground";
    case 'dlc': return "bg-orange-700/50 text-foreground";
    case 'inGamePurchase': return "bg-blue-600/50 text-foreground";
    case 'disc': return "bg-blue-400/50 text-foreground";
    case 'physical': return "bg-yellow-400/50 text-foreground";
    case 'subscription': return "bg-red-800/50 text-foreground";
    default: return "bg-muted text-foreground";
  }
}

const getTransactionTypeStyle = (type: string) => {
  switch(type) {
    case 'subscription': return "bg-purple-900/50 text-foreground";
    case 'one-time': return "bg-muted text-foreground";
    default: return "bg-muted text-foreground";
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
    // console.log(`Clicked on ${item.title} payment`);
    onClick?.();
  };

  const renderIcon = () => {
    // For subscriptions, use the digital location icon

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
      className="flex items-center justify-between p-4 hover:bg-muted/50 cursor-pointer transition-colors rounded-lg"
      onClick={handleClick}
    >
      <div className="flex items-center gap-1">
        <span className="text-muted-foreground w-16">
          {dateDisplay}
        </span>
        <div className="flex items-center gap-3">
          <div className="h-9 w-9 flex items-center justify-center">
            {renderIcon()}
          </div>
          <div className="flex flex-col">
            <span
              className="text-foreground truncate text-wrap max-w-[120px] lg:max-w-full"
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
        <span className="text-foreground w-24 text-right">{formatCurrency(item.amount)}</span>
      </div>
    </div>
  );
});
