import { useMemo } from "react";

// Recharts components
import { Bar, BarChart, CartesianGrid, LabelList, XAxis, YAxis } from "recharts"

// Types
import { chartConfig } from "./itemsByPlatformChard.const"
import type { PlatformItem } from "./itemsByPlatformCard.types"

// ShadCN UI components
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/shared/components/ui/card"
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/shared/components/ui/chart"

// Utils
import { normalizePlatformName } from '@/features/dashboard/lib/utils/normalizePlatformName';

// Icons
import { TrendingUp } from "lucide-react"

type ItemsByPlatformCardProps = {
  platformList: PlatformItem[];
  newItemsThisMonth: number;
};

export function ItemsByPlatformCard({
  platformList,
  newItemsThisMonth,
}: ItemsByPlatformCardProps) {

  const chartData = useMemo(() => {
    // Sort by item count descending to show the most popular platforms at the top
    const sortedData = [...platformList].sort((a, b) => b.itemCount - a.itemCount);

    return sortedData.map(item => {
      const { displayName } = normalizePlatformName(item.platform);
      return {
        ...item,
        platform: displayName, // Use the clean name for the legend
      };
    });
  }, [platformList]);

  const totalItemCount = useMemo(() => {
    return platformList.reduce((sum, platform) => sum + platform.itemCount, 0);
  }, [platformList]);

  return (
    <Card className="flex flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>Games by Platform</CardTitle>
      </CardHeader>
      <CardContent className="flex-1 pb-0">
        <ChartContainer
          config={chartConfig}
          className="mx-auto aspect-square max-h-[400px]"
        >
          <BarChart
            accessibilityLayer
            data={chartData}
            layout="vertical"
            margin={{ left: 0, right: 30, top: 10, bottom: 10 }}
          >
            <CartesianGrid horizontal={false} />
            <YAxis
              dataKey="platform"
              type="category"
              tickLine={false}
              tickMargin={10}
              axisLine={false}
              width={120}
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
            />
            <XAxis dataKey="itemCount" type="number" hide />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent indicator="line" />}
            />
            <Bar
              dataKey="itemCount"
              fill="hsl(var(--chart-1))"
              radius={4}
              barSize={25}
            >
              <LabelList
                dataKey="itemCount"
                position="right"
                offset={8}
                className="fill-foreground"
                fontSize={12}
              />
            </Bar>
          </BarChart>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col gap-2 text-sm">
        <div className="flex items-center gap-2 font-medium leading-none">
          {newItemsThisMonth} new titles this month <TrendingUp className="h-4 w-4" />
        </div>
        <div className="leading-none text-muted-foreground">
          Showing total <span className="font-bold text-white">{totalItemCount}</span> titles across <span className="font-bold text-white">{platformList.length}</span> platforms
        </div>
      </CardFooter>
    </Card>
  )
}
