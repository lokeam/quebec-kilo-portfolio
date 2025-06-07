import { useCallback, useEffect } from 'react';
import { Button } from '@/shared/components/ui/button';
import { LayoutGrid, LayoutList, Sheet } from 'lucide-react';
import { Input } from '@/shared/components/ui/input';

// Components
import { FilterDropdown } from '@/shared/components/ui/FilterDropdown/FilterDropdown';

// Hooks
import { useFilterCheckboxes } from '@/shared/components/ui/FilterDropdown/useFilterCheckboxes';
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';

// Types
import { PhysicalLocationType } from '@/types/domain/location-types';
import { SublocationType } from '@/types/domain/location-types';

interface FilterOption {
  key: string;
  label: string;
}

interface PhysicalLocationsToolbarProps {
  sublocationTypes: FilterOption[];
  parentTypes: FilterOption[];
}

export function PhysicalLocationsToolbar({ sublocationTypes, parentTypes }: PhysicalLocationsToolbarProps) {
  const {
    viewMode,
    setViewMode,
    setSearchQuery,
    setSublocationTypeFilters,
    setParentLocationTypeFilters,
  } = useOnlineServicesStore();

  const sublocationTypeFilter = useFilterCheckboxes(
    sublocationTypes.map(option => option.key)
  );

  const parentTypeFilter = useFilterCheckboxes(
    parentTypes.map(option => option.key)
  );

  const handleSearchChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setSearchQuery(event.target.value.toLowerCase());
    },
    [setSearchQuery]
  );

  useEffect(() => {
    const selectedSublocationTypes = Object.entries(sublocationTypeFilter.checkboxes)
      .filter(([, isChecked]) => isChecked === true)
      .map(([key]) => key);

    setSublocationTypeFilters(selectedSublocationTypes);
  }, [sublocationTypeFilter.checkboxes, setSublocationTypeFilters]);

  useEffect(() => {
    const selectedParentTypes = Object.entries(parentTypeFilter.checkboxes)
      .filter(([, isChecked]) => isChecked === true)
      .map(([key]) => key);

    setParentLocationTypeFilters(selectedParentTypes);
  }, [parentTypeFilter.checkboxes, setParentLocationTypeFilters]);

  return (
    <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
      <div className="flex flex-wrap items-center gap-3">
        <Input
          className="w-[300px]"
          id="filterServices"
          placeholder="Filter Services..."
          onChange={handleSearchChange}
        />
        <span className="text-sm text-gray-500">Filter by</span>

        <FilterDropdown
          label="Storage Type"
          options={sublocationTypes}
          width="230px"
          {...sublocationTypeFilter}
        />

        <FilterDropdown
          label="Property Type"
          options={parentTypes}
          width="180px"
          {...parentTypeFilter}
        />
      </div>

      <div className="flex items-center gap-3">
        <span className="text-sm text-gray-500">View</span>
        <div className="flex bg-black rounded-md p-1 gap-1">
          <Button
            variant={viewMode === 'grid' ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setViewMode('grid')}
          >
            <LayoutGrid className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === 'list' ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setViewMode('list')}
          >
            <LayoutList className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === 'table' ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setViewMode('table')}
          >
            <Sheet className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
