// Components
import { PageMain } from '@/shared/components/layout/page-main';
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { NotificationList } from '@/features/dashboard/components/organisms/NotificationsPage/NotificationsList';
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
        {/* <Timeline>
          {notificationsMockData.map((notification, index) => (
            <TimelineItem key={`${index}-${notification.notificationTitle}`}>
              <TimelineSeparator>
                {index !== 0 && <TimelineConnector />}
                <TimelineDot />
                {index !== notificationsMockData.length - 1 && <TimelineConnector />}
              </TimelineSeparator>
              <TimelineContent>
                <h3>{notification.notificationTitle}</h3>
                <p>{notification.notificationMsg}</p>
              </TimelineContent>
            </TimelineItem>
          ))}
        </Timeline> */}
      </div>
    </PageMain>
  );
}
