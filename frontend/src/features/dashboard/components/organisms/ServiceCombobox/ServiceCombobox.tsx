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
import { FALLBACK_SERVICES } from "@/core/api/services/digitalServices.service"
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services'
import type { ServiceType } from '@/shared/constants/service.constants'
import type { ServiceTierName } from '@/features/dashboard/lib/types/online-services/tiers'
import { BILLING_CYCLES } from '@/shared/constants/payment'

interface ServiceComboboxProps {
  onServiceSelect: (service: OnlineService) => void;
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
  const { data: apiServices = [], isLoading } = useDigitalServicesCatalog();

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
    console.log('Original Service:', service);
    setSelectedService(service);
    setOpen(false);
    // Transform DigitalServiceItem into OnlineService
    const onlineService: OnlineService = {
      id: service.id,
      name: service.name,
      label: service.name,
      logo: service.logo,
      type: service.isSubscriptionService ? 'subscription' : 'online' as ServiceType,
      isSubscriptionService: service.isSubscriptionService || false,
      status: 'active',
      url: '#',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      features: [],
      tier: {
        currentTier: 'standard' as ServiceTierName,
        availableTiers: [{
          name: 'standard' as ServiceTierName,
          features: [],
          id: 'standard',
          isDefault: true
        }]
      },
      billing: service.isSubscriptionService ? {
        cycle: BILLING_CYCLES.MONTHLY,
        fees: {
          monthly: '0'
        },
        paymentMethod: 'Generic'
      } : undefined
    };
    console.log('Transformed Service:', onlineService);
    onServiceSelect(onlineService);
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
            onValueChange={(value) => {
              console.log("Search query changed:", value);
              setSearchQuery(value);
            }}
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
                filteredServices.map((service: DigitalServiceItem) => (
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
                ))
              )}
            </CommandGroup>
          </div>
        </Command>
      </PopoverContent>
    </Popover>
  )
}