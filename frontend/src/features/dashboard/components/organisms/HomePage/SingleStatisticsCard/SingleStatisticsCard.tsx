import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from '@/shared/components/ui/card';
import { SINGLE_STATISTICS_ICONS } from './constants/singleStatCardIcons';
import type { StatCardIconType } from './types/singleStatCard.types';

interface SingleStatisticsCardProps {
  title: string;
  value: number;
  lastUpdated: string;
  icon: string;
  size: 'sm' | 'lg';
};

const CARD_SIZE_MAP = {
  sm: 'md:col-span-1',
  lg: 'lg:col-span-2 flex flex-col h-full',
};

export const SingleStatisticsCard = ({
  title,
  value,
  lastUpdated,
  icon,
  size = 'sm',
}: SingleStatisticsCardProps) => {
    const IconComponent = SINGLE_STATISTICS_ICONS[icon as StatCardIconType];

  return (
    <Card className={`col-span-full ${CARD_SIZE_MAP[size]}`}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-1">
        <CardTitle className="text-md font-medium">{title}</CardTitle>
        <IconComponent className="h-[2rem] w-[2rem]" />
      </CardHeader>
      <CardContent>
        <p className="text-3xl font-bold">{value}</p>
        <p className="text-sm text-muted-foreground">{`from ${lastUpdated} ago`}</p>
      </CardContent>
    </Card>
  );
};
