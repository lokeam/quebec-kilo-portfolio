import { useState } from 'react';

// Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer' ;
import { PhysicalLocationFormSingle } from '@/features/dashboard/components/organisms/MediaStoragePage/PhysicalLocationFormSingle/PhysicalLocationFormSingle';
import { PhysicalLocationDrawerList } from '@/features/dashboard/components/organisms/MediaStoragePage/PhysicalLocationDrawerList/PhysicalLocationDrawerList';
import { MediaStoragePageAccordion } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/MediaStoragePageAccordion';

// Hooks
import { usePhysicalLocations } from '@/core/api/hooks/usePhysicalLocations';
import { useDigitalLocations } from '@/core/api/hooks/useDigitalLocations';
import { useStorageAnalytics } from '@/core/api/hooks/useAnalyticsData';

export function MediaStoragePageContent() {
  const [addLocationOpen, setAddLocationOpen] = useState<boolean>(false);
  const [editLocationOpen, setEditLocationOpen] = useState<boolean>(false);

  // Fetch physical locations with the hook
  const { data: physicalLocationsData, isLoading: isLoadingPhysicalLocations } = usePhysicalLocations();

  // Fetch digital locations with the hook
  const { data: digitalLocations, isLoading: isLoadingDigitalLocations } = useDigitalLocations();

  // Fetch storage analytics
  const { data: analyticsData } = useStorageAnalytics();

  // Ensure we always have arrays for locations, even if data is still loading
  const safePhysicalLocations = physicalLocationsData || [];
  const safeDigitalLocations = digitalLocations || [];

  // Extract metadata from analytics
  const meta = {
    counts: {
      locations: {
        total: (analyticsData?.data?.storage?.totalPhysicalLocations || 0) + (analyticsData?.data?.storage?.totalDigitalLocations || 0),
        physical: analyticsData?.data?.storage?.totalPhysicalLocations || 0,
        digital: analyticsData?.data?.storage?.totalDigitalLocations || 0,
      },
      items: {
        total: 0,
        physical: 0,
        digital: 0,
        byLocation: {},
      },
    },
    lastUpdated: new Date(),
    version: '1.0',
  };

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
        locationData={safeDigitalLocations}
        title="Digital Locations"
        meta={meta}
        type="digital"
        isLoading={isLoadingDigitalLocations}
      />
    </PageMain>
  );
}
