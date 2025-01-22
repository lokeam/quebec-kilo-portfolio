
import { Check } from "lucide-react"
import { Badge } from "@/shared/components/ui/badge"
import type { SpendTrackingMediaType, SpendTrackingService } from "@/features/dashboard/lib/types/service.types"

import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';

export interface RecurringExpensesItemProps extends SpendTrackingService {
  onClick?: () => void;
  isWatchedByResizeObserver?: boolean;
}

export function RecurringExpensesItem({
  day,
  month,
  title,
  billingCycle,
  name,
  spendType,
  amount,
  isPaid,
  mediaType = 'subscription',
  onClick,
  isWatchedByResizeObserver,
}: RecurringExpensesItemProps) {
  const { games, physicalMedia, digitalMedia } = useDomainMaps();


  const getLogoOrIcon = (name: string, mediaType: SpendTrackingMediaType) => {
    if (!mediaType) return null;
    let IconComponent = null;

    switch (mediaType) {
      case 'subscription': {
        const LogoComponent = games[name];
        return LogoComponent ? <LogoComponent className="w-full h-full object-contain" /> : null;
      }
      case 'dlc':
        IconComponent = digitalMedia[mediaType];
        return IconComponent ? <IconComponent className="w-full h-full" /> : null;
      case 'inGamePurchase':
        IconComponent = digitalMedia[mediaType];
        return IconComponent ? <IconComponent className="w-full h-full" /> : null;
      case 'disc':
          IconComponent = physicalMedia[mediaType];
          return IconComponent ? <IconComponent className="w-full h-full" /> : null;
      case 'hardware':
        IconComponent = physicalMedia[mediaType];
        return IconComponent ? <IconComponent className="w-full h-full" /> : null;
      default:
        return null;
    }
   };
  return (
    <div
      className="flex items-center justify-between p-4 hover:bg-slate-800/50 cursor-pointer transition-colors rounded-lg"
      {...(isWatchedByResizeObserver ? { 'data-card-sentinel': true } : {})}
      onClick={() => {
        console.log(`Clicked on ${title} payment`)
        onClick?.()
      }}
    >
      <div className="flex items-center gap-1">
        <span className="text-slate-400 w-16">
          {month} {day}
        </span>
        <div className="flex items-center gap-3">
          <div className="h-9 w-9 flex items-center justify-center">
            {getLogoOrIcon(name, mediaType)}
          </div>
          <div className="flex flex-col">
            <span
              className="text-slate-200 truncate overflow-hidden"
              style={{
                maxWidth: 'var(--label-max-width)',
                display: 'block',
              }}
            >{title}</span>
            <span className="text-slate-400">{billingCycle}</span>
          </div>
        </div>
      </div>
      <div className="flex items-center gap-4">
        <Badge
          variant="secondary"
          className={`
            hidden md:inline-flex
            ${spendType === "subscription" ? "bg-purple-900/50 text-purple-200" : ""}
            ${spendType === "one-time" ? "bg-slate-700/50 text-slate-200" : ""}
          `}
        >
          {spendType}
        </Badge>
        {
          mediaType === "subscription" ? null : (
            <Badge
              variant="secondary"
              className={`
                ${mediaType === "hardware" ? "bg-green-700/50 text-slate-200" : ""}
                ${mediaType === "dlc" ? "bg-orange-700/50 text-slate-200" : ""}
                ${mediaType === "inGamePurchase" ? "bg-blue-600/50 text-slate-200" : ""}
                ${mediaType === "disc" ? "bg-blue-400/50 text-slate-200" : ""}
                ${mediaType === "physical" ? "bg-yellow-400/50 text-slate-200" : ""}
              `}
            >
              {mediaType}
            </Badge>
          )
        }

        <span className="text-slate-200 w-24 text-right">${amount}</span>
        {isPaid && <Check className="h-4 w-4 text-green-500" />}
      </div>
    </div>
  )
}

