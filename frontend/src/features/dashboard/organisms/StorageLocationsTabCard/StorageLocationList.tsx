import { Avatar, AvatarFallback, AvatarImage } from '@/shared/components/ui/avatar';
import type { DigitalStorageService, PhysicalStorageLocation } from './storageLocationsTabCard.types';

type StorageLocationListProps = {
  services: DigitalStorageService[] | PhysicalStorageLocation[];
};

export function StorageLocationList({ services }: StorageLocationListProps) {
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
              <p className="text-sm font-medium leading-none">{service.name}</p>
            </div>
            <div className="font-medium">{service.itemsStored} items stored</div>
          </div>
        </div>
      ))}
    </div>
  );
}
