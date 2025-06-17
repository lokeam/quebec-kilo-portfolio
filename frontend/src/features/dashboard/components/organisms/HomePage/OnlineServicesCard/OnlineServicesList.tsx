import { Badge } from '@/shared/components/ui/badge';
import type { ServiceListProps } from './onlineServicesCard.types';
import SVGLogo from "@/shared/components/ui/LogoMap/LogoMap";
import type { LogoName } from "@/shared/components/ui/LogoMap/LogoMap";
import { IconCloudDataConnection } from '@tabler/icons-react';

export function ServiceList({ digitalLocations }: ServiceListProps) {
  return (
    <div className="space-y-6">
      {digitalLocations.map((digitalLocation, index) => (
        <div key={index} className="flex items-center gap-4">
          <div className="h-9 w-9 flex items-center justify-center">
            {digitalLocation.logo ? (
              <SVGLogo
                domain="games"
                name={digitalLocation.logo as LogoName<'games'>}
                className="w-full h-full object-contain"
              />
            ) : (
              <IconCloudDataConnection className="w-full h-full" />
            )}
          </div>

          <div className="flex flex-1 flex-wrap items-center justify-between">
            <div className="space-y-1">
              <p className="text-md font-medium leading-none">{digitalLocation.name}</p>
              <a className="text-sm text-muted-foreground" href={digitalLocation.url} target="_blank" rel="noopener noreferrer">{digitalLocation.url}</a>
            </div>
            <div className="flex items-center gap-2">
              { digitalLocation.billingCycle && <Badge variant="outline">{digitalLocation.billingCycle}</Badge> }
              <div className="font-medium">{digitalLocation.monthlyFee ? `$${digitalLocation.monthlyFee}` : 'FREE'}</div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
