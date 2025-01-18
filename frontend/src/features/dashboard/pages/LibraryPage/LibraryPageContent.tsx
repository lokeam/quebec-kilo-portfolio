import { useEffect } from 'react';

// Components
import { LibraryPageToolbar } from '@/features/dashboard/components/organisms/LibraryPage/LibraryPageToolbar/LibraryPageToolbar';
import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';
import { LibraryLayoutContainer } from '@/features/dashboard/components/templates/LibraryLayoutContainer';

// Utils + Hooks
import { useLibraryStore } from '@/features/dashboard/lib/stores/libraryStore';

// Mock Data
import { libraryPageMockData } from './LibraryPage.mockdata';

export function LibraryPageContent() {
  const { viewMode, setGames } = useLibraryStore();

  // Set games in store when page mounts
  useEffect(() => {
    setGames(libraryPageMockData);
  }, [setGames]);

  return (
    <LibraryLayoutContainer
      viewMode={viewMode}
      EmptyPage={NoResultsFound}
      Toolbar={LibraryPageToolbar}
      title="All Games"
    />
  );
}
