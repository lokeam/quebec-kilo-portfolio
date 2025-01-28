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
import { useDebounce } from "@/shared/hooks/useDebounce"
import { useAvailableServices } from "@/features/dashboard/lib/hooks/useAvailableService"

interface ServiceComboboxProps {
  onServiceSelect: (service: any) => void;
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
  const [open, setOpen] = React.useState(false)
  const [searchQuery, setSearchQuery] = React.useState("")
  const debouncedSearchQuery = useDebounce(searchQuery, 150)

  // Fetch all services when no search query is present
  const { availableServices: allServices, isLoading: isLoadingAll } = useAvailableServices("")
  // Fetch filtered services when there is a search query
  const { availableServices: searchedServices, isLoading: isLoadingSearch } =
    useAvailableServices(debouncedSearchQuery)

  const [selectedService, setSelectedService] = React.useState<any>(null)

  // Determine which services to display based on search state
  const displayedServices = React.useMemo(() => {
    if (debouncedSearchQuery) {
      return searchedServices
    }
    return allServices
  }, [debouncedSearchQuery, searchedServices, allServices])

  const isLoading = debouncedSearchQuery ? isLoadingSearch : isLoadingAll

  const handleSelect = React.useCallback((service: any) => {
    setSelectedService(service)
    setOpen(false)
    onServiceSelect(service)
  }, [onServiceSelect])

  return (
    <Popover open={open} onOpenChange={setOpen}>
    <PopoverTrigger asChild>
      <Button
        variant="outline"
        role="combobox"
        aria-expanded={open}
        className={cn("w-full justify-between", className)}
      >
        {selectedService ? selectedService.label : placeholder}
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
          onValueChange={setSearchQuery}
        />
        <div className="command-scroll-area max-h-[200px] overflow-y-auto">
          <div
            cmdk-list=""
            style={{
              "--cmdk-list-height": `${Math.min(displayedServices.length * 40, 300)}px`
            } as React.CSSProperties}
          >
            {isLoading ? (
              <CommandEmpty className="text-sm text-gray-500 text-center p-2">Loading services...</CommandEmpty>
            ) : displayedServices.length === 0 ? (
              <CommandEmpty>{emptyMessage}</CommandEmpty>
            ) : (
              <CommandGroup>
                {displayedServices.map((service) => (
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
                    {service.label}
                  </CommandItem>
                ))}
              </CommandGroup>
            )}
          </div>
        </div>
      </Command>
    </PopoverContent>
    </Popover>
  )
}