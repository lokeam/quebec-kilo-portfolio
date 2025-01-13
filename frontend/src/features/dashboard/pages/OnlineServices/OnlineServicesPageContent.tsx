import { useEffect } from 'react';

// Components
import { OnlineServicesToolbar } from '@/features/dashboard/components/organisms/OnlineServicesToolbar/OnlineServicesToolbar';
import { SingleOnlineServiceCard } from '@/features/dashboard/components/organisms/SingleOnlineServiceCard/SingleOnlineServiceCard';
import { OnlineServicesTable } from '@/features/dashboard/components/organisms/OnlineServicesTable/OnlineServicesTable';
import { AddNewServiceDialog } from '@/features/dashboard/components/organisms/AddNewServiceDialog/AddNewServiceDialog';
import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';
import { ServiceListContainer } from '@/features/dashboard/components/templates/ServiceListContainer';

// Utils + Hooks
import { useCardLabelWidth } from '@/features/dashboard/components/organisms/SingleOnlineServiceCard/useCardLabelWidth';
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';

// Mock Data
import { onlineServicesPageMockData } from './onlineServicesPage.mockdata';
import { useFilteredServices } from '@/features/dashboard/lib/hooks/useFilteredServices';
import { OnlineServicesEmptyPage } from '@/features/dashboard/pages/OnlineServices/OnlineServicesEmptyPage';

export function OnlineServicesPageContent() {
  const { viewMode } = useOnlineServicesStore();
  const filteredServices = useFilteredServices(onlineServicesPageMockData?.services);
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
    setServices(onlineServicesPageMockData?.services);
  }, [setServices])

  console.log('checking filteredServices', filteredServices);

  return (
    <ServiceListContainer
      services={filteredServices || []}
      totalServices={onlineServicesPageMockData?.totalServices || 0}
      viewMode={viewMode}
      title="Online Services"
      EmptyPage={OnlineServicesEmptyPage}
      NoResultsFound={NoResultsFound}
      AddNewDialog={AddNewServiceDialog}
      Toolbar={OnlineServicesToolbar}
      Table={OnlineServicesTable}
      Card={SingleOnlineServiceCard}
    />
  );
}
