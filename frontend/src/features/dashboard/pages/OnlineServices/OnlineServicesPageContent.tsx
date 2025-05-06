import { useCallback, useState } from 'react';

// Template Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';

// Components
import { SingleOnlineServiceCard } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/SingleOnlineServiceCard';
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';
import { OnlineServicesToolbar } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServicesToolbar/OnlineServicesToolbar';
import { OnlineServicesTable } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServicesTable/OnlineServicesTable';
import { OnlineServiceForm } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServiceForm/OnlineServiceForm';
import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';

// API Hooks and Utilities
import { useGetAllDigitalLocations, useDeleteDigitalLocation } from '@/core/api/queries/useDigitalLocations';
import { useStorageAnalytics } from '@/core/api/hooks/useAnalyticsData';
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { useCardLabelWidth } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/useCardLabelWidth';
import { useFilteredServices } from '@/features/dashboard/lib/hooks/useFilteredServices';

// Types
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import type { SelectableItem } from '@/shared/components/ui/ResponsiveCombobox/ResponsiveCombobox';
import { SERVICE_STATUS_CODES } from '@/shared/constants/service.constants';

export function OnlineServicesPageContent() {
  const [addServiceOpen, setAddServiceOpen] = useState<boolean>(false);
  const [editServiceOpen, setEditServiceOpen] = useState<boolean>(false);
  const [selectedService, setSelectedService] = useState<OnlineService | null>(null);
  const viewMode = useOnlineServicesStore((state) => state.viewMode);

  // Fetch digital locations using our hook - now returns OnlineService[] directly
  const { data: services = [], isLoading, error } = useGetAllDigitalLocations();

  // Get filtered services using the unified hook
  const filteredServices = useFilteredServices(services);

  // Set up the delete mutation
  const deleteServiceMutation = useDeleteDigitalLocation();

  // Handler to delete a service
  const handleDeleteService = useCallback((serviceId: string) => {
    deleteServiceMutation.mutate(serviceId);
  }, [deleteServiceMutation]);

  // Handler to edit a service
  const handleEditService = useCallback((service: OnlineService) => {
    setSelectedService(service);
    setEditServiceOpen(true);
  }, []);

  // Handler to close edit drawer
  const handleCloseEditDrawer = useCallback(() => {
    setEditServiceOpen(false);
    setSelectedService(null);
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

  // Fetch storage analytics in parallel
  const {
    data: analyticsData,
    isLoading: analyticsLoading,
    error: analyticsError
  } = useStorageAnalytics();

  const handleCloseAddDrawer = useCallback(() => {
    setAddServiceOpen(false);
  }, []);

  // Analytics summary component
  const renderAnalyticsSummary = () => {
    if (analyticsLoading) {
      return <div className="text-sm text-gray-500">Loading analytics...</div>;
    }

    if (analyticsError) {
      return null; // Don't show anything if analytics fail
    }

    const storageStats = analyticsData?.storage;

    if (!storageStats) {
      return null;
    }

    return (
      <div className="mb-4 p-4 bg-blue-50 border border-blue-100 rounded-md">
        <h3 className="text-sm font-semibold text-blue-800">Storage Analytics</h3>
        <div className="mt-2 grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-600">Digital Locations:</span>{' '}
            <span className="font-medium">{storageStats.totalDigitalLocations}</span>
          </div>
          <div>
            <span className="text-gray-600">Physical Locations:</span>{' '}
            <span className="font-medium">{storageStats.totalPhysicalLocations}</span>
          </div>
        </div>
      </div>
    );
  };

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

    if (filteredServices.length === 0) {
      return <NoResultsFound />;
    }

    if (viewMode === 'table') {
      return <OnlineServicesTable services={filteredServices} onDelete={handleDeleteService} onEdit={handleEditService} />;
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
              data-card-sentinel={index === 0}
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

        <div className='flex items-center space-x-2'>
          <DrawerContainer
            open={addServiceOpen}
            onOpenChange={setAddServiceOpen}
            triggerAddLocation="Add Digital Service"
            title="Digital Service"
            description="Tell us about your digital service."
          >
            <OnlineServiceForm onClose={handleCloseAddDrawer} />
          </DrawerContainer>
        </div>
      </PageHeadline>

      {/* Analytics Summary */}
      {renderAnalyticsSummary()}

      {/* Digital Services Display Section */}
      <div className="mt-6">
        <OnlineServicesToolbar />
        <div className="mt-4 space-y-4">
          {renderContent()}
        </div>
      </div>

      {/* Edit Service Drawer */}
      <DrawerContainer
        open={editServiceOpen}
        onOpenChange={setEditServiceOpen}
        title="Edit Digital Service"
        description="Update your digital service details."
      >
        {selectedService && (
          <OnlineServiceForm
            serviceData={{
              id: selectedService.id,
              service: selectedService,
              expenseType: selectedService.billing?.cycle,
              cost: selectedService.billing?.fees?.monthly ? parseFloat(selectedService.billing.fees.monthly.replace('$', '')) : 0,
              billingPeriod: selectedService.billing?.cycle,
              nextPaymentDate: selectedService.billing?.renewalDate ? new Date(`${selectedService.billing.renewalDate.month} ${selectedService.billing.renewalDate.day}`) : undefined,
              paymentMethod: typeof selectedService.billing?.paymentMethod === 'string'
                ? { id: selectedService.billing.paymentMethod, displayName: selectedService.billing.paymentMethod }
                : selectedService.billing?.paymentMethod as SelectableItem | undefined,
              is_active: selectedService.status === SERVICE_STATUS_CODES.ACTIVE
            }}
            isEditing={true}
            onClose={handleCloseEditDrawer}
            onDelete={handleDeleteService}
          />
        )}
      </DrawerContainer>
    </PageMain>
  );
}