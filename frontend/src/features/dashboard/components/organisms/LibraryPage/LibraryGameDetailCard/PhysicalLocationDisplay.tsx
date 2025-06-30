import { memo } from 'react';
import { PhysicalLocationIcon } from '@/features/dashboard/lib/utils/getPhysicalLocationIcon';
import { SublocationIcon } from '@/features/dashboard/lib/utils/getSublocationIcon';
import { formatPlatformsForDisplay } from '@/features/dashboard/lib/utils/libraryCardUtils';
import type { PhysicalLocationDisplayData } from '@/features/dashboard/lib/utils/libraryCardUtils';
import type { LocationIconBgColor } from '@/types/domain/location-types';

interface PhysicalLocationDisplayProps {
  location: PhysicalLocationDisplayData;
  maxPlatforms?: number;
}

export const PhysicalLocationDisplay = memo(({
  location,
  maxPlatforms = 2
}: PhysicalLocationDisplayProps) => {
  const { displayPlatforms, overflowCount } = formatPlatformsForDisplay(location.platforms, maxPlatforms);

  return (
    <div className="flex items-center gap-2 mb-4">
      <div className="flex items-center gap-1">
        <PhysicalLocationIcon
          type={location.parentLocationType}
          bgColor={location.parentLocationBgColor as LocationIconBgColor}
        />
        <span className="text-xs text-white/80">{location.parentLocationName}</span>
      </div>
      <div className="flex items-center gap-1">
        <SublocationIcon
          type={location.sublocationType}
          bgColor={location.parentLocationBgColor as LocationIconBgColor}
        />
        <span className="text-xs text-white/80">{location.sublocationName}</span>
      </div>
      <div className="flex flex-wrap gap-1">
        {displayPlatforms.map((platform, index) => (
          <span
            key={`${platform}-${index}`}
            className="text-xs bg-white/20 px-1 rounded text-white"
          >
            {platform}
          </span>
        ))}
        {overflowCount > 0 && (
          <span className="text-xs text-white/60">
            +{overflowCount} more
          </span>
        )}
      </div>
    </div>
  );
});

PhysicalLocationDisplay.displayName = 'PhysicalLocationDisplay';