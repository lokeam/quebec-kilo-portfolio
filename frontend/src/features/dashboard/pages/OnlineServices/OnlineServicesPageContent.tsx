import { useState, useCallback } from 'react';

// Template Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';

// Components
import { SingleOnlineServiceCard } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/SingleOnlineServiceCard';
import { OnlineServicesToolbar } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServicesToolbar/OnlineServicesToolbar';
import { OnlineServicesTable } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServicesTable/OnlineServicesTable';
import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';
import { OnlineServiceForm } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServiceForm/OnlineServiceForm';

// API Hooks and Utilities
import { useStorageAnalytics } from '@/core/api/queries/analyticsData.queries';
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { useCardLabelWidth } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/useCardLabelWidth';
import { useFilteredServices } from '@/features/dashboard/lib/hooks/useFilteredServices';

// Types
import type { DigitalLocation } from '@/types/domain/online-service';

// Skeleton Components
const TableSkeleton = () => (
  <div className="w-full">
    <div className="h-[72px] border rounded-md mb-2 bg-slate-100 animate-pulse" />
    <div className="space-y-2">
      {[1, 2, 3].map((i) => (
        <div key={i} className="h-[72px] border rounded-md bg-slate-100 animate-pulse" />
      ))}
    </div>
  </div>
);

const CardSkeleton = () => (
  <div className="grid grid-cols-1 gap-4 md:grid-cols-2 2xl:grid-cols-3">
    {[1, 2, 3].map((i) => (
      <div
        key={i}
        className="h-[100px] border rounded-md bg-gradient-to-b from-slate-100 to-slate-200 animate-pulse"
      />
    ))}
  </div>
);

export function OnlineServicesPageContent() {
  const [addServiceOpen, setAddServiceOpen] = useState<boolean>(false);
  const viewMode = useOnlineServicesStore((state) => state.viewMode);

  // Fetch digital locations using analytics
  const { data: storageData, isLoading, error } = useStorageAnalytics();

  // Get filtered services using the unified hook
  const services = storageData?.digitalLocations || [];
  const filteredServices = useFilteredServices(services);


  const handleAddService = useCallback((service: DigitalLocation) => {
    console.log('Add service:', service);
    setAddServiceOpen(false);
  }, []);

  const handleCloseAddDrawer = useCallback((data: DigitalLocation) => {
    console.log("Closing drawer after adding service: ", data)
    setAddServiceOpen(false);
  }, [])

  // Handler to delete a service
  const handleDeleteService = useCallback((serviceId: string) => {
    // TODO: Implement delete functionality
    console.log('Delete service:', serviceId);
  }, []);

  // Handler to edit a service
  const handleEditService = useCallback((service: DigitalLocation) => {
    // TODO: Implement edit functionality
    console.log('Edit service:', service);
  }, []);

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

  // Render content based on the selected view mode
  const renderContent = () => {
    if (isLoading) {
      return viewMode === 'table' ? <TableSkeleton /> : <CardSkeleton />;
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

    if (filteredServices.length === 0) {
      return <NoResultsFound />;
    }

    if (viewMode === 'table') {
      return (
        <OnlineServicesTable
          services={filteredServices}
          onDelete={handleDeleteService}
          onEdit={handleEditService}
        />
      );
    }

    return (
      <div className="p-4 border rounded-md">
        <h2 className="text-lg font-semibold">Digital Services</h2>
        <p className="text-gray-500 mb-4">{
          filteredServices.length === 1
            ? '1 service found'
            : `${filteredServices.length} services found`
        }</p>
        <div className={`grid grid-cols-1 gap-4 ${
          viewMode === 'grid' ? 'md:grid-cols-2 2xl:grid-cols-3' : ''
        }`}>
          {filteredServices.map((service, index) => (
            <SingleOnlineServiceCard
              key={`${service.id}-${index}`}
              service={service}
              onDelete={handleDeleteService}
              onEdit={() => handleEditService(service)}
              isWatchedByResizeObserver={index === 0}
            />
          ))}
        </div>
      </div>
    );
  };

  return (
    <PageMain>
      <PageHeadline>
        <div className="flex items-center">
          <h1 className='text-2xl font-bold tracking-tight'>Online Services</h1>
        </div>

        <div className="flex items-center space-x-2">
          <DrawerContainer
            open={addServiceOpen}
            onOpenChange={setAddServiceOpen}
            triggerAddLocation="Add Digital Service"
            title="Digital Service"
            description="Tell us about the service you want to add"
          >
            <OnlineServiceForm
              onSuccess={handleAddService}
              onClose={handleCloseAddDrawer}
              defaultValues={{
                service: null,
                cost: 0,
                is_active: true,
              }}
              buttonText="Add Service"
            />

          </DrawerContainer>

        </div>
      </PageHeadline>

      {/* Digital Services Display Section */}
      <div className="mt-6">
        <OnlineServicesToolbar />
        <div className="mt-4 space-y-4">
          {renderContent()}
        </div>
      </div>
    </PageMain>
  );
}