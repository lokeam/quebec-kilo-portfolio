import { useState, useCallback, useMemo } from 'react';

// Template Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';

// Components
import { SinglePhysicalLocationCard } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/SinglePhysicalLocationCard/SinglePhysicalLocationCard';
import { SingleSublocationCard } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/SingleSublocationCard/SingleSublocationCard';
import { PhysicalLocationsToolbar } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/PhysicalLocationsToolbar/PhysicalLocationsToolbar';
import { PhysicalLocationsTable } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/PhysicalLocationsTable/PhysicalLocationsTable';
import { PhysicalLocationForm } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/PhysicalLocationsForm/PhysicalLocationsForm';
import { SublocationForm } from '@/features/dashboard/components/organisms/SublocationForm/SublocationForm';
import { PhysicalLocationsPageSkeleton } from '@/features/dashboard/pages/PhysicalLocations/PhysicalLocationsPageSkeleton';

import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';

// API Hooks and Utilities
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { useCardLabelWidth } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/useCardLabelWidth';
import { usePhysicalLocationFilters } from '@/features/dashboard/hooks/usePhysicalLocationFilters';

// Types
import type { LocationsBFFPhysicalLocationResponse, LocationsBFFSublocationResponse } from '@/types/domain/physical-location';
import { useGetPhysicalLocationsBFFResponse, useDeletePhysicalLocation } from '@/core/api/queries/physicalLocation.queries';


// Empty States
const NoDataState = () => (
  <div className="p-4 border rounded-md">
    <p className="text-gray-500">No physical location found. Add a location to get started.</p>
  </div>
);

const NoFilterResultsState = () => (
  <div className="p-4 border rounded-md">
    <p className="text-gray-500">No locations match your current filters. Try adjusting your search or filters.</p>
  </div>
);

const NoSublocationsState = () => (
  <div className="p-4 border rounded-md">
    <p className="text-gray-500">No sublocations found. Add a sublocation to get started.</p>
  </div>
);

// Error State
const ErrorState = ({ error }: { error: unknown }) => (
  <div className="p-4 border border-red-300 bg-red-50 rounded-md">
    <p className="text-red-500">
      Error loading physical location data: {error instanceof Error ? error.message : 'An unknown error occurred'}
    </p>
    {/*
    <p className="text-red-500">
      Error loading physical location data: {error instanceof Error ? error.message : 'An unknown error occurred'}
    </p>
    */}
  </div>
);

// Helper function to get physical locations without sublocations
const getPhysicalLocationsWithoutSublocations = (
  physicalLocations: LocationsBFFPhysicalLocationResponse[],
  sublocations: LocationsBFFSublocationResponse[]
) => {
  const parentIdsWithSublocations = new Set(
    sublocations.map(sublocation => sublocation.parentLocationId)
  );

  return physicalLocations.filter(
    location => !parentIdsWithSublocations.has(location.physicalLocationId)
  );
};

