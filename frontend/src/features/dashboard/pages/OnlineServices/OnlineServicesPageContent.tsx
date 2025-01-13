import { useEffect } from 'react';
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';

// Components
import { OnlineServicesToolbar } from '@/features/dashboard/components/organisms/OnlineServicesToolbar/OnlineServicesToolbar';
import { SingleOnlineServiceCard } from '@/features/dashboard/components/organisms/SingleOnlineServiceCard/SingleOnlineServiceCard';
import { OnlineServicesTable } from '@/features/dashboard/components/organisms/OnlineServicesTable/OnlineServicesTable';
import { AddNewServiceDialog } from '@/features/dashboard/components/organisms/AddNewServiceDialog/AddNewServiceDialog';

// Utils + Hooks
import { useCardLabelWidth } from '@/features/dashboard/components/organisms/SingleOnlineServiceCard/useCardLabelWidth';
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { ViewModes } from '@/features/dashboard/lib/stores/onlineServicesStore';

// Mock Data
import { onlineServicesPageMockData } from './onlineServicesPage.mockdata';
import { useFilteredServices } from '@/features/dashboard/lib/hooks/useFilteredServices';
import { OnlineServicesEmptyPage } from '@/features/dashboard/pages/OnlineServices/OnlineServicesEmptyPage';

export function OnlineServicesPageContent() {
  const { viewMode } = useOnlineServicesStore();
  const filteredServices = useFilteredServices(onlineServicesPageMockData);
  const setServices = useOnlineServicesStore((state) => state.setServices);

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

  useEffect(() => {
    setServices(onlineServicesPageMockData);
  }, [setServices])

  console.log('checking filteredServices', filteredServices);

  // If there are no services, show the empty page
  if (filteredServices.length === 0) {
    return (
      <PageMain>
        <PageHeadline>
          <div className='flex items-center'>
            <h1 className='text-2xl font-bold tracking-tight'>Online Services</h1>
          </div>
        </PageHeadline>
        <OnlineServicesEmptyPage />
      </PageMain>
    );
  }

  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>Online Services</h1>
        </div>
        <div className='flex items-center space-x-2'>
          <AddNewServiceDialog />
        </div>
      </PageHeadline>

      <OnlineServicesToolbar />

      {viewMode === ViewModes.TABLE ? (
        <OnlineServicesTable services={filteredServices} />
      ) : (
        <div className={`grid grid-cols-1 gap-4 ${
          viewMode === ViewModes.GRID
            ? 'md:grid-cols-2 2xl:grid-cols-3'
            : ''
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
    </PageMain>
  );
}
