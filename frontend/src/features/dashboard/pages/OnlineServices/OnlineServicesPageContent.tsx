import { useState, useEffect } from 'react';

// Template Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';

// Components
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';
import { OnlineServicesToolbar } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServicesToolbar/OnlineServicesToolbar';

// Mock Data and Utilities
import { onlineServicesPageMockData } from './onlineServicesPage.mockdata';
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';

// Types
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import type { ServiceTierName } from '@/features/dashboard/lib/types/online-services/tiers';
import type { ServiceStatusCode, ServiceType } from '@/shared/constants/service.constants';
import { OnlineServiceForm } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServiceForm/OnlineServiceForm';

// Add this interface above the transformService function
interface RawOnlineService {
  id: string;
  name: string;
  logo?: string;
  url?: string;
  type: string;
  status: string;
  tier: {
    name: string;
    features: string[];
  };
  billing?: {
    cycle?: string;
    fees?: {
      monthly?: string;
      quarterly?: string;
      annual?: string;
    };
    paymentMethod?: string;
    renewalDate?: {
      month: string;
      day: string;
    };
  };
}

// Helper function to transform raw service data
const transformService = (service: RawOnlineService): OnlineService => ({
  ...service,
  logo: service.logo || 'default-logo',
  url: service.url || '#',
  type: service.type as ServiceType,
  status: service.status as ServiceStatusCode,
  features: service.tier.features || [],
  label: service.name,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
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
    cycle: service.billing.cycle || 'NA',
    fees: {
      monthly: service.billing.fees?.monthly || '0'
    },
    paymentMethod: service.billing.paymentMethod || 'Generic',
    renewalDate: service.billing.renewalDate ? {
      month: service.billing.renewalDate.month,
      day: Number(service.billing.renewalDate.day)
    } : undefined
  } : undefined
});

export function OnlineServicesPageContent() {
  const [addServiceOpen, setAddServiceOpen] = useState<boolean>(false);
  const [editServiceOpen, setEditServiceOpen] = useState<boolean>(false);
  const setServices = useOnlineServicesStore((state) => state.setServices);

  // Transform and load services to the store
  useEffect(() => {
    setServices(onlineServicesPageMockData?.services.map(transformService));
  }, [setServices]);

  return (
    <PageMain>
      <PageHeadline>
        <div className="flex items-center">
          <h1 className='text-2xl font-bold tracking-tight'>Online Services</h1>
        </div>

        <div className='flex items-center space-x-2'>
          {/* Add Digital Service Button */}
          <DrawerContainer
            open={addServiceOpen}
            onOpenChange={setAddServiceOpen}
            triggerAddLocation="Add Digital Service"
            title="Digital Service"
            description="Tell us about your digital service."
          >
            {/* Replace with actual form component when available */}
            <OnlineServiceForm />
          </DrawerContainer>

          {/* Edit Digital Service Button */}
          <DrawerContainer
            open={editServiceOpen}
            onOpenChange={setEditServiceOpen}
            triggerEditLocation="Edit Digital Services"
            title="Edit Digital Services"
            description="Edit your digital services"
          >
            {/* Replace with actual service list component when available */}
            <div>Digital Service List will go here</div>
          </DrawerContainer>
        </div>
      </PageHeadline>

      {/* Digital Services Display Section */}
      <div className="mt-6">
        {/* Add a toolbar for filtering/searching services */}
        <OnlineServicesToolbar />

        {/* Add accordion or other display components for services */}
        <div className="mt-4 space-y-4">
          {/* Replace with actual accordion component when available */}
          <div className="p-4 border rounded-md">
            <h2 className="text-lg font-semibold">Digital Services</h2>
            <p className="text-gray-500">Your digital services will be displayed here</p>
          </div>
        </div>
      </div>
    </PageMain>
  );
}