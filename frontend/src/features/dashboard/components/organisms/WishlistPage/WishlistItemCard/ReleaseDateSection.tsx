import { memo } from 'react';
import { SupportedItemPlatforms } from './SupportedItemPlatforms';
import type { PlatformCategory } from '@/shared/types/platform';

interface ReleaseDateSectionProps {
  platform: string;
  releaseDate: string;
  hasAndroidVersion: boolean;
  hasIOSVersion: boolean;
  hasMacOSVersion: boolean;
}

export const ReleaseDateSection = memo(({
  platform,
  releaseDate,
  hasAndroidVersion,
  hasIOSVersion,
  hasMacOSVersion
}: ReleaseDateSectionProps) => {
  return (
    <div className="flex items-center gap-2 text-sm text-gray-400 mt-2">
      <SupportedItemPlatforms
        platform={platform as PlatformCategory}
        hasAndroidVersion={hasAndroidVersion}
        hasIOSVersion={hasIOSVersion}
        hasMacOSVersion={hasMacOSVersion}
      />
      {releaseDate}
    </div>
  );
}, (prevProps, nextProps) => {
  return (
    prevProps.platform === nextProps.platform &&
    prevProps.releaseDate === nextProps.releaseDate &&
    prevProps.hasAndroidVersion === nextProps.hasAndroidVersion &&
    prevProps.hasIOSVersion === nextProps.hasIOSVersion &&
    prevProps.hasMacOSVersion === nextProps.hasMacOSVersion
  );
});
