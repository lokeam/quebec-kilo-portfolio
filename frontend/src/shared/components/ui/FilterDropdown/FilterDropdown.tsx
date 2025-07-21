import { Button } from '@/shared/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu';
import { ChevronDown } from '@/shared/components/ui/icons';

interface FilterOption {
  key: string;
  label: string;
};

interface FilterDropdownProps {
  label: string;
  options: readonly FilterOption[];
  width?: string;
  checkboxes: Record<string, boolean>;
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  onCheckboxChange: (key: string) => (checked: boolean) => void;
  onClearAll: () => void;
  areAllUnchecked: boolean;
};

export function FilterDropdown({
  label,
  options,
  width = '140px',
  checkboxes,
  isOpen,
  onOpenChange,
  onCheckboxChange,
  onClearAll,
  areAllUnchecked,
}: FilterDropdownProps) {
  return (
    <DropdownMenu open={isOpen} onOpenChange={onOpenChange}>
      <DropdownMenuTrigger
        asChild
        aria-label={`Filter online services by ${label}`}
        aria-expanded={isOpen}
      >
        <Button variant="outline" className={`w-[${width}] justify-between`}>
          {label}
          <ChevronDown className="h-4 w-4 opacity-50" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        className={`w-[${width}]`}
        onInteractOutside={() => onOpenChange(false)}
        onEscapeKeyDown={(e) => e.preventDefault()}
      >
        <DropdownMenuGroup>
          {options.map(({ key, label }) => (
            <DropdownMenuCheckboxItem
              key={key}
              checked={checkboxes[key]}
              onCheckedChange={onCheckboxChange(key)}
              onSelect={(e) => e.preventDefault()}
            >
              {label}
            </DropdownMenuCheckboxItem>
          ))}
          <DropdownMenuSeparator />
          <DropdownMenuCheckboxItem
            checked={areAllUnchecked}
            onCheckedChange={onClearAll}
            onSelect={(e) => e.preventDefault()}
          >
            Clear All
          </DropdownMenuCheckboxItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

