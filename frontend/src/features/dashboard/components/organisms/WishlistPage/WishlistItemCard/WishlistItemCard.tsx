import { memo, useMemo, useCallback, useRef, useReducer } from 'react';
import { isEqual } from 'lodash';

// ShadCN Components
import { Button } from "@/shared/components/ui/button"
import { Card } from "@/shared/components/ui/card"

// Components
import { PriceSection } from './PriceSection';
import { RatingSection } from './RatingSection';
import { TagsSection } from './TagsSection';
import { ReleaseDateSection } from './ReleaseDateSection';

// Hooks + Utils
import { useElementBreakpoint } from '@/shared/hooks/useElementBreakpoint'
import { visibilityReducer } from './visibilityReducer';
import { toast } from 'sonner';

// Icons
import { IconX } from '@tabler/icons-react'

// Types
import type { CardVisibility } from './WishlistCardItem.types';


export type WishlistItemCardProps = {
  id: string;
  title: string;
  thumbnailUrl: string;
  tags: string[];
  releaseDate: string;
  rating: {
    positive: number;
    negative: number;
    totalReviews: number;
  };
  price: {
    original: number;
    discounted?: number | undefined | null;
    discountPercentage?: number | undefined | null;
    vendor: string;
  };
  platform: string;
  hasMacOSVersion?: boolean | undefined | null;
  hasAndroidVersion?: boolean | undefined | null;
  hasIOSVersion?: boolean | undefined | null;
  index: number;
};


// Handlers for wishlist removal
const createRemoveFromWishlistToast = (title: string) => {
  toast(`${title} successfully removed from wishlist`, {
    className: 'bg-green-500 text-white',
    duration: 2500,
  });
};

export const WishlistItemCard = memo(({
  id,
  title,
  thumbnailUrl,
  tags,
  releaseDate,
  rating,
  platform,
  price,
  index,
  hasAndroidVersion,
  hasIOSVersion,
  hasMacOSVersion
}: WishlistItemCardProps) => {
  // Edit: Move visibility state to useReducer for better performance
  const [visibility, dispatch] = useReducer(visibilityReducer, {
    showTags: true,
    showRating: true,
    showReleaseDate: true,
    showMoreDeals: true,
    stackPriceContent: false,
  });

  const cardRef = useRef<HTMLDivElement>(null);

  const handleRemoveFromWishlist = useCallback(() => {
    // TODO: Wire-up handler for wishlist removal to query the backend
    createRemoveFromWishlistToast(title);
  }, [title]);

  const selector = useMemo(() =>
    `[data-wishlist-item="${index}-${id}"]`,
    [index, id]
  );

  // Visibility reducer to handle breakpoint related visibility changes
  const setVisibilityCallback = useCallback((value: CardVisibility) => {
    dispatch({ type: 'SET_VISIBILITY', payload: value });
  }, []);

  const breakpointRules = useMemo(() => [
    {
      breakpoint: 865,
      value: {
        showTags: false,
        showRating: false,
        showReleaseDate: false,
        showMoreDeals: true,
        stackPriceContent: false,
      },
    },
    {
      breakpoint: 630,
      value: {
        showTags: false,
        showRating: false,
        showReleaseDate: false,
        showMoreDeals: false,
        stackPriceContent: true,
      }
    }
  ], []);

  const defaultValue = useMemo(() => ({
    showTags: true,
    showRating: true,
    showReleaseDate: true,
    showMoreDeals: true,
    stackPriceContent: false,
  }), []);

  useElementBreakpoint({
    selector,
    breakpointRules,
    defaultValue,
    onBreakpointChange: setVisibilityCallback
  });

  // Memoize the card content to prevent unnecessary re-renders
  const cardContent = useMemo(() => (
    <div className="flex-1 min-w-0 flex flex-col justify-between gap-2">
      <div className="flex justify-between items-start gap-4">
        <div className="space-y-2 flex-1 min-w-0">
          <h2 className="text-2xl font-semibold">{title}</h2>
          {visibility.showTags && <TagsSection tags={tags} />}
        </div>

        <Button
          size="icon"
          className="text-gray-400 bg-gray-700 hover:text-white hover:bg-red-800"
          onClick={handleRemoveFromWishlist}
        >
          <IconX className="h-8 w-8" />
        </Button>
      </div>

      <div className="flex items-end justify-between mt-auto">
        <div className="flex flex-col justify-between h-full gap-1">
          {visibility.showReleaseDate && (
            <ReleaseDateSection
              platform={platform}
              releaseDate={releaseDate}
              hasAndroidVersion={hasAndroidVersion ?? false}
              hasIOSVersion={hasIOSVersion ?? false}
              hasMacOSVersion={hasMacOSVersion ?? false}
            />
          )}
          {visibility.showRating && <RatingSection {...rating} />}
        </div>

        <PriceSection
          price={price}
          showMoreDeals={visibility.showMoreDeals}
          stackPriceContent={visibility.stackPriceContent}
        />
      </div>
    </div>
  ), [
    title,
    tags,
    visibility,
    handleRemoveFromWishlist,
    platform,
    releaseDate,
    hasAndroidVersion,
    hasIOSVersion,
    hasMacOSVersion,
    rating,
    price
  ]);


  return (
    <Card
      ref={cardRef}
      data-wishlist-item={`${index}-${id}`}
      className={`bg-[#1a1b1f] text-white overflow-hidden w-full ${index % 2 === 0 ? 'mb-2' : ''}`}
    >
      <div className="flex gap-4 p-4 h-[190px]">
        <div className="relative shrink-0 h-full">
          <img
            src={thumbnailUrl}
            alt={title}
            width={292}
            height={136}
            className="rounded-sm w-[292px] h-full object-cover max-lg:w-[100px]"
          />
        </div>
        {cardContent}
      </div>
    </Card>
  );
}, (prevProps, nextProps) => {
  return (
    prevProps.id === nextProps.id &&
    prevProps.index === nextProps.index &&
    prevProps.title === nextProps.title &&
    prevProps.thumbnailUrl === nextProps.thumbnailUrl &&
    prevProps.platform === nextProps.platform &&
    prevProps.releaseDate === nextProps.releaseDate &&
    isEqual(prevProps.price, nextProps.price) &&
    isEqual(prevProps.rating, nextProps.rating) &&
    isEqual(prevProps.tags, nextProps.tags) &&
    prevProps.hasAndroidVersion === nextProps.hasAndroidVersion &&
    prevProps.hasIOSVersion === nextProps.hasIOSVersion &&
    prevProps.hasMacOSVersion === nextProps.hasMacOSVersion
  );
});
