import { memo } from 'react';
import { DigitalLocationIcon } from '@/features/dashboard/lib/utils/getDigitalLocationIcon';
import { formatPlatformsForDisplay } from '@/features/dashboard/lib/utils/libraryCardUtils';
import type { DigitalLocationDisplayData } from '@/features/dashboard/lib/utils/libraryCardUtils';

interface DigitalLocationDisplayProps {
  location: DigitalLocationDisplayData;
  maxPlatforms?: number;
}

export const DigitalLocationDisplay = memo(({
  location,
  maxPlatforms = 2
}: DigitalLocationDisplayProps) => {
  const { displayPlatforms, overflowCount } = formatPlatformsForDisplay(location.platforms, maxPlatforms);

  return (
    <div className="flex items-center gap-2">
      <div className="flex items-center gap-1">
        <DigitalLocationIcon
          name={location.digitalLocationName}
          className="h-7 w-7"
        />
        <span className="text-xs text-white/80">{location.digitalLocationName}</span>
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

DigitalLocationDisplay.displayName = 'DigitalLocationDisplay';