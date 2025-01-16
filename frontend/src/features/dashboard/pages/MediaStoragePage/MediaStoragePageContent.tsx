import { useState } from 'react';

// Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer' ;
import { MediaPageLocationForm } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaPageLocationForm/MediaPageLocationForm';
import { MediaPageSublocationForm } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaPageSublocationForm/MediaPageSublocationForm';
import { MediaStoragePageAccordion } from '@/features/dashboard/components/organisms/MediaStoragePageAccordion/MediaStoragePageAccordion';

// Mock Data
import { mediaStoragePageMockData } from '@/features/dashboard/pages/MediaStoragePage/mediaStoragePage.mockdata';

// Mock data - You might want to move this to a separate file
const physicalLocationsData = [
  {
    name: "Physical Location 1",
    subLocations: [
      {
        title: "Sub Location A",
        description: "Description for Sub Location A",
        src: "/placeholder.svg",
        ctaText: "Learn More",
      },
    ],
  },
  // Add more locations as needed
];

const digitalLocationsData = [
  {
    name: "Digital Location 1",
    subLocations: [
      {
        title: "Digital Sub Location A",
        description: "Description for Digital Sub Location A",
        src: "/placeholder.svg",
        ctaText: "Access",
      },
    ],
  },
  // Add more locations as needed
];

export function MediaStoragePageContent() {
  const [open, setOpen] = useState<boolean>(false)


  const { data: { physicalLocations, digitalLocations }, meta } = mediaStoragePageMockData;
  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>Media Storage Page</h1>
        </div>

        {/* Add Physical Location Button */}
        <div className='flex items-center space-x-2'>
          <DrawerContainer
            open={open}
            onOpenChange={setOpen}
            triggerText="Add Physical Location"
            title="Physical Location"
            description="Tell us about where your games are stored."
          >
            <MediaPageSublocationForm onSuccess={() => setOpen(false)} />
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
