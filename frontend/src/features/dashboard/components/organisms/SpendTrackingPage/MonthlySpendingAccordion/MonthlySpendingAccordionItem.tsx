import { memo } from 'react';

// Hooks + Utils
import { LogoOrIcon } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/LogoOrIcon';

// Icons
import { Check } from "lucide-react"
import { Badge } from "@/shared/components/ui/badge"

// Types
import type { SpendTrackingService } from "@/features/dashboard/lib/types/service.types"

export interface MonthlySpendingAccordionItemProps extends SpendTrackingService {
  onClick?: () => void;
}

const BADGE_STYLES = {
  mediaType: {
    hardware: "bg-green-700/50 text-slate-200",
    dlc: "bg-orange-700/50 text-slate-200",
    inGamePurchase: "bg-blue-600/50 text-slate-200",
    disc: "bg-blue-400/50 text-slate-200",
    physical: "bg-yellow-400/50 text-slate-200",
    subscription: "bg-red-800/50 text-slate-200"
  },
  spendType: {
    subscription: "bg-purple-900/50 text-purple-200",
    "one-time": "bg-slate-700/50 text-slate-200"
  }
} as const;
const SUBSCRIPTION_MEDIA = 'subscription';

// Tightly coupled render fn
function SpendingBadge({
  variant = "secondary",
  className,
  children
}: {
  variant?: "secondary";
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <Badge variant={variant} className={className}>
      {children}
    </Badge>
  );
}


export const MemoizedMonthlySpendingAccordionItem = memo(function MonthlySpendingAccordionItem({
  day,
  month,
  title,
  billingCycle,
  name,
  spendType,
  amount,
  isPaid,
  mediaType = SUBSCRIPTION_MEDIA,
  onClick,
}: MonthlySpendingAccordionItemProps) {
  const handleClick = () => {
    console.log(`Clicked on ${title} payment`)
    onClick?.()
  };

  return (
    <div
      className="flex items-center justify-between p-4 hover:bg-slate-800/50 cursor-pointer transition-colors rounded-lg"
      onClick={handleClick}
    >
      <div className="flex items-center gap-1">
        <span className="text-slate-400 w-16">
          {month} {day}
        </span>
        <div className="flex items-center gap-3">
          <div className="h-9 w-9 flex items-center justify-center">
            <LogoOrIcon name={name} mediaType={mediaType} />
          </div>
          <div className="flex flex-col">
            <span
              className="text-slate-200 truncate text-wrap max-w-[120px] lg:max-w-full"
            >{title}</span>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-4">
      {mediaType === SUBSCRIPTION_MEDIA && (
          <SpendingBadge className={BADGE_STYLES.mediaType.subscription}>
            {billingCycle}
          </SpendingBadge>
        )}

        <SpendingBadge className={`hidden md:inline-flex ${BADGE_STYLES.spendType[spendType]}`}>
          {spendType}
        </SpendingBadge>

        {mediaType !== SUBSCRIPTION_MEDIA && (
          <SpendingBadge className={BADGE_STYLES.mediaType[mediaType]}>
            {mediaType}
          </SpendingBadge>
        )}
        <span className="text-slate-200 w-24 text-right">${amount}</span>
        {isPaid && <Check className="h-4 w-4 text-green-500" />}
      </div>
    </div>
  )
});
