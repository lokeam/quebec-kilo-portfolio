// Components
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { LibraryMediaItem } from '@/features/dashboard/components/organisms/LibraryPage/LibraryMediaItem/LibraryMediaItem';

// Mock Data
import { libraryPageMockData } from './LibraryPage.mockdata';

export function LibraryPageContent() {
  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>Library Page</h1>
        </div>
      </PageHeadline>

      <div className="flex h-full w-full flex-wrap content-start">
          {libraryPageMockData.map((game) => (
            <LibraryMediaItem
              key={game.id}
              href={`steam://rungameid/${game.id}`}
              imageUrl={game.image}
            />
          ))}
      </div>

    </PageMain>
  );
}
