import { useState } from 'react';
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/shared/components/ui/dialog';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Plus } from 'lucide-react';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { SingleOnlineServiceCard } from '@/features/dashboard/organisms/SingleOnlineServiceCard/SingleOnlineServiceCard';
import { useDebounce } from '@/shared/hooks/useDebounce';
import { useAvailableServices } from '@/features/dashboard/hooks/useAvailableService';

export function AddNewServiceDialog() {
  const [searchQuery, setSearchQuery] = useState('');
  const debouncedSearchQuery = useDebounce(searchQuery, 400);
  const { availableServices, isLoading, error } = useAvailableServices(debouncedSearchQuery);

  return (
    <Dialog>
      {/* Trigger & button */}
      <DialogTrigger asChild>
        <Button>
          <Plus className="h-4 w-4" />
          New Service
        </Button>
      </DialogTrigger>

      {/* Dialog Modal + Content */}
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add a New Service</DialogTitle>
        </DialogHeader>

        <div className="grid gap-4 py-4">
          <Input
            placeholder="Search for a service by name"
            value={searchQuery}
            onChange={(event) => setSearchQuery(event.target.value)}
            className="w-full"
          />
          <div className="flex flex-col space-y-2">
            { isLoading ?(
                Array(3).fill(0).map((_, index) => (
                <div key={index} className="p-4 border rounded-lg">
                  <div className="flex items-center space-x-4">
                    <Skeleton className="h-12 w-12 rounded" />
                    <div className="space-y-2">
                      <Skeleton className="h-4 w-[200px]" />
                      <Skeleton className="h-4 w-[100px]" />
                    </div>
                  </div>
                </div>
              ))
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
                      console.log('route to add service page')
                      // Handle service selection
                    }}
                  />
                ))
              )
            }
          </div>
        </div>

        <DialogFooter>
          <Button>Manually Add a Service</Button>
        </DialogFooter>
      </DialogContent>

    </Dialog>
  )
}