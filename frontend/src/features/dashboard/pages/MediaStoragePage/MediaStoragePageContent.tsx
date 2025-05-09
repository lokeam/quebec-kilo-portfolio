import { useStorageAnalytics } from '@/core/api/queries/analyticsData.queries';
import { MediaStoragePageAccordion } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/MediaStoragePageAccordion';

export function MediaStoragePageContent() {
  const { data: storageData, isLoading, error } = useStorageAnalytics();

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error loading storage data</div>;
  }

  if (!storageData) {
    return <div>No storage data available</div>;
  }

  return (
    <div className="space-y-6">
      <MediaStoragePageAccordion
        title="Physical Locations"
        locationData={storageData.physicalLocations}
        meta={storageData.metadata}
        type="physical"
        isLoading={isLoading}
      />

      <MediaStoragePageAccordion
        title="Digital Locations"
        locationData={storageData.digitalLocations}
        meta={storageData.metadata}
        type="digital"
        isLoading={isLoading}
      />
    </div>
  );
}
