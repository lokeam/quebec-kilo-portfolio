import { useMemo, useCallback, useState } from "react"
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
import { useDigitalServicesCatalog } from "@/core/api/queries/useDigitalServicesCatalog"
import type { DigitalServiceItem } from "@/core/api/services/digitalServices.service"

// Fallback data in case the API fails
const FALLBACK_SERVICES: DigitalServiceItem[] = [
  { id: 'amazonluna', name: 'Amazon Luna', logo: 'amazon' },
  { id: 'applearcade', name: 'Apple Arcade', logo: 'apple' },
  { id: 'blizzard', name: 'Blizzard Battle.net', logo: 'blizzard' },
  { id: 'ea', name: 'EA Play', logo: 'ea' },
  { id: 'epicgames', name: 'Epic Games', logo: 'epicgames' },
  { id: 'fanatical', name: 'Fanatical', logo: 'fanatical' },
  { id: 'gog', name: 'GOG', logo: 'gog' },
  { id: 'googleplaypass', name: 'Google Play Pass', logo: 'google' },
  { id: 'greenmangaming', name: 'Green Man Gaming', logo: 'greenmangaming' },
  { id: 'humblebundle', name: 'Humble Bundle', logo: 'humblebundle' },
  { id: 'itchio', name: 'itch.io', logo: 'itchio' },
  { id: 'meta', name: 'Meta', logo: 'meta' },
  { id: 'netflix', name: 'Netflix', logo: 'netflix' },
  { id: 'nintendo', name: 'Nintendo', logo: 'nintendo' },
  { id: 'nvidia', name: 'NVIDIA', logo: 'nvidia' },
  { id: 'primegaming', name: 'Prime Gaming', logo: 'prime' },
  { id: 'playstation', name: 'PlayStation Network', logo: 'ps' },
  { id: 'shadow', name: 'Shadow', logo: 'shadow' },
  { id: 'ubisoft', name: 'Ubisoft', logo: 'ubisoft' },
  { id: 'xboxlive', name: 'Xbox Live', logo: 'xbox' },
  { id: 'xboxgamepass', name: 'Xbox Game Pass', logo: 'xbox' }
];

interface ServiceComboboxProps {
  onServiceSelect: (service: DigitalServiceItem) => void;
  placeholder?: string;
  emptyMessage?: string;
  className?: string;
}

export function ServiceCombobox({
  onServiceSelect,
  placeholder = "Search services...",
  emptyMessage = "No service found.",
  className
}: ServiceComboboxProps) {
  const [open, setOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState("")

  // Use TanStack Query to fetch services with fallback
  const { data: apiServices = [], isLoading, isError } = useDigitalServicesCatalog();

  // Use fallback data if API returns empty array
  const services = apiServices.length > 0 ? apiServices : FALLBACK_SERVICES;

  const [selectedService, setSelectedService] = useState<DigitalServiceItem | null>(null);

  console.log("ServiceCombobox rendering with", services.length, "services");

  // Filter services based on search query
  const filteredServices = useMemo(() => {
    if (!searchQuery) return services;
    return services.filter(service =>
      service.name.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [searchQuery, services]);

  const handleSelect = useCallback((service: DigitalServiceItem) => {
    setSelectedService(service);
    setOpen(false);
    onServiceSelect(service);
  }, [onServiceSelect]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className={cn("w-full justify-between", className)}
          onClick={() => {
            console.log("Trigger clicked, open:", !open);
            setOpen(!open);
          }}
        >
          {selectedService ? selectedService.name : placeholder}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        className="w-full p-0"
        align="start"
        side="bottom"
        sideOffset={5}
        sticky="always"
        avoidCollisions={false}
      >
        <Command className="overflow-hidden">
          <CommandInput
            placeholder={placeholder}
            value={searchQuery}
            onValueChange={(value) => {
              console.log("Search query changed:", value);
              setSearchQuery(value);
            }}
          />
          <div className="max-h-[300px] overflow-y-auto">
            {isLoading ? (
              <CommandEmpty className="text-sm text-gray-500 text-center p-2">Loading services...</CommandEmpty>
            ) : filteredServices.length === 0 ? (
              <CommandEmpty>{emptyMessage}</CommandEmpty>
            ) : (
              <CommandGroup>
                {filteredServices.map((service) => (
                  <CommandItem
                    key={service.id}
                    value={service.name}
                    onSelect={() => {
                      console.log("Service selected:", service.name);
                      handleSelect(service);
                    }}
                  >
                    <Check
                      className={cn(
                        "mr-2 h-4 w-4",
                        selectedService?.id === service.id ? "opacity-100" : "opacity-0"
                      )}
                    />
                    {service.name}
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