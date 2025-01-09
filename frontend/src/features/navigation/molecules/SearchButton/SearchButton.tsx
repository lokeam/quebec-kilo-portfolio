import { useState, useEffect, useCallback, useMemo } from 'react';

// Components
import { SearchResult } from './SearchResult';

// Shadcn Components
import { Button } from "@/shared/components/ui/button"
import {
  Dialog,
  DialogOverlay,
  DialogTrigger,
  DialogContent,
  DialogPortal,
  DialogTitle,
} from '@/shared/components/ui/dialog';
import { VisuallyHidden } from '@radix-ui/react-visually-hidden';
import { Input } from '@/shared/components/ui/input';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { cn } from '@/shared/components/ui/utils';

// Hooks
import { useSearchGames } from '@/features/dashboard/hooks/useSearchGames';
import { useDebounce } from '@/shared/hooks/useDebounce';

// Icons
import { SearchIcon } from 'lucide-react';

// TODO: Replace memoization with virtualized list after wiring up backend
export function SearchButton() {
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const debouncedSearchQuery = useDebounce(searchQuery, 400);
  const { games, isLoading, error } = useSearchGames(debouncedSearchQuery);

  // Memoize the handleOpenChange callback
  const handleOpenChange = useCallback((open: boolean) => {
    setIsOpen(open);
    if (!open) {
      setSearchQuery('');
    }
  }, []);

  // Memoize the onChange handler
  const handleSearchQueryChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(event.target.value);
  }, []);

  // Move the useEffect outside of render cycle
  useEffect(() => {
    if (games?.length > 0) {
      setSearchResults(games);
    }
  }, [games]);

  // Memoize displaySearchResults calculation
  const displaySearchResults = useMemo(() => {
    if (games?.length > 0) return games;
    if (searchResults.length > 0) return searchResults;
    return [];
  }, [games, searchResults]);

  // Memoize the loading skeleton
  const LoadingSkeleton = useMemo(() => (
    Array(3).fill(0).map((_, index) => (
      <div key={index} className="p-4 border rounded-lg">
        <div className="flex items-center space-x-4">
          <Skeleton className="h-12 w-12 rounded" />
          <div className="space-y-2">
            <Skeleton className="h-4 w-[200px]" />
            <Skeleton className="h-4 w-[100px]" />
          </div>
        </div>
      </div>
    ))
  ), []);

  return (
    <Dialog open={isOpen} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button variant="outline" size="icon">
          <SearchIcon />
        </Button>
      </DialogTrigger>

      <DialogPortal>
        <DialogOverlay className="fixed inset-0 bg-black/50 data-[state=open]:animate-fadeIn data-[state=closed]:animate-fadeOut" />

        {/* Container div to control exact positioning */}
        <div className="fixed inset-0 overflow-hidden pt-20">
          <div className="flex items-start justify-center">
            <DialogContent className={cn(
              "fixed left-[50%] top-[5%] z-50 w-[90vw] max-w-[940px] -translate-x-[50%] translate-y-0 bg-background rounded-lg shadow-lg",
              // Remove ALL animation classes except simple fade
              "transition-opacity duration-200",
              "opacity-0 data-[state=open]:opacity-100"
            )}>
              <div className="flex flex-col h-[calc(100vh-120px)]">
                <VisuallyHidden>
                  <DialogTitle>Search for a game by name</DialogTitle>
                </VisuallyHidden>

                {/* Search input - fixed at top */}
                <div className="shrink-0 p-4 border-b">
                  <Input
                    placeholder="Search for a game by name"
                    value={searchQuery}
                    onChange={handleSearchQueryChange}
                    className="w-full"
                  />
                </div>

                {/* Scrollable results container */}
                <div className="flex-1 overflow-y-auto">
                  <div className="p-4 space-y-2">
                {/* Search results content */}
                  {
                    isLoading ? LoadingSkeleton :
                    error ? (
                      <div className="text-red-500 p-4">
                        We're having trouble with search. Please try again later.
                      </div>
                    ) : displaySearchResults.length === 0 ? (
                      <div className="text-muted-foreground p-4 text-center">
                        {debouncedSearchQuery ? 'No games found' : 'Start typing to search for games'}
                      </div>
                    ) : (
                      displaySearchResults.map((game, index) => (
                        <SearchResult
                          key={`${game.name}-${index}`}
                          title={game.name}
                          imageUrl={game.coverImage}
                          isInLibrary={game.isInLibrary}
                        />
                      ))
                    )
                  }
                  </div>
                </div>

              </div>
            </DialogContent>
          </div>
        </div>

      </DialogPortal>
    </Dialog>
  )
};
