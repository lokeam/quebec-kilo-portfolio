import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
import { Button } from '@/shared/components/ui/button';
import { Plus } from 'lucide-react';

import { Skeleton } from '@/shared/components/ui/skeleton';

// Mock Data
import { onlineServicesPageMockData } from './onlineServicesPage.mockdata';

export function OnlineServicesPageContent() {
  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items center'>
          <h1 className='text-2xl font-bold tracking-tight'>Online Services</h1>
        </div>
        <div className='flex items-center space-x-2'>
          <Button>
            <Plus className="h-4 w-4" />
            New Service
          </Button>
        </div>
      </PageHeadline>

      {/* <PageGrid columns={{ sm: 1, md: 2, lg: 3 }}> */}
      <div className="grid grid-cols-1 lg:grid-cols-4 gap-2">

        <div className="lg:col-span-3 grid grid-cols-1 md:grid-cols-2 xl:grid-cols-2 gap-4">
          {onlineServicesPageMockData.map((service, index) => (
            <Skeleton
              key={`${service.name}-${index}`}
              className="w-full h-[100px]"
            />
          ))}
        </div>

        <div className="space-y-4">
          <Skeleton className="w-full h-[200px]" />
          <Skeleton className="w-full h-[300px]" />
        </div>

        </div>
      {/* </PageGrid> */}
    </PageMain>
  );
}
