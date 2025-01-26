// Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { NotificationList } from '@/features/dashboard/components/organisms/NotificationsPage/NotificationsList';

// Mock data
import { notificationsMockData } from '../../components/organisms/NotificationsPage/notificationsPage.mockdata';

export function NotificationsPageContent() {
  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>Notifications</h1>
        </div>
      </PageHeadline>

      <div className='flex flex-col gap-4'>
        <NotificationList initialNotifications={notificationsMockData} />
      </div>
    </PageMain>
  );
}
