import { useState } from 'react';

// Components
import { Card } from '@/shared/components/ui/card';
import { Button } from '@/shared/components/ui/button';
import { Badge } from '@/shared/components/ui/badge';
import { ImageWithFallback } from '@/shared/components/ui/ImageWithFallback/ImageWithFallback';

// Type Refactor
import type { Game } from '@/types/domain/game';

// Legacy Types
//import type { WishlistItem } from '@/features/dashboard/lib/types/wishlist/base';

// NOTE: Commented out because we don't have a way to add to library or wishlist yet
// Replace these legacy hooks with the new ones when we have them
// import { useAddToLibrary, useAddToWishlist } from '@/core/api/queries/useLibraryMutations';

// Icons
import { LibraryBig } from '@/shared/components/ui/icons';
import { IconHeart } from '@/shared/components/ui/icons';

// Hooks
//import { toast } from 'sonner';

// New
type SearchResultProps = {
  game: Game;
  onAction?: () => void; // Callback to close dialog
  onAddToLibrary?: (game: Game) => void;
}

// Helper function to format release date
const formatReleaseDate = (timestamp?: number): string => {
  if (!timestamp) return '';
  const date = new Date(timestamp * 1000); // Convert Unix timestamp to milliseconds
  return `(${date.getFullYear()})`;
};

// Helper function to get badge variant based on game type
const getGameTypeBadgeVariant = (normalizedText?: string): 'default' | 'secondary' | 'destructive' | 'outline' => {
  switch (normalizedText) {
    case 'main':
      return 'default';
    case 'dlc':
      return 'secondary';
    case 'expansion':
      return 'destructive';
    default:
      return 'outline';
  }
};

export function SearchResult({
  game,
  onAction,
  onAddToLibrary,
}: SearchResultProps) {
  const [imageLoaded, setImageLoaded] = useState(false);
  // const addToLibrary = useAddToLibrary();
  // const addToWishList = useAddToWishlist();

  // Debug the game object and firstReleaseDate
  console.log('--------------------------------');
  console.log('SearchResult onAddToLibrary: ', onAddToLibrary);
  console.log('--------------------------------');
  console.log('Game object:', game);
  console.log('First release date:', game.firstReleaseDate);
  console.log('Formatted date:', game.firstReleaseDate ? formatReleaseDate(game.firstReleaseDate) : 'No date');

  const handleAddToLibrary = () => {
    console.log('Adding to library fired for selected game:', game);

    if (onAddToLibrary) {
      onAddToLibrary(game);
    }

    // Optinally call onAction to close dialog or handle other side effect
    if (onAction) onAction();

    // Move this to mutation success
    // toast( `${game.name} successfully added to library`,{
    //   description: 'You can now access it in your library',
    //   className: 'bg-green-500 text-white',
    //   duration: 2500,
    // });
  };

  const handleAddToWishlist = () => {
    console.log('Adding to wishlist', game);

    if (onAction) onAction();

    // toast(`${game.name} successfully added to wishlist`, {
    //   description: 'You can now access it in your wishlist',
    //   className: 'bg-green-500 text-white',
    //   duration: 2500,
    // });
  };

  const showLibraryButton = !game.isInLibrary;
  // Mods are free, so we don't need to show the wishlist button
  const showWishlistButton = !game.isInWishlist && game.gameType?.normalizedText !== 'mod';

  return (
    <Card className="relative flex items-center transition-all duration-200 bg-[#2A2A2A] hover:bg-[#E5E5E5] group overflow-hidden opacity-0 data-[loaded=true]:opacity-100" data-loaded={imageLoaded}>
      { game.isInLibrary && (
        <div className="absolute left-0 top-0 bottom-0 z-10">
          <div className="flex h-full items-center">
            <div className="flex items-center bg-[#1A9FFF] h-[34px] overflow-hidden transition-all duration-200 group-hover:w-[120px] w-[34px] rounded-r-md">
              <div className="flex items-center gap-2 px-2 w-[120px]">
                <LibraryBig className="w-5 h-5 text-black shrink-0" />
                <span className="text-sm font-medium text-black whitespace-nowrap">IN LIBRARY</span>
              </div>
            </div>
          </div>
        </div>
      )}
      <div className="shrink-0 p-2">
        <div className="relative w-24 md:w-32 h-[130px] p-2">
          <ImageWithFallback
            src={game.coverUrl}
            alt={`cover image for ${game.name}`}
            width={292}
            height={130}
            className="rounded-sm w-full h-full object-cover"
            onLoad={() => setImageLoaded(true)}
            igdbSize="cover_big"
          />
        </div>
      </div>

      <div className="flex flex-1 min-w-0 items-center pr-2">
        <div className="flex-1 min-w-0"> {/* nested min-w-0 for text truncation */}
          <div className="flex items-center gap-2">
            <h3 className="font-medium text-white text-wrap max-w-[140px] max-h-[48px] md:max-w-full md:max-h-unset truncate">
              {game.name}
              {game.firstReleaseDate && (
                <time className="text-gray-400 text-sm font-normal ml-1" dateTime={new Date(game.firstReleaseDate * 1000).toISOString()}>
                  {formatReleaseDate(game.firstReleaseDate)}
                </time>
              )}
            </h3>
          </div>
          {game.platforms && game.platforms.length > 0 && (
            <div className="text-gray-400 text-sm mt-1">
              {game.platforms.map((platform, index) => (
                <span key={platform.id}>
                  {platform.name}
                  {index < (game.platforms?.length || 0) - 1 && (
                    <span className="mx-1">/</span>
                  )}
                </span>
              ))}
            </div>
          )}
          {game.gameType && (
            <div className="text-gray-400 text-sm mt-3">
              <Badge variant={getGameTypeBadgeVariant(game.gameType.normalizedText)}>
                {game.gameType.displayText}
              </Badge>
            </div>
          )}
        </div>

        {
          (
            <div className="flex shrink-0 gap-1 mt-1 ml-2">
            {/* Add to Library Button */}
              {showLibraryButton && (
                <Button variant="outline" onClick={handleAddToLibrary}>
                  <LibraryBig className="w-5 h-5" />
                  <span className="hidden md:block text-sm font-medium text-white whitespace-nowrap">Add to library</span>
                </Button>
              )}
            {/* Add to Wishlist Button */}
              {/* {showWishlistButton && (
                <Button variant="outline" onClick={handleAddToWishlist}>
                  <IconHeart className="w-5 h-5" />
                  <span className="hidden md:block text-sm font-medium text-white whitespace-nowrap">Add to wishlist</span>
                </Button>
              )} */}
            </div>
          )
        }
      </div>
    </Card>
  );
}
