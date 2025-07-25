import { useState } from 'react';

// Components
import { SearchResult } from '@/features/dashboard/components/organisms/GameSearchAndSelectDialog/SearchResult';
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';
import { AddGameToLibraryForm } from '@/features/dashboard/components/organisms/GameSearchAndSelectDialog/AddGameToLibraryForm/AddGameToLibraryForm';

// Types
import type { SearchResult as SearchResultType } from '@/types/domain/search';
import type { Game } from '@/types/domain/game';
import type { AddGameFormPhysicalLocationsResponse, AddGameFormDigitalLocationsResponse } from '@/types/domain/search';

interface ResultsSectionProps {
  results?: SearchResultType[];
  isLoading: boolean;
  error: Error | null;
  onSelect: (gameId: string) => void;
  physicalLocations?: AddGameFormPhysicalLocationsResponse[];
  digitalLocations?: AddGameFormDigitalLocationsResponse[];
}

export function ResultsSection({
  results,
  isLoading,
  error,
  onSelect,
  physicalLocations,
  digitalLocations,
}: ResultsSectionProps) {
  const [selectedGameWithDrawerOpen, setSelectedGameWithDrawerOpen] = useState<Game | null>(null);

  const handleAddGameToLibrary = (game: Game) => setSelectedGameWithDrawerOpen(game);
  const handleCloseDrawer = () => setSelectedGameWithDrawerOpen(null);

  // console.log('DEBUG: handleCloseDrawer called');
  if (isLoading) {
    return (
      <div className="space-y-4">
        {Array.from({ length: 3 }).map((_, index) => (
          <div key={index} className="h-24 animate-pulse bg-muted rounded-lg" />
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

  // console.log('ResultsSection, physicalLocations', physicalLocations);
  // console.log('ResultsSection, digitalLocations', digitalLocations);

  return (
    <div className="space-y-4">
      {results.map((result) => (
        <SearchResult
          key={result.game.id}
          game={result.game}
          onAction={() => onSelect(String(result.game.id))}
          onAddToLibrary={handleAddGameToLibrary}
        />
      ))}

      <DrawerContainer
        open={!!selectedGameWithDrawerOpen}
        onOpenChange={open => {
          // console.log('Drawer onOpenChanged called with: ', open);
          // console.log('🔍 DEBUG: Drawer onOpenChange:', {
          //   open,
          //   selectedGameWithDrawerOpen,
          //   hasOnClose: !!handleCloseDrawer
          // });
          if (!open) handleCloseDrawer();
        }}
        title=""
        description=""
      >
        {
          selectedGameWithDrawerOpen && (
            <AddGameToLibraryForm
              selectedGame={selectedGameWithDrawerOpen}
              physicalLocations={physicalLocations ?? []}
              digitalLocations={digitalLocations ?? []}
              onClose={handleCloseDrawer}
            />
          )
        }
      </DrawerContainer>
    </div>
  );
}