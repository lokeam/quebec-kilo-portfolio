import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
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

      {/* Online Services Toolbar Skeleton */}
      <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
        <div className="flex flex-wrap items-center gap-3">
          <Skeleton className="h-10 w-[300px]" />
          <Skeleton className="h-4 w-16" />
          <Skeleton className="h-10 w-[140px]" />
          <Skeleton className="h-10 w-[180px]" />
        </div>

        <div className="flex items-center gap-3">
          <Skeleton className="h-4 w-8" />
          <div className="flex rounded-md p-1 gap-1">
            <Skeleton className="h-8 w-8 rounded" />
            <Skeleton className="h-8 w-8 rounded" />
            <Skeleton className="h-8 w-8 rounded" />
          </div>
        </div>
      </div>

      {/* Online Services Grid Skeleton */}
      <div className="p-4 border rounded-md">
        <Skeleton className="h-6 w-48 mb-2" />
        <Skeleton className="h-4 w-32 mb-4" />
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 2xl:grid-cols-3">
          {Array.from({ length: 12 }).map((_, i) => (
            <Skeleton
              key={i}
              className="h-[100px] w-full border rounded-lg"
            />
          ))}
        </div>
      </div>
    </PageMain>
  )
}