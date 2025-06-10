// Components
import { InfoSection } from "@/features/dashboard/components/organisms/LibraryPage/LibraryMediaListItem/InfoSection"
import { CoverImage } from "@/shared/components/ui/CoverImage/CoverImage"
import { LibraryCountIcon } from "@/features/dashboard/lib/utils/getLibraryItemCountIcon"

// Shadcn UI components
import { Card } from "@/shared/components/ui/card"

// Hooks
import { cn } from "@/shared/components/ui/utils"

// Icons
import { useLocationIcons } from '@/features/dashboard/lib/hooks/useLocationIcons';
import { IconFileFilled } from "@tabler/icons-react"
import { LibraryItemContextMenu } from '../LibraryMediaListItem/LibraryItemContextMenu';

// Types
import type { GamePlatformLocation, GamePlatformLocationResponse, GameType } from '@/types/domain/library-types';

interface LibraryMediaItemProps {
  index?: number;
  id: number;
  steamHref?: string
  name: string;
  coverUrl: string;
  firstReleaseDate: number;
  rating: number;
  themeNames: string[] | null;
  isInLibrary: boolean;
  isInWishlist: boolean;
  gameType: GameType;
  favorite: boolean;
  gamesByPlatformAndLocation?: GamePlatformLocationResponse[];
  onRemoveFromLibrary?: () => void;
}

export function LibraryMediaItem({
  steamHref,
  coverUrl,
  name,
  gamesByPlatformAndLocation = [],
  onRemoveFromLibrary = () => {},
}: LibraryMediaItemProps) {
  const { locationIcon, subLocationIcon } = useLocationIcons({
    gamesByPlatformAndLocation,
    selectedIndex: 0 // For now, just show the first platform/location
  });

  // Get the first platform/location for display
  const selectedLocation = gamesByPlatformAndLocation?.[0];

  const content = (
    <div className={cn(
      "group relative aspect-[11/15] w-full bg-black shadow-[0px_3px_10px_rgba(0,0,0,0.9)]",
      "overflow-hidden cursor-pointer mx-2.5",
      "[transform:perspective(450px)_rotateX(0deg)_scale(1)_translateY(0px)]",
      "[transform-origin:top_50%] transition-[transform,box-shadow] duration-200 ease-in",
      "hover:[transform:perspective(450px)_rotateX(5deg)_scale(1.05)_translateY(-4px)]",
      "hover:shadow-[0px_8px_20px_rgba(0,0,0,0.9)]"
    )}>
      <Card className="w-full h-full rounded-none overflow-hidden">
        <CoverImage
          src={coverUrl}
          size="cover_big"
          alt={name}
          className="w-full h-full"
        />
        <div className="absolute top-2 right-2 z-10">
          <LibraryCountIcon
            count={gamesByPlatformAndLocation.length}
            className="h-6 w-6 text-white drop-shadow-[0_2px_2px_rgba(0,0,0,0.8)]"
          />
        </div>
        <div className="card-gradient absolute left-0 top-[-35%] h-full w-full opacity-10 transition-all duration-400 group-hover:top-0 group-hover:opacity-15"
          style={{
            background: 'linear-gradient(30deg, rgba(0, 0, 0, 0), rgba(0, 0, 0, 0) 50%, rgb(255,255,255) 55%)'
          }}
        />
      </Card>
      <div
        className="absolute bottom-0 left-0 right-0 p-5 pb-2.5 text-xl text-white
                   transform translate-y-full transition-transform duration-200 ease-out
                   group-hover:translate-y-0 backdrop-blur-[5px] bg-black bg-opacity-50 space-y-2"
      >
        {/* Game Location */}
        {selectedLocation && (
          <>
            <InfoSection
              icon={locationIcon}
              label={selectedLocation.Type === 'physical' ? "Location" : "Service"}
              value={selectedLocation.LocationName ?? ""}
              hasStackedContent={false}
              isMobile={false}
            />

            {/* Game Sublocation */}
            <InfoSection
              icon={subLocationIcon}
              label="Sublocation"
              value={selectedLocation.SublocationName ?? ""}
              isVisible={!!selectedLocation.SublocationName && !!selectedLocation.SublocationType}
              hasStackedContent={false}
              isMobile={false}
              isCardView={true}
            />

            {/* Disk Size - Only show for digital games */}
            {selectedLocation.Type === 'digital' && (
              <InfoSection
                icon={<IconFileFilled className="h-7 w-7" />}
                label="Disk Size"
                value="0 GB" // TODO: Add disk size to the API response
                hasStackedContent={false}
                isMobile={false}
              />
            )}
          </>
        )}
      </div>
      <div
        data-testid="library-media-item-shine"
        className="absolute top-0 left-0 w-[200%] h-[300px]
                   bg-gradient-to-b from-white/80 to-transparent
                   [transform:translate3d(0px,0px,0px)_rotate(45deg)]
                   [transform-origin:top_right] transition-transform duration-300 ease-out
                   group-hover:[transform:translate3d(0px,100px,0px)_rotate(45deg)]"
      />
    </div>
  );

  return (
    <LibraryItemContextMenu onRemoveFromLibrary={onRemoveFromLibrary}>
      <a href={steamHref} className="w-full sm:w-1/2 md:w-1/3 lg:w-1/4 xl:w-1/6 p-2">
        {content}
      </a>
    </LibraryItemContextMenu>
  )
}
