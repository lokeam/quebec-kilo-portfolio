import { LibraryBig } from 'lucide-react';
import { Card } from '@/shared/components/ui/card';
import { Skeleton } from '@/shared/components/ui/skeleton';

type SearchResultProps = {
  title: string;
  imageUrl: string;
  isInLibrary: boolean;
}

export function SearchResult({
  title,
  // imageUrl,
 // price,
  isInLibrary = false,
}: SearchResultProps) {
  return (
    <Card className="relative flex items-center gap-4 p-4 transition-all duration-200 bg-[#2A2A2A] hover:bg-[#E5E5E5] group">
      {isInLibrary && (
        <div className="absolute left-0 top-0 bottom-0 z-10">
          <div className="flex h-full items-center">
            <div className="flex items-center bg-[#1A9FFF] h-[34px] overflow-hidden transition-all duration-200 group-hover:w-[120px] w-[34px] rounded-r-md">
              <div className="flex items-center gap-2 px-2 w-[120px]">
                <LibraryBig className="w-5 h-5 text-black shrink-0" />
                <span className="text-sm font-medium text-black whitespace-nowrap">IN LIBRARY</span>
              </div>
            </div>
          </div>
        </div>
      )}
      <div className="relative w-16 h-20 shrink-0">
        <div className="relative aspect-[4/3] md:aspect-square">
          <Skeleton className="absolute inset-0 w-full h-full rounded-lg" />
        </div>
      </div>
      <div className="flex flex-col">
        <h3 className="font-medium text-white group-hover:text-black">
          {title}
        </h3>
        {/* <p className="text-white/80 group-hover:text-black/80">
          ${price.toFixed(2)}
        </p> */}
      </div>
    </Card>
  );
}
