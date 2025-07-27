// Shadcn UI Components / Charts
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { type ChartConfig, ChartContainer, ChartTooltip } from '@/shared/components/ui/chart';
import { Bar, BarChart, CartesianGrid, LabelList, XAxis } from 'recharts';

// Icons
import { TrendingUp } from '@/shared/components/ui/icons';

// Utils
import { formatCurrency, type CurrencyAmount } from "@/features/dashboard/lib/utils/formatCurrency";

// Types
import type { SingleMonthlyExpenditureBFFResponse } from '@/types/domain/spend-tracking';

const chartConfig = {
  expenditure: {
    label: "Expenditure",
    color: "hsl(var(--chart-1))",
  },
} satisfies ChartConfig;

// TODO: Normalize Card colors + styling post design a/b test
export function TotalAnnualSpendingCard({
  totalAnnualSpending
}: {
  totalAnnualSpending: {
    dateRange: string;
    monthlyExpenditures: SingleMonthlyExpenditureBFFResponse[];
    medianMonthlyCost: number;
  }
}) {
  const { dateRange, monthlyExpenditures, medianMonthlyCost } = totalAnnualSpending;

  return (
    <Card className="col-span-full lg:col-span-2">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-xl font-bold text-foreground">Net This Year</CardTitle>
        {/* <button className="inline-flex items-center text-xs text-muted-foreground hover:text-foreground">
          DOWNLOAD REPORT
          <ArrowUpRight className="ml-1 h-3 w-3" />
        </button> */}
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          <div className="text-md font-bold">{dateRange}</div>
        </div>
          <ChartContainer config={chartConfig}>
            <BarChart
              accessibilityLayer
              data={monthlyExpenditures}
              margin={{
                top: 20,
                right: 10,
                left: 10,
              }}
            >
              <CartesianGrid vertical={false} />
              <XAxis
                dataKey="month"
                tickLine={false}
                tickMargin={10}
                axisLine={false}
                tickFormatter={(value) => value.slice(0, 3)}
              />
              <ChartTooltip
                cursor={false}
                content={({ active, payload }) => {
                  if (!active || !payload?.length) return null;

                  const item = payload[0];
                  const capitalizedName = typeof item.name === 'string'
                    ? item.name.charAt(0).toUpperCase() + item.name.slice(1)
                    : 'Expenditure';
                  return (
                    <div className="rounded-lg border border-border/50 bg-background px-2.5 py-1.5 text-xs shadow-xl">
                      <div className="flex items-center gap-2">
                        <div className="h-2.5 w-2.5 rounded-[2px] bg-blue-500" />
                        <span className="text-muted-foreground">
                          {capitalizedName}
                        </span>
                        <span className="font-mono font-medium tabular-nums text-foreground">
                          {formatCurrency(item.value as CurrencyAmount)}
                        </span>
                      </div>
                    </div>
                  );
                }}
              />
              <Bar
                dataKey="expenditure"
                fill="var(--color-expenditure)"
                radius={8}
                barSize={20}
                minPointSize={4}
              >
                <LabelList
                  position="top"
                  offset={12}
                  className="fill-foreground"
                  fontSize={12}
                  formatter={(value: CurrencyAmount) => formatCurrency(value)}
                />
              </Bar>
            </BarChart>
          </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col items-start gap-2 text-sm">
        <div className="flex gap-2 font-medium leading-none">
          Average {formatCurrency(medianMonthlyCost as CurrencyAmount)} <TrendingUp className="h-4 w-4" />
        </div>
        <div className="leading-none text-muted-foreground">
          Showing all monthly expenditures for the last 12 months
        </div>
      </CardFooter>
    </Card>
  );
}
