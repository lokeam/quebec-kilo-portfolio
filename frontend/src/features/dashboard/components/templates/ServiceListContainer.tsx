
// Components
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';

// Utils + Hooks
import { ViewModes } from '@/features/dashboard/lib/stores/onlineServicesStore';

// Types
import { type ComponentType } from 'react';
import type { Service } from '@/features/dashboard/lib/types/service.types';

interface ServiceListContainerProps<T extends Service> {
  services: T[];
  totalServices: number;
  viewMode: typeof ViewModes[keyof typeof ViewModes];
  title: string;
  EmptyPage: ComponentType;
  AddNewDialog: ComponentType;
  Toolbar: ComponentType;
  Table: ComponentType<{ services: T[] }>;
  Card: ComponentType<T & { isWatchedByResizeObserver: boolean }>;
  NoResultsFound: ComponentType;
}

export function ServiceListContainer<T extends Service>({
  services,
  totalServices,
  viewMode,
  title,
  EmptyPage,
  AddNewDialog,
  Toolbar,
  Table,
  Card,
}: ServiceListContainerProps<T>) {
  /* True empty state - first time user zero services */
  if (totalServices === 0) {
    return (
      <PageMain>
        <PageHeadline>
          <div className='flex items-center'>
            <h1 className='text-2xl font-bold tracking-tight'>{title}</h1>
          </div>
        </PageHeadline>
        <EmptyPage />
      </PageMain>
    );
  }
  const renderContent = () => {
    if (services.length === 0) {
      return <NoResultsFound />;
    }

    if (viewMode === ViewModes.TABLE) {
      return <Table services={services} />;
    }

    return (
      <div className={`grid grid-cols-1 gap-4 ${
        viewMode === ViewModes.GRID ? 'md:grid-cols-2 2xl:grid-cols-3' : ''
      }`}>
        {services.map((service, index) => (
          <Card
            key={`${service.name}-${index}`}
            {...service}
            isWatchedByResizeObserver={index === 0}
          />
        ))}
      </div>
    );
  };

  /* Standard layout with toolbar for all non-empty states */
  return (
    <PageMain>

      <PageHeadline>
        <div className='flex items-center'>
          <h1 className='text-2xl font-bold tracking-tight'>{title}</h1>
        </div>
        <div className='flex items-center space-x-2'>
          <AddNewDialog />
        </div>
      </PageHeadline>

      <Toolbar />

      {renderContent()}

    </PageMain>
  );
}
