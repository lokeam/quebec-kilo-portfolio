import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from '@/shared/components/ui/card';
import { SINGLE_STATISTICS_ICONS } from './constants/singleStatCardIcons';
import type { StatCardIconType } from './types/singleStatCard.types';

interface EmptySingleStatisticsCardProps {
  title: string;
  value: number;
  lastUpdated: string;
  icon: string;
};

  export const EmptySingleStatisticsCard = ({
  title,
  icon,
}: EmptySingleStatisticsCardProps) => {
    const IconComponent = SINGLE_STATISTICS_ICONS[icon as StatCardIconType];

  return (
    <Card className="col-span-full md:col-span-1">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-1">
        <CardTitle className="text-md font-medium">{title}</CardTitle>
        <IconComponent className="h-[2rem] w-[2rem]" />
      </CardHeader>
      <CardContent>
        <p className="text-xl font-bold">Your total count of games will appear here</p>
        <p className="text-sm text-muted-foreground">Since the last time you updated your library</p>
      </CardContent>
    </Card>
  );
};
