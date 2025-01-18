import { CONSOLE_PLATFORMS } from '@/features/dashboard/lib/constants/filterOptions/library/platform.filterOptions';
import { FilterDropdown } from '@/shared/components/ui/FilterDropdown';

export function LibraryPageToolbar() {
  return <FilterDropdown options={CONSOLE_PLATFORMS} />;
}