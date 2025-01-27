import { useEffect } from 'react';

// Components
import { OnlineServicesToolbar } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServicesToolbar/OnlineServicesToolbar';
import { SingleOnlineServiceCard } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/SingleOnlineServiceCard';
import { OnlineServicesTable } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServicesTable/OnlineServicesTable';
import { AddNewServiceDialog } from '@/features/dashboard/components/organisms/OnlineServicesPage/AddNewServiceDialog/AddNewServiceDialog';
import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';
import { ServiceListContainer } from '@/features/dashboard/components/templates/ServiceListContainer';

// Utils + Hooks
import { useCardLabelWidth } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/useCardLabelWidth';
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';

// Types
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';

// Mock Data
import { onlineServicesPageMockData } from './onlineServicesPage.mockdata';
import { useFilteredServices } from '@/features/dashboard/lib/hooks/useFilteredServices';
import { OnlineServicesEmptyPage } from '@/features/dashboard/pages/OnlineServices/OnlineServicesEmptyPage';
import type { ServiceType } from '@/shared/constants/service.constants';

export function OnlineServicesPageContent() {
  const { viewMode } = useOnlineServicesStore();
  const filteredServices = useFilteredServices(onlineServicesPageMockData?.services.map(service => ({
    ...service,
    type: service.type as ServiceType,
    tier: {
      currentTier: service.tier.name,
      availableTiers: [{ name: service.tier.name, features: service.tier.features }]
    }
  })) as OnlineService[]);
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
    setServices(onlineServicesPageMockData?.services.map(service => ({
      ...service,
      type: service.type as ServiceType,
      tier: {
        currentTier: service.tier.name,
        availableTiers: [{ name: service.tier.name, features: service.tier.features }]
      }
    })) as OnlineService[]);
  }, [setServices])

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
      hasCSSGridLayout={true}
    />
  );
}
