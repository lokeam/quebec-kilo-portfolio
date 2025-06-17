//import { Avatar, AvatarFallback, AvatarImage } from '@/shared/components/ui/avatar';
//import type { DigitalStorageService, PhysicalStorageLocation } from './storageLocationsTabCard.types';
import { SublocationIcon } from '@/features/dashboard/lib/utils/getSublocationIcon';
import { DigitalLocationIcon } from '@/features/dashboard/lib/utils/getDigitalLocationIcon';
import type { StorageLocationTabCardPhysicalItem, StorageLocationTabCardDigitalItem } from './StorageLocationsTabCard';

type StorageLocationListProps = {
  isPhysical: boolean;
  services: StorageLocationTabCardPhysicalItem[] | StorageLocationTabCardDigitalItem[];
};

const physicalLocationListItem = (service: StorageLocationTabCardPhysicalItem, index: number) => {
  return (
    <div key={index} className="flex items-center gap-4">
    <SublocationIcon type={service.sublocationType} bgColor={service.parentLocationBgColor || 'gray'} />
      <div className="flex flex-1 flex-wrap items-center justify-between">
        <div className="space-y-1">
          <p className="text-sm font-medium leading-none">{service.sublocationName}</p>
        </div>
        <div className="font-medium">{service.storedItems} items stored</div>
      </div>
    </div>
  )
};

const digitalLocationListItem = (service: StorageLocationTabCardDigitalItem, index: number) => {
  return (
    <div key={index} className="flex items-center gap-4">
      <DigitalLocationIcon name={service.logo} className="w-9 h-9" />
      <div className="flex flex-1 flex-wrap items-center justify-between">
        <div className="space-y-1">
          <p className="text-sm font-medium leading-none">{service.name}</p>
        </div>
        <div className="font-medium">{service.storedItems} items stored</div>
      </div>
    </div>
  )
};

export function StorageLocationList({ services, isPhysical = false }: StorageLocationListProps) {

  if (isPhysical) {
    return (
      <div className="space-y-6">
        {(services as StorageLocationTabCardPhysicalItem[]).map((service, index) =>
          physicalLocationListItem(service, index)
        )}
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {(services as StorageLocationTabCardDigitalItem[]).map((service, index) =>
        digitalLocationListItem(service, index)
      )}
    </div>
  );
}
