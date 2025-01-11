// import { Heart } from 'lucide-react'
import { IconX } from '@tabler/icons-react'
import { Badge } from "@/shared/components/ui/badge"
import { Button } from "@/shared/components/ui/button"
import { Card } from "@/shared/components/ui/card"
import { useState } from 'react'
import { useElementBreakpoint } from '@/shared/hooks/useElementBreakpoint'

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
    discounted: number;
    discountPercentage: number;
    vendor: string;
  };
  platform: string;
  index: number;
};

interface CardVisibility {
  showTags: boolean;
  showRating: boolean;
  showReleaseDate: boolean;
}

export function WishlistItemCard({
  title,
  thumbnailUrl,
  tags,
  releaseDate,
  rating,
  price,
  index
}: WishlistItemCardProps) {
  const [visibility, setVisibility] = useState<CardVisibility>({
    showTags: true,
    showRating: true,
    showReleaseDate: true
  });

  useElementBreakpoint({
    selector: '[data-wishlist-item]',
    breakpointRules: [
      {
        breakpoint: 865,
        value: {
          showTags: false,
          showRating: false,
          showReleaseDate: false
        }
      }
    ],
    defaultValue: {
      showTags: true,
      showRating: true,
      showReleaseDate: true
    },
    onBreakpointChange: setVisibility
  });

  return (
    <Card
      data-wishlist-item
      className={`bg-[#1a1b1f] text-white overflow-hidden w-full ${index % 2 === 0 ? 'mb-2' : ''}`}
    >
      <div className="flex gap-4 p-4 h-[190px]">
        {/* Game Thumbnail */}
        <div className="relative shrink-0 h-full">
          <img
            src={thumbnailUrl}
            alt={title}
            width={292}
            height={136}
            className="rounded-sm w-[292px] h-full object-cover max-lg:w-[100px]"
          />
        </div>

        {/* Game Info */}
        <div className="flex-1 min-w-0 flex flex-col justify-between gap-2">
          <div className="flex justify-between items-start gap-4">
            <div className="space-y-2 flex-1 min-w-0">
              <h2 className="text-2xl font-semibold">{title}</h2>

              {/* Tags */}
              {visibility.showTags && (
                <div className="flex flex-wrap gap-2">
                  {tags.map((tag, index) => (
                    <Badge
                      key={index}
                      variant="secondary"
                      className="bg-[#42464e] hover:bg-[#42464e] rounded-none"
                    >
                      {tag}
                    </Badge>
                  ))}
                </div>
              )}
            </div>

            {/* Wishlist Button */}
            <Button
              size="icon"
              className="text-gray-400 bg-gray-700 hover:text-white hover:bg-red-800"
            >
              <IconX className="h-8 w-8" />
            </Button>
          </div>

          {/* Bottom Row */}
          <div className="flex items-end justify-between mt-auto">
            {/* Release Date and Rating */}
            <div className="flex flex-col gap-1">
              {/* Release Date */}
              {visibility.showReleaseDate && (
                <div className="flex items-center gap-2 text-sm text-gray-400">
                  <svg
                    viewBox="0 0 16 16"
                    className="w-4 h-4"
                    fill="currentColor"
                  >
                    <path d="M0 4v8h16V4H0zm15 7H1V7h14v4z"/>
                  </svg>
                  {releaseDate}
                </div>
              )}

              {/* Rating Bar */}
              {visibility.showRating && (
                <div className="space-y-1">
                  <div className="h-2 w-32 bg-gray-700 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-blue-500 rounded-full"
                      style={{ width: `${rating.positive}%` }}
                    />
                  </div>
                  <div className="flex items-center gap-1 text-xs text-gray-400">
                    <span className="text-blue-500">{rating.positive}%</span>
                    <span className="text-red-500">{rating.negative}%</span>
                    <span className="ml-1">
                      {rating.totalReviews.toLocaleString()} User Reviews
                    </span>
                  </div>
                </div>
              )}
            </div>

            {/* Price and CTA */}
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <div className="bg-[#94d933] text-black font-bold hover:bg-[#567b27] rounded-sm py-2 px-3 mr-3">
                  -{price.discountPercentage}%
                </div>
                <div className="flex flex-col">
                  <div className="text-gray-400 line-through text-sm">
                    ${price.original.toFixed(2)}
                  </div>
                  <div className="text-lg font-bold">
                    ${price.discounted.toFixed(2)}
                  </div>
                </div>
              </div>
              <div className="flex flex-col gap-2">
                <Button className="bg-[#4c6b22] hover:bg-[#567b27] text-white">
                  Purchase on {price.vendor}
                </Button>
                <Button className="bg-[#492fef] hover:bg-[#632de1] text-white">
                  See more deals
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Card>
  );
}

