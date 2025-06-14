import { useEffect } from 'react';

// Components
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { LibraryPageToolbar } from '@/features/dashboard/components/organisms/LibraryPage/LibraryPageToolbar/LibraryPageToolbar';
import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';
import { LibraryMediaItem } from '@/features/dashboard/components/organisms/LibraryPage/LibraryMediaItem/LibraryMediaItem';
import { MemoizedLibraryMediaListItem } from '@/features/dashboard/components/organisms/LibraryPage/LibraryMediaListItem/LibraryMediaListItem';

// Utils + Hooks
import { useLibraryStore } from '@/features/dashboard/lib/stores/libraryStore';
import { useLibraryTitle } from '@/features/dashboard/lib/hooks/useLibraryTitle';
import { useFilteredLibraryItems } from '@/features/dashboard/lib/hooks/useFilteredLibraryItems';
import { useGetLibraryPageBFFResponse } from '@/core/api/queries/gameLibrary.queries';
import type { LibraryGameItemResponse } from '@/types/domain/library-types';

export function LibraryPageContent() {
  const { viewMode, setGames } = useLibraryStore();
  const services = useLibraryStore((state) => state.userGames);
  const platformFilter = useLibraryStore((state) => state.platformFilter);
  const searchQuery = useLibraryStore((state) => state.searchQuery);

  const {
    data: bffResponse,
    isLoading,
    error
  } = useGetLibraryPageBFFResponse();

  const libraryItems = bffResponse?.libraryItems ?? [];
  const recentlyAdded = bffResponse?.recentlyAdded ?? [];

  // Set games in store when page mounts or when library items change
  useEffect(() => {
    console.log('🔍 DEBUG: LibraryPageContent useEffect:', {
      libraryItemsLength: libraryItems.length,
      storeGamesLength: services.length
    });

    // Only update if we have new library items
    if (libraryItems.length > 0) {
      setGames(libraryItems);
    }
  }, [libraryItems, setGames]);

  useEffect(() => {
    console.log('🔍 DEBUG: LibraryPageContent useEffect:', {
      libraryItemsLength: libraryItems.length,
      storeGamesLength: services.length,
      bffResponse
    });
  }, [libraryItems, services, bffResponse]);

  /* Combined filtering for both platform and title search */
  const filteredServices = useFilteredLibraryItems(services, platformFilter, searchQuery);

  /* Pass filtered data to title hook */
  const { title, countText } = useLibraryTitle({
    baseTitle: 'All Games',
    filteredCount: filteredServices.length,
    platformFilter,
  });

  // Handle loading state
  if (isLoading) {
    return (
      <PageMain>
        <PageHeadline>
          <div className='flex items-center'>
            <h1 className='text-3xl font-bold tracking-tight'>Loading...</h1>
          </div>
        </PageHeadline>
      </PageMain>
    );
  }

  // Handle error state
  if (error) {
    return (
      <PageMain>
        <PageHeadline>
          <div className='flex items-center'>
            <h1 className='text-3xl font-bold tracking-tight'>Error loading library</h1>
          </div>
        </PageHeadline>
      </PageMain>
    );
  }

  /* Guard clause empty state - first time user zero services */
  if (services.length === 0) {
    return (
      <PageMain>
        <PageHeadline>
          <div className='flex items-center'>
            <h1 className='text-2xl font-bold tracking-tight'>{title}</h1>
            <span className='text-[20px] text-gray-500 ml-1'>{countText}</span>
          </div>
        </PageHeadline>
        <NoResultsFound />
      </PageMain>
    );
  }

  /* Render content based on view mode */
  const renderContent = () => {

    console.log('🔍 DEBUG: LibraryPageContent render:', {
      filteredServicesLength: filteredServices.length,
      viewMode
    });


    if (services.length === 0) {
      return <NoResultsFound />;
    }

    const CardComponent = viewMode === 'grid'
      ? LibraryMediaItem
      : MemoizedLibraryMediaListItem;

    return (
      <div className="flex h-full w-full flex-wrap content-start">
        {filteredServices.map((item: LibraryGameItemResponse, index) => (
          <CardComponent
            key={`${item.name}-${index}`}
            index={index}
            id={item.id}
            name={item.name}
            coverUrl={item.coverUrl}
            firstReleaseDate={item.firstReleaseDate}
            rating={item.rating}
            themeNames={item.themeNames ?? null}
            isInLibrary={item.isInLibrary}
            isInWishlist={item.isInWishlist}
            gameType={item.gameType}
            favorite={item.favorite}
            gamesByPlatformAndLocation={item.gamesByPlatformAndLocation}
            onRemoveFromLibrary={() => {}}
          />
        ))}
      </div>
    );
  };

  return (
    <PageMain>
      <LibraryPageToolbar />
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-3xl font-bold tracking-tight'>
            {title}
            <span className='text-[20px] text-gray-500 ml-1'>{countText}</span>
          </h1>
        </div>
      </PageHeadline>
      {renderContent()}
    </PageMain>
  );
}
