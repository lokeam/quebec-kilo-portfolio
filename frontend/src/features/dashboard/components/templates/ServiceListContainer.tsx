
// Components
import { PageHeadline } from '@/shared/components/layout/page-headline';
import { PageMain } from '@/shared/components/layout/page-main';
import { NoResultsFound } from '@/features/dashboard/components/molecules/NoResultsFound';

// Utils + Hooks
import type { ViewMode } from '@/shared/constants/viewModes';

// Types
import { type ComponentType } from 'react';
import type { DigitalLocation } from '@/types/domain/digital-location';

interface ServiceListContainerProps {
  services: DigitalLocation[];
  totalServices: number;
  viewMode: ViewMode;
  title: string;
  EmptyPage: ComponentType;
  AddNewDialog: ComponentType;
  Toolbar: ComponentType;
  Table: ComponentType<{ services: DigitalLocation[] }>;
  Card: ComponentType<DigitalLocation & { isWatchedByResizeObserver: boolean }>;
  NoResultsFound: ComponentType;
  containerClassName?: string;
  hasCSSGridLayout?: boolean;
}

export function ServiceListContainer({
  services,
  totalServices,
  viewMode,
  title,
  EmptyPage,
  AddNewDialog,
  Toolbar,
  Table,
  Card,
  containerClassName = 'grid grid-cols-1 gap-4',
  hasCSSGridLayout = true,
}: ServiceListContainerProps) {
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

    if (viewMode === 'table') {
      return <Table services={services} />;
    }

    return (
      <div className={`${containerClassName} ${
        viewMode === 'grid'
        && hasCSSGridLayout
        ? 'md:grid-cols-2 2xl:grid-cols-3'
        : ''
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
