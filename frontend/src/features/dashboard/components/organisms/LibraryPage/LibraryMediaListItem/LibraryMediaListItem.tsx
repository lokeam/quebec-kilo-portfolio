import { memo, useReducer, useCallback, useMemo, useRef } from 'react';

// Components
import { InfoSection } from '@/features/dashboard/components/organisms/LibraryPage/LibraryMediaListItem/InfoSection';

// ShadCN Components
import { Button } from '@/shared/components/ui/button';

// Hooks + Utils
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';
import { useElementBreakpoint } from '@/shared/hooks/useElementBreakpoint';
import { visibilityReducer } from '@/features/dashboard/components/organisms/WishlistPage/WishlistItemCard/visibilityReducer';
import { toast } from 'sonner';

// Types
import type { CardVisibility } from '@/features/dashboard/components/organisms/WishlistPage/WishlistItemCard/WishlistCardItem.types';

// Icons
import { ImageWithFallback } from '@/shared/components/ui/ImageWithFallback/ImageWithFallback';
import { Settings } from 'lucide-react';
import {
  IconFileFilled,
  IconStar,
  IconStarFilled,
  IconCloudDataConnection,
  IconDevicesPc,
  IconDeviceGamepad,
} from '@tabler/icons-react';

// Constants
import { LIBRARY_MEDIA_ITEM_BREAKPOINT_RULES } from '@/features/dashboard/lib/constants/dashboard.constants';

interface LibraryMediaListItemProps {
  index: number;
  title: string;
  coverImage?: string;
  lastPlayed?: string;
  favorite?: boolean;
  physicalLocation?: string;
  physicalLocationType?: string;
  physicalSublocation?: string;
  physicalSublocationType?: string;
  digitalLocation?: string;
  diskSize?: string;
  platformVersion?: string;
  onFavorite?: () => void;
  onSettings?: () => void;
}

const createAddToFavoritesToast = (title: string) => {
  toast(`${title} successfully added to favorites`, {
    className: 'bg-green-500 text-white',
    duration: 2500,
  });
};

function LibraryMediaListItem({
  index,
  title,
  coverImage,
  favorite = false,
  physicalLocation = "",
  physicalLocationType = "",
  physicalSublocation = "",
  physicalSublocationType = "",
  digitalLocation = "",
  platformVersion = "",
  diskSize = "",
  onFavorite,
  onSettings
}: LibraryMediaListItemProps) {
  const cardRef = useRef<HTMLDivElement>(null);

  const handleAddToFavorites = useCallback(() => {
    createAddToFavoritesToast(title);
  }, [title]);

  /* Memoize selector for useElementBreakpoint hook */
  const selector = useMemo(() =>
    `[data-library-item="${index}-${title}"]`,
    [index, title]
  );

  const {
    location: locationIcons,
    sublocation: sublocationIcons
  } = useDomainMaps();

  const locationIcon = useMemo(() => {
    if (physicalLocation && physicalLocationType) {
      const IconComponent = locationIcons[physicalLocationType.toLowerCase() as keyof typeof locationIcons];
      return IconComponent ? <IconComponent className="h-6 w-6" /> : null;
    }

    if (digitalLocation) {
      return <IconCloudDataConnection className="h-6 w-6" />
    }
    return null;
  }, [physicalLocation, physicalLocationType, digitalLocation, locationIcons])

  const subLocationIcon = useMemo(() => {
    if (physicalSublocationType) {
      const IconComponent = sublocationIcons[physicalSublocationType.toLowerCase() as keyof typeof sublocationIcons];
      return IconComponent ? <IconComponent className="h-6 w-6" /> : null;
    }

    return null;
  }, [physicalSublocationType, sublocationIcons]);

  // Visibility reducer to handle breakpoint related visibility changes
  const [visibility, dispatch] = useReducer(visibilityReducer, {
    showTags: true,
    showRating: true,
    showReleaseDate: true,
    showMoreDeals: true,
    stackPriceContent: true,
    showLocationInfo: true,
    stackInfoContent: false,
    isMobile: false,
  });

  const setVisibilityCallback = useCallback((value: CardVisibility) => {
    dispatch({ type: 'SET_VISIBILITY', payload: value });
  }, []);

  // Memoize expensive computations
  const platformIcon = useMemo(() => {
    return platformVersion === "PC" ?
      <IconDevicesPc className="h-6 w-6" /> :
      <IconDeviceGamepad className="h-6 w-6" />
  }, [platformVersion]);

  const defaultValue = useMemo(() => ({
    showTags: true,
    showRating: true,
    showReleaseDate: true,
    showMoreDeals: true,
    stackPriceContent: false,
    showLocationInfo: true,
    stackInfoContent: false,
  }), []);

  useElementBreakpoint({
    selector,
    breakpointRules: LIBRARY_MEDIA_ITEM_BREAKPOINT_RULES,
    defaultValue,
    onBreakpointChange: setVisibilityCallback
  });

  return (
    <div
      className={`flex items-center gap-4 w-full rounded-lg border bg-card p-4 text-card-foreground shadow-sm overflow-x-hidden ${index % 2 === 0 ? 'my-2' : ''}`}
      data-library-item={`${index}-${title}`}
      ref={cardRef}
    >
      {/* Game Cover */}
      <div className="h-16 w-28 flex-shrink-0">
        <ImageWithFallback
          src={coverImage}
          alt={title}
          className="h-full w-full rounded-md object-cover"
        />
      </div>

      {/* Game Info */}
      <div className="flex flex-1 items-center justify-between">
        <div className="space-y-1">
          <h3 className="font-semibold mb-2">{title}</h3>

          <div className="flex gap-8 text-sm text-muted-foreground">

            {/* Game Platform */}
            <InfoSection
              icon={platformIcon}
              label="Platform"
              value={platformVersion}
              hasStackedContent={visibility.stackInfoContent}
              isMobile={visibility.isMobile}
            />

            {/* Game Location */}
            <InfoSection
              icon={locationIcon}
              label={physicalLocation ? "Location" : "Service"}
              value={physicalLocation || digitalLocation}
              hasStackedContent={visibility.stackInfoContent}
              isMobile={visibility.isMobile}
            />

            {/* Game Sublocation */}
            <InfoSection
              icon={subLocationIcon}
              label="Sublocation"
              value={physicalSublocation}
              isVisible={!!physicalSublocation && !!physicalSublocationType}
              hasStackedContent={visibility.stackInfoContent}
              isMobile={visibility.isMobile}
            />

            {/* Game Disk Size */}
            <InfoSection
              icon={<IconFileFilled className="h-6 w-6" />}
              label="Disk Size"
              value={diskSize}
              hasStackedContent={visibility.stackInfoContent}
              isMobile={visibility.isMobile}
            />

          </div>
        </div>

        {/* Actions */}
        <div className={`flex items-center gap-2 ${visibility.isMobile ? 'flex-col' : ''}`}>
          <Button
            variant={favorite ? "default" : "secondary"}
            size="sm"
            onClick={handleAddToFavorites}
            className={`transition-colors hover:bg-[#5bf563] ${!favorite ? 'bg-muted' : ''}`}
          >
            {favorite ? (
              <IconStarFilled className="h-4 w-4" />
            ) : (
              <IconStar className="h-4 w-4" />
            )}
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={onSettings}
          >
            <Settings className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}

export const MemoizedLibraryMediaListItem = memo(LibraryMediaListItem);
