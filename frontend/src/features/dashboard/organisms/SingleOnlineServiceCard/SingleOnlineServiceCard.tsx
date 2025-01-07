import { Card } from "@/shared/components/ui/card"
import { ChevronDown } from 'lucide-react'
import { IconCloudDataConnection } from '@tabler/icons-react';
import SVGLogo from "@/shared/components/ui/LogoMap/LogoMap";
import type { LogoName } from "@/shared/components/ui/LogoMap/LogoMap";

type SingleOnlineServiceCardProps = {
  name: string;
  label: string;
  logo: string;
  tier: string;
  price: string;
  billingCycle: string;
  currency: string;
  isWatchedByResizeObserver: boolean;
}

export function SingleOnlineServiceCard({
  label,
  logo,
  tier,
  price,
  billingCycle,
  isWatchedByResizeObserver,
}: SingleOnlineServiceCardProps) {
  console.log('logo: ', logo);

  const isServiceFree = price === 'FREE';

  // Create a mapping for special cases
  const logoNameMap: Record<string, string> = {
    'greenmanlogo': 'greenman',
    'primegaminglogo': 'prime',
    'netflixgameslogo': 'netflix',
    'geforcelogo': 'nvidia',
    'eaplaylogo': 'ea',
    'metaquestlogo': 'meta',
    'amazonlunalogo': 'luna'
  };

  // Get the correct logo name using the mapping or fallback to simple replacement
  const logoName = logoNameMap[logo] || logo?.replace('logo', '');
  const hasValidLogo = Boolean(logoName);

  return (
    <Card
      className={`w-full max-w-lg min-h-[100px] max-h-[100px] p-4 bg-gradient-to-b from-slate-900 to-slate-950 border-slate-800 ${isWatchedByResizeObserver ? 'w-full' : 'max-w-lg'}`}
      {...(isWatchedByResizeObserver ? { 'data-card-sentinel': true } : {})}
    >
      <div className="flex items-center justify-between h-full">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 text-white flex items-center justify-center">
            {hasValidLogo ? (
              <SVGLogo
                domain="games"
                name={logoName as LogoName<'games'>}
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
            {!isServiceFree && tier && (
              <span className="text-xs text-muted-foreground">{tier}</span>
            )}
          </div>
        </div>
        <div className="flex items-center gap-1 text-sm">
          {!isServiceFree && (
            <>
              <span className="font-medium text-white">{price}</span>
              <span className="text-muted-foreground text-xs">/ {billingCycle}</span>
            </>
          )}
          <ChevronDown className="h-4 w-4 text-muted-foreground ml-1" />
        </div>
      </div>
    </Card>
  )
}

