import { memo } from 'react';

// Components
import { MemoizedDashboardBadge } from '@/features/dashboard/components/molecules/DashboardBadge/DashboardBadge';

// Hooks + Utils
import { LogoOrIcon } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/LogoOrIcon';
import { useFormattedDate } from '@/features/dashboard/lib/hooks/useFormattedDate';

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

// Local Type Definitions
interface BaseSpendItem {
  id: string;
  title: string;
  amount: number;
  spendTransactionType: 'subscription' | 'one-time';
  paymentMethod: string;
  mediaType: string;
  serviceName?: {
    id: string;
    displayName: string;
  };
  createdAt: number;
  updatedAt: number;
  isActive: boolean;
}

interface SubscriptionSpend extends BaseSpendItem {
  spendTransactionType: 'subscription';
  billingCycle: string;
  nextBillingDate: number;
  yearlySpending: Array<{
    year: number;
    amount: number;
  }>;
}

interface OneTimeSpend extends BaseSpendItem {
  spendTransactionType: 'one-time';
  isDigital: boolean;
  isWishlisted: boolean;
  purchaseDate: number;
}

interface MonthlySpendingAccordionItemProps {
  item: SubscriptionSpend | OneTimeSpend;
  onClick?: () => void;
}

// Type Guards
const isSubscriptionSpend = (item: SubscriptionSpend | OneTimeSpend): item is SubscriptionSpend => {
  return item.spendTransactionType === 'subscription';
};

export const MemoizedMonthlySpendingAccordionItem = memo(function MonthlySpendingAccordionItem({
  item,
  onClick,
}: MonthlySpendingAccordionItemProps) {
  const { spendTransactionType } = item;
  const nextBillingDate = isSubscriptionSpend(item) ? item.nextBillingDate : undefined;
  const purchaseDate = !isSubscriptionSpend(item) ? item.purchaseDate : undefined;
  const dateDisplay = useFormattedDate(
    spendTransactionType,
    nextBillingDate,
    purchaseDate
  );

  const handleClick = () => {
    console.log(`Clicked on ${item.title} payment`);
    onClick?.();
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
              name={item.serviceName?.id ?? ''}
              mediaType={item.mediaType as 'subscription' | 'dlc' | 'inGamePurchase' | 'disc' | 'hardware'}
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
            <MemoizedDashboardBadge className={getMediaTypeStyle('subscription')}>
              {item.billingCycle}
            </MemoizedDashboardBadge>
        )}

        <MemoizedDashboardBadge
          className={`hidden md:inline-flex ${getTransactionTypeStyle(spendTransactionType)}`}>
          {spendTransactionType}
        </MemoizedDashboardBadge>

        {isSubscriptionSpend(item) && (
          <MemoizedDashboardBadge className={getMediaTypeStyle(item.mediaType)}>
            {item.mediaType}
          </MemoizedDashboardBadge>
        )}
        <span className="text-slate-200 w-24 text-right">${item.amount}</span>
      </div>
    </div>
  );
});
