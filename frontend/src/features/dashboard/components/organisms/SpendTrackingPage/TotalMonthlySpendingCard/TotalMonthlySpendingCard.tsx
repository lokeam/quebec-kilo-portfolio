// Icons
import { ArrowUpRight, TrendingDown } from '@/shared/components/ui/icons';

// Shadcn UI Components / Charts
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Cell, Pie, PieChart, ResponsiveContainer } from 'recharts';

// Types
import type { SpendingCategoryBFFResponse } from '@/types/domain/spend-tracking';

// Constants
import {
  GRAPH_CATEGORY_COLORS,
  GRAPH_CATEGORY_DISPLAY_NAMES
} from '@/shared/constants/graphCategoryColors';


export function TotalMonthlySpendingCard({
  totalMonthlySpending
}: {
  totalMonthlySpending: {
    currentMonthTotal: number;
    lastMonthTotal: number;
    percentageChange: number;
    comparisonDateRange: string;
    spendingCategories: SpendingCategoryBFFResponse[];
  }
}) {
  const {
    currentMonthTotal,
    lastMonthTotal,
    percentageChange,
    comparisonDateRange,
    spendingCategories
  } = totalMonthlySpending;

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(value);
  };

  return (
    <Card className="col-span-full lg:col-span-2">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-xl font-bold text-foreground">Net this month</CardTitle>
        <button className="inline-flex items-center text-xs text-muted-foreground hover:text-foreground">
          DOWNLOAD REPORT
          <ArrowUpRight className="ml-1 h-3 w-3" />
        </button>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          <div>
            <div className="text-2xl font-bold text-green-600 dark:text-green-400">{formatCurrency(currentMonthTotal)}</div>
            <div className="inline-flex items-center text-sm">
              <TrendingDown className="mr-1 h-4 w-4 text-red-500" />
              <span className="text-red-500">{percentageChange}%</span>
              <span className="ml-1 text-muted-foreground">
                vs {formatCurrency(lastMonthTotal)} in {comparisonDateRange}
              </span>
            </div>
          </div>

          <div className="h-[150px] mt-4">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={spendingCategories}
                  dataKey="value"
                  nameKey="name"
                  startAngle={180}
                  endAngle={0}
                  cx="50%"
                  cy="100%"
                  outerRadius={120}
                  innerRadius={90}
                >
                  {spendingCategories.map((entry) => (
                    <Cell
                      key={entry.name}
                      fill={GRAPH_CATEGORY_COLORS[entry.name]}
                      strokeWidth={0}
                    />
                  ))}
                </Pie>
              </PieChart>
            </ResponsiveContainer>
          </div>

          <div className="grid grid-cols-2 gap-4 pt-4">
            {spendingCategories.map((item) => (
              <div key={item.name} className="flex items-center">
                <div
                  className="h-3 w-3 rounded-full mr-2"
                  style={{ backgroundColor: GRAPH_CATEGORY_COLORS[item.name] }}
                />
                <div className="flex flex-col">
                  <span className="text-sm">{GRAPH_CATEGORY_DISPLAY_NAMES[item.name] ?? item.name}</span>
                  <span className="text-sm text-muted-foreground">
                    {formatCurrency(item.value)}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
