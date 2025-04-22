import { useCallback, useState, useEffect } from 'react';

// Template Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';

// Components
import { SingleOnlineServiceCard } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/SingleOnlineServiceCard';
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';
import { OnlineServicesToolbar } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServicesToolbar/OnlineServicesToolbar';
import { OnlineServicesTable } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServicesTable/OnlineServicesTable';

// API Hooks and Utilities
import { useDigitalLocations } from '@/core/api/hooks/useDigitalLocations';
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { useCardLabelWidth } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/useCardLabelWidth';

// Types
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import type { ServiceTierName } from '@/features/dashboard/lib/types/online-services/tiers';
import type { ServiceStatusCode, ServiceType } from '@/shared/constants/service.constants';
import type { DigitalLocation } from '@/features/dashboard/lib/types/media-storage/digital-location.types';
import { OnlineServiceForm } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServiceForm/OnlineServiceForm';

// Helper function to transform digital location data to online service format
const transformDigitalLocationToService = (location: DigitalLocation): OnlineService => ({
  id: location.id,
  name: location.name,
  logo: 'default-logo',
  url: location.url || '#',
  type: location.service_type as ServiceType,
  status: location.is_active ? 'active' as ServiceStatusCode : 'inactive' as ServiceStatusCode,
  features: [],
  label: location.name,
  createdAt: location.created_at,
  updatedAt: location.updated_at,
  isSubscriptionService: !!location.subscription,
  tier: {
    currentTier: 'Basic' as ServiceTierName,
    availableTiers: [{
      name: 'Basic' as ServiceTierName,
      features: [],
      id: `tier-basic`,
      isDefault: true
    }]
  },
  billing: location.subscription ? {
    cycle: location.subscription.billing_cycle || 'NA',
    fees: {
      monthly: location.subscription.cost_per_cycle.toString() || '0'
    },
    paymentMethod: location.subscription.payment_method || 'Generic',
    renewalDate: {
      month: new Date(location.subscription.next_payment_date).toLocaleString('default', { month: 'long' }),
      day: new Date(location.subscription.next_payment_date).getDate()
    }
  } : undefined
});

export function OnlineServicesPageContent() {
  const [addServiceOpen, setAddServiceOpen] = useState<boolean>(false);
  const [editServiceOpen, setEditServiceOpen] = useState<boolean>(false);
  const setServices = useOnlineServicesStore((state) => state.setServices);
  const services = useOnlineServicesStore((state) => state.services);
  const viewMode = useOnlineServicesStore((state) => state.viewMode);

  // Set up card label width for responsive design
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

  // Fetch digital locations using our hook
  const { data: digitalLocations, isLoading, error } = useDigitalLocations();

  // Transform digital locations to online services format and update store
  useEffect(() => {
    if (digitalLocations) {
      const transformedServices = digitalLocations.map(transformDigitalLocationToService);
      setServices(transformedServices);
    }
  }, [digitalLocations, setServices]);

  const handleCloseAddDrawer = useCallback(() => {
    setAddServiceOpen(false);
  }, []);

  // Render content based on the selected view mode
  const renderContent = () => {
    if (isLoading) {
      return (
        <div className="p-4 border rounded-md">
          <p className="text-gray-500">Loading digital services...</p>
        </div>
      );
    }

    if (error) {
      return (
        <div className="p-4 border border-red-300 bg-red-50 rounded-md">
          <p className="text-red-500">Error loading digital services</p>
        </div>
      );
    }

    if (services.length === 0) {
      return (
        <div className="p-4 border rounded-md">
          <p className="text-gray-500">No digital services found. Add a service to get started.</p>
        </div>
      );
    }

    if (viewMode === 'table') {
      return <OnlineServicesTable services={services} />;
    }

    // Grid or list view (cards)
    return (
      <div className={`p-4 border rounded-md ${
        viewMode === 'grid' ? 'grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4' : ''
      }`}>
        <h2 className="text-lg font-semibold">Digital Services</h2>
        <p className="text-gray-500">{services.length} services found</p>

        <ul className={`mt-2 ${viewMode === 'grid' ? '' : 'space-y-2'}`}>
          {services.map((service, index) => (
            <SingleOnlineServiceCard
              key={`${service.id}-${index}`}
              {...service}
              isWatchedByResizeObserver={index === 0}
            />
          ))}
        </ul>
      </div>
    );
  };

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
            <OnlineServiceForm onClose={handleCloseAddDrawer} />
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

        {/* Loading, Error and Data States */}
        <div className="mt-4 space-y-4">
          {renderContent()}
        </div>
      </div>
    </PageMain>
  );
}