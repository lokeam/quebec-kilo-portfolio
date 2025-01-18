
import { useState, useMemo } from 'react';

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
import { useLibraryGames } from '@/features/dashboard/lib/stores/libraryStore';
import { usePlatformSearch } from '@/features/dashboard/lib/hooks/usePlatformSearch';

// Icons
import { Check, ChevronsUpDown } from 'lucide-react';


interface PlatformComboboxProps {
  onPlatformSelect?: (platform: string) => void;
}

export function PlatformCombobox({
  onPlatformSelect,
}: PlatformComboboxProps) {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState("");
  const { availablePlatforms, handleSearch } = usePlatformSearch();


  const handlePlatformSelect = (platform: string) => {
    setValue(platform);
    setOpen(false);
    onPlatformSelect?.(platform);
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
          {value || "Select platform..."}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0">

      <Command className="rounded-lg border border-input bg-background">
        <Command.Input
          placeholder="Search platform..."
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
                  heading={manufacturer.toUpperCase()}
                  className="text-sm font-medium text-muted-foreground px-2 py-1.5"
                >
                  {platforms.map((platform) => (
                    <Command.Item
                      key={platform.key}
                      value={[
                        platform.key,
                        platform.label.toLowerCase(),
                        ...platform.searchTerms
                      ].join(' ')}
                      onSelect={() => handlePlatformSelect(platform.key)}
                      className="relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none aria-selected:bg-accent aria-selected:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 hover:bg-accent hover:text-accent-foreground"
                    >
                      {platform.label}
                      <Check
                        className={`ml-auto h-4 w-4 ${
                          value === platform.label ? "opacity-100" : "opacity-0"
                        }`}
                      />
                    </Command.Item>
                  ))}
                </Command.Group>
              ))}
            </Command.List>
          </ScrollArea>
        </Command>

        {/* <Command>
          <Command.Input placeholder="Search platform..." />
          <ScrollArea className="h-[200px]">
            <Command.List>
              <Command.Empty>No platform found.</Command.Empty>
                {Object.entries(availablePlatforms).map(([manufacturer, platforms]) => (
                  <Command.Group
                    key={manufacturer}
                    heading={manufacturer.toUpperCase()}
                  >
                    {platforms.map((platform) => (
                      <Command.Item
                        key={platform.key}
                        value={platform.key}
                        onSelect={() => handlePlatformSelect(platform.key)}
                      >
                        {platform.label}
                        <Check
                          className={`ml-auto h-4 w-4 ${
                            value === platform.label ? "opacity-100" : "opacity-0"
                          }`}
                        />
                      </Command.Item>
                    ))}
                  </Command.Group>
                ))}
            </Command.List>
          </ScrollArea>

        </Command> */}
      </PopoverContent>
    </Popover>
  )
}
