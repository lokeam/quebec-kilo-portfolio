import * as React from "react"
import { Check, ChevronsUpDown } from "lucide-react"
import { cn } from "@/shared/components/ui/utils"
import { Button } from "@/shared/components/ui/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from "@/shared/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/shared/components/ui/popover"

export interface SelectableItem {
  id: string;
  displayName: string;
  [key: string]: string | number | boolean | undefined;
}

interface SimpleComboboxProps<T extends SelectableItem> {
  onSelect: (item: T) => void;
  items: T[];
  placeholder?: string;
  emptyMessage?: string;
  className?: string;
  label?: string;
  initialValue?: string;
}

export function SimpleCombobox<T extends SelectableItem>({
  onSelect,
  items,
  placeholder = "Search...",
  emptyMessage = "No items found.",
  className,
  label = "Select an option",
  initialValue
}: SimpleComboboxProps<T>) {
  const [open, setOpen] = React.useState(false)
  const [searchQuery, setSearchQuery] = React.useState("")

  // Find the selected item based on initialValue
  const selectedItem = React.useMemo(() => {
    if (!initialValue) return null;
    return items.find(item => item.id.toLowerCase() === initialValue.toLowerCase()) || null;
  }, [items, initialValue]);

  // Filter items based on search query
  const displayedItems = React.useMemo(() => {
    if (!searchQuery) return items;
    return items.filter(item =>
      item.displayName.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [items, searchQuery]);

  const handleSelect = React.useCallback((item: T) => {
    onSelect(item);
    setOpen(false);
    setSearchQuery("");
  }, [onSelect]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className={cn("w-full justify-between", className)}
          aria-label={label}
        >
          {selectedItem ? selectedItem.displayName : placeholder}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-full p-0" align="start">
        <Command>
          <CommandInput
            placeholder={placeholder}
            value={searchQuery}
            onValueChange={setSearchQuery}
          />
          <div
            className="max-h-[300px] overflow-y-auto"
            onWheel={(e) => {
              e.stopPropagation();
            }}
          >
            {displayedItems.length === 0 ? (
              <CommandEmpty>{emptyMessage}</CommandEmpty>
            ) : (
              <CommandGroup>
                {displayedItems.map((item) => (
                  <CommandItem
                    key={item.id}
                    value={item.id}
                    onSelect={() => handleSelect(item)}
                  >
                    <Check
                      className={cn(
                        "mr-2 h-4 w-4",
                        selectedItem?.id === item.id ? "opacity-100" : "opacity-0"
                      )}
                    />
                    {item.displayName}
                  </CommandItem>
                ))}
              </CommandGroup>
            )}
          </div>
        </Command>
      </PopoverContent>
    </Popover>
  )
}