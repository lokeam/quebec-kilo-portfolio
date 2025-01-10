import { useState, useEffect, useCallback, useMemo } from 'react'
import { Button } from "@/shared/components/ui/button"
import { DialogTrigger } from '@/shared/components/ui/dialog'
import { SearchIcon } from 'lucide-react'
import { SearchDialog, SearchDialogSkeleton } from '@/shared/components/ui/SearchDialog'
import { SearchResult } from '@/features/dashboard/organisms/AddItemSearchDialog/SearchResult'
import { useSearchGames } from '@/features/dashboard/hooks/useSearchGames'
import { useDebounce } from '@/shared/hooks/useDebounce'
import type { Game } from '@/types/types/domain.types'

export function AddItemSearchDialog() {
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [searchResults, setSearchResults] = useState<Game[]>([]);
  const debouncedSearchQuery = useDebounce(searchQuery, 400);
  const { games, isLoading, error } = useSearchGames(debouncedSearchQuery);

  const handleOpenChange = useCallback((open: boolean) => {
    setIsOpen(open)
    if (!open) {
      setSearchQuery('')
    }
  }, []);

  const handleSearchQueryChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(event.target.value)
  }, []);

  useEffect(() => {
    let isMounted = true;

    if (games?.length > 0 && isMounted) {
      setSearchResults(games)
    }

    return () => {
      isMounted = false;
    }
  }, [games]);

  const displaySearchResults = useMemo(() => {
    // Early return for undefined/empty cases
    if (!games && !searchResults.length) return [];

    // Prefer games over searchResults when available
    return games?.length ? games : searchResults;
  }, [games, searchResults]);

  return (
    <SearchDialog
      open={isOpen}
      onOpenChange={handleOpenChange}
      searchQuery={searchQuery}
      onSearchChange={handleSearchQueryChange}
      searchPlaceholder="Search for a game by name"
      dialogTitle="Search Games"
      hideHeader={true}
      trigger={
        <DialogTrigger asChild>
          <Button variant="outline" size="icon">
            <SearchIcon />
          </Button>
        </DialogTrigger>
      }
    >
      {isLoading ? (
        <SearchDialogSkeleton />
      ) : error ? (
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
      )}
    </SearchDialog>
  );
}
