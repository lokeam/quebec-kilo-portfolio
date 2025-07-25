// Template
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { PageGrid } from '@/shared/components/layout/page-grid';

// Components
import { SettingsForm } from '@/features/dashboard/pages/SettingsPage/SettingsForm';

export function SettingsPageContent() {
  // console.log(`SettingsPageContent`);

  return (
    <PageMain>
      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>Settings</h1>
        </div>
      </PageHeadline>

      <PageGrid>
        <SettingsForm />
      </PageGrid>
    </PageMain>
  )
}

