import { memo } from 'react';
import { truncateWarningText } from '@/features/dashboard/lib/utils/truncateWarningText';
import type {
  LocationsBFFPhysicalLocationResponse,
  LocationsBFFSublocationResponse
} from '@/types/domain/physical-location';

interface PhysicalLocationDeleteWarningProps {
  location: LocationsBFFPhysicalLocationResponse;
  associatedItems: LocationsBFFSublocationResponse[];
}

export const PhysicalLocationDeleteWarning = memo(({
  location,
  associatedItems
}: PhysicalLocationDeleteWarningProps) => {
  const childSublocations = associatedItems.filter(
    sublocation => sublocation.parentLocationId === location.physicalLocationId
  );
  const hasSublocations = childSublocations.length > 0;

  if (!hasSublocations) {
    return (
      <p>
        Are you sure you want to delete this location? This action cannot be undone.
      </p>
    );
  }

  return (
    <div className="space-y-2">
      <p>
        Are you sure to want to delete this location? You will be deleting the following child sublocations stored here:
      </p>
      <ul className="list-disc pl-4 space-y-1">
        {childSublocations.slice(0, 5).map((sublocation) => (
          <li key={sublocation.sublocationId} className="text-sm">
            "{truncateWarningText(sublocation.sublocationName)}": ({sublocation.storedItems} games)
          </li>
        ))}
        {childSublocations.length > 5 && (
          <li className="text-sm italic">
            Complete sublocation listing for this location available on either the physical locations or media storage pages
          </li>
        )}
      </ul>
      <p className="mt-2">
        Some of these sublocations may have unique copies of games that you don't have anywhere else.
        You will *also* delete those games from your library.
      </p>
      <p className="mt-2">
        Are you sure that you want to delete this location? This action cannot be undone.
      </p>
    </div>
  );
});

PhysicalLocationDeleteWarning.displayName = 'PhysicalLocationDeleteWarning';