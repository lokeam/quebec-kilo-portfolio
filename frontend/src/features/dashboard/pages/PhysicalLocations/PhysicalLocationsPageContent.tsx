import { useState, useCallback } from 'react';

// Template Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';

// Components
import { SinglePhysicalLocationCard } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/SinglePhysicalLocationCard/SinglePhysicalLocationCard';
import { SingleSublocationCard } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/SingleSublocationCard/SingleSublocationCard';
import { PhysicalLocationsToolbar } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/PhysicalLocationsToolbar/PhysicalLocationsToolbar';
import { PhysicalLocationsTable } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/PhysicalLocationsTable/PhysicalLocationsTable';
import { PhysicalLocationForm } from '@/features/dashboard/components/organisms/MediaStoragePage/PhysicalLocationForm/PhysicalLocationForm';
import { SublocationForm } from '@/features/dashboard/components/organisms/SublocationForm/SublocationForm';

import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';

// API Hooks and Utilities
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { useCardLabelWidth } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/useCardLabelWidth';
import { usePhysicalLocationFilters } from '@/features/dashboard/hooks/usePhysicalLocationFilters';

// Types
import type { LocationsBFFPhysicalLocationResponse, LocationsBFFSublocationResponse } from '@/types/domain/physical-location';
import { useGetPhysicalLocationsBFFResponse, useDeletePhysicalLocation } from '@/core/api/queries/physicalLocation.queries';

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

