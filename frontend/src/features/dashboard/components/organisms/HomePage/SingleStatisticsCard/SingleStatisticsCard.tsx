// ShadCN UI components
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from '@/shared/components/ui/card';

// Constants
import { SINGLE_STATISTICS_ICONS } from './constants/singleStatCardIcons';

// Types
import type { StatCardIconType } from './types/singleStatCard.types';

// Utils
import { formatCurrency } from '@/features/dashboard/lib/utils/formatCurrency';
import { formatMonthsAgo } from '@/features/dashboard/lib/utils/formatMonthsAgo';

interface SingleStatisticsCardProps {
  stats: {
    title: string;
    icon: string;
    value: number;
    secondaryValue?: number;
    lastUpdated: number;
  }
};

const CARD_SIZE_MAP = {
  sm: 'md:col-span-1',
  lg: 'lg:col-span-2 flex flex-col h-full',
};

const renderCardContentValue = (
  value: number,
  icon: string,
  secondaryValue?: number,
) => {
  switch (icon) {
    case 'package':
      return (
        <>
          {`${secondaryValue}`}
          <span className="mx-1 text-sm">sub</span>
          {secondaryValue &&
            <>
              <span className="mx-1 text-sm">|</span>
              <span className="mx-1 text-md">{value}</span>
              <span className="mx-1 text-sm">parent</span>
            </>
          }
        </>
      );
    case 'coin':
      //return formatCurrency(value);
      return (
        <>
          {`${formatCurrency(value)}`}
          <span className="mx-1 text-sm">this month</span>
        </>
      )
    default:
      return value.toLocaleString();
  }
}

export const SingleStatisticsCard = ({
  stats,
}: SingleStatisticsCardProps) => {
    const { title, value, lastUpdated, icon, secondaryValue = 0 } = stats;
    const IconComponent = SINGLE_STATISTICS_ICONS[icon as StatCardIconType];
    const formattedLastUpdatedTimestamp = formatMonthsAgo(lastUpdated);

  return (
    <Card className={`col-span-full ${CARD_SIZE_MAP['sm']} overflow-hidden`}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-1">
        <CardTitle className="text-md font-medium">{title}</CardTitle>
        <IconComponent className="h-[2rem] w-[2rem]" />
      </CardHeader>
      <CardContent>
        <p className="text-3xl font-bold">{renderCardContentValue(value, icon, secondaryValue)}</p>
        <p className="text-sm text-muted-foreground">{formattedLastUpdatedTimestamp}</p>
      </CardContent>
    </Card>
  );
};
