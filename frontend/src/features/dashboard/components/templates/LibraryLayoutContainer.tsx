import { type ReactNode } from 'react';

// Components
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';
import { LibraryMediaItem } from '@/features/dashboard/components/organisms/LibraryPage/LibraryMediaItem/LibraryMediaItem';
import { MemoizedLibraryMediaListItem } from '@/features/dashboard/components/organisms/LibraryPage/LibraryMediaListItem/LibraryMediaListItem';

// Utils + Hooks
import { ViewModes, useLibraryGames } from '@/features/dashboard/lib/stores/libraryStore';

// Types
import { type ComponentType } from 'react';

interface LibraryLayoutProps {
  viewMode: typeof ViewModes[keyof typeof ViewModes];
  EmptyPage: ComponentType;
  Toolbar: ComponentType;
  title: ReactNode;
}

export function LibraryLayoutContainer({ viewMode, EmptyPage, Toolbar, title }: LibraryLayoutProps) {
  const services = useLibraryGames();

  console.log('services', services);

  /* True empty state - first time user zero services */
  if (services.length === 0) {
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>{title}</h1>
        </div>
      </PageHeadline>
      <EmptyPage />
    </PageMain>
  }

  const renderContent = () => {
    if (services.length === 0) {
      return <NoResultsFound />;
    }

    const CardComponent = viewMode === ViewModes.GRID ? LibraryMediaItem : MemoizedLibraryMediaListItem;

    console.log('services', services);
    return (
      <div className="flex h-full w-full flex-wrap content-start">
        {services.map((item, index) => (
          <CardComponent
            key={`${item.id}-${index}`}
            index={index}
            {...item}
          />
        ))}
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
            <span className='text-[20px] text-gray-500 ml-1'>({services.length} games)</span>
          </h1>
        </div>
      </PageHeadline>

      {renderContent()}
    </PageMain>
  );
}