export function PhysicalLocationsPageContent() {
  const [addPhysicalLocationOpen, setAddPhysicalLocationOpen] = useState<boolean>(false);
  const [addSublocationOpen, setAddSublocationOpen] = useState<boolean>(false);
  const [editServiceOpen, setEditServiceOpen] = useState<boolean>(false);
  const [serviceBeingEdited, setServiceBeingEdited] = useState<LocationsBFFPhysicalLocationResponse | LocationsBFFSublocationResponse | null>(null);
  const [selectedParentLocation, setSelectedParentLocation] = useState<LocationsBFFPhysicalLocationResponse | null>(null);
  const [isSelectingParentLocation, setIsSelectingParentLocation] = useState<boolean>(false);

  const viewMode = useOnlineServicesStore((state) => state.viewMode);
  const { searchQuery, sublocationTypeFilters, parentLocationTypeFilters } = useOnlineServicesStore();

  // Fetch physical locations using BFF
  const { data: storageData, isLoading, error } = useGetPhysicalLocationsBFFResponse();

  // Get filter options from BFF data
  const filterOptions = usePhysicalLocationFilters(storageData);

  // Memoize lowercase search query
  const memoizedSearchQuery = useMemo(() => searchQuery.toLowerCase(), [searchQuery]);

  // Enhanced edit handlers
  const handleEditService = useCallback((location: LocationsBFFPhysicalLocationResponse | LocationsBFFSublocationResponse) => {
    setServiceBeingEdited(location);
    setEditServiceOpen(true);
  }, []);

  // Sublocation creation when clicking on a SinglePhysicalLocationCard's Add Sublocation button
  const handleAddSublocation = useCallback((location: LocationsBFFPhysicalLocationResponse) => {
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

  // Memoize filter functions
  const filterSublocation = useCallback((sublocation: LocationsBFFSublocationResponse) => {
    const matchesSearch = sublocation.sublocationName.toLowerCase().includes(memoizedSearchQuery);
    const matchesSublocationType = sublocationTypeFilters.length === 0 ||
      sublocationTypeFilters.includes(sublocation.sublocationType);
    const matchesParentType = parentLocationTypeFilters.length === 0 ||
      parentLocationTypeFilters.includes(sublocation.parentLocationType);
    return matchesSearch && matchesSublocationType && matchesParentType;
  }, [memoizedSearchQuery, sublocationTypeFilters, parentLocationTypeFilters]);

  const filterPhysicalLocation = useCallback((location: LocationsBFFPhysicalLocationResponse) => {
    const matchesSearch = location.name.toLowerCase().includes(memoizedSearchQuery);
    const matchesType = parentLocationTypeFilters.length === 0 ||
      parentLocationTypeFilters.includes(location.physicalLocationType);
    return matchesSearch && matchesType;
  }, [memoizedSearchQuery, parentLocationTypeFilters]);

  // Memoize filtered results
  const filteredSublocations = useMemo(() => {
    if (!storageData?.sublocations) return [];

    // Early return if no filters are active
    if (searchQuery === '' &&
        sublocationTypeFilters.length === 0 &&
        parentLocationTypeFilters.length === 0) {
      return storageData.sublocations;
    }

    return storageData.sublocations.filter(filterSublocation);
  }, [
    storageData?.sublocations,
    filterSublocation,
    searchQuery,
    sublocationTypeFilters.length,
    parentLocationTypeFilters.length
  ]);

  const filteredPhysicalLocations = useMemo(() => {
    if (!storageData?.physicalLocations) return [];

    // Early return if no filters are active
    if (searchQuery === '' && parentLocationTypeFilters.length === 0) {
      return storageData.physicalLocations;
    }

    return storageData.physicalLocations.filter(filterPhysicalLocation);
  }, [
    storageData?.physicalLocations,
    filterPhysicalLocation,
    searchQuery,
    parentLocationTypeFilters.length
  ]);

  // Render content based on the selected view mode
  const renderContent = () => {
    // Loading state
    if (isLoading) {
      return <PhysicalLocationsPageSkeleton />;
    }

    // Error state
    if (error) {
      return <ErrorState error={error} />;
    }

    // No data state
    if (!storageData?.physicalLocations || storageData.physicalLocations.length === 0) {
      return <NoDataState />;
    }

    // No results after filtering
    if (filteredSublocations.length === 0 && filteredPhysicalLocations.length === 0) {
      return <NoFilterResultsState />;
    }

    if (viewMode === 'table') {
      // Get physical locations without sublocations
      const physicalLocationsWithoutSublocations = getPhysicalLocationsWithoutSublocations(
        filteredPhysicalLocations,
        filteredSublocations
      );

      return (
        <PhysicalLocationsTable
          sublocationRows={filteredSublocations}
          physicalLocationRows={physicalLocationsWithoutSublocations}
          onEdit={handleEditService}
          onDelete={handleDeleteLocation}
        />
      );
    }

    return (
      <>
        <div className="p-4 border rounded-md">
          <h2 className="text-lg font-semibold">Physical Locations</h2>
          <p className="text-gray-500 mb-4">{
            filteredPhysicalLocations.length === 1
              ? '1 location found'
              : `${filteredPhysicalLocations.length} locations found`
          }</p>
          <div className={`grid grid-cols-1 gap-4 ${
            viewMode === 'grid' ? 'md:grid-cols-2 2xl:grid-cols-3' : ''
          }`}>
            {filteredPhysicalLocations.map((location, index) => (
              <SinglePhysicalLocationCard
                key={location.physicalLocationId}
                location={location}
                sublocations={filteredSublocations}
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
            filteredSublocations.length === 1
              ? '1 sublocation found'
              : `${filteredSublocations.length} sublocations found`
          }</p>
          <div className={`grid grid-cols-1 gap-4 ${
            viewMode === 'grid' ? 'md:grid-cols-2 2xl:grid-cols-3' : ''
          }`}>
            {filteredSublocations.length === 0 ? (
              <NoSublocationsState />
            ) : (
              filteredSublocations.map((location, index) => (
                <SingleSublocationCard
                  key={location.sublocationId}
                  location={location}
                  onEdit={handleEditService}
                  isWatchedByResizeObserver={index === 0}
                />
              ))
            )}
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