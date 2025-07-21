import { memo } from 'react';
import { BoxIcon } from '@/shared/components/ui/icons';
import { IconCloudDataConnection } from '@/shared/components/ui/icons';
import type { PhysicalLocationResponse, DigitalLocationResponse } from '@/types/domain/library-types';
import {
  extractPhysicalLocationData,
  extractDigitalLocationData
} from '@/features/dashboard/lib/utils/libraryCardUtils';
import { PhysicalLocationDisplay } from './PhysicalLocationDisplay';
import { DigitalLocationDisplay } from './DigitalLocationDisplay';

interface LibraryGameDetailCardInfoSectionProps {
  physicalLocations?: PhysicalLocationResponse[];
  digitalLocations?: DigitalLocationResponse[];
  totalPhysicalVersions?: number;
  totalDigitalVersions?: number;
  isVisible?: boolean;
  isMobile?: boolean;
  hasStackedContent?: boolean;
}

export const LibraryGameDetailCardInfoSection = memo(({
  physicalLocations = [],
  digitalLocations = [],
  totalPhysicalVersions = 0,
  totalDigitalVersions = 0,
  isVisible = true,
  isMobile = false,
  hasStackedContent = false,
}: LibraryGameDetailCardInfoSectionProps) => {

  if (!isVisible || isMobile) return null;

  // Extract location data for display
  const physicalLocationData = extractPhysicalLocationData(physicalLocations);
  const digitalLocationData = extractDigitalLocationData(digitalLocations);

  return (
    <div className={`flex flex-col gap-2 ${
      hasStackedContent ? 'flex-col max-w-[70px] overflow-x-hidden' : ''
    }`}>
      {/* Physical Versions */}
      {totalPhysicalVersions > 0 && (
        <div className="flex flex-col gap-2 w-96">
          <div className="flex flex-row items-center justify-between text-foreground">
            <div className="text-lg font-semibold leading-none tracking-tight">
              Copies in Physical Storage
              <span className="mr-2 text-sm uppercase"> {` (${totalPhysicalVersions})`}</span>
            </div>
            <BoxIcon className="h-7 w-7" />
          </div>

          {/* Physical Location Displays */}
          {physicalLocationData.map((location, index) => (
            <PhysicalLocationDisplay
              key={`physical-${location.parentLocationName}-${index}`}
              location={location}
              maxPlatforms={2}
            />
          ))}
        </div>
      )}

      {/* Digital Versions */}
      {totalDigitalVersions > 0 && (
        <div className="flex flex-col gap-2 w-96">
          <div className="flex flex-row items-center justify-between text-foreground">
            <div className="text-lg font-semibold leading-none tracking-tight">
              Copies in Digital Storage
              <span className="mr-2 text-xs uppercase"> {` (${totalDigitalVersions})`}</span>
            </div>
            <IconCloudDataConnection className="h-8 w-8" />
          </div>

          {/* Digital Location Displays */}
          {digitalLocationData.map((location, index) => (
            <DigitalLocationDisplay
              key={`digital-${location.digitalLocationName}-${index}`}
              location={location}
              maxPlatforms={2}
            />
          ))}
        </div>
      )}
    </div>
  );
}, (prevProps, nextProps) => {
  return prevProps.hasStackedContent === nextProps.hasStackedContent &&
    prevProps.isMobile === nextProps.isMobile &&
    prevProps.isVisible === nextProps.isVisible &&
    prevProps.totalPhysicalVersions === nextProps.totalPhysicalVersions &&
    prevProps.totalDigitalVersions === nextProps.totalDigitalVersions &&
    JSON.stringify(prevProps.physicalLocations) === JSON.stringify(nextProps.physicalLocations) &&
    JSON.stringify(prevProps.digitalLocations) === JSON.stringify(nextProps.digitalLocations);
});

LibraryGameDetailCardInfoSection.displayName = 'LibraryGameDetailCardInfoSection';

/*
LEGACY Implementation for reference
import { memo, type ReactNode } from 'react';

interface LibraryGameDetailCardInfoSectionProps {
  icon: ReactNode;
  label: string;
  value: string;
  hasStackedContent?: boolean;
  isVisible?: boolean;
  isMobile?: boolean;
  isCardView?: boolean;
};

export const LibraryGameDetailCardInfoSection = memo(({
  icon,
  label,
  value,
  hasStackedContent = false,
  isVisible = true,
  isMobile = false,
  isCardView = false,
}: LibraryGameDetailCardInfoSectionProps) => {

  if (!isVisible || isMobile) return null;

  return (
    <div className={`flex flex-row items-center gap-2 ${
      hasStackedContent ? 'flex-col max-w-[70px] overflow-x-hidden' : ''
    }`}>
      {icon}
      <div className={`flex flex-col ${isCardView ? 'ml-[5px]' : ''}`}>
        <span className={`mr-2 text-xs uppercase ${hasStackedContent ? 'hidden' : ''}`}>{label}</span>
        <span className={`text-sm text-white ${hasStackedContent ? 'max-w-[70px]' : 'max-w-[105px]'} overflow-x-hidden truncate`}>{value.charAt(0).toUpperCase() + value.slice(1)}</span>
      </div>
    </div>
  );
}, (prevProps, nextProps) => {
  return prevProps.hasStackedContent === nextProps.hasStackedContent &&
    prevProps.isMobile === nextProps.isMobile &&
    prevProps.isVisible === nextProps.isVisible &&
    prevProps.label === nextProps.label &&
    prevProps.value === nextProps.value &&
    prevProps.icon === nextProps.icon;
});

*/
