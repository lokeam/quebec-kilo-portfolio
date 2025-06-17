import * as React from "react"
import { TrendingUp } from "lucide-react"
import { Label, Pie, PieChart } from "recharts"
import { chartConfig } from "./itemsByPlatformChard.const"
import type { PlatformItem } from "./itemsByPlatformCard.types"
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
  ChartLegend,
  ChartLegendContent,
} from "@/shared/components/ui/chart"

type ItemsByPlatformCardProps = {
  totalItemCount: number;
  platformList: PlatformItem[];
  newItemsThisMonth: number;
};

export function ItemsByPlatformCard({
  totalItemCount,
  platformList,
  newItemsThisMonth,
}: ItemsByPlatformCardProps) {

  // Add fill colors to platform data using chartConfig
  const platformDataWithColors = platformList.map(item => ({
    ...item,
    fill: (chartConfig[item.platform as keyof typeof chartConfig]?.color) || 'hsl(var(--chart-pc))'
  }));

  return (
    <Card className="flex flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>Games by Platform</CardTitle>
      </CardHeader>
      <CardContent className="flex-1 pb-0">
        <ChartContainer
          config={chartConfig}
          className="mx-auto aspect-square max-h-[250px]"
        >
          <PieChart>
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent hideLabel />}
            />
            <Pie
              data={platformDataWithColors}
              dataKey="itemCount"
              nameKey="platform"
              innerRadius={60}
              strokeWidth={5}
            >
              <Label
                content={({ viewBox }) => {
                  if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                    return (
                      <text
                        x={viewBox.cx}
                        y={viewBox.cy}
                        textAnchor="middle"
                        dominantBaseline="middle"
                      >
                        <tspan
                          x={viewBox.cx}
                          y={viewBox.cy}
                          className="fill-foreground text-3xl font-bold"
                        >
                          {totalItemCount.toLocaleString()}
                        </tspan>
                        <tspan
                          x={viewBox.cx}
                          y={(viewBox.cy || 0) + 24}
                          className="fill-muted-foreground"
                        >
                          Titles
                        </tspan>
                      </text>
                    )
                  }
                }}
              />
            </Pie>
            <ChartLegend
              content={<ChartLegendContent nameKey="platform" />}
              className="-translate-y-2 flex-wrap gap-2 [&>*]:basis-1/4 [&>*]:justify-center"
            />
          </PieChart>
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
