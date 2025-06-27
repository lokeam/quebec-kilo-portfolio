import { memo } from 'react';

// Components
import { MemoizedDashboardBadge } from '@/features/dashboard/components/molecules/DashboardBadge/DashboardBadge';
import { SpendingItemYearGrid } from './SpendingItemYearGrid';
import { SpendingItemPaymentDetails } from './SpendingItemPaymentDetails';

// Shadcn UI Components
import { Card, CardContent, CardHeader } from '@/shared/components/ui/card';

// Hooks + Utils
import { useFormattedDate } from '@/features/dashboard/lib/hooks/useFormattedDate';
import { useSpendingData } from '@/features/dashboard/lib/hooks/useSpendingData';
import { MediaIcon } from '@/features/dashboard/lib/utils/getMediaIcon';
import { DigitalLocationIcon } from '@/features/dashboard/lib/utils/getDigitalLocationIcon';

// Types
import type { SpendingItemBFFResponse, SingleYearlyTotalBFFResponse } from '@/types/domain/spend-tracking';
import { MediaCategory } from '@/types/domain/spend-tracking';

// Icons
import { PaymentIcon } from 'react-svg-credit-card-payment-icons/dist';
//import { CreditCard } from 'lucide-react';

interface MonthlySpendingItemDetailsProps {
  item: SpendingItemBFFResponse;
  oneTimeTotal: SingleYearlyTotalBFFResponse[];
}

type PaymentMethodType = "Alipay" | "Amex" | "Code" | "CodeFront" | "Diners" | "Discover" | "Elo" | "Generic" | "Hiper" | "Hipercard" | "Jcb" | "Maestro" | "Mastercard" | "Mir" | "Paypal" | "Unionpay" | "Visa";

export const MonthlySpendingItemDetails = memo(function MonthlySpendingItemDetails({
  item,
  oneTimeTotal,
}: MonthlySpendingItemDetailsProps) {
  const { spendingData, title, isSubscription } = useSpendingData(item, oneTimeTotal);

  const dateDisplay = useFormattedDate(
    item.spendTransactionType,
    item.nextBillingDate,
    item.purchaseDate
  );

  const renderIcon = () => {
    // For subscriptions, use the digital location icon
    if (item.mediaType === MediaCategory.SUBSCRIPTION) {
      return (
        <DigitalLocationIcon
          name={item.provider}
          className="h-14 w-14"
        />
      );
    }

    // For other media types, use the media icon
    return (
      <MediaIcon
        mediaType={item.mediaType}
        className="h-14 w-14"
      />
    );
  };

  return (
    <Card className="bg-[#0A0A0A] text-white border-none mb-4">
      <CardHeader className="space-y-1.5">
        <div className="flex flex-row items-center justify-between space-y-4">
          <div className="flex flex-col">
            <div className="flex flex-row gap-4">
              <MemoizedDashboardBadge
                variant="outline"
                className="bg-purple-900/50 text-purple-300 border-purple-700 w-auto"
                data-testid="media-type-badge"
              >
                {item.mediaType}
              </MemoizedDashboardBadge>
            </div>

            {/* Provider Logo / Item Icon*/}
            <div className="h-14 w-14 flex items-center justify-center my-2">
              {renderIcon()}
            </div>
            <h2 className="text-xl font-semibold">{item.title}</h2>
          </div>

          <SpendingItemPaymentDetails
            amount={item.amount}
            date={dateDisplay}
            isSubscription={isSubscription}
          />
        </div>
      </CardHeader>
      <CardContent className="space-y-8">
        {isSubscription && item.billingCycle && (
          <div>
            <h3 className="text-lg font-semibold mb-4">Subscription details</h3>
            <div className="flex xs:flex-col flex-row gap-4">
              <MemoizedDashboardBadge
                variant="outline"
                className="bg-blue-900/50 text-blue-300 border-blue-700"
              >
                {item.billingCycle}
              </MemoizedDashboardBadge>
              <MemoizedDashboardBadge
                variant="outline"
                className="bg-green-900/50 text-green-300 border-green-700"
              >
                {new Intl.NumberFormat('en-US', {
                  style: 'currency',
                  currency: 'USD'
                }).format(Number(item.amount))}
              </MemoizedDashboardBadge>
            </div>
          </div>
        )}

        {/* Yearly Spending */}
        <div>
          <h3 className="text-lg font-semibold mb-4">{title}</h3>
          <SpendingItemYearGrid data={spendingData} />
        </div>

        {/* Payment Method */}
        <div>
          <h3 className="text-lg font-semibold mb-4">Payment method</h3>
          <div className="flex items-center space-x-4">
            <div className="w-12 h-12 bg-gray-800 rounded-lg flex items-center justify-center">
              {/* <CreditCard className="w-6 h-6 text-gray-400" /> */}
              <PaymentIcon
                type={item.paymentMethod as PaymentMethodType}
                format="flatRounded"
              />
            </div>
            <div>
              <div className="font-semibold">{item.paymentMethod}</div>
              <div className="text-sm text-gray-400">
                {new Intl.NumberFormat('en-US', {
                  style: 'currency',
                  currency: 'USD'
                }).format(Number(item.amount))}
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
});
