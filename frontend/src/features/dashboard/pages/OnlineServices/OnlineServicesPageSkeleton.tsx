import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';
import { Skeleton } from '@/shared/components/ui/skeleton';

export function OnlineServicesPageSkeleton() {
  return (
    <PageMain>
      <PageHeadline>
        <Skeleton className="h-10 w-[320px]" />
        <div className="flex items-center space-x-2">
          <Skeleton className="h-10 w-[200px]" />
        </div>
      </PageHeadline>

      <PageGrid>
        {Array.from({ length: 12 }).map((_, i) => (
          <div key={i} className="md:col-span-1">
            <Skeleton className="h-[100px] w-full" />
          </div>
        ))}
      </PageGrid>
    </PageMain>

  )
}