import { memo, useReducer, useCallback, useMemo, useRef } from 'react';

// Components
import { InfoSection } from '@/features/dashboard/components/organisms/LibraryPage/LibraryMediaListItem/InfoSection';
import { CoverImage } from "@/shared/components/ui/CoverImage/CoverImage";
import { Badge } from "@/shared/components/ui/badge";

// ShadCN Components
import { Button } from '@/shared/components/ui/button';

// Hooks + Utils
import { useLocationIcons } from '@/features/dashboard/lib/hooks/useLocationIcons';
import { useElementBreakpoint } from '@/shared/hooks/useElementBreakpoint';
import { visibilityReducer } from '@/features/dashboard/components/organisms/WishlistPage/WishlistItemCard/visibilityReducer';
//import { toast } from 'sonner';

// Types
import type { CardVisibility } from '@/features/dashboard/lib/types/wishlist/cards';
import type { GamePlatformLocation, GamePlatformLocationResponse, GameType } from '@/types/domain/library-types';

// Icons
import {
  IconFileFilled,
  IconStar,
  IconStarFilled,
  IconDevicesPc,
  IconDeviceGamepad,
} from '@tabler/icons-react';

// Constants
import { LIBRARY_MEDIA_ITEM_BREAKPOINT_RULES } from '@/features/dashboard/lib/constants/dashboard.constants';
import { MemoizedMediaListItemDropDownMenu } from './MediaListItemDropDownMenu';
import { showToast } from '@/shared/components/ui/TanstackMutationToast/showToast';

interface LibraryMediaListItemProps {
  index: number;
  id: number;
  name: string;
  coverUrl: string;
  firstReleaseDate: number;
  rating: number;
  themeNames: string[] | null;
  isInLibrary: boolean;
  isInWishlist: boolean;
  gameType: GameType;
  favorite: boolean;
  gamesByPlatformAndLocation: GamePlatformLocationResponse[];
  onFavorite?: () => void;
  onSettings?: () => void;
  onRemoveFromLibrary?: () => void;
}

const createAddToFavoritesToast = (title: string) => {
  showToast({
    message: `${title} successfully added to favorites`,
    variant: 'success',
    duration: 2500,
  });
};

function LibraryMediaListItem({
  index,
  name,
  coverUrl,
  favorite = false,
  gamesByPlatformAndLocation = [],
  onRemoveFromLibrary = () => {},
}: LibraryMediaListItemProps) {
  const cardRef = useRef<HTMLDivElement>(null);

  const handleAddToFavorites = useCallback(() => {
    createAddToFavoritesToast(name);
  }, [name]);

  /* Memoize selector for useElementBreakpoint hook */
  const selector = useMemo(() =>
    `[data-library-item="${index}-${name}"]`,
    [index, name]
  );

  const { locationIcon, subLocationIcon } = useLocationIcons({
    gamesByPlatformAndLocation,
    selectedIndex: 0
  });

  // Get the first platform/location for display
  const selectedLocation = gamesByPlatformAndLocation[0];

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
    return selectedLocation?.PlatformName === "PC" ?
      <IconDevicesPc className="h-8 w-8 mt-[-4px]" /> :
      <IconDeviceGamepad className="h-7 w-7" />
  }, [selectedLocation?.PlatformName]);

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
      className={`flex items-center gap-4 w-full rounded-lg border bg-card p-4 text-card-foreground shadow-sm overflow-x-hidden ${
        index % 2 === 0 ? 'my-2' : ''}`
      }
      data-library-item={`${index}-${name}`}
      ref={cardRef}
    >
      {/* Game Cover */}
      <div className="h-16 w-28 flex-shrink-0">
        <CoverImage
          src={coverUrl}
          size="cover_small"
          alt={name}
          className="h-full w-full rounded-md"
        />
      </div>

      {/* Game Info */}
      <div className="flex flex-1 items-center justify-between">
        <div className="space-y-1">
          <div className="flex items-center gap-2 mb-2">
            <h3 className="font-semibold">{name}</h3>
            {gamesByPlatformAndLocation.length > 1 && (
              <Badge
                variant="secondary"
                className="bg-purple-500/20 text-purple-500 hover:bg-purple-500/30"
              >
                {gamesByPlatformAndLocation.length} versions
              </Badge>
            )}
          </div>

          <div className="flex gap-8 text-sm text-muted-foreground">

            {/* Game Platform */}
            <InfoSection
              icon={platformIcon}
              label="Platform"
              value={selectedLocation?.PlatformName ?? ""}
              hasStackedContent={visibility.stackInfoContent}
              isMobile={visibility.isMobile}
            />

            {/* Game Location */}
            <InfoSection
              icon={locationIcon}
              label={selectedLocation?.Type === 'physical' ? "Location" : "Service"}
              value={selectedLocation?.LocationName ?? ""}
              hasStackedContent={visibility.stackInfoContent}
              isMobile={visibility.isMobile}
            />

            {/* Game Sublocation */}
            <InfoSection
              icon={subLocationIcon}
              label="Sublocation"
              value={selectedLocation?.SublocationName ?? ""}
              isVisible={!!selectedLocation?.SublocationName && !!selectedLocation?.SublocationType}
              hasStackedContent={visibility.stackInfoContent}
              isMobile={visibility.isMobile}
            />

            {/* Game Disk Size - Only show for digital games */}
            {selectedLocation?.Type === 'digital' && (
              <InfoSection
                icon={<IconFileFilled className="h-7 w-7" />}
                label="Disk Size"
                value="0 GB" // TODO: Add disk size to the API response
                hasStackedContent={visibility.stackInfoContent}
                isMobile={visibility.isMobile}
              />
            )}

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
          <MemoizedMediaListItemDropDownMenu onRemoveFromLibrary={onRemoveFromLibrary} />
        </div>
      </div>
    </div>
  );
}

export const MemoizedLibraryMediaListItem = memo(LibraryMediaListItem);
