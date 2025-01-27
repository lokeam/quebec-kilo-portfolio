import { useCallback } from 'react';

// ShadCN UI Components
import { Button } from '@/shared/components/ui/button';
import { LayoutGrid, LayoutList } from 'lucide-react';
import { Input } from '@/shared/components/ui/input';

// Components
import { PlatformCombobox } from '@/features/dashboard/components/organisms/LibraryPage/LibraryPagePlatformCombobox/LibraryPagePlatformCombobox';

// Hooks
import {
  useLibrarySetSearchQuery,
  useLibraryViewMode,
  useLibrarySetViewMode,
} from '@/features/dashboard/lib/stores/libraryStore';

// Constants
import { featureViewModes } from '@/shared/constants/viewModes';

export function LibraryPageToolbar() {
  const currentViewMode = useLibraryViewMode();
  const setCurrentViewMode = useLibrarySetViewMode();
  const setSearchQuery = useLibrarySetSearchQuery();

  const handleSearchChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setSearchQuery(event.target.value.toLowerCase());
  }, [setSearchQuery]);

  return (
    <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
      <div className="flex flex-wrap items-center gap-3">
        <Input
          className="w-[300px]"
          id="filterServices"
          placeholder="Filter games in your library..."
          onChange={handleSearchChange}
        />
        <span className="text-sm text-gray-500">Filter by</span>

        <PlatformCombobox />

      </div>

      <div className="flex items-center gap-3">
        <span className="text-sm text-gray-500">View</span>
        <div className="flex bg-black rounded-md p-1 gap-1">
        <Button
            variant={currentViewMode === featureViewModes.library.allowed[0] ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setCurrentViewMode(featureViewModes.library.allowed[0])}
          >
            <LayoutGrid className="h-4 w-4" />
          </Button>
          <Button
            variant={currentViewMode === featureViewModes.library.allowed[1] ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setCurrentViewMode(featureViewModes.library.allowed[1])}
          >
            <LayoutList className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
