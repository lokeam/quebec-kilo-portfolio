import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { Skeleton } from '@/shared/components/ui/skeleton';

export function LibraryPageSkeleton() {
  return (
    <PageMain>
      {/* Library Page Toolbar Skeleton */}
      <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
        <div className="flex flex-wrap items-center gap-3">
          <Skeleton className="h-10 w-[300px]" />
          <Skeleton className="h-4 w-16" />
          <Skeleton className="h-10 w-[200px]" />
          <Skeleton className="h-10 w-[200px]" />
        </div>

        <div className="flex items-center gap-3">
          <Skeleton className="h-4 w-8" />
          <div className="flex rounded-md p-1 gap-1">
            <Skeleton className="h-8 w-8 rounded" />
            <Skeleton className="h-8 w-8 rounded" />
          </div>
        </div>
      </div>

      <PageHeadline>
        <div className="flex items-center">
          <Skeleton className="h-9 w-[200px]" />
          <Skeleton className="h-6 w-[100px] ml-1" />
        </div>
      </PageHeadline>

      {/* Library Games Content Skeleton */}
      <div className="flex h-full w-full flex-wrap content-start gap-4">
        {Array.from({ length: 12 }).map((_, i) => (
          <Skeleton
            key={i}
            className="h-[280px] w-[200px] border rounded-lg"
          />
        ))}
      </div>
    </PageMain>
  )
}