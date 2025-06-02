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
  Drawer,
  DrawerContent,
  DrawerTrigger,
} from "@/shared/components/ui/drawer"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/shared/components/ui/popover"
import { useMediaQuery } from "@/shared/hooks/useMediaQuery"

export interface SelectableItem {
  id: string;
  displayName: string;
  [key: string]: string | number | boolean | undefined;
}

interface ResponsiveComboboxProps<T extends SelectableItem> {
  onSelect: (item: T) => void;
  items: T[];
  placeholder?: string;
  emptyMessage?: string;
  className?: string;
  label?: string;
  initialValue?: string;
}

function ComboboxContent<T extends SelectableItem>({
  displayedItems,
  emptyMessage,
  selectedItem,
  handleSelect,
  searchQuery,
  setSearchQuery,
  placeholder,
  label,
}: {
  displayedItems: T[];
  emptyMessage: string;
  selectedItem: T | null;
  handleSelect: (item: T) => void;
  searchQuery: string;
  setSearchQuery: (value: string) => void;
  placeholder: string;
  label: string;
}) {
  return (
    <Command className="overflow-visible" aria-label={label}>
      <CommandInput
        placeholder={placeholder}
        value={searchQuery}
        onValueChange={setSearchQuery}
      />
      <div
        className="max-h-[300px] overflow-y-auto"
        role="listbox"
        aria-label={`${label} options`}
        onWheel={(e) => e.stopPropagation()}
      >
        {displayedItems.length === 0 ? (
          <CommandEmpty>{emptyMessage}</CommandEmpty>
        ) : (
          <CommandGroup className="overflow-visible">
            {displayedItems.map((item) => (
              <CommandItem
                aria-selected={selectedItem?.id === item.id}
                key={item.id}
                value={item.id}
                role="option"
                onSelect={() => handleSelect(item)}
              >
                <Check
                  aria-hidden="true"
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
  )
}

export function ResponsiveCombobox<T extends SelectableItem>({
  onSelect,
  items,
  placeholder = "Search...",
  emptyMessage = "No items found.",
  className,
  label = "Select an option",
  initialValue
}: ResponsiveComboboxProps<T>) {
  const [open, setOpen] = React.useState(false)
  const isDesktop = useMediaQuery("(min-width: 768px)")
  const [searchQuery, setSearchQuery] = React.useState("")
  const [selectedItem, setSelectedItem] = React.useState<T | null>(null)

  // Find initial item if provided
  const initialItem = React.useMemo(() => {
    if (!initialValue) return null;
    return items.find(item =>
      item.id.toLowerCase() === initialValue.toLowerCase() ||
      item.displayName.toLowerCase() === initialValue.toLowerCase()
    ) || null;
  }, [initialValue, items]);

  // Set initial item when it's found
  React.useEffect(() => {
    if (initialItem && !selectedItem) {
      setSelectedItem(initialItem);
    }
  }, [initialItem, selectedItem]);

  const displayedItems = React.useMemo(() => {
    if (!searchQuery) return items;
    return items.filter(item =>
      item.displayName.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [items, searchQuery]);

  const handleSelect = React.useCallback((item: T) => {
    setSelectedItem(item)
    setOpen(false)
    onSelect(item)
  }, [onSelect])

  const commonTriggerButton = (
    <Button
      variant="outline"
      role="combobox"
      aria-expanded={open}
      aria-haspopup="listbox"
      aria-controls={open ? "combobox-content" : undefined}
      aria-label={label}
      className={cn("w-full justify-between", className)}
    >
      {selectedItem ? selectedItem.displayName : placeholder}
      <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
    </Button>
  )

  const commonContent = (
    <ComboboxContent
      displayedItems={displayedItems}
      emptyMessage={emptyMessage}
      selectedItem={selectedItem}
      handleSelect={handleSelect}
      searchQuery={searchQuery}
      setSearchQuery={setSearchQuery}
      placeholder={placeholder}
      label={label}
    />
  )

  if (isDesktop) {
    return (
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          {commonTriggerButton}
        </PopoverTrigger>
        <PopoverContent
          className="w-[300px] p-0"
          align="start"
          side="bottom"
          sideOffset={5}
          sticky="always"
          id="combobox-content"
        >
          {commonContent}
        </PopoverContent>
      </Popover>
    )
  }

  return (
    <Drawer open={open} onOpenChange={setOpen}>
      <DrawerTrigger asChild>
        {commonTriggerButton}
      </DrawerTrigger>
      <DrawerContent>
        <div className="mt-4 border-t p-4 w-[300px]">
          {commonContent}
        </div>
      </DrawerContent>
    </Drawer>
  )
}

ResponsiveCombobox.displayName = 'ResponsiveCombobox'
