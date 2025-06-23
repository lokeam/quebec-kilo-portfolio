// ShadCN Recharts Components
import { Bar, BarChart, XAxis } from "recharts"

// ShadCN Components
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

// Utils
import { formatCurrency, type CurrencyAmount } from "@/features/dashboard/lib/utils/formatCurrency";

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


type MonthlySpendingCardProps = {
  mediaTypeDomains: string[];
  monthlyExpenditures: object[];
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
                  className="w-[220px] p-2"
                  formatter={(value, name) => (
                    <>
                      <span
                        className="inline-block h-2.5 w-2.5 rounded-[2px] mr-2 align-middle"
                        style={{
                          background: chartConfig[name as keyof typeof chartConfig]?.color || "#8884d8",
                        }}
                      />
                      {chartConfig[name as keyof typeof chartConfig]?.label || name}
                      <div className="ml-auto flex items-baseline gap-0.5 font-mono font-medium tabular-nums text-foreground">
                        {formatCurrency(value as CurrencyAmount)}
                      </div>
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
