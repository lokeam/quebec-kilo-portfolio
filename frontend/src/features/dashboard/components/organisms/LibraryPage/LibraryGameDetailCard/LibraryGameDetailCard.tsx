import { memo, useReducer, useCallback, useMemo, useRef } from 'react';

// Components
import { LibraryGameDetailCardInfoSection } from '@/features/dashboard/components/organisms/LibraryPage/LibraryGameDetailCard/LibraryGameDetailCardInfoSection';
import { CoverImage } from "@/shared/components/ui/CoverImage/CoverImage";
import { Badge } from "@/shared/components/ui/badge";

// ShadCN Components
import { Button } from '@/shared/components/ui/button';

// Hooks + Utils
import { useElementBreakpoint } from '@/shared/hooks/useElementBreakpoint';
import { visibilityReducer } from '@/features/dashboard/components/organisms/WishlistPage/WishlistItemCard/visibilityReducer';
//import { toast } from 'sonner';

// Types
import type { CardVisibility } from '@/features/dashboard/lib/types/wishlist/cards';
import type { GameType, PhysicalLocationResponse, DigitalLocationResponse } from '@/types/domain/library-types';

// Icons
import {
  IconStar,
  IconStarFilled,
} from '@tabler/icons-react';

// Constants
import { LIBRARY_MEDIA_ITEM_BREAKPOINT_RULES } from '@/features/dashboard/lib/constants/dashboard.constants';
import { MemoizedMediaListItemDropDownMenu } from './MediaListItemDropDownMenu';
import { showToast } from '@/shared/components/ui/TanstackMutationToast/showToast';
import { formatReleaseDate } from '@/features/dashboard/lib/utils/libraryCardUtils';

interface LibraryGameDetailCardProps {
  index: number;
  id: number;
  name: string;
  coverUrl: string;
  firstReleaseDate?: number;
  rating: number;
  themeNames: string[] | null;
  isInLibrary: boolean;
  isInWishlist: boolean;
  gameType: GameType;
  favorite: boolean;
  physicalLocations: PhysicalLocationResponse[];
  digitalLocations: DigitalLocationResponse[];
  onFavorite?: () => void;
  onSettings?: () => void;
  onRemoveFromLibrary?: () => void;
  totalDigitalVersions?: number;
  totalPhysicalVersions?: number;
}

const createAddToFavoritesToast = (title: string) => {
  showToast({
    message: `${title} successfully added to favorites`,
    variant: 'success',
    duration: 2500,
  });
};

function LibraryGameDetailCard({
  index,
  name,
  coverUrl,
  favorite = false,
  firstReleaseDate,
  physicalLocations = [],
  digitalLocations = [],
  onRemoveFromLibrary = () => {},
  totalPhysicalVersions = 0,
  totalDigitalVersions = 0,
}: LibraryGameDetailCardProps) {
  const cardRef = useRef<HTMLDivElement>(null);

  const handleAddToFavorites = useCallback(() => {
    createAddToFavoritesToast(name);
  }, [name]);

  /* Memoize selector for useElementBreakpoint hook */
  const selector = useMemo(() =>
    `[data-library-item="${index}-${name}"]`,
    [index, name]
  );

  // Calculate total locations for the count badge
  const totalLocations = physicalLocations.length + digitalLocations.length;

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
      className={`flex items-start gap-4 w-full rounded-lg border bg-card p-4 text-card-foreground shadow-sm overflow-x-hidden ${
        index % 2 === 0 ? 'my-2' : ''}`
      }
      data-library-item={`${index}-${name}`}
      ref={cardRef}
    >
      {/* Game Cover */}
      <div className="h-[341px] w-64 flex-shrink-0">
        <CoverImage
          src={coverUrl}
          size="cover_big"
          alt={name}
          className="h-full w-full rounded-md"
        />
      </div>

      {/* Game Info */}
      <div className="flex flex-1 items-start justify-between">
        <div className="">
          <div className="flex items-center gap-2 mb-2">
            <h3 className="text-2xl font-semibold">{name}</h3>
            {totalLocations > 1 && (
              <Badge
                variant="secondary"
                className="bg-purple-500/20 text-purple-500 hover:bg-purple-500/30"
              >
                {totalLocations} versions
              </Badge>
            )}
          </div>
          <div className="flex flex-row gap-2 mb-6">Release Date: {formatReleaseDate(firstReleaseDate || 0)}</div>
          {/* DO NOT DELETE THIS COMMENT: NON FUNCTIONAL REQUIREMENT ITEMS START */}
          {/* All of these require either updated db adapter logic, db columns, db  */}
          {/* <div className="flex flex-row gap-2">Genres: (List of all genres as classified by IGDB)</div>
          <div className="flex flex-row gap-2">Published on the following platforms: (List of all platforms as classified by IGDB)</div>
          <div className="flex flex-row gap-2">If this title is not an IGDB gameType main game, show specific badge denoting IGDB game type () | Game Title</div> */}
          {/* DO NOT DELETE THIS COMMENT: NON FUNCTIONAL REQUIREMENT ITEMS END */}

          <div className="flex gap-8 text-sm text-muted-foreground">

            {/* Game Platform */}
            <LibraryGameDetailCardInfoSection
              physicalLocations={physicalLocations}
              digitalLocations={digitalLocations}
              totalPhysicalVersions={totalPhysicalVersions}
              totalDigitalVersions={totalDigitalVersions}
              hasStackedContent={visibility.stackInfoContent}
              isMobile={visibility.isMobile}
            />
          </div>
        </div>

        {/* Actions */}
        <div className={`flex items-start gap-2 ${visibility.isMobile ? 'flex-col' : ''}`}>
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

export const MemoizedLibraryGameDetailCard = memo(LibraryGameDetailCard);
