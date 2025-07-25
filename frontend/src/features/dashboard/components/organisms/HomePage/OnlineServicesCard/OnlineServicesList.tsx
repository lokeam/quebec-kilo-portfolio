// ShadCN UI components
import { Badge } from '@/shared/components/ui/badge';

// Utils
import { DigitalLocationIcon } from '@/features/dashboard/lib/utils/getDigitalLocationIcon';

// Types
import type { ServiceListProps } from './onlineServicesCard.types';

// Icons
import { IconCloudDataConnection } from '@/shared/components/ui/icons';
import { formatServicePrice } from '@/features/dashboard/lib/utils/online-service-status';

export function ServiceList({ digitalLocations }: ServiceListProps) {

  // console.log('OnlineServicesList digitalLocations - ', digitalLocations);

  return (
    <div className="space-y-6">
      {digitalLocations.map((digitalLocation, index) => (
        <div key={index} className="flex items-center gap-4">
          <div className="h-9 w-9 flex items-center justify-center">
            {digitalLocation.name ? (
              <DigitalLocationIcon
                name={digitalLocation.name}
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
              <div className="font-medium">{formatServicePrice(digitalLocation.monthlyFee)}</div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
