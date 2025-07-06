import { useCallback, useEffect } from 'react';
import { Button } from '@/shared/components/ui/button';
import { LayoutGrid, LayoutList } from 'lucide-react';
import { Input } from '@/shared/components/ui/input';

// Components
import { FilterDropdown } from '@/shared/components/ui/FilterDropdown/FilterDropdown';

// Hooks
import { useFilterCheckboxes } from '@/shared/components/ui/FilterDropdown/useFilterCheckboxes';
import { useLibraryStore } from '@/features/dashboard/lib/stores/libraryStore';

// Constants
import { featureViewModes } from '@/shared/constants/viewModes';

interface FilterOption {
  key: string;
  label: string;
}

interface LibraryPageToolbarProps {
  platforms: FilterOption[];
  locations: FilterOption[];
}

export function LibraryPageToolbar({ platforms, locations }: LibraryPageToolbarProps) {
  const {
    viewMode,
    setViewMode,
    setSearchQuery,
    setPlatformFilters,
    setLocationFilters,
  } = useLibraryStore();

  const platformFilter = useFilterCheckboxes(
    platforms.map(option => option.key)
  );

  const locationFilter = useFilterCheckboxes(
    locations.map(option => option.key)
  );

  const handleSearchChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setSearchQuery(event.target.value.toLowerCase());
    },
    [setSearchQuery]
  );

  useEffect(() => {
    const selectedPlatforms = Object.entries(platformFilter.checkboxes)
      .filter(([, isChecked]) => isChecked === true)
      .map(([key]) => key);

    setPlatformFilters(selectedPlatforms);
  }, [platformFilter.checkboxes, setPlatformFilters]);

  useEffect(() => {
    const selectedLocations = Object.entries(locationFilter.checkboxes)
      .filter(([, isChecked]) => isChecked === true)
      .map(([key]) => key);

    setLocationFilters(selectedLocations);
  }, [locationFilter.checkboxes, setLocationFilters]);

  return (
    <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
      <div className="flex flex-wrap items-center gap-3">
        <Input
          className="w-[300px]"
          id="filterGames"
          placeholder="Filter games in your library..."
          onChange={handleSearchChange}
        />
        <span className="text-sm text-muted-foreground">Filter by</span>

        <FilterDropdown
          label="Platform"
          options={platforms}
          width="200px"
          {...platformFilter}
        />

        <FilterDropdown
          label="Location"
          options={locations}
          width="200px"
          {...locationFilter}
        />
      </div>

      <div className="flex items-center gap-3">
        <span className="text-sm text-muted-foreground">View</span>
        <div className="flex bg-muted rounded-md p-1 gap-1">
          <Button
            variant={viewMode === featureViewModes.library.allowed[0] ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setViewMode(featureViewModes.library.allowed[0])}
          >
            <LayoutGrid className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === featureViewModes.library.allowed[1] ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setViewMode(featureViewModes.library.allowed[1])}
          >
            <LayoutList className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
