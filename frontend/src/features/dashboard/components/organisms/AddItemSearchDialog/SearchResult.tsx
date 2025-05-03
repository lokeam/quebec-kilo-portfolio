
// Components
import { Card } from '@/shared/components/ui/card';
import { Button } from '@/shared/components/ui/button';
import { ImageWithFallback } from '@/shared/components/ui/ImageWithFallback/ImageWithFallback';

// Type Refactor
import type { Game } from '@/types/game';

// Legacy Types
//import type { WishlistItem } from '@/features/dashboard/lib/types/wishlist/base';


import { useAddToLibrary, useAddToWishlist } from '@/core/api/queries/useLibraryMutations';

// Icons
import { LibraryBig } from 'lucide-react';
import { IconHeart } from '@tabler/icons-react';

// Hooks
//import { toast } from 'sonner';

// New
type SearchResultProps = {
  game: Game;
  onAction?: () => void; // Callback to close dialog
}


export function SearchResult({ game, onAction}: SearchResultProps) {
  const addToLibrary = useAddToLibrary();
  const addToWishList = useAddToWishlist();

  const handleAddToLibrary = () => {
    addToLibrary.mutate({
      id: Number(game.id),
      name: game.name,
      cover_url: game.coverUrl,
      rating: game.rating ? Number(game.rating) : undefined,
      theme_names: game.themeNames ? [...game.themeNames] : undefined,
    })

    if (onAction) onAction();

    // Move this to mutation success
    // toast( `${game.name} successfully added to library`,{
    //   description: 'You can now access it in your library',
    //   className: 'bg-green-500 text-white',
    //   duration: 2500,
    // });
  };

  const handleAddToWishlist = () => {
    addToWishList.mutate({
      id: Number(game.id),
      name: game.name,
      cover_url: game.coverUrl,
      rating: game.rating ? Number(game.rating) : undefined,
      theme_names: game.themeNames ? [...game.themeNames] : undefined,
    })

    if (onAction) onAction();

    // toast(`${game.name} successfully added to wishlist`, {
    //   description: 'You can now access it in your wishlist',
    //   className: 'bg-green-500 text-white',
    //   duration: 2500,
    // });
  };

  const showLibraryButton = !game.isInLibrary;
  const showWishlistButton = !game.isInWishlist;

  console.log('SearchResult', game);
  return (
    <Card className="relative flex items-center transition-all duration-200 bg-[#2A2A2A] hover:bg-[#E5E5E5] group overflow-hidden">
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
        <div className="relative w-24 md:w-32 p-2">
          <ImageWithFallback
            src={game.coverUrl}
            alt={`cover image for ${game.name}`}
            width={292}
            height={120}
            className="rounded-sm w-[140px] h-full object-cover"
          />
        </div>
      </div>

      <div className="flex flex-1 min-w-0 items-center pr-2">
        <div className="flex-1 min-w-0"> {/* nested min-w-0 for text truncation */}
          <h3 className="font-medium text-white text-wrap max-w-[140px] max-h-[48px] md:max-w-full md:max-h-unset truncate">
            {game.name}
          </h3>
        </div>

        {
          (
            <div className="flex shrink-0 gap-1 mt-1 ml-2">
              {showLibraryButton && (
                <Button variant="outline" onClick={handleAddToLibrary}>
                  <LibraryBig className="w-5 h-5" />
                  <span className="hidden md:block text-sm font-medium text-white whitespace-nowrap">Add to library</span>
                </Button>
              )}
              {showWishlistButton && (
                <Button variant="outline" onClick={handleAddToWishlist}>
                  <IconHeart className="w-5 h-5" />
                  <span className="hidden md:block text-sm font-medium text-white whitespace-nowrap">Add to wishlist</span>
                </Button>
              )}
            </div>
          )
        }
      </div>
    </Card>
  );
}
