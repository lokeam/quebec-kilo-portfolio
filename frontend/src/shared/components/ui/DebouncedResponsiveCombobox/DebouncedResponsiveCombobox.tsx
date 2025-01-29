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
import { useDebounce } from "@/shared/hooks/useDebounce"
import { useMediaQuery } from "@/shared/hooks/useMediaQuery"

export interface SelectableItem {
  id: string;
  displayName: string;
  [key: string]: any;
}

interface DebouncedResponsiveComboboxProps<T extends SelectableItem> {
  onSelect: (item: T) => void;
  items?: T[];
  fetchItems?: (query: string) => Promise<T[]>;
  placeholder?: string;
  emptyMessage?: string;
  className?: string;
  label?: string;
}

function ComboboxContent<T extends SelectableItem>({
  isLoading,
  displayedItems,
  emptyMessage,
  selectedItem,
  handleSelect,
  searchQuery,
  setSearchQuery,
  placeholder,
  label,
}: {
  isLoading: boolean;
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
    <Command className="overflow-hidden" aria-label={label}>
      <CommandInput
        placeholder={placeholder}
        value={searchQuery}
        onValueChange={setSearchQuery}
      />
      <div
        className="command-scroll-area max-h-[300px] overflow-y-auto"
        role="listbox"
        aria-label={`${label} options`}
      >
        {isLoading ? (
          <CommandEmpty className="text-sm text-gray-500 text-center p-2">
            Loading...
          </CommandEmpty>
        ) : displayedItems.length === 0 ? (
          <CommandEmpty>{emptyMessage}</CommandEmpty>
        ) : (
          <CommandGroup>
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

export function DebouncedResponsiveCombobox<T extends SelectableItem>({
  onSelect,
  items,
  fetchItems,
  placeholder = "Search...",
  emptyMessage = "No items found.",
  className,
  label = "Select an option",
}: DebouncedResponsiveComboboxProps<T>) {
  const [open, setOpen] = React.useState(false)
  const isDesktop = useMediaQuery("(min-width: 768px)")
  const [searchQuery, setSearchQuery] = React.useState("")
  const [selectedItem, setSelectedItem] = React.useState<T | null>(null)
  const [displayedItems, setDisplayedItems] = React.useState<T[]>([])
  const [isLoading, setIsLoading] = React.useState(false)

  const debouncedSearchQuery = useDebounce(searchQuery, 150)

  // Handle static items (constants)
  React.useEffect(() => {
    if (items) {
      const filtered = items.filter(item =>
        item.displayName.toLowerCase().includes(debouncedSearchQuery.toLowerCase())
      )
      setDisplayedItems(filtered)
    }
  }, [items, debouncedSearchQuery])

  // Handle API-based items
  React.useEffect(() => {
    async function fetchData() {
      if (!fetchItems) return

      setIsLoading(true)
      try {
        const results = await fetchItems(debouncedSearchQuery)
        setDisplayedItems(results)
      } catch (error) {
        console.error('Error fetching items:', error)
        setDisplayedItems([])
      } finally {
        setIsLoading(false)
      }
    }

    if (fetchItems) {
      fetchData()
    }
  }, [fetchItems, debouncedSearchQuery])

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
      isLoading={isLoading}
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
          className="w-full p-0"
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
        <div className="mt-4 border-t p-4">
          {commonContent}
        </div>
      </DrawerContent>
    </Drawer>
  )
}

DebouncedResponsiveCombobox.displayName = 'DebouncedResponsiveCombobox'
