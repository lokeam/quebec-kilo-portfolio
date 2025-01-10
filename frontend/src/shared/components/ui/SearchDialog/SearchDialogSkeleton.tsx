import { Skeleton } from '@/shared/components/ui/skeleton'

export function SearchDialogSkeleton() {
  return (
    <>
      {Array(3).fill(0).map((_, index) => (
        <div key={index} className="p-4 border rounded-lg">
          <div className="flex items-center space-x-4">
            <Skeleton className="h-12 w-12 rounded" />
            <div className="space-y-2">
              <Skeleton className="h-4 w-[200px]" />
              <Skeleton className="h-4 w-[100px]" />
            </div>
          </div>
        </div>
      ))}
    </>
  );
}
