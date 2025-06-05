import { useState, useCallback } from 'react';

// Template Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';

// Components
import { SinglePhysicalLocationCard } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/SinglePhysicalLocationCard/SinglePhysicalLocationCard';
import { PhysicalLocationsToolbar } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/PhysicalLocationsToolbar/PhysicalLocationsToolbar';
import { PhysicalLocationsTable } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/PhysicalLocationsTable/PhysicalLocationsTable';

import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';

// API Hooks and Utilities
import { useStorageAnalytics } from '@/core/api/queries/analyticsData.queries';
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { useCardLabelWidth } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/useCardLabelWidth';
import { useFilteredPhysicalLocations } from '@/features/dashboard/lib/hooks/useFilteredPhysicalLocations';

// Types
import type { PhysicalLocation } from '@/types/domain/physical-location';
import type { SublocationRowData } from '@/core/api/adapters/analytics.adapter';

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
  const [addServiceOpen, setAddServiceOpen] = useState<boolean>(false);
  const [editServiceOpen, setEditServiceOpen] = useState<boolean>(false);
  const [serviceBeingEdited, setServiceBeingEdited] = useState<PhysicalLocation | null>(null);

  const viewMode = useOnlineServicesStore((state) => state.viewMode);

  // Fetch physical locations using analytics
  const { data: storageData, isLoading, error } = useStorageAnalytics();

  // Get filtered locations using the new hook
  const filteredSublocationRows = useFilteredPhysicalLocations(storageData?.sublocationRows || []);

  console.log('[DEBUG] PhysicalLocationsPageContent: Render', {
    isLoading,
    hasError: !!error,
    sublocationRowsCount: storageData?.sublocationRows.length || 0,
    filteredSublocationRowsCount: filteredSublocationRows.length,
    viewMode,
    timestamp: new Date().toISOString()
  });

  // Enhanced edit handlers
  const handleEditService = useCallback((sublocation: SublocationRowData) => {
    console.log('[DEBUG] PhysicalLocationsPageContent: Edit service', {
      sublocationId: sublocation.sublocationId,
      parentLocationId: sublocation.parentLocationId,
      timestamp: new Date().toISOString()
    });

    // Find the parent location
    const parentLocation = storageData?.physicalLocations.find(loc =>
      loc.sublocations?.some(sub => sub.id === sublocation.sublocationId)
    );
    if (parentLocation) {
      setServiceBeingEdited(parentLocation);
      setEditServiceOpen(true);
    }
  }, [storageData]);

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
    console.log('[DEBUG] PhysicalLocationsPageContent: renderContent', {
      isLoading,
      hasError: !!error,
      sublocationRowsCount: storageData?.sublocationRows.length || 0,
      filteredSublocationRowsCount: filteredSublocationRows.length,
      viewMode,
      timestamp: new Date().toISOString()
    });

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

    if (storageData?.sublocationRows.length === 0) {
      return (
        <div className="p-4 border rounded-md">
          <p className="text-gray-500">No physical location found. Add a location to get started.</p>
        </div>
      );
    }

    if (filteredSublocationRows.length === 0) {
      return <NoResultsFound />;
    }

    if (viewMode === 'table') {
      return (
        <PhysicalLocationsTable
          sublocationRows={filteredSublocationRows}
          onEdit={handleEditService}
        />
      );
    }

    return (
      <div className="p-4 border rounded-md">
        <h2 className="text-lg font-semibold">Physical Locations</h2>
        <p className="text-gray-500 mb-4">{
          filteredSublocationRows.length === 1
            ? '1 location found'
            : `${filteredSublocationRows.length} locations found`
        }</p>
        <div className={`grid grid-cols-1 gap-4 ${
          viewMode === 'grid' ? 'md:grid-cols-2 2xl:grid-cols-3' : ''
        }`}>
          {filteredSublocationRows.map((sublocation, index) => (
            <SinglePhysicalLocationCard
              key={sublocation.sublocationId}
              location={sublocation}
              onEdit={handleEditService}
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
          <h1 className='text-2xl font-bold tracking-tight'>Physical Locations</h1>
        </div>

        <div className="flex items-center space-x-2">
          {/* Add Location Drawer */}
          <DrawerContainer
            open={addServiceOpen}
            onOpenChange={setAddServiceOpen}
            triggerAddLocation="Add Physical Location"
            title="Physical Location"
            description="Tell us about the location you want to add"
            triggerBtnIcon="location"
          >
            <p>Add Physical Location Form</p>
          </DrawerContainer>

          {/* Edit Location Drawer */}
          <DrawerContainer
            open={editServiceOpen}
            onOpenChange={setEditServiceOpen}
            title="Edit Physical Location"
            description="Update your location details"
          >
            {serviceBeingEdited && (
              <p>Edit Physical Location Form</p>
            )}
          </DrawerContainer>
        </div>
      </PageHeadline>

      {/* Physical Locations Display Section */}
      <div className="mt-6">
        <PhysicalLocationsToolbar />
        <div className="mt-4 space-y-4">
          {renderContent()}
        </div>
      </div>
    </PageMain>
  );
}