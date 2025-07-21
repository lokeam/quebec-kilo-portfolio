import { memo } from 'react';
import { BoxIcon } from '@/shared/components/ui/icons';
import { IconCloudDataConnection } from '@/shared/components/ui/icons';
import type { PhysicalLocationResponse, DigitalLocationResponse } from '@/types/domain/library-types';
import {
  extractPhysicalPlatformNamesFromLocations,
  extractDigitalPlatformNamesFromLocations,
  formatPlatformsForDisplay
} from '@/features/dashboard/lib/utils/libraryCardUtils';

interface LibraryGameCardInfoSectionProps {
  isVisible?: boolean;
  totalDigitalVersions?: number;
  totalPhysicalVersions?: number;
  physicalLocations?: PhysicalLocationResponse[];
  digitalLocations?: DigitalLocationResponse[];
}

export const LibraryGameCardInfoSection = memo(({
  physicalLocations = [],
  digitalLocations = [],
  totalPhysicalVersions = 0,
  totalDigitalVersions = 0,
  isVisible = true,
}: LibraryGameCardInfoSectionProps) => {

  if (!isVisible) return null;

  // Extract platform names using utility functions
  const physicalPlatforms = extractPhysicalPlatformNamesFromLocations(physicalLocations);
  const digitalPlatforms = extractDigitalPlatformNamesFromLocations(digitalLocations);

  // Format platforms for display (limit to 2 per section to avoid overflow)
  const { displayPlatforms: displayPhysicalPlatforms, overflowCount: physicalOverflowCount } =
    formatPlatformsForDisplay(physicalPlatforms, 3);

  const { displayPlatforms: displayDigitalPlatforms, overflowCount: digitalOverflowCount } =
    formatPlatformsForDisplay(digitalPlatforms, 3);

  return (
    <div className="flex flex-col space-y-2">
      {/* Physical Copies Section */}
      {totalPhysicalVersions > 0 && (
        <div className="flex flex-col space-y-1">
          <div className="flex flex-row items-center justify-between">
            <div className="text-sm font-semibold leading-none tracking-tight">
              Physical Copies
              <span className="mr-2 text-xs uppercase"> {` (${totalPhysicalVersions})`}</span>
            </div>
            <BoxIcon className="h-7 w-7" />
          </div>

          {/* Physical Platforms List */}
          {displayPhysicalPlatforms.length > 0 && (
            <div className="flex flex-wrap gap-1 ml-2">
              {displayPhysicalPlatforms.map((platform: string, index: number) => (
                <span
                  key={`physical-${platform}-${index}`}
                  className="text-xs bg-white/20 px-1 rounded text-white"
                >
                  {platform}
                </span>
              ))}
              {physicalOverflowCount > 0 && (
                <span className="text-xs text-white/60">
                  +{physicalOverflowCount} more
                </span>
              )}
            </div>
          )}
        </div>
      )}

      {/* Digital Copies Section */}
      {totalDigitalVersions > 0 && (
        <div className="flex flex-col space-y-1">
          <div className="flex flex-row items-center justify-between">
            <div className="text-sm font-semibold leading-none tracking-tight">
              Digital Copies
              <span className="mr-2 text-xs uppercase"> {` (${totalDigitalVersions})`}</span>
            </div>
            <IconCloudDataConnection className="h-7 w-7" />
          </div>

          {/* Digital Platforms List */}
          {displayDigitalPlatforms.length > 0 && (
            <div className="flex flex-wrap gap-1 ml-2">
              {displayDigitalPlatforms.map((platform: string, index: number) => (
                <span
                  key={`digital-${platform}-${index}`}
                  className="text-xs bg-white/20 px-1 rounded text-white"
                >
                  {platform}
                </span>
              ))}
              {digitalOverflowCount > 0 && (
                <span className="text-xs text-white/60">
                  +{digitalOverflowCount} more
                </span>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}, (prevProps, nextProps) => {
  return prevProps.totalPhysicalVersions === nextProps.totalPhysicalVersions &&
    prevProps.totalDigitalVersions === nextProps.totalDigitalVersions &&
    prevProps.isVisible === nextProps.isVisible &&
    JSON.stringify(prevProps.physicalLocations) === JSON.stringify(nextProps.physicalLocations) &&
    JSON.stringify(prevProps.digitalLocations) === JSON.stringify(nextProps.digitalLocations);
});