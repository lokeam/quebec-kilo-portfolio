import { useState, useCallback } from 'react';

// Shadcn Components
import { Button } from '@/shared/components/ui/button';
import { DialogTrigger } from '@/shared/components/ui/dialog';

// Components
import { SearchDialog, SearchDialogSkeleton } from '@/shared/components/ui/SearchDialog';
import { ResultsSection } from './ResultsSection/ResultsSection';
import { ActionsSection } from './ActionsSection/ActionsSection';

// Hooks
import { useDebounce } from '@/shared/hooks/useDebounce';
import { useGameSearch } from '@/core/api/queries/gameSearch.queries';

// Icons
import { SearchIcon } from 'lucide-react';

export function GameSearchAndSelectDialog() {
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [selectedGames, setSelectedGames] = useState<Set<string>>(new Set());

  const debouncedSearchQuery = useDebounce(searchQuery, 400);
  const { data, isLoading, error } = useGameSearch({
    query: debouncedSearchQuery,
    filters: {},
    sortBy: 'rating',
    sortOrder: 'desc'
  });

  const handleOpenChange = useCallback((open: boolean) => {
    setIsOpen(open);
    if (!open) {
      setSearchQuery('');
      setSelectedGames(new Set());
    }
  }, []);

  const handleSearchQueryChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(event.target.value);
  }, []);

  const handleGameSelect = useCallback((gameId: string) => {
    setSelectedGames(prev => {
      const next = new Set(prev);
      if (next.has(gameId)) {
        next.delete(gameId);
      } else {
        next.add(gameId);
      }
      return next;
    });
  }, []);

  const handleConfirm = useCallback(() => {
    // TODO: Implement confirmation logic
    setIsOpen(false);
  }, []);

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
      <div className="flex flex-col gap-4">
        {isLoading ? (
          <SearchDialogSkeleton />
        ) : error ? (
          <div className="text-red-500 p-4">
            We're having trouble with search. Please try again later.
          </div>
        ) : data?.results.length === 0 ? (
          <div className="text-muted-foreground p-4 text-center">
            {debouncedSearchQuery ? 'No games found' : 'Start typing to search for games'}
          </div>
        ) : (
          <>
            <ResultsSection
              results={data?.results}
              isLoading={isLoading}
              error={error instanceof Error ? error : null}
              onSelect={handleGameSelect}
            />
            <ActionsSection
              selectedGames={selectedGames}
              onConfirm={handleConfirm}
            />
          </>
        )}
      </div>
    </SearchDialog>
  );
}