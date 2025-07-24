import { useEffect, useMemo, useCallback } from 'react';

// Components
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { LibraryPageToolbar } from '@/features/dashboard/components/organisms/LibraryPage/LibraryPageToolbar/LibraryPageToolbar';
import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';
import { LibraryGameCard } from '@/features/dashboard/components/organisms/LibraryPage/LibraryGameCard/LibraryGameCard';
import { MemoizedLibraryGameDetailCard } from '@/features/dashboard/components/organisms/LibraryPage/LibraryGameDetailCard/LibraryGameDetailCard';
import { LibraryPageSkeleton } from '@/features/dashboard/components/organisms/LibraryPage/LibraryPageSkeleton';

// Utils + Hooks
import { useLibraryStore } from '@/features/dashboard/lib/stores/libraryStore';
import { useGetLibraryPageBFFResponse } from '@/core/api/queries/gameLibrary.queries';
import { useLibraryFilters } from '@/features/dashboard/lib/hooks/useLibraryFilters';
import { useShowConditionalIntroToasts } from '@/features/dashboard/hooks/intro-toasts/useShowConditionalIntroToasts';

// Types
import type { LibraryGameItemRefactoredResponse } from '@/types/domain/library-types';

export function LibraryPageContent() {
  const { viewMode, setGames } = useLibraryStore();
  const services = useLibraryStore((state) => state.userGames);
  const { searchQuery, platformFilters, locationFilters } = useLibraryStore();

  const {
    data: bffResponse,
    isLoading,
    error
  } = useGetLibraryPageBFFResponse();

  const libraryItems = useMemo(() => bffResponse?.libraryItems ?? [], [bffResponse?.libraryItems]);

  // Get filter options from library data
  const filterOptions = useLibraryFilters(libraryItems);

  // Show intro toast for adding games to library
  const hasGamesInLibrary = Boolean(services && services.length > 0);
  useShowConditionalIntroToasts(7, hasGamesInLibrary);



  // Set games in store when page mounts or when library items change
  useEffect(() => {
    // console.log('🔍 DEBUG: LibraryPageContent useEffect:', {
    //   libraryItemsLength: libraryItems.length,
    //   storeGamesLength: services.length
    // });

    // Only update if we have new library items
    if (libraryItems.length > 0) {
      setGames(libraryItems);
    }
  }, [libraryItems, setGames, services.length]);


  // Memoize filter function
  const filterGame = useCallback((game: LibraryGameItemRefactoredResponse) => {
    // Search filter - match against game name
    const matchesSearch = !searchQuery ||
      game.name.toLowerCase().includes(searchQuery.toLowerCase());

    // Platform filter - match against any of the selected platforms
    const matchesPlatform = platformFilters.length === 0 ||
      [
        // Check physical locations
        ...(game.physicalLocations || []).flatMap(location =>
          (location.gamePlatformVersions || []).map(platform => platform.platformName)
        ),
        // Check digital locations
        ...(game.digitalLocations || []).flatMap(location =>
          (location.gamePlatformVersions || []).map(platform => platform.platformName)
        )
      ].some(platformName => platformFilters.includes(platformName));

    // Location filter - match against any of the selected locations
    const matchesLocation = locationFilters.length === 0 ||
      [
        // Check physical locations
        ...(game.physicalLocations || []).flatMap(location => [
          location.sublocationName || '',
          location.parentLocationName || ''
        ]),
        // Check digital locations
        ...(game.digitalLocations || []).map(location => location.digitalLocationName || '')
      ].some(locationName => locationFilters.includes(locationName));

    return matchesSearch && matchesPlatform && matchesLocation;
  }, [searchQuery, platformFilters, locationFilters]);

  // Memoize filtered results
  const filteredServices = useMemo(() => {
    if (!services || services.length === 0) return [];

    // Early return if no filters are active
    if (searchQuery === '' && platformFilters.length === 0 && locationFilters.length === 0) {
      return services;
    }

    return services.filter(filterGame);
  }, [services, filterGame, searchQuery, platformFilters.length, locationFilters.length]);

  // Handle loading state
  if (isLoading) {
    return <LibraryPageSkeleton />;
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
            <h1 className='text-2xl font-bold tracking-tight'>All Games</h1>
            <span className='text-[20px] text-gray-500 ml-1'>(0 games)</span>
          </div>
        </PageHeadline>
        <NoResultsFound />
      </PageMain>
    );
  }

  /* Render content based on view mode */
  const renderContent = () => {
    // console.log('🔍 DEBUG: LibraryPageContent render:', {
    //   servicesLength: services.length,
    //   filteredServicesLength: filteredServices.length,
    //   searchQuery,
    //   platformFilters,
    //   locationFilters,
    //   viewMode
    // });

    if (filteredServices.length === 0) {
      return <NoResultsFound />;
    }

    const CardComponent = viewMode === 'grid'
      ? LibraryGameCard
      : MemoizedLibraryGameDetailCard;


    return (
      <div className="flex h-full w-full flex-wrap content-start">
        {filteredServices.map((item: LibraryGameItemRefactoredResponse, index) => (
          <CardComponent
            key={`${item.name}-${index}`}
            index={index}
            id={item.id}
            name={item.name}
            coverUrl={item.coverUrl}
            firstReleaseDate={item.firstReleaseDate}
            rating={0}
            themeNames={item.genreNames ?? null}
            isInLibrary={true}
            isInWishlist={item.isInWishlist}
            gameType={item.gameType}
            favorite={item.favorite}
            physicalLocations={item.physicalLocations}
            digitalLocations={item.digitalLocations}
            onRemoveFromLibrary={() => {}}
            totalDigitalVersions={item.totalDigitalVersions}
            totalPhysicalVersions={item.totalPhysicalVersions}
          />
        ))}
      </div>
    );
  };

  const title = 'All Games';

  return (
    <PageMain>
      <LibraryPageToolbar
        platforms={filterOptions.platforms}
        locations={filterOptions.locations}
      />
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-3xl font-bold tracking-tight'>
            {title}
            <span className='text-[20px] text-gray-500 ml-1'>({filteredServices.length} {filteredServices.length === 1 ? 'game' : 'games'})</span>
          </h1>
        </div>
      </PageHeadline>
      {renderContent()}
    </PageMain>
  );
}
