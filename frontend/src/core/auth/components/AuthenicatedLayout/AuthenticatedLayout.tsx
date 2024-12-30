
import { TopNavigation } from '@/features/navigation/organisms/TopNav/TopNavigation';
import { NotifyButton } from '@/features/navigation/molecules/NotifyButton/NotifyButton';
import { SearchButton } from '@/features/navigation/molecules/SearchButton/SearchButton';
import { AppSidebar } from '@/features/navigation/organisms/SideNav/AppSidebar';
import { AvatarDropDownMenu } from '@/features/navigation/molecules/AvatarDropDownMenu/AvatarDropDownMenu';
import { cn } from '@/shared/components/ui/utils';

//import { OfflineBanner } from '@/core/network-status/components/OfflineBanner';

import { Outlet } from 'react-router-dom';

export default function AuthenticatedLayout() {
  return (
    <>
      {/* Shadcn UI Sidebar */}
      <AppSidebar />

      <div
        id="content"
        className={cn(
          'max-w-full w-full ml-auto',
          'peer-data-[state=collapsed]:w-[calc(100%-var(--sidebar-width-icon)-1rem)]',
          'peer-data-[state=expanded]:w-[calc(100%-var(--sidebar-width))]',
          'transition-[width] ease-linear duration-200',
          'h-svh flex flex-col',
          'group-data-[scroll-locked=1]/body:h-full',
          'group-data-[scroll-locked=1]/body:has-[main.fixed-main]:h-svh'
        )}
      >
        {/* Network Status Provider - Offline Banner */}
        {/* <OfflineBanner /> */}

        {/* Top Navigation w/ Sidebar Toggle */}
          <TopNavigation>
            <div className='ml-auto flex items-center space-x-4'>
              {/* Search Bar */}
              <SearchButton />

              {/* Notifications */}
              <NotifyButton />

              {/* User Login Avatar */}
              <AvatarDropDownMenu />
            </div>
          </TopNavigation>

          <Outlet />
      </div>
    </>
  );
};
