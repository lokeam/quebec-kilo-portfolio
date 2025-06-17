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

const chartConfig = {
  oneTimePurchase: {
    label: "One Time Purchase",
    color: "hsl(var(--chart-1))",
  },
  dlc: {
    label: "DLC",
    color: "hsl(var(--chart-2))",
  },
  inGamePurchase: {
    label: "In Game Purchase",
    color: "hsl(var(--chart-3))",
  },
  hardware: {
    label: "Hardware",
    color: "hsl(var(--chart-4))",
  },
  subscription: {
    label: "Subscription",
    color: "hsl(var(--chart-5))",
  },
} satisfies ChartConfig;

type MediaTypeDomain = "games" | "movies" | "oneTimePurchase" | "hardware" | "dlc" | "inGamePurchase" | "subscription";
type MonthlyExpenditureItem = {
  date: string;
  oneTimePurchase: number;
  hardware: number;
  dlc: number;
  inGamePurchase: number;
  subscription: number;
}

type MonthlySpendingCardProps = {
  mediaTypeDomains: MediaTypeDomain[];
  monthlyExpenditures: MonthlyExpenditureItem[];
};

export function MonthlySpendingCard({
  mediaTypeDomains,
  monthlyExpenditures,
}: MonthlySpendingCardProps) {

  // Filter to only valid domains present in chartConfig
  const validDomains = mediaTypeDomains.filter((domain) => chartConfig[domain as keyof typeof chartConfig]);
  if (process.env.NODE_ENV === 'development') {
    const invalidDomains = mediaTypeDomains.filter((domain) => !chartConfig[domain as keyof typeof chartConfig]);
    if (invalidDomains.length > 0) {
      // eslint-disable-next-line no-console
      console.warn('Some mediaTypeDomains are not present in chartConfig:', invalidDomains);
    }
  }

  return (
    <Card className="col-span-full lg:col-span-2 flex flex-col h-full">
      <CardHeader>
        <CardTitle>Monthly Game Spending</CardTitle>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <BarChart
            accessibilityLayer
            data={monthlyExpenditures}
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
            {validDomains.map((domain, idx) => (
              <Bar
                key={domain}
                dataKey={domain}
                stackId="a"
                fill={chartConfig[domain as keyof typeof chartConfig].color}
                radius={idx === validDomains.length - 1 ? [4, 4, 0, 0] : [0, 0, 0, 0]}
              />
            ))}
            <ChartTooltip
              content={
                <ChartTooltipContent
                  hideLabel
                  className="w-[180px]"
                  formatter={(value, name) => (
                    <>
                      <div
                        className="h-2.5 w-2.5 shrink-0 rounded-[2px] bg-[--color-bg]"
                        style={{
                          "--color-bg": chartConfig[name as keyof typeof chartConfig]?.color || "#8884d8",
                        } as React.CSSProperties}
                      />
                      {chartConfig[name as keyof typeof chartConfig]?.label || name}
                      <div className="ml-auto flex items-baseline gap-0.5 font-mono font-medium tabular-nums text-foreground">
                        ${value}
                      </div>
                    </>
                  )}
                  // Optionally, you can add a total below if desired
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
