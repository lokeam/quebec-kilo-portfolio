

// Components
import { PageHeadline } from "@/shared/components/layout/page-headline";
import { PageMain } from "@/shared/components/layout/page-main";
import { WishlistItemCard } from '@/features/dashboard/organisms/WishlistItemCard/WishlistItemCard';

// Mock Data
import { wishlistPageMockData } from '@/features/dashboard/pages/WishListPage/WishListPage.mockdata';

export function WishListPageContent() {

  const hasPcItems = wishlistPageMockData?.pc.length > 0;
  const hasConsoleItems = wishlistPageMockData?.console.length > 0;
  const hasMobileItems = wishlistPageMockData?.mobile.length > 0;

  return (
    <PageMain>
      <PageHeadline>
        <h1 className="text-2xl font-bold tracking-tight">Wish List Page</h1>
      </PageHeadline>

      <div className="flex h-full w-full flex-wrap content-start">

        {hasPcItems && <h2 className="text-2xl font-bold tracking-tight my-4">PC Wishlist</h2>}
        {hasPcItems && wishlistPageMockData.pc.map((item, index) => (
          <WishlistItemCard key={`${item.id}-${index}`} {...item} index={index} />
        ))}

        {hasConsoleItems && <h2 className="text-2xl font-bold tracking-tight mt-8 mb-4">Console Wishlist</h2>}
        {hasConsoleItems && wishlistPageMockData.console.map((item, index) => (
          <WishlistItemCard key={`${item.id}-${index}`} {...item} index={index} />
        ))}

        {hasMobileItems && <h2 className="text-2xl font-bold tracking-tight mt-8 mb-4">Mobile Wishlist</h2>}
        {hasMobileItems && wishlistPageMockData.mobile.map((item, index) => (
          <WishlistItemCard key={`${item.id}-${index}`} {...item} index={index} />
        ))}
      </div>

      <h2 className="text-2xl font-bold tracking-tight">Console Wishlist</h2>
    </PageMain>
  );
}