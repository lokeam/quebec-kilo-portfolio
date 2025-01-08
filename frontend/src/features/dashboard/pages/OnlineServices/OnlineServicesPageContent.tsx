import { useState } from 'react';
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { Button } from '@/shared/components/ui/button';
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/shared/components/ui/dialog';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Plus } from 'lucide-react';

// Components
import { OnlineServicesToolbar } from '@/features/dashboard/organisms/OnlineServicesToolbar/OnlineServicesToolbar';
import { SingleOnlineServiceCard } from '@/features/dashboard/organisms/SingleOnlineServiceCard/SingleOnlineServiceCard';
import { OnlineServicesTable } from '@/features/dashboard/organisms/OnlineServicesTable/OnlineServicesTable';

// Utils + Hooks
import { useDebounce } from '@/shared/hooks/useDebounce';
import { useCardLabelWidth } from '@/features/dashboard/organisms/SingleOnlineServiceCard/useCardLabelWidth';
import { useOnlineServicesStore } from '@/features/dashboard/stores/onlineServicesStore';
import { useAvailableServices } from '@/features/dashboard/hooks/useAvailableService';
import { ViewModes } from '@/features/dashboard/stores/onlineServicesStore';

// Mock Data
import { onlineServicesPageMockData } from './onlineServicesPage.mockdata';
import { useFilteredServices } from '@/features/dashboard/hooks/useFilteredServices';
import { Input } from '@/shared/components/ui/input';

export function OnlineServicesPageContent() {
  const [searchQuery, setSearchQuery] = useState('');
  const debouncedSearchQuery = useDebounce(searchQuery, 400);
  const { availableServices, isLoading, error } = useAvailableServices(debouncedSearchQuery);

  const {
    viewMode,
  } = useOnlineServicesStore();
  const filteredServices = useFilteredServices(onlineServicesPageMockData);

  useCardLabelWidth({
    selectorAttribute: '[data-card-sentinel]',
    breakpoints: {
      narrow: 320,
      medium: 360
    },
    widths: {
      narrow: '60px',
      medium: '100px',
      wide: '200px'
    }
  });

  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>Online Services</h1>
        </div>
        <div className='flex items-center space-x-2'>
            <Dialog>
              <DialogTrigger>
              <Button>
                <Plus className="h-4 w-4" />
                New Service
            </Button>
          </DialogTrigger>
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
        </div>
      </PageHeadline>

      <OnlineServicesToolbar />

      {viewMode === ViewModes.TABLE ? (
        <OnlineServicesTable services={filteredServices} />
      ) : (
        <div className={`grid grid-cols-1 gap-4 ${
          viewMode === ViewModes.GRID
            ? 'md:grid-cols-2 2xl:grid-cols-3'
            : ''
        }`}>
          {filteredServices.length > 0 ? (
            filteredServices.map((service, index) => (
              <SingleOnlineServiceCard
                key={`${service.name}-${index}`}
                {...service}
                isWatchedByResizeObserver={index === 0}
              />
            ))
          ) : (
            <div className="col-span-full">
              <p>No online services found</p>
            </div>
          )}
        </div>
      )}

      <div className="mt-4 space-y-4">
        <Skeleton className="w-full h-[200px]" />
        <Skeleton className="w-full h-[300px]" />
      </div>
    </PageMain>
  );
}
