import { memo, useMemo } from 'react';
import SVGLogo from '@/shared/components/ui/LogoMap/LogoMap';
import type { PlatformCategory } from '@/shared/types/platform';

type SupportedItemPlatformsProps = {
  platform: PlatformCategory;
  hasAndroidVersion: boolean;
  hasIOSVersion: boolean;
  hasMacOSVersion: boolean;
}

export const SupportedItemPlatforms = memo(({
  platform,
  hasAndroidVersion,
  hasIOSVersion,
  hasMacOSVersion
}: SupportedItemPlatformsProps) => {

  console.log('supported item platforms', platform, hasAndroidVersion, hasIOSVersion, hasMacOSVersion);

  // Object literal pattern for platform-specific render logic
  const platformIconMap = useMemo(() => ({
    PC: () => [
      <SVGLogo key="pc" domain="platforms" name="pc" className="w-8 h-8" />,
      hasMacOSVersion && <SVGLogo key="macos" domain="platforms" name="macos" className="w-8 h-8" />
    ],
    Console: () => [
      <SVGLogo key="pc" domain="platforms" name="console" className="w-8 h-8" />
    ],
    Mobile: () => [
      hasAndroidVersion && (
        <SVGLogo key="android" domain="platforms" name="android" className="w-8 h-8" />
      ),
      hasIOSVersion && (
        <SVGLogo key="ios" domain="platforms" name="ios" className="w-8 h-8" />
      ),
    ],
    // Handle individual mobile platforms explicitly
    Android: () => [
      <SVGLogo key="android" domain="platforms" name="android" className="w-8 h-8" />
    ],
    iOS: () => [
      <SVGLogo key="ios" domain="platforms" name="ios" className="w-8 h-8" />
    ]
  }), [hasAndroidVersion, hasIOSVersion, hasMacOSVersion]);

  // Get icons for current platform and remove falsy values
  const platformIcons = useMemo(() =>
    platformIconMap[platform]?.()?.filter(Boolean) ?? [],
    [platform, platformIconMap]
  );

  if (platformIcons.length === 0) return null;

  return (
    <div className="flex items-center gap-2">
      {platformIcons}
    </div>
  )
}, (prevProps, nextProps) => {
  return prevProps.platform === nextProps.platform &&
    prevProps.hasAndroidVersion === nextProps.hasAndroidVersion &&
    prevProps.hasIOSVersion === nextProps.hasIOSVersion &&
    prevProps.hasMacOSVersion === nextProps.hasMacOSVersion;
});

SupportedItemPlatforms.displayName = 'SupportedItemPlatforms';