export function PhysicalLocationsPageContent() {
  const [addPhysicalLocationOpen, setAddPhysicalLocationOpen] = useState<boolean>(false);
  const [addSublocationOpen, setAddSublocationOpen] = useState<boolean>(false);
  const [editServiceOpen, setEditServiceOpen] = useState<boolean>(false);
  const [serviceBeingEdited, setServiceBeingEdited] = useState<LocationsBFFPhysicalLocationResponse | LocationsBFFSublocationResponse | null>(null);
  const [selectedParentLocation, setSelectedParentLocation] = useState<LocationsBFFPhysicalLocationResponse | null>(null);
  const [isSelectingParentLocation, setIsSelectingParentLocation] = useState<boolean>(false);

  const viewMode = useOnlineServicesStore((state) => state.viewMode);

  // Fetch physical locations using BFF
  const { data: storageData, isLoading, error } = useGetPhysicalLocationsBFFResponse();

  // Get filter options from BFF data
  const filterOptions = usePhysicalLocationFilters(storageData);

  // Enhanced edit handlers
  const handleEditService = useCallback((location: LocationsBFFPhysicalLocationResponse | LocationsBFFSublocationResponse) => {
    setServiceBeingEdited(location);
    setEditServiceOpen(true);
  }, []);

  // Sublocation creation when clicking on a SinglePhysicalLocationCard's Add Sublocation button
  const handleAddSublocation = useCallback((location: LocationsBFFPhysicalLocationResponse) => {
    console.log('Adding sublocation to:', {
      physicalLocationId: location.physicalLocationId,
      name: location.name,
      physicalLocationType: location.physicalLocationType,
      bgColor: location.bgColor
    });
    setSelectedParentLocation(location);
    setAddSublocationOpen(true);
  }, []);

  // Sublocation creation when clicking Page Content's Add Sublocation button
  const handleStartSublocationCreation = useCallback(() => {
    setIsSelectingParentLocation(true);
    setSelectedParentLocation(null);
    setAddSublocationOpen(true);
  }, []);

  const handleParentLocationSelect = useCallback((location: LocationsBFFPhysicalLocationResponse) => {
    setSelectedParentLocation(location);
    setIsSelectingParentLocation(false);
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

  // Handle form submission success
  const handleFormSuccess = useCallback(() => {
    setAddPhysicalLocationOpen(false);
    setEditServiceOpen(false);
    setAddSublocationOpen(false);
    setServiceBeingEdited(null);
    setSelectedParentLocation(null);
    setIsSelectingParentLocation(false);
  }, []);

  // Add the mutation hook near the top of the component
  const deleteMutation = useDeletePhysicalLocation();

  const handleDeleteLocation = useCallback((locationId: string) => {
    deleteMutation.mutate(locationId);
  }, [deleteMutation]);

  // Render content based on the selected view mode
  const renderContent = () => {
    if (isLoading) {
      return viewMode === 'table' ? <TableSkeleton /> : <CardSkeleton />;
    }

    if (error) {
      return (
        <div className="p-4 border border-red-300 bg-red-50 rounded-md">
          <p className="text-red-500">Error loading physical location data</p>
        </div>
      );
    }

    if (!storageData?.physicalLocations || storageData.physicalLocations.length === 0) {
      return (
        <div className="p-4 border rounded-md">
          <p className="text-gray-500">No physical location found. Add a location to get started.</p>
        </div>
      );
    }

    if (viewMode === 'table') {
      return (
        <PhysicalLocationsTable
          sublocationRows={storageData.sublocations}
          onEdit={handleEditService}
        />
      );
    }

    return (
      <>
        <div className="p-4 border rounded-md">
          <h2 className="text-lg font-semibold">Physical Locations</h2>
          <p className="text-gray-500 mb-4">{
            storageData.physicalLocations.length === 1
              ? '1 location found'
              : `${storageData.physicalLocations.length} locations found`
          }</p>
          <div className={`grid grid-cols-1 gap-4 ${
            viewMode === 'grid' ? 'md:grid-cols-2 2xl:grid-cols-3' : ''
          }`}>
            {storageData.physicalLocations.map((location, index) => (
              <SinglePhysicalLocationCard
                key={location.physicalLocationId}
                location={location}
                sublocations={storageData.sublocations}
                onEdit={handleEditService}
                onDelete={handleDeleteLocation}
                onAddSublocation={handleAddSublocation}
                isWatchedByResizeObserver={index === 0}
              />
            ))}
          </div>
        </div>

        <div className="p-4 border rounded-md">
          <h2 className="text-lg font-semibold">Sublocations</h2>
          <p className="text-gray-500 mb-4">{
            storageData.sublocations.length === 1
              ? '1 sublocation found'
              : `${storageData.sublocations.length} sublocations found`
          }</p>
          <div className={`grid grid-cols-1 gap-4 ${
            viewMode === 'grid' ? 'md:grid-cols-2 2xl:grid-cols-3' : ''
          }`}>
            {storageData.sublocations.map((location, index) => (
              <SingleSublocationCard
                key={location.sublocationId}
                location={location}
                onEdit={handleEditService}
                isWatchedByResizeObserver={index === 0}
              />
            ))}
          </div>
        </div>
      </>
    );
  };

  // Render drawer content based on the current state
  const renderDrawerContent = () => {
    if (isSelectingParentLocation) {
      return (
        <div
          className="space-y-6"
          role="region"
          aria-label="Parent location selection"
        >
          {/* Header Section */}
          <div className="space-y-2">
            <h2
              className="text-lg font-semibold"
              id="selection-header"
            >
              Select a Parent Location
            </h2>
            <p
              className="text-sm text-gray-500"
              id="selection-description"
            >
              Choose the physical location where you want to add a sublocation
            </p>
          </div>

          {/* Location List Section */}
          <div
            className="space-y-4"
            role="listbox"
            aria-labelledby="selection-header"
            aria-describedby="selection-description"
          >
            {storageData?.physicalLocations.length === 0 ? (
              <div
                className="p-4 border rounded-md bg-gray-50"
                role="alert"
              >
                <p className="text-gray-500">
                  No physical locations found. Please add a physical location first.
                </p>
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-4">
                {storageData?.physicalLocations.map((location) => (
                  <SinglePhysicalLocationCard
                    key={location.physicalLocationId}
                    location={location}
                    sublocations={storageData.sublocations}
                    isSelectionMode={true}
                    onSelect={handleParentLocationSelect}
                    isSelected={selectedParentLocation?.physicalLocationId === location.physicalLocationId}
                  />
                ))}
              </div>
            )}
          </div>

          {/* Selection Instructions */}
          <div
            className="p-4 border rounded-md bg-blue-50"
            role="status"
            aria-live="polite"
          >
            <p className="text-sm text-blue-700">
              Click on a location to select it as the parent for your new sublocation
            </p>
          </div>
        </div>
      );
    }

    if (selectedParentLocation) {
      return (
        <SublocationForm
          parentLocation={selectedParentLocation}
          onSuccess={handleFormSuccess}
          buttonText="Add Sublocation"
        />
      );
    }

    return null;
  };

  return (
    <PageMain>
      <PageHeadline>
        <div className="flex items-center">
          <h1 className='text-2xl font-bold tracking-tight'>Physical Locations</h1>
        </div>

        <div className="flex items-center space-x-2">
          {/* Add Location Drawer */}
          <DrawerContainer
            open={addPhysicalLocationOpen}
            onOpenChange={setAddPhysicalLocationOpen}
            triggerAddLocation="Add Physical Location"
            title="Physical Location"
            description="Tell us about the location you want to add"
            triggerBtnIcon="location"
          >
            <PhysicalLocationForm
              onSuccess={handleFormSuccess}
              buttonText="Add Location"
            />
          </DrawerContainer>

          {/* Add Sublocation Drawer */}
          {
            storageData?.physicalLocations && storageData.physicalLocations.length > 0 && (
              <DrawerContainer
                open={addSublocationOpen}
                onOpenChange={setAddSublocationOpen}
                triggerAddLocation="Add Sublocation"
                title={isSelectingParentLocation ? "Select Parent Location" : "Add Sublocation"}
                description={isSelectingParentLocation ? "Choose a parent location for your new sublocation" : "Tell us about the sublocation you want to add"}
                triggerBtnIcon="location"
                onTriggerClick={handleStartSublocationCreation}
              >
                {renderDrawerContent()}
              </DrawerContainer>
            )
          }

          {/* Edit Location Drawer */}
          <DrawerContainer
            open={editServiceOpen}
            onOpenChange={setEditServiceOpen}
            title="Edit Physical Location"
            description="Update your location details"
          >
            {serviceBeingEdited && (
              <PhysicalLocationForm
                onSuccess={handleFormSuccess}
                buttonText="Update Location"
                isEditing={true}
                locationData={{
                  id: 'sublocationId' in serviceBeingEdited
                    ? serviceBeingEdited.sublocationId
                    : serviceBeingEdited.physicalLocationId,
                  name: 'sublocationName' in serviceBeingEdited
                    ? serviceBeingEdited.sublocationName
                    : serviceBeingEdited.name,
                  locationType: 'sublocationType' in serviceBeingEdited
                    ? serviceBeingEdited.sublocationType
                    : serviceBeingEdited.physicalLocationType,
                  bgColor: 'parentLocationBgColor' in serviceBeingEdited
                    ? serviceBeingEdited.parentLocationBgColor
                    : serviceBeingEdited.bgColor,
                  mapCoordinates: typeof serviceBeingEdited.mapCoordinates === 'string'
                    ? serviceBeingEdited.mapCoordinates
                    : serviceBeingEdited.mapCoordinates?.coords,
                  createdAt: serviceBeingEdited.createdAt ? new Date(serviceBeingEdited.createdAt) : undefined,
                  updatedAt: serviceBeingEdited.updatedAt ? new Date(serviceBeingEdited.updatedAt) : undefined
                }}
                onClose={handleFormSuccess}
              />
            )}
          </DrawerContainer>
        </div>
      </PageHeadline>

      {/* Physical Locations Display Section */}
      <div className="mt-6">
        <PhysicalLocationsToolbar
          sublocationTypes={filterOptions.sublocationTypes}
          parentTypes={filterOptions.parentTypes}
        />
        <div className="mt-4 space-y-4">
          {renderContent()}
        </div>
      </div>
    </PageMain>
  );
}