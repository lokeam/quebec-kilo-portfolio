import { memo } from 'react';
import { Card } from "@/shared/components/ui/card"
import { Power } from 'lucide-react'
import { IconCloudDataConnection, IconCalendarDollar } from '@tabler/icons-react';
import SVGLogo from "@/shared/components/ui/LogoMap/LogoMap";
import type { LogoName } from "@/shared/components/ui/LogoMap/LogoMap";
import { Badge } from "@/shared/components/ui/badge";

interface SingleOnlineServiceCardProps {
  name: string;
  label: string;
  logo: string;
  tierName: string;
  monthlyFee: string;
  isActive: boolean;
  renewalMonth: string;
  isWatchedByResizeObserver?: boolean;
  onClick?: () => void;
}

export const SingleOnlineServiceCard = memo(({
  label,
  logo,
  tierName,
  monthlyFee,
  renewalMonth,
  isActive,
  isWatchedByResizeObserver,
}: SingleOnlineServiceCardProps) => {

  const isServiceFree = monthlyFee === 'FREE';
  const hasValidLogo = Boolean(logo);
  const date = new Date();
  const currentMonth = date.toLocaleString('default', { month: 'long' });
  const isRenewalMonth = currentMonth === renewalMonth;

  return (
    <Card
      className={`flex relative cursor-pointer w-full min-h-[100px] max-h-[100px] p-4 bg-gradient-to-b from-slate-900 to-slate-950 border-slate-800`}
      {...(isWatchedByResizeObserver ? { 'data-card-sentinel': true } : {})}
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
            >{label}</span>
            {!isServiceFree && (
              <span className="text-xs text-muted-foreground">{tierName || 'Standard subscription'}</span>
            )}
          </div>
        </div>
        <div className="flex items-center gap-1 text-sm shrink-0">
          {!isServiceFree && (
            <>
              <span className="font-medium text-white">{monthlyFee}</span>
              <span className="text-muted-foreground text-xs">/ 1 mo</span>
            </>
          )}
          {
            isActive && (
              <Power className="h-5 w-5 ml-1 text-green-500" />
            )
          }
          {
            !isActive && isRenewalMonth && (
              <Badge variant="default" className="ml-1 bg-red-900 absolute top-3 right-2">
                <IconCalendarDollar className="h-5 w-5 ml-1 text-white" />
                <span className="ml-1">Renews this month</span>
              </Badge>

            )
          }
        </div>
      </div>
    </Card>
  )
});
