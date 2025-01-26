import { type ReactNode } from 'react';

// Components
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';
import { LibraryMediaItem } from '@/features/dashboard/components/organisms/LibraryPage/LibraryMediaItem/LibraryMediaItem';
import { MemoizedLibraryMediaListItem } from '@/features/dashboard/components/organisms/LibraryPage/LibraryMediaListItem/LibraryMediaListItem';

// Utils + Hooks
import {
  ViewModes,
  useLibraryGames,
  useLibraryPlatformFilter,
} from '@/features/dashboard/lib/stores/libraryStore';
import { useLibrarySearchQuery } from '@/features/dashboard/lib/stores/libraryStore';
import { useLibraryTitle } from '@/features/dashboard/lib/hooks/useLibraryTitle';

// Types
import { type ComponentType } from 'react';
import { useFilteredLibraryItems } from '../../lib/hooks/useFilteredLibraryItems';

interface LibraryLayoutProps {
  viewMode: typeof ViewModes[keyof typeof ViewModes];
  EmptyPage: ComponentType;
  Toolbar: ComponentType;
  title: ReactNode;
  baseTitle?: string;
}

export function LibraryLayoutContainer({
  viewMode,
  EmptyPage,
  Toolbar,
  baseTitle,
}: LibraryLayoutProps) {

  /* Grab data from Zustand store */
  const services = useLibraryGames();
  const platformFilter = useLibraryPlatformFilter();
  const searchQuery = useLibrarySearchQuery();

  /* Combined filtering for both platform and title search */
  const filteredServices = useFilteredLibraryItems(services, platformFilter, searchQuery);

  console.log('filteredServices', filteredServices);


  /* Pass filtered data to title hook */
  const { title, countText } = useLibraryTitle({
    baseTitle: baseTitle || '',
    filteredCount: filteredServices.length,
    platformFilter,
  });

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
        <EmptyPage />
      </PageMain>
    );
  }

  /* Render content based on view mode */
  const renderContent = () => {
    if (services.length === 0) {
      return <NoResultsFound />;
    }

    const CardComponent = viewMode === ViewModes.GRID
      ? LibraryMediaItem
      : MemoizedLibraryMediaListItem;

    return (
      <div className="flex h-full w-full flex-wrap content-start">
        {filteredServices.map((item, index) => (
          <CardComponent
            key={`${item.title}-${index}`}
            {...item}
          />
        ))
        }
      </div>
    );
  };

  return (
    <PageMain>
      <Toolbar />
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
