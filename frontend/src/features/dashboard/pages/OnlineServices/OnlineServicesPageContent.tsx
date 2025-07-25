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

// Skeleton Components
import { OnlineServicesPageSkeleton } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageSkeleton';

// API Hooks and Utilities
import { useShowConditionalIntroToasts } from '@/features/dashboard/hooks/intro-toasts/useShowConditionalIntroToasts';
//import { useStorageAnalytics } from '@/core/api/queries/analyticsData.queries';


import {
  useGetDigitalLocationsBFFResponse,
  useDeleteDigitalLocation
} from '@/core/api/queries/digitalLocation.queries';

import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { useCardLabelWidth } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/useCardLabelWidth';
import { useFilteredServices } from '@/features/dashboard/lib/hooks/useFilteredServices';
import { useOnlineServicesFilters } from '@/features/dashboard/lib/hooks/useOnlineServicesFilters';

// Types
import type { DigitalLocation } from '@/types/domain/digital-location';
import type { FormValues } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServiceForm/OnlineServiceForm';

// Utils
const DEFAULT_FORM_VALUES = {
  name: '',
  isActive: true,
  url: '',
  isSubscriptionService: false,
  billingCycle: '',
  costPerCycle: 0,
  anchorDate: new Date(),
  paymentMethod: '',
};

const transformServiceToFormValues = (service: DigitalLocation): FormValues => {
  // Ensure we have a valid service object
  if (!service) {
    return DEFAULT_FORM_VALUES;
  }

  return {
    name: service.name || '',
    isActive: service.isActive ?? true,
    url: service.url || '',
    isSubscriptionService: service.isSubscription ?? false,
    billingCycle: service.billingCycle || '',
    costPerCycle: service.costPerCycle || 0,
    anchorDate: service.nextPaymentDate
      ? new Date(service.nextPaymentDate)
      : new Date(),
    paymentMethod: service.paymentMethod || '',
  };
};

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
  // Add state for edit mode
  const [addServiceOpen, setAddServiceOpen] = useState<boolean>(false);
  const [editServiceOpen, setEditServiceOpen] = useState<boolean>(false);
  const [serviceBeingEdited, setServiceBeingEdited] = useState<DigitalLocation | null>(null);

  const viewMode = useOnlineServicesStore((state) => state.viewMode);

  // Show intro toast for adding digital locations
  const { data: storageData, isLoading, error } = useGetDigitalLocationsBFFResponse();
  const deleteDigitalLocation = useDeleteDigitalLocation();
  const hasDigitalLocations = Boolean(storageData && storageData.length > 0);
  useShowConditionalIntroToasts(4, !hasDigitalLocations);

  //const storageData = {digitalLocations: []};
  // console.log('[DEBUG] digitalLocationsBFFResponse:', storageData);

  // Get filtered services using the unified hook
  const services = storageData || [];
  const filteredServices = useFilteredServices(services);

  // Get filter options from BFF data
  const filterOptions = useOnlineServicesFilters(services);

  // data: FormValues
  const handleAddService = useCallback(() => {
    // console.log('Add service:', data);
    setAddServiceOpen(false);
  }, []);

  const handleCloseAddDrawer = useCallback(() => {
    setAddServiceOpen(false);
  }, []);

  // Enhanced edit handlers
  const handleEditService = useCallback((service: DigitalLocation) => {
    setServiceBeingEdited(service);
    setEditServiceOpen(true);
  }, []);

  const handleCloseEditDrawer = useCallback(() => {
    setEditServiceOpen(false);
    setServiceBeingEdited(null);
  }, []);

  // data: FormValues
  const handleEditSubmit = useCallback(() => {
    //console.log('Edit service:', data);
    setEditServiceOpen(false);
    setServiceBeingEdited(null);
  }, []);

  const handleDeleteService = useCallback((serviceId: string) => {
    deleteDigitalLocation.mutate(serviceId);
  }, [deleteDigitalLocation]);

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

  if (isLoading) {
    return <OnlineServicesPageSkeleton />
  }

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
              onEdit={handleEditService}
              onDelete={handleDeleteService}
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
          {/* Add Service Drawer */}
          <DrawerContainer
            open={addServiceOpen}
            onOpenChange={setAddServiceOpen}
            triggerAddLocation="Add Digital Service"
            title="Digital Service"
            description="Tell us about the service you want to add"
            triggerBtnIcon="digital"
          >
            <OnlineServiceForm
              onSuccess={handleAddService}
              onClose={handleCloseAddDrawer}
              buttonText="Add Service"
              isEditMode={false}
            />
          </DrawerContainer>

          {/* Edit Service Drawer */}
          <DrawerContainer
            open={editServiceOpen}
            onOpenChange={setEditServiceOpen}
            title="Edit Digital Service"
            description="Update your service details"
          >
            {serviceBeingEdited && (
              <OnlineServiceForm
                onSuccess={handleEditSubmit}
                onClose={handleCloseEditDrawer}
                buttonText="Edit Service"
                initialValues={transformServiceToFormValues(serviceBeingEdited)}
                isEditMode={true}
                serviceId={serviceBeingEdited.id}
              />
            )}
          </DrawerContainer>
        </div>
      </PageHeadline>

      {/* Digital Services Display Section */}
      <div className="mt-6">
        <OnlineServicesToolbar
          paymentMethods={filterOptions.paymentMethods}
          billingCycles={filterOptions.billingCycles}
        />
        <div className="mt-4 space-y-4">
          {renderContent()}
        </div>
      </div>
    </PageMain>
  );
}