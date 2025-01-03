import { Bar, BarChart, XAxis } from "recharts"

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/shared/components/ui/card"
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/shared/components/ui/chart"
import type { Domain, MonthlyDomainExpenditure } from "./monthlySpendingCard.types"

const chartConfig = {
  games: {
    label: "Games",
    color: "hsl(var(--chart-1))",
  },
  movies: {
    label: "Movies",
    color: "hsl(var(--chart-2))",
  },
} satisfies ChartConfig;


type MonthlySpendingCardProps = {
  domains: Domain[];
  spendingByMonth: MonthlyDomainExpenditure[];
};

export function MonthlySpendingCard({
  domains,
  spendingByMonth,
}: MonthlySpendingCardProps) {

  console.log(domains, spendingByMonth);

  return (
    <Card className="col-span-full lg:col-span-2 flex flex-col h-full">
      <CardHeader>
        <CardTitle>Monthly Spending Across Games and Movies</CardTitle>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <BarChart
            accessibilityLayer
            data={spendingByMonth}
            margin={{ top: 20, right: 20, bottom: 20, left: 20 }}
          >
            <XAxis
              dataKey="date"
              tickLine={false}
              tickMargin={10}
              axisLine={false}
              tickFormatter={(value) => {
                return new Date(value).toLocaleDateString("en-US", {
                  month: "short",
                })
              }}
            />
            <Bar
              dataKey="games"
              stackId="a"
              fill="var(--color-games)"
              radius={[0, 0, 4, 4]}
            />
            <Bar
              dataKey="movies"
              stackId="a"
              fill="var(--color-movies)"
              radius={[4, 4, 0, 0]}
            />
            <ChartTooltip
              content={
                <ChartTooltipContent
                  hideLabel
                  className="w-[180px]"
                  formatter={(value, name, item, index) => (
                    <>
                      <div
                        className="h-2.5 w-2.5 shrink-0 rounded-[2px] bg-[--color-bg]"
                        style={
                          {
                            "--color-bg": `var(--color-${name})`,
                          } as React.CSSProperties
                        }
                      />
                      {chartConfig[name as keyof typeof chartConfig]?.label ||
                        name}
                      <div className="ml-auto flex items-baseline gap-0.5 font-mono font-medium tabular-nums text-foreground">
                        ${value}
                      </div>
                      {/* Add after the last item */}
                      {index === 1 && (
                        <div className="mt-1.5 flex basis-full items-center border-t pt-1.5 text-xs font-medium text-foreground">
                          Total
                          <div className="ml-auto flex items-baseline gap-0.5 font-mono font-medium tabular-nums text-foreground">
                            ${item.payload.games + item.payload.movies}
                          </div>
                        </div>
                      )}
                    </>
                  )}
                />
              }
              cursor={false}
              defaultIndex={1}
            />
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
