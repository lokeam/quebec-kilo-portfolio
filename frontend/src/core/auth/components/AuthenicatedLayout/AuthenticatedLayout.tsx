// Components
import { TopNavigation } from '@/features/navigation/organisms/TopNav/TopNavigation';
import { AppSidebar } from '@/features/navigation/organisms/SideNav/AppSidebar';
import { AvatarDropDownMenu } from '@/features/navigation/molecules/AvatarDropDownMenu/AvatarDropDownMenu';
import { OfflineBanner } from '@/core/network-status/components/OfflineBanner';
import { GameSearchAndSelectDialog } from '@/features/dashboard/components/organisms/GameSearchAndSelectDialog/GameSearchAndSelectDialog';
import { FeedbackButton } from '@/shared/components/ui/feedback-button/FeedbackButton';

// Utils
import { Outlet } from 'react-router-dom';
import { cn } from '@/shared/components/ui/utils';

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
        <OfflineBanner />

        {/* Top Navigation w/ Sidebar Toggle */}
          <TopNavigation>
            <div className='ml-auto flex items-center space-x-4'>
              {/* Search Bar */}
              <GameSearchAndSelectDialog />

              {/* Notifications */}
              {/* <NotifyPopover /> */}

              {/* Feedback Button */}
              <FeedbackButton variant="outline" size="sm" />

              {/* User Login Avatar */}
              <AvatarDropDownMenu />
            </div>
          </TopNavigation>

          {/* Main content area */}
          <Outlet />
      </div>
    </>
  );
};
