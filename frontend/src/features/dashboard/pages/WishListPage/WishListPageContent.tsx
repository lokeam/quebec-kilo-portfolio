import { useMemo } from 'react';

// Components
import { PageHeadline } from "@/shared/components/layout/page-headline";
import { PageMain } from "@/shared/components/layout/page-main";
import { WishlistItemCard } from '@/features/dashboard/organisms/WishlistItemCard/WishlistItemCard';

// Mock Data
import { wishlistPageMockData } from '@/features/dashboard/pages/WishListPage/WishListPage.mockdata';

export function WishListPageContent() {
  const mockData = useMemo(() => wishlistPageMockData, []);

  // Memoize boolean checks for each platform section
  const { hasPcItems, hasConsoleItems, hasMobileItems } = useMemo(() => {
    return {
      hasPcItems: mockData?.pc.length > 0,
      hasConsoleItems: mockData?.console.length > 0,
      hasMobileItems: mockData?.mobile.length > 0
    }
  }, [mockData.pc, mockData.console, mockData.mobile]);

  // Memoize the wishlist item renderers
  const pcItems = useMemo(() => {
    return hasPcItems && mockData.pc.map((item, index) => (
      <WishlistItemCard key={`${item.id}-${index}-${item.title}`} {...item} index={index} />
    ));
  }, [hasPcItems, mockData.pc]);

  const consoleItems = useMemo(() => {
    return hasConsoleItems && mockData.console.map((item, index) => (
      <WishlistItemCard key={`${item.id}-${index}-${item.title}`} {...item} index={index} />
    ));
  }, [hasConsoleItems, mockData.console]);

  const mobileItems = useMemo(() => {
    return hasMobileItems && mockData.mobile.map((item, index) => (
      <WishlistItemCard key={`${item.id}-${index}-${item.title}`} {...item} index={index} />
    ));
  }, [hasMobileItems, mockData.mobile]);


  return (
    <PageMain>
      <PageHeadline>
        <h1 className="text-2xl font-bold tracking-tight">Wish List Page</h1>
      </PageHeadline>

      <div className="flex h-full w-full flex-wrap content-start">

        {hasPcItems && <h2 className="text-2xl font-bold tracking-tight my-4">PC Wishlist</h2>}
        {pcItems}

        {hasConsoleItems && <h2 className="text-2xl font-bold tracking-tight mt-8 mb-4">Console Wishlist</h2>}
        {consoleItems}

        {hasMobileItems && <h2 className="text-2xl font-bold tracking-tight mt-8 mb-4">Mobile Wishlist</h2>}
        {mobileItems}
      </div>

    </PageMain>
  );
}
