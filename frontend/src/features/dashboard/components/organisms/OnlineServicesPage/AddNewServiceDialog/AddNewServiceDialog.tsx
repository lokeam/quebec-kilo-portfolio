import { useState } from 'react';
import { DialogTrigger } from '@/shared/components/ui/dialog';
import { Button } from '@/shared/components/ui/button';
import { Plus } from 'lucide-react';
import { SearchDialog, SearchDialogSkeleton } from '@/shared/components/ui/SearchDialog';
import { SingleOnlineServiceCard } from '@/features/dashboard/components/organisms/OnlineServicesPage/SingleOnlineServiceCard/SingleOnlineServiceCard';
import { useDebounce } from '@/shared/hooks/useDebounce';
import { useAvailableServices } from '@/features/dashboard/lib/hooks/useAvailableService';

export function AddNewServiceDialog() {
  const [isOpen, setIsOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const debouncedSearchQuery = useDebounce(searchQuery, 400);
  const { availableServices, isLoading, error } = useAvailableServices(debouncedSearchQuery);

  const handleOpenChange = (open: boolean) => {
    setIsOpen(open);
    if (!open) {
      setSearchQuery('');
    }
  };

  const handleSearchQueryChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(event.target.value);
  };

  return (
    <SearchDialog
      open={isOpen}
      onOpenChange={handleOpenChange}
      searchQuery={searchQuery}
      onSearchChange={handleSearchQueryChange}
      searchPlaceholder="Search for a service by name"
      dialogTitle="Add a New Service"
      hideHeader={false}
      trigger={
        <DialogTrigger asChild>
          <Button>
            <Plus className="h-4 w-4" />
            New Service
          </Button>
        </DialogTrigger>
      }
      footer={
        <Button>Manually Add a Service</Button>
      }
    >
      {isLoading ? (
        <SearchDialogSkeleton />
      ) : error ? (
        <div className="text-red-500 p-4">
          Error loading services. Please try again later.
        </div>
      ) : availableServices.length === 0 ? (
        <div className="text-muted-foreground p-4 text-center">
          {debouncedSearchQuery ? 'No services found' : 'Start typing to search for services'}
        </div>
      ) : (
        availableServices.map((service, index) => (
          <SingleOnlineServiceCard
            key={`${service.name}-${index}`}
            {...service}
            onClick={() => {
              console.log('route to add service page');
              // Handle service selection
            }}
          />
        ))
      )}
    </SearchDialog>
  );
}
