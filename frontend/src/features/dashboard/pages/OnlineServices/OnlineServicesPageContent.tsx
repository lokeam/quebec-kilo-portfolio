import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Plus } from 'lucide-react';

// Components
import { OnlineServicesToolbar } from '@/features/dashboard/organisms/OnlineServicesToolbar/OnlineServicesToolbar';
import { SingleOnlineServiceCard } from '@/features/dashboard/organisms/SingleOnlineServiceCard/SingleOnlineServiceCard';
import { OnlineServicesTable } from '@/features/dashboard/organisms/OnlineServicesTable/OnlineServicesTable';

// Utils + Hooks
import { useCardLabelWidth } from '@/features/dashboard/organisms/SingleOnlineServiceCard/useCardLabelWidth';
import { useOnlineServicesStore } from '@/features/dashboard/stores/onlineServicesStore';
import { ViewModes } from '@/features/dashboard/stores/onlineServicesStore';

// Mock Data
import { onlineServicesPageMockData } from './onlineServicesPage.mockdata';
import { useFilteredServices } from '@/features/dashboard/hooks/useFilteredServices';

export function OnlineServicesPageContent() {
  const { viewMode } = useOnlineServicesStore();
  const filteredServices = useFilteredServices(onlineServicesPageMockData);

  useCardLabelWidth({
    selectorAttribute: '[data-card-sentinel]',
    breakpoints: {
      narrow: 320,
      medium: 360
    },
    widths: {
      narrow: '60px',
      medium: '100px',
      wide: '200px'
    }
  });

  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>Online Services</h1>
        </div>
        <div className='flex items-center space-x-2'>
          <Button>
            <Plus className="h-4 w-4" />
            New Service
          </Button>
        </div>
      </PageHeadline>

      <OnlineServicesToolbar />

      {viewMode === ViewModes.TABLE ? (
        <OnlineServicesTable services={filteredServices} />
      ) : (
        <div className={`grid grid-cols-1 gap-4 ${
          viewMode === ViewModes.GRID
            ? 'md:grid-cols-2 2xl:grid-cols-3'
            : 'md:grid-cols-1'
        }`}>
          {filteredServices.length > 0 ? (
            filteredServices.map((service, index) => (
              <SingleOnlineServiceCard
                key={`${service.name}-${index}`}
                {...service}
                isWatchedByResizeObserver={index === 0}
              />
            ))
          ) : (
            <div className="col-span-full">
              <p>No online services found</p>
            </div>
          )}
        </div>
      )}

      <div className="mt-4 space-y-4">
        <Skeleton className="w-full h-[200px]" />
        <Skeleton className="w-full h-[300px]" />
      </div>
    </PageMain>
  );
}
