import { Avatar, AvatarFallback, AvatarImage } from '@/shared/components/ui/avatar';
import { Badge } from '@/shared/components/ui/badge';
import type { ServiceListProps } from './onlineServicesCard.types';

export function ServiceList({ services }: ServiceListProps) {
  return (
    <div className="space-y-6">
      {services.map((service, index) => (
        <div key={index} className="flex items-center gap-4">
          <Avatar className="h-9 w-9">
            <AvatarImage src={service.avatar} alt={`${service.name} avatar`} />
            <AvatarFallback>{service.name.slice(0, 2)}</AvatarFallback>
          </Avatar>
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