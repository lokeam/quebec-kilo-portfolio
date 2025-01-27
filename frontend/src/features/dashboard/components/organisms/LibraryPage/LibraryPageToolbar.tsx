import { useMemo } from 'react';
import { CONSOLE_PLATFORMS } from '@/features/dashboard/lib/constants/filter-options/library/platform.filterOptions';
import { FilterDropdown } from '@/shared/components/ui/FilterDropdown/FilterDropdown';

export function LibraryPageToolbar() {
  const options = useMemo(() => Object.values(CONSOLE_PLATFORMS).flat(), []);
  return (<FilterDropdown
    options={options}
    label="Platforms"
    checkboxes={{}}
    isOpen={false}
    onOpenChange={() => {}}
    onCheckboxChange={(): ((checked: boolean) => void) => () => {}}
    onClearAll={() => {}}

    areAllUnchecked={false}
  />);
}

