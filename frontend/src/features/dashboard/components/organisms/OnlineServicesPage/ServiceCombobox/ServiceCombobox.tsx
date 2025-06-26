import { useMemo, useCallback, useState, useEffect } from 'react';

// Shadcn UI Components
import { Button } from '@/shared/components/ui/button';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from '@/shared/components/ui/command';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/shared/components/ui/popover';

// Icons
import { Check, ChevronsUpDown } from 'lucide-react';

// Utils
import { cn } from '@/shared/components/ui/utils';

// Query Hooks
import { useGetDigitalServicesCatalog } from '@/core/api/queries/digitalServicesCatalog.queries';

// Types
import type { DigitalServiceItem } from '@/types/domain/digital-location';

interface ServiceComboboxProps {
  onServiceSelect: (service: {
    name: string;
    url: string;
    isSubscriptionService: boolean;
  }) => void;
  placeholder?: string;
  emptyMessage?: string;
  className?: string;
  initialValue?: string;
}

export function ServiceCombobox({
  onServiceSelect,
  placeholder = "Search services...",
  emptyMessage = "No service found.",
  className,
  initialValue
}: ServiceComboboxProps) {
  const [open, setOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState("")
  const [selectedService, setSelectedService] = useState<DigitalServiceItem | null>(null)

  // Use TanStack Query to fetch services
  const { data: services = [], isLoading } = useGetDigitalServicesCatalog();

  // Find initial service if provided
  const initialService = useMemo(() => {
    if (!initialValue) return null;
    return services.find(service => service.name === initialValue) || null;
  }, [initialValue, services]);

  // Set initial service when it's found
  useEffect(() => {
    if (initialService && !selectedService) {
      setSelectedService(initialService);
    }
  }, [initialService, selectedService]);

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
    onServiceSelect({
      name: service.name,
      url: service.url,
      isSubscriptionService: service.isSubscriptionService
    });
  }, [onServiceSelect]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className={cn("w-full justify-between", className)}
        >
          {selectedService ? selectedService.name : placeholder}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        className="w-[300px] p-0"
        align="start"
        side="bottom"
        sideOffset={5}
        sticky="always"
        avoidCollisions={false}
      >
        <Command className="overflow-visible">
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
            <CommandGroup className="overflow-visible">
              {isLoading ? (
                <CommandEmpty className="text-sm text-gray-500 text-center p-2">Loading services...</CommandEmpty>
              ) : filteredServices.length === 0 ? (
                <CommandEmpty>{emptyMessage}</CommandEmpty>
              ) : (
                filteredServices.map((service) => (
                  <CommandItem
                    key={service.id}
                    value={service.name}
                    onSelect={() => handleSelect(service)}
                  >
                    <Check
                      className={cn(
                        "mr-2 h-4 w-4",
                        selectedService?.id === service.id ? "opacity-100" : "opacity-0"
                      )}
                    />
                    {service.name}
                  </CommandItem>
                ))
              )}
            </CommandGroup>
          </div>
        </Command>
      </PopoverContent>
    </Popover>
  )
}