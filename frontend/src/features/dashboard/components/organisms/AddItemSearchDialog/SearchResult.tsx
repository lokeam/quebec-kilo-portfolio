
// Components
import { Card } from '@/shared/components/ui/card';
import { Button } from '@/shared/components/ui/button';
import { ImageWithFallback } from '@/shared/components/ui/ImageWithFallback/ImageWithFallback';

// Icons
import { LibraryBig } from 'lucide-react';
import { IconHeart } from '@tabler/icons-react';

// Hooks
import { toast } from 'sonner';


type SearchResultProps = {
  title: string;
  imageUrl: string | null | undefined;
  isInLibrary: boolean;
}

export function SearchResult({
  title,
  imageUrl,
 // price,
  isInLibrary = false,
}: SearchResultProps) {

  const handleAddToLibrary = () => {
    toast( `${title} successfully added to library`,{
      description: 'You can now access it in your library',
      className: 'bg-green-500 text-white',
      duration: 2500,
    });
  };

  const handleAddToWishlist = () => {
    toast(`${title} successfully added to wishlist`, {
      description: 'You can now access it in your wishlist',
      className: 'bg-green-500 text-white',
      duration: 2500,
    });
  };

  return (
    <Card className="relative flex items-center transition-all duration-200 bg-[#2A2A2A] hover:bg-[#E5E5E5] group overflow-hidden">
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
      <div className="shrink-0 p-2">
        <div className="relative w-24 md:w-32 p-2">
          <ImageWithFallback
            src={imageUrl}
            alt={title}
            width={292}
            height={120}
            className="rounded-sm w-[140px] h-full object-cover"
          />
        </div>
      </div>

      <div className="flex flex-1 min-w-0 items-center pr-2">
        <div className="flex-1 min-w-0"> {/* nested min-w-0 for text truncation */}
          <h3 className="font-medium text-white text-wrap max-w-[140px] max-h-[48px] md:max-w-full md:max-h-unset truncate">
            {title}
          </h3>
        </div>

        {
          !isInLibrary && (
            <div className="flex shrink-0 gap-1 mt-1 ml-2">
              <Button variant="outline" onClick={handleAddToLibrary}>
                <LibraryBig className="w-5 h-5" />
                <span className="hidden md:block text-sm font-medium text-white whitespace-nowrap">Add to library</span>
              </Button>
              <Button variant="outline" onClick={handleAddToWishlist}>
                <IconHeart className="w-5 h-5" />
                <span className="hidden md:block text-sm font-medium text-white whitespace-nowrap">Add to wishlist</span>
              </Button>
            </div>
          )
        }
      </div>
    </Card>
  );
}
