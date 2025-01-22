
// Shadcn UI Components / Charts
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { type ChartConfig, ChartContainer, ChartTooltip, ChartTooltipContent } from '@/shared/components/ui/chart';
import { Bar, BarChart, CartesianGrid, LabelList, XAxis } from 'recharts';

// Icons
import { ArrowUpRight, TrendingUp } from 'lucide-react';

// Make up some random expenses between 80 and 500
const generateExpenditure = () => Math.floor(Math.random() * (500 - 80 + 1)) + 80;

// TODO: Sync all this hardcoded data with mockdata file
const chartData = [
  { month: "Jan", expenditure: generateExpenditure() },
  { month: "Feb", expenditure: generateExpenditure() },
  { month: "Mar", expenditure: generateExpenditure() },
  { month: "Apr", expenditure: generateExpenditure() },
  { month: "May", expenditure: generateExpenditure() },
  { month: "Jun", expenditure: generateExpenditure() },
  { month: "Jul", expenditure: generateExpenditure() },
  { month: "Aug", expenditure: generateExpenditure() },
  { month: "Sep", expenditure: generateExpenditure() },
  { month: "Oct", expenditure: generateExpenditure() },
  { month: "Nov", expenditure: generateExpenditure() },
  { month: "Dec", expenditure: generateExpenditure() },
];

// Use median monthly expenses for better measure of central tendency
const medianExpenditure = () => {
  const sorted = [...chartData].sort((a, b) => a.expenditure - b.expenditure)
  const middle = Math.floor(sorted.length / 2)
  return sorted.length % 2 === 0
    ? Math.round((sorted[middle - 1].expenditure + sorted[middle].expenditure) / 2)
    : sorted[middle].expenditure
};

const chartConfig = {
  expenditure: {
    label: "Expenditure",
    color: "hsl(var(--chart-1))",
  },
} satisfies ChartConfig;

// TODO: Normalize Card colors + styling post design a/b test
export function TotalAnnualSpendingCard() {
  const medianMonthlyCost = medianExpenditure();

  return (
    <Card className="dark col-span-full lg:col-span-2">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-xl font-bold text-[#c0c4cf]">Net This Year</CardTitle>
        <button className="inline-flex items-center text-xs text-[#8B949E] hover:text-white">
          DOWNLOAD REPORT
          <ArrowUpRight className="ml-1 h-3 w-3" />
        </button>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          <div className="text-md font-bold">January 2024 - January 2025</div>
        </div>
          <ChartContainer config={chartConfig}>
            <BarChart
              accessibilityLayer
            data={chartData}
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
            <ChartTooltip cursor={false} content={<ChartTooltipContent hideLabel />} />
            <Bar dataKey="expenditure" fill="var(--color-expenditure)" radius={8} barSize={20}>
              <LabelList position="top" offset={12} className="fill-foreground" fontSize={12} />
            </Bar>
          </BarChart>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col items-start gap-2 text-sm">
        <div className="flex gap-2 font-medium leading-none">
          Average ${medianMonthlyCost} <TrendingUp className="h-4 w-4" />
        </div>
        <div className="leading-none text-muted-foreground">
          Showing all monthly expenditures for the last 12 months
        </div>
      </CardFooter>
    </Card>
  );
}
