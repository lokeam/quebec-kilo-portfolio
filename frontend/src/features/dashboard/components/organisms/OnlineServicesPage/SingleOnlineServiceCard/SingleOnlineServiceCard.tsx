import { memo } from 'react';
import { Card } from "@/shared/components/ui/card"
import { Power } from 'lucide-react'
import { IconCloudDataConnection, IconCalendarDollar } from '@tabler/icons-react';
import SVGLogo from "@/shared/components/ui/LogoMap/LogoMap";
import type { LogoName } from "@/shared/components/ui/LogoMap/LogoMap";
import { Badge } from "@/shared/components/ui/badge";
import {
  isServiceFree,
  formatCurrency,
  isRenewalMonth,
} from '@/features/dashboard/lib/utils/online-service-status';
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import { SERVICE_STATUS_CODES } from '@/shared/constants/service.constants';

interface SingleOnlineServiceCardProps extends OnlineService {
  isWatchedByResizeObserver?: boolean;
  onClick?: () => void;
}

export const SingleOnlineServiceCard = memo(({
  label,
  logo,
  status,
  billing,
  tier,
  isWatchedByResizeObserver,
  onClick
}: SingleOnlineServiceCardProps) => {
  const currentTierDetails = tier?.availableTiers?.find(t =>
    t?.name?.toLowerCase() === tier?.currentTier?.toLowerCase()
  );

  const hasValidLogo = Boolean(logo);
  const isFree = isServiceFree({ billing } as OnlineService);
  const showRenewalBadge = status !== SERVICE_STATUS_CODES.ACTIVE &&
    !isFree &&
    isRenewalMonth({ billing } as OnlineService);

  return (
    <Card
      className="flex relative cursor-pointer w-full min-h-[100px] max-h-[100px] p-4 bg-gradient-to-b from-slate-900 to-slate-950 border-slate-800"
      {...(isWatchedByResizeObserver ? { 'data-card-sentinel': true } : {})}
      onClick={onClick}
    >
      <div className="flex items-center justify-between min-w-0 w-full">
        <div className="flex items-center gap-3 min-w-0">
          <div className="w-10 h-10 shrink-0 text-white flex items-center justify-center">
            {hasValidLogo ? (
              <SVGLogo
                domain="games"
                name={logo as LogoName<'games'>}
                className="w-full h-full object-contain"
              />
            ) : (
              <IconCloudDataConnection className="w-full h-full" />
            )}
          </div>
          <div className="flex flex-col">
            <span
              className="font-medium text-sm text-white truncate overflow-hidden"
              style={{
                maxWidth: 'var(--label-max-width)',
                display: 'block',
              }}
            >
              {label}
            </span>
            {!isFree && currentTierDetails && (
              <div className="flex flex-col">
                <span className="text-xs text-muted-foreground">
                {currentTierDetails.name}
                </span>
                {tier.maxDevices && (
                  <span className="text-xs text-muted-foreground">
                    Up to {tier.maxDevices} devices
                  </span>
                )}
              </div>
            )}
          </div>
        </div>
        <div className="flex items-center gap-1 text-sm shrink-0">
          {isFree ? (
              <span className="font-medium text-white">
                FREE
              </span>
            ) : (
            <>
              <span className="font-medium text-white">
                {formatCurrency(billing.fees.monthly)}
              </span>
              <span className="text-muted-foreground text-xs">/ 1 mo</span>
            </>
          )}
          {status === SERVICE_STATUS_CODES.ACTIVE && (
            <Power className="h-5 w-5 ml-1 text-emerald-500" />
          )}
          {showRenewalBadge && (
            <Badge
              variant="default"
              className="ml-1 bg-red-900 absolute top-3 right-2"
            >
              <IconCalendarDollar className="h-5 w-5 ml-1 text-white" />
              <span className="ml-1">Renews this month</span>
            </Badge>
          )}
        </div>
      </div>
    </Card>
  );
});

SingleOnlineServiceCard.displayName = 'SingleOnlineServiceCard';
