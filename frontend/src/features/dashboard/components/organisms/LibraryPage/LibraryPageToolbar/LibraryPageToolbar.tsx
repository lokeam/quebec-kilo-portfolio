import { Button } from '@/shared/components/ui/button';
import { LayoutGrid, LayoutList } from 'lucide-react';
import { Input } from '@/shared/components/ui/input';

// Components
import { PlatformCombobox } from '@/features/dashboard/components/organisms/LibraryPage/LibraryPagePlatformCombobox/LibraryPagePlatformCombobox';

// Hooks
import { useLibraryViewMode, useLibrarySetViewMode, ViewModes } from '@/features/dashboard/lib/stores/libraryStore';

/*
  Filters:

  - Platform
  - Online Service
  - Alphabetical by Title
  - Size
  - Date Added

*/


export function LibraryPageToolbar() {
  const currentViewMode = useLibraryViewMode();
  const setCurrentViewMode = useLibrarySetViewMode();

  return (
    <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
      <div className="flex flex-wrap items-center gap-3">
        <Input
          className="w-[300px]"
          id="filterServices"
          placeholder="Filter games in your library..."
          onChange={() => {console.log('input filter change')}}
        />
        <span className="text-sm text-gray-500">Filter by</span>

        <PlatformCombobox />

      </div>

      <div className="flex items-center gap-3">
        <span className="text-sm text-gray-500">View</span>
        <div className="flex bg-black rounded-md p-1 gap-1">
          <Button
            variant={currentViewMode === ViewModes.GRID ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setCurrentViewMode(ViewModes.GRID)}
          >
            <LayoutGrid className="h-4 w-4" />
          </Button>
          <Button
            variant={currentViewMode === ViewModes.LIST ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setCurrentViewMode(ViewModes.LIST)}
          >
            <LayoutList className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
