import { useCallback, useEffect, useMemo } from 'react';
import { Button } from '@/shared/components/ui/button';
import { LayoutGrid, LayoutList, Sheet } from 'lucide-react';
import { Input } from '@/shared/components/ui/input';

// Components
import { FilterDropdown } from '@/shared/components/ui/FilterDropdown/FilterDropdown';

// Hooks
import { useFilterCheckboxes } from '@/shared/components/ui/FilterDropdown/useFilterCheckboxes';
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { useStorageAnalytics } from '@/core/api/queries/analyticsData.queries';

// Types
import { PhysicalLocationType } from '@/types/domain/location-types';
import { SublocationType } from '@/types/domain/location-types';

export function PhysicalLocationsToolbar() {
  const {
    viewMode,
    setViewMode,
    setSearchQuery,
    setSublocationTypeFilters,
    setParentLocationTypeFilters,
  } = useOnlineServicesStore();

  const { data: storageData } = useStorageAnalytics();

  // Generate filter options from the data
  const filterOptions = useMemo(() => {
    if (!storageData?.sublocationRows) return { sublocationTypes: [], parentTypes: [] };

    // Get unique sublocation types and format them for display
    const uniqueSublocationTypes = Array.from(new Set(
      storageData.sublocationRows.map(row => row.sublocationType)
    ))
    .filter((type): type is SublocationType =>
      Object.values(SublocationType).includes(type as SublocationType)
    )
    .map(type => ({
      key: type,
      label: type.charAt(0).toUpperCase() + type.slice(1) // Capitalize first letter
    }));

    // Get unique physical location types and format them for display
    const uniqueParentTypes = Array.from(new Set(
      storageData.sublocationRows.map(row => row.parentLocationType)
    ))
    .filter((type): type is PhysicalLocationType =>
      Object.values(PhysicalLocationType).includes(type as PhysicalLocationType)
    )
    .map(type => ({
      key: type,
      label: type.charAt(0).toUpperCase() + type.slice(1) // Capitalize first letter
    }));

    return {
      sublocationTypes: uniqueSublocationTypes,
      parentTypes: uniqueParentTypes
    };
  }, [storageData?.sublocationRows]);

  const sublocationTypeFilter = useFilterCheckboxes(
    filterOptions.sublocationTypes.map(option => option.key)
  );

  const parentTypeFilter = useFilterCheckboxes(
    filterOptions.parentTypes.map(option => option.key)
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
          options={filterOptions.sublocationTypes}
          width="230px"
          {...sublocationTypeFilter}
        />

        <FilterDropdown
          label="Property Type"
          options={filterOptions.parentTypes}
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
