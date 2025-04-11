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
import type { ServiceTierName } from '@/features/dashboard/lib/types/online-services/tiers';
import type { ServiceStatusCode } from '@/shared/constants/service.constants';

// Mock Data
import { onlineServicesPageMockData } from './onlineServicesPage.mockdata';
import { useFilteredServices } from '@/features/dashboard/lib/hooks/useFilteredServices';
import { OnlineServicesEmptyPage } from '@/features/dashboard/pages/OnlineServices/OnlineServicesEmptyPage';
import type { ServiceType } from '@/shared/constants/service.constants';

// Type for raw service data
interface RawServiceData {
  id: string;
  name: string;
  label: string;
  status: string;
  logo?: string;
  type: string;
  features: string[];
  createdAt: string;
  updatedAt: string;
  url?: string;
  tier: {
    name: string;
    features: string[];
  };
  billing?: {
    cycle: string;
    fees: {
      monthly: string;
      quarterly: string;
      annual: string;
    };
    renewalDate?: {
      month: string;
      day: string; // Note: string in raw data
    };
    paymentMethod: string;
  };
}

// This is complete bullshit that Typescript forces me to create this casting:
const transformService = (service: RawServiceData): OnlineService => ({
  ...service,
  logo: service.logo || 'default-logo',
  url: service.url || '#',
  type: service.type as ServiceType,
  status: service.status as ServiceStatusCode,
  tier: {
    currentTier: service.tier.name as ServiceTierName,
    availableTiers: [{
      name: service.tier.name as ServiceTierName,
      features: service.tier.features,
      id: `tier-${service.tier.name.toLowerCase().replace(/\s+/g, '-')}`,
      isDefault: true
    }]
  },
  billing: service.billing ? {
    // Ensure required fields have values
    cycle: service.billing.cycle || 'NA',
    fees: service.billing.fees || { monthly: '0', quarterly: '0', annual: '0' },
    paymentMethod: service.billing.paymentMethod || 'Generic',
    // Optional
    renewalDate: service.billing.renewalDate ? {
      month: service.billing.renewalDate.month,
      day: Number(service.billing.renewalDate.day)
    } : undefined
  } : undefined
});

export function OnlineServicesPageContent() {
  const { viewMode } = useOnlineServicesStore();
  const filteredServices = useFilteredServices(
    onlineServicesPageMockData?.services.map(transformService)
  );
  const setServices = useOnlineServicesStore((state) => state.setServices);

  useCardLabelWidth({
    selectorAttribute: '[data-card-sentinel]',
    breakpoints: {
      narrow: 320,
      medium: 360
    },
    widths: {
      narrow: '120px',
      medium: '140px',
      wide: '200px'
    }
  });

  useEffect(() => {
    setServices(onlineServicesPageMockData?.services.map(transformService));
  }, [setServices]);

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
