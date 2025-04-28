import { useState } from 'react';

// Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer' ;
import { PhysicalLocationFormSingle } from '@/features/dashboard/components/organisms/MediaStoragePage/PhysicalLocationFormSingle/PhysicalLocationFormSingle';
import { PhysicalLocationDrawerList } from '@/features/dashboard/components/organisms/MediaStoragePage/PhysicalLocationDrawerList/PhysicalLocationDrawerList';
import { MediaStoragePageAccordion } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/MediaStoragePageAccordion';

// Mock Data - keeping for digital locations until that's implemented
import { mediaStoragePageMockData } from '@/features/dashboard/pages/MediaStoragePage/mediaStoragePage.mockdata';

// Hooks
import { usePhysicalLocations } from '@/core/api/hooks/usePhysicalLocations';

export function MediaStoragePageContent() {
  const [addLocationOpen, setAddLocationOpen] = useState<boolean>(false);
  const [editLocationOpen, setEditLocationOpen] = useState<boolean>(false);

  // Fetch physical locations with the hook
  const { data: physicalLocationsData, isLoading: isLoadingPhysicalLocations } = usePhysicalLocations();

  // Use mock data for digital locations until that's implemented
  const { data: { digitalLocations }, meta } = mediaStoragePageMockData;

  // Ensure we always have an array for physical locations, even if data is still loading
  const safePhysicalLocations = physicalLocationsData || [];

  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>Media Storage Page</h1>
        </div>


        <div className='flex items-center space-x-2'>
          {/* Add Physical Location Button */}
          <DrawerContainer
            open={addLocationOpen}
            onOpenChange={setAddLocationOpen}
            triggerAddLocation="Add Physical Location"
            title="Physical Location"
            description="Tell us about where your games are stored."
          >
            <PhysicalLocationFormSingle onSuccess={() => setAddLocationOpen(false)} />
          </DrawerContainer>

          {/* Edit Physical Location Button */}
          <DrawerContainer
            open={editLocationOpen}
            onOpenChange={setEditLocationOpen}
            triggerEditLocation="Edit Physical Location"
            title="Edit Locations"
            description="Edit your physical locations and sublocations"
          >
            <PhysicalLocationDrawerList
              onSuccess={() => setEditLocationOpen(false)}
              locationData={safePhysicalLocations}
            />
          </DrawerContainer>
        </div>

      </PageHeadline>

      {/* Physical Locations Accordion */}
      <MediaStoragePageAccordion
        locationData={safePhysicalLocations}
        title="Physical Locations"
        meta={meta}
        type="physical"
        isLoading={isLoadingPhysicalLocations}
      />

      {/* Digital Location Accordion */}
      <MediaStoragePageAccordion
        locationData={digitalLocations}
        title="Digital Locations"
        meta={meta}
        type="digital"
      />
    </PageMain>
  );
}
