import { useState } from 'react';
import type { DropdownMenuCheckboxItemProps } from '@radix-ui/react-dropdown-menu';

type Checked = DropdownMenuCheckboxItemProps['checked'];

/**
 * Hook for managing filter dropdown checkbox state.
 * Used to:
 * - Ensure Shadcn UI Dropdown Menu items correctly display checkbox state
 * - Force Shadcn UI Dropdown Menu to remain open after checkbox click
 * - Provide a way to clear all checkboxes
 *
 * @param keys - Array of string identifiers for filter options
 * @returns Object containing:
 * - checkboxes: Record of checkbox states
 * - isOpen: Dropdown open state
 * - onOpenChange: Handler for dropdown open/close
 * - onCheckboxChange: Handler for checkbox changes
 * - onClearAll: Handler to uncheck all boxes
 * - areAllUnchecked: Boolean indicating if all checkboxes are unchecked
 *
 * @example
 * ```tsx
 * const filterControls = useFilterCheckboxes(['option1', 'option2']);
 *
 * <FilterDropdown {...filterControls} />
 * ```
 */
export function useFilterCheckboxes(keys: string[]) {
  const [isOpen, setIsOpen] = useState(false);
  const [checkboxes, setCheckboxes] = useState<Record<string, boolean>>(
    keys.reduce((acc, key) => ({ ...acc, [key]: false }), {})
  );

  const onOpenChange = (open: boolean) => {
    setIsOpen(open);
  };

  const onCheckboxChange = (key: string) => (checked: Checked) => {
    setCheckboxes(prev => ({
      ...prev,
      [key]: checked === true,
    }));
  };

  const onClearAll = () => {
    setCheckboxes(
      Object.keys(checkboxes).reduce((acc, key) => ({
        ...acc,
        [key]: false,
      }), {})
    );
  };

  const areAllUnchecked = Object.values(checkboxes).every(value => !value);

  return {
    checkboxes,
    isOpen,
    onOpenChange,
    onCheckboxChange,
    onClearAll,
    areAllUnchecked,
  } as const;
}
