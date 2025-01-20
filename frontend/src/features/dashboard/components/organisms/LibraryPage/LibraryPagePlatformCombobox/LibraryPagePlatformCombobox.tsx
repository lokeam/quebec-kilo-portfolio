
import { useState } from 'react';

// Shadcn UI Componets
import { Button } from '@/shared/components/ui/button';
import { Command } from 'cmdk';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/shared/components/ui/popover';
import { ScrollArea } from '@/shared/components/ui/scroll-area';

// Hooks, Utils + Stores
import { useLibrarySetPlatformFilter } from '@/features/dashboard/lib/stores/libraryStore';
import { usePlatformSearch } from '@/features/dashboard/lib/hooks/usePlatformSearch';

// Icons
import { Check, ChevronsUpDown, X } from 'lucide-react';


interface PlatformComboboxProps {
  onPlatformSelect?: (platform: string) => void;
}

export function PlatformCombobox({
  onPlatformSelect,
}: PlatformComboboxProps) {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState('');
  const setPlatformFilter = useLibrarySetPlatformFilter();
  const { availablePlatforms, handleSearch } = usePlatformSearch();

  const handlePlatformSelect = (platform: string) => {
    // Only update if the value has changed
    if (platform !== value) {
      setValue(platform);
      setOpen(false);
      onPlatformSelect?.(platform);
      setPlatformFilter(platform);
      handleSearch(''); // Reset search when selection is made
    }
  }

  const handleClear = (event: React.MouseEvent<HTMLDivElement>) => {
    event.stopPropagation();
    setValue('');
    setPlatformFilter('');
    onPlatformSelect?.('');
    setOpen(false);
    handleSearch(''); // Reset search when cleared
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-[200px] justify-between"
        >
          {value || "All platforms..."}
          {value ? (
            <div onClick={handleClear}>
              <X className="ml-2 h-4 w-4 shrink-0 opacity-50 hover:opacity-100" />
            </div>
          ) : (
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0">
        <Command className="rounded-lg border border-input bg-background">
          <Command.Input
            autoFocus
            placeholder="Search platforms..."
            onValueChange={handleSearch}
            value={value}
            className="h-9 px-3 py-2 text-sm bg-transparent focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
          />
          <ScrollArea className="h-[200px]">
            <Command.List className="px-1">
              <Command.Empty className="py-6 text-center text-sm">
                No platform found.
              </Command.Empty>
              {Object.entries(availablePlatforms).map(([manufacturer, platforms]) => (
                <Command.Group
                  key={manufacturer}
                  heading={`Developer: ${manufacturer.toUpperCase()}`}
                  className="text-sm text-gray-500"
                >
                  {platforms.map((platform) => (
                    <Command.Item
                      key={platform.key}
                      value={platform.key}
                      onSelect={handlePlatformSelect}
                      className="relative flex cursor-default gap-2 select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none data-[disabled=true]:pointer-events-none data-[selected='true']:bg-accent data-[selected=true]:text-accent-foreground data-[disabled=true]:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0"
                    >
                      <div className="flex items-center justify-between w-full text-md text-white">
                        {platform.label}
                        <Check
                          className={`ml-auto h-4 w-4 ${
                            value === platform.key ? "opacity-100" : "opacity-0"
                          }`}
                        />
                      </div>
                    </Command.Item>
                  ))}
                </Command.Group>
              ))}
            </Command.List>
          </ScrollArea>
        </Command>
      </PopoverContent>
    </Popover>
  );
}