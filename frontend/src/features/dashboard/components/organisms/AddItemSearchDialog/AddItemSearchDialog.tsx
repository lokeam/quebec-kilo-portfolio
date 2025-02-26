import { useState, useCallback } from 'react';

// Shadcn Components
import { Button } from '@/shared/components/ui/button';
import { DialogTrigger } from '@/shared/components/ui/dialog';

// Components
import { SearchDialog, SearchDialogSkeleton } from '@/shared/components/ui/SearchDialog';
import { SearchResult } from '@/features/dashboard/components/organisms/AddItemSearchDialog/SearchResult';

// Hooks
import { useDebounce } from '@/shared/hooks/useDebounce';
import { useMediaItemSearch } from '@/core/api/queries/useMediaItemSearch';

// Icons
import { SearchIcon } from 'lucide-react';

// Types
//import type { Game } from '@/types/types/domain.types';

export function AddItemSearchDialog() {
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [isOpen, setIsOpen] = useState<boolean>(false);

  const debouncedSearchQuery = useDebounce(searchQuery, 400);
  const { data: games, isLoading, error } = useMediaItemSearch(debouncedSearchQuery);

  const handleOpenChange = useCallback((open: boolean) => {
    setIsOpen(open)
    if (!open) setSearchQuery('')
  }, []);

  const handleSearchQueryChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(event.target.value)
  }, []);

  const handleAction = () => {
    setIsOpen(false);
  }

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
      ) : games?.length === 0 ? (
        <div className="text-muted-foreground p-4 text-center">
          {debouncedSearchQuery ? 'No games found' : 'Start typing to search for games'}
        </div>
      ) : (
        games?.map((game, index) => (
          <SearchResult
            key={`${game?.name}-${index}`}
            // name={game?.name ?? ''}
            game={game}
            onAction={handleAction}
            // cover_url={game?.cover_url ?? ''}
            // is_in_library={game?.is_in_library ?? false}
            // is_in_wishlist={game?.is_in_wishlist ?? false}
          />
        ))
      )}
    </SearchDialog>
  );
}
