import { useState } from 'react';

// Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer' ;
import { PhysicalLocationFormSingle } from '@/features/dashboard/components/organisms/MediaStoragePage/PhysicalLocationFormSingle/PhysicalLocationFormSingle';
import { PhysicalLocationDrawerList } from '@/features/dashboard/components/organisms/MediaStoragePage/PhysicalLocationDrawerList/PhysicalLocationDrawerList';
import { MediaStoragePageAccordion } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/MediaStoragePageAccordion';

// Mock Data
import { mediaStoragePageMockData } from '@/features/dashboard/pages/MediaStoragePage/mediaStoragePage.mockdata';

export function MediaStoragePageContent() {

  const [addLocationOpen, setAddLocationOpen] = useState<boolean>(false);
  const [editLocationOpen, setEditLocationOpen] = useState<boolean>(false);

  const { data: { physicalLocations, digitalLocations }, meta } = mediaStoragePageMockData;
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
            {/* <PhysicalLocationFormList
              locationData={physicalLocations}
              onSuccess={() => setEditLocationOpen(false)}
            /> */}
            <PhysicalLocationDrawerList
              onSuccess={() => setEditLocationOpen(false)}
              locationData={physicalLocations}
            />
          </DrawerContainer>
        </div>

      </PageHeadline>

      {/* Physical Locations Accordion */}
      <MediaStoragePageAccordion
        locationData={physicalLocations}
        title="Physical Locations"
        meta={meta}
        type="physical"
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
