import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from '@/shared/components/ui/card';
import { SINGLE_STATISTICS_ICONS } from './constants/singleStatCardIcons';
import type { StatCardIconType } from './types/singleStatCard.types';
import { formatMonthsAgo } from '@/features/dashboard/lib/utils/formatMonthsAgo';


interface SingleStatisticsCardProps {
  // title: string;
  // value: number;
  // lastUpdated: number;
  // icon: string;
  // size: 'sm' | 'lg';
  stats: {
    title: string;
    icon: string;
    value: number;
    lastUpdated: number;
  }
};

const CARD_SIZE_MAP = {
  sm: 'md:col-span-1',
  lg: 'lg:col-span-2 flex flex-col h-full',
};

export const SingleStatisticsCard = ({
  stats,
}: SingleStatisticsCardProps) => {
    const { title, value, lastUpdated, icon } = stats;
    const IconComponent = SINGLE_STATISTICS_ICONS[icon as StatCardIconType];
    const formattedLastUpdatedTimestamp = formatMonthsAgo(lastUpdated);

  return (
    <Card className={`col-span-full ${CARD_SIZE_MAP['sm']}`}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-1">
        <CardTitle className="text-md font-medium">{title}</CardTitle>
        <IconComponent className="h-[2rem] w-[2rem]" />
      </CardHeader>
      <CardContent>
        <p className="text-3xl font-bold">{value}</p>
        <p className="text-sm text-muted-foreground">{formattedLastUpdatedTimestamp}</p>
      </CardContent>
    </Card>
  );
};
