import { memo } from 'react';

// Shadcn UI Components
import { Badge } from '@/shared/components/ui/badge';
import { Card, CardContent, CardHeader } from '@/shared/components/ui/card';

// Hooks + Utils
import { LogoOrIcon } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/LogoOrIcon';

// Icons
import { CreditCard } from 'lucide-react';

// Types
import type { SpendTrackingService } from "@/features/dashboard/lib/types/service.types";

interface MonthlySpendingItemDetailsProps {
  item: SpendTrackingService;
};

export const MonthlySpendingItemDetails = memo(function MonthlySpendingItemDetails({
  item : {
    name,
    day,
    month,
    // year // <--- TODO: Swap this out with array of years
    title,
    amount,
    billingCycle,
    mediaType = 'subscription',
    paymentMethod,
  }
}: MonthlySpendingItemDetailsProps) {

  const calculateYearlyAmount = (amount: string, cycle: string): number => {
    const baseAmount = Number.parseFloat(amount)
    switch (cycle) {
      case "1 year":
        return baseAmount
      case "3 months":
        return baseAmount * 4
      case "1 month":
        return baseAmount * 12
      default:
        return baseAmount
    }
  };

  const yearlyAmount = calculateYearlyAmount(amount ?? '0', billingCycle ?? '1 month');


  return (
    <Card className="bg-[#0A0A0A] text-white border-none mb-4">
      <CardHeader className="space-y-1.5">
        <div className="flex flex-row items-center justify-between space-y-4">

          <div className="flex flex-col">
            <div className="flex flex-row gap-4">
              <Badge
                variant="outline"
                className="bg-purple-900/50 text-purple-300 border-purple-700 w-auto"
                data-testid="media-type-badge"
              >
                {mediaType.charAt(0).toUpperCase() + mediaType.slice(1)}
              </Badge>
            </div>

            <div className="h-14 w-14 flex items-center justify-center my-2">
              <LogoOrIcon name={name} mediaType={mediaType} />
            </div>
            <h2 className="text-xl font-semibold">{title}</h2>
          </div>

          <div className="text-right">
          <div className="text-sm text-gray-400">NEXT PAYMENT</div>
            <div className="text-2xl font-bold">${amount}</div>
            <div className="text-sm text-gray-400">
              <span>Due <span className="font-bold text-white ml-1">{month} {day}</span></span>
            </div>
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-8">
        <div>
          <h3 className="text-lg font-semibold mb-4">Subscription details</h3>
          <div className="flex xs:flex-col flex-row gap-4">
            <Badge variant="outline" className="bg-blue-900/50 text-blue-300 border-blue-700">
              {billingCycle}
            </Badge>
            <Badge variant="outline" className="bg-green-900/50 text-green-300 border-green-700">
              ${amount}
            </Badge>
          </div>
        </div>

        <div>
          <h3 className="text-lg font-semibold mb-4">Total spent per year</h3>
          <div className="space-y-4 sm:space-y-0 sm:grid sm:grid-cols-3 sm:gap-4">
            <div data-spending-year-col1 className="text-gray-400">2023</div>
            <div data-spending-year-col2 className="text-gray-400">2024</div>
            <div data-spending-year-col3 className="text-gray-400">2025</div>
            <div data-spending-year-col1 className="font-semibold">${yearlyAmount.toFixed(2)}</div>
            <div data-spending-year-col2 className="font-semibold">${(yearlyAmount * 1.1).toFixed(2)}</div>
            <div data-spending-year-col3 className="font-semibold">${(yearlyAmount * 1.2).toFixed(2)}</div>
          </div>
        </div>

        <div>
          <h3 className="text-lg font-semibold mb-4">Payment method</h3>
          <div className="flex items-center space-x-4">
            <div className="w-12 h-12 bg-gray-800 rounded-lg flex items-center justify-center">
              <CreditCard className="w-6 h-6 text-gray-400" />
            </div>
            <div>
              <div className="font-semibold">{paymentMethod}</div>
              <div className="text-sm text-gray-400">${amount}</div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
});
