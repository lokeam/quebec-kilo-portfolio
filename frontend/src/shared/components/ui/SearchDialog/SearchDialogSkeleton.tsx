import { Skeleton } from '@/shared/components/ui/skeleton'

export function SearchDialogSkeleton() {
  return (
    <>
      {Array(9).fill(0).map((_, index) => (
        <Skeleton
          key={index}
          className="h-[146px] w-full border rounded-lg"
        />
      ))}
    </>
  );
}
