// Icons
import { ArrowUpRight, TrendingDown } from 'lucide-react';

// Shadcn UI Components / Charts
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Cell, Pie, PieChart, ResponsiveContainer } from 'recharts';

// Types
import type { SpendingCategoryBFFResponse } from '@/types/domain/spend-tracking';
import { MediaCategory } from '@/types/domain/spend-tracking';

// TODO: Sync all this hardcoded data with mockdata file
const mediaTypeColors: Record<MediaCategory, string> = {
  [MediaCategory.HARDWARE]: "#3FB950",
  [MediaCategory.DLC]: "#F85149",
  [MediaCategory.IN_GAME_PURCHASE]: "#FFA657",
  [MediaCategory.SUBSCRIPTION]: "#A371F7",
  [MediaCategory.PHYSICAL]: "#2F81F7",
  [MediaCategory.DISC]: "#DB61A2",
};

// TODO: Normalize Card colors + styling post design a/b test
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
    <Card className="col-span-full lg:col-span-2 bg-[#0B0F13] border-[#1D2127] text-white">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-xl font-bold text-[#c0c4cf]">Net this month</CardTitle>
        <button className="inline-flex items-center text-xs text-[#8B949E] hover:text-white">
          DOWNLOAD REPORT
          <ArrowUpRight className="ml-1 h-3 w-3" />
        </button>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          <div>
            <div className="text-2xl font-bold text-[#3FB950]">{formatCurrency(currentMonthTotal)}</div>
            <div className="inline-flex items-center text-sm">
              <TrendingDown className="mr-1 h-4 w-4 text-red-500" />
              <span className="text-red-500">{percentageChange}%</span>
              <span className="ml-1 text-[#8B949E]">
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
                      fill={mediaTypeColors[entry.name as MediaCategory]}
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
                  style={{ backgroundColor: mediaTypeColors[item.name as MediaCategory] }}
                />
                <div className="flex flex-col">
                  <span className="text-sm">{item.name}</span>
                  <span className="text-sm text-[#8B949E]">
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
