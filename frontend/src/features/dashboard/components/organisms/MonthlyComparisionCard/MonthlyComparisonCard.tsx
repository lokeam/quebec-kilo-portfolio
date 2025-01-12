import { TrendingDown, TrendingUp } from 'lucide-react'
import { Label, PolarRadiusAxis, RadialBar, RadialBarChart } from "recharts"

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/shared/components/ui/card"
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/shared/components/ui/chart"

interface MonthlyComparisonCardProps {
  currentMonth: string;
  previousMonth: string;
  currentMonthTotal: number;
  previousMonthTotal: number;
  chartData: {
    previousMonthTotal: number;
    currentMonthTotal: number;
  };
}

const chartConfig = {
  previousMonthTotal: {
    label: "Last:",
    color: "hsl(var(--chart-1))",
  },
  currentMonthTotal: {
    label: "Current:",
    color: "hsl(var(--chart-2))",
  },
} satisfies ChartConfig

export function MonthlyComparisonCard({
  currentMonth,
  previousMonth,
  currentMonthTotal,
  previousMonthTotal,
  chartData,
}: MonthlyComparisonCardProps) {
  const percentageChange = ((currentMonthTotal - previousMonthTotal) / previousMonthTotal) * 100
  const isIncrease = percentageChange > 0

  return (
    <Card className="flex flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>Monthly Comparison</CardTitle>
        <CardDescription>{previousMonth} - {currentMonth}</CardDescription>
      </CardHeader>
      <CardContent className="flex flex-1 items-center pb-0">
        <ChartContainer
          config={chartConfig}
          className="mx-auto aspect-square w-full max-w-[250px]"
        >
          <RadialBarChart
            data={[chartData]}
            endAngle={180}
            innerRadius={80}
            outerRadius={130}
          >
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent hideLabel />}
            />
            <PolarRadiusAxis tick={false} tickLine={false} axisLine={false}>
              <Label
                content={({ viewBox }) => {
                  if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                    return (
                      <text x={viewBox.cx} y={viewBox.cy} textAnchor="middle">
                        <tspan
                          x={viewBox.cx}
                          y={(viewBox.cy || 0) - 16}
                          className="fill-foreground text-2xl font-bold"
                        >
                          ${currentMonthTotal.toLocaleString()}
                        </tspan>
                        <tspan
                          x={viewBox.cx}
                          y={(viewBox.cy || 0) + 4}
                          className="fill-muted-foreground"
                        >
                          {currentMonth.charAt(0).toUpperCase() + currentMonth.slice(1)}
                        </tspan>
                      </text>
                    )
                  }
                }}
              />
            </PolarRadiusAxis>
            <RadialBar
              dataKey="previousMonthTotal"
              stackId="a"
              cornerRadius={5}
              fill="var(--color-previousMonthTotal)"
              className="stroke-transparent stroke-2"
            />
            <RadialBar
              dataKey="currentMonthTotal"
              fill="var(--color-currentMonthTotal)"
              stackId="a"
              cornerRadius={5}
              className="stroke-transparent stroke-2"
            />
          </RadialBarChart>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col gap-2 text-sm">
        <div className="flex items-center gap-2 font-medium leading-none">
          {isIncrease ? "Trending up" : "Trending down"} by{" "}
          <span className={isIncrease ? "text-red-500" : "text-green-500"}>
            {Math.abs(percentageChange).toFixed(1)}%
          </span>{" "}
          this month{" "}
          {isIncrease ? (
            <TrendingUp className="h-4 w-4" />
          ) : (
            <TrendingDown className="h-4 w-4" />
          )}
        </div>
        <div className="leading-none text-muted-foreground">
          Comparing service fees for the last 2 months
        </div>
      </CardFooter>
    </Card>
  )
}

