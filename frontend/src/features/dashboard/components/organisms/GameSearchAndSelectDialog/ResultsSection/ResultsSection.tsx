import { useCallback } from 'react';

// Shadcn Components
import { Card } from '@/shared/components/ui/card';
import { Button } from '@/shared/components/ui/button';
import { ImageWithFallback } from '@/shared/components/ui/ImageWithFallback/ImageWithFallback';

// Icons
import { LibraryBig } from 'lucide-react';
import { IconHeart } from '@tabler/icons-react';

// Types
import type { SearchResult } from '@/types/domain/search';

interface ResultsSectionProps {
  results?: SearchResult[];
  isLoading: boolean;
  error: Error | null;
  onSelect: (gameId: string) => void;
}

export function ResultsSection({ results, isLoading, error, onSelect }: ResultsSectionProps) {
  const handleAddToLibrary = useCallback((gameId: number) => {
    // TODO: Implement add to library
    onSelect(String(gameId));
  }, [onSelect]);

  const handleAddToWishlist = useCallback((gameId: number) => {
    // TODO: Implement add to wishlist
    onSelect(String(gameId));
  }, [onSelect]);

  if (isLoading) {
    return (
      <div className="space-y-4">
        {Array.from({ length: 3 }).map((_, index) => (
          <Card key={index} className="h-24 animate-pulse bg-muted" />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-red-500 p-4">
        Error loading results: {error.message}
      </div>
    );
  }

  if (!results) return null;

  return (
    <div className="space-y-4">
      {results.map((result) => {
        const { game } = result;
        const showLibraryButton = !game.isInLibrary;
        const showWishlistButton = !game.isInWishlist;

        return (
          <Card key={game.id} className="relative flex items-center transition-all duration-200 bg-[#2A2A2A] hover:bg-[#E5E5E5] group overflow-hidden">
            {game.isInLibrary && (
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
              <div className="flex-1 min-w-0">
                <h3 className="font-medium text-white text-wrap max-w-[140px] max-h-[48px] md:max-w-full md:max-h-unset truncate">
                  {game.name}
                </h3>
              </div>

              <div className="flex shrink-0 gap-1 mt-1 ml-2">
                {showLibraryButton && (
                  <Button variant="outline" onClick={() => handleAddToLibrary(game.id)}>
                    <LibraryBig className="w-5 h-5" />
                    <span className="hidden md:block text-sm font-medium text-white whitespace-nowrap">Add to library</span>
                  </Button>
                )}
                {showWishlistButton && (
                  <Button variant="outline" onClick={() => handleAddToWishlist(game.id)}>
                    <IconHeart className="w-5 h-5" />
                    <span className="hidden md:block text-sm font-medium text-white whitespace-nowrap">Add to wishlist</span>
                  </Button>
                )}
              </div>
            </div>
          </Card>
        );
      })}
    </div>
  );
}