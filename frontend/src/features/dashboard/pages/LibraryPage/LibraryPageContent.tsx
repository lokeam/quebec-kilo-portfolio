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
import { useGetAllLibraryGames } from '@/core/api/queries/gameLibrary.queries';

export function LibraryPageContent() {
  const { viewMode, setGames } = useLibraryStore();
  const services = useLibraryStore((state) => state.userGames);
  const platformFilter = useLibraryStore((state) => state.platformFilter);
  const searchQuery = useLibraryStore((state) => state.searchQuery);

  const { data: libraryGames } = useGetAllLibraryGames();

  // Set games in store when page mounts
  useEffect(() => {
    if (libraryGames) {
      setGames(libraryGames);
    }
  }, [libraryGames, setGames]);

  /* Combined filtering for both platform and title search */
  const filteredServices = useFilteredLibraryItems(services, platformFilter, searchQuery);

  /* Pass filtered data to title hook */
  const { title, countText } = useLibraryTitle({
    baseTitle: 'All Games',
    filteredCount: filteredServices.length,
    platformFilter,
  });

  console.log('filteredServices', filteredServices);


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
    if (services.length === 0) {
      return <NoResultsFound />;
    }

    const CardComponent = viewMode === 'grid'
      ? LibraryMediaItem
      : MemoizedLibraryMediaListItem;

    return (
      <div className="flex h-full w-full flex-wrap content-start">
        {filteredServices.map((item, index) => (
          <CardComponent
            key={`${item.name}-${index}`}
            index={index}
            id={item.id}
            name={item.name}
            coverUrl={item.coverUrl}
            firstReleaseDate={item.firstReleaseDate}
            rating={item.rating}
            themeNames={item.themeNames}
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
