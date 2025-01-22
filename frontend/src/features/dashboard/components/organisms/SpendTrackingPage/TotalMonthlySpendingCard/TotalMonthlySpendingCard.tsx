// Icons
import { ArrowUpRight, TrendingDown } from 'lucide-react';

// Shadcn UI Components / Charts
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Cell, Pie, PieChart, ResponsiveContainer } from 'recharts';


// TODO: Sync all this hardcoded data with mockdata file
const spendingData = [
  { name: 'Hardware', value: 534.04 },
  { name: 'Dlc', value: 267.02 },
  { name: 'InGamePurchase', value: 178.01 },
  { name: 'Subscription', value: 356.02 },
  { name: 'Physical', value: 267.02 },
  { name: 'Disc', value: 178.01 },
];

const mediaTypeColors: Record<string, string> = {
  Hardware: "#3FB950",
  Dlc: "#F85149",
  InGamePurchase: "#FFA657",
  Subscription: "#A371F7",
  Physical: "#2F81F7",
  Disc: "#DB61A2",
};

// TODO: Normalize Card colors + styling post design a/b test
export function TotalMonthlySpendingCard() {
  const netThisMonth = "$1,784.04";
  const percentChange = "-20.91";
  const netLastMonth = "$2,255.92";

  return (
    <Card className="col-span-full lg:col-span-2 w-full max-w-2xl bg-[#0B0F13] border-[#1D2127] text-white">
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
            <div className="text-2xl font-bold text-[#3FB950]">{netThisMonth}</div>
            <div className="inline-flex items-center text-sm">
              <TrendingDown className="mr-1 h-4 w-4 text-red-500" />
              <span className="text-red-500">{percentChange}%</span>
              <span className="ml-1 text-[#8B949E]">
                vs {netLastMonth} in Dec 1 - Dec 22, 2024
              </span>
            </div>
          </div>

          <div className="h-[150px] mt-4">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={spendingData}
                  dataKey="value"
                  nameKey="name"
                  startAngle={180}
                  endAngle={0}
                  cx="50%"
                  cy="100%"
                  outerRadius={120}
                  innerRadius={90}
                >
                  {spendingData.map((entry) => (
                    <Cell
                      key={entry.name}
                      fill={mediaTypeColors[entry.name]}
                      strokeWidth={0}
                    />
                  ))}
                </Pie>
              </PieChart>
            </ResponsiveContainer>
          </div>

          <div className="grid grid-cols-2 gap-4 pt-4">
            {spendingData.map((item) => (
              <div key={item.name} className="flex items-center">
                <div
                  className="h-3 w-3 rounded-full mr-2"
                  style={{ backgroundColor: mediaTypeColors[item.name] }}
                />
                <div className="flex flex-col">
                  <span className="text-sm">{item.name}</span>
                  <span className="text-sm text-[#8B949E]">
                    ${item.value.toFixed(2)}
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
