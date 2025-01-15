import { Badge } from '@/shared/components/ui/badge';
import type { ServiceListProps } from './onlineServicesCard.types';
import SVGLogo from "@/shared/components/ui/LogoMap/LogoMap";
import type { LogoName } from "@/shared/components/ui/LogoMap/LogoMap";
import { IconCloudDataConnection } from '@tabler/icons-react';

export function ServiceList({ services }: ServiceListProps) {
  return (
    <div className="space-y-6">
      {services.map((service, index) => (
        <div key={index} className="flex items-center gap-4">
          <div className="h-9 w-9 flex items-center justify-center">
            {service.logo ? (
              <SVGLogo
                domain="games"
                name={service.logo as LogoName<'games'>}
                className="w-full h-full object-contain"
              />
            ) : (
              <IconCloudDataConnection className="w-full h-full" />
            )}
          </div>

          <div className="flex flex-1 flex-wrap items-center justify-between">
            <div className="space-y-1">
              <p className="text-md font-medium leading-none">{service.name}</p>
              <a className="text-sm text-muted-foreground" href={service.url} target="_blank" rel="noopener noreferrer">{service.url}</a>
            </div>
            <div className="flex items-center gap-2">
              { service.plan && <Badge variant="outline">{service.plan}</Badge> }
              <div className="font-medium">{service.monthlyFee}</div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
