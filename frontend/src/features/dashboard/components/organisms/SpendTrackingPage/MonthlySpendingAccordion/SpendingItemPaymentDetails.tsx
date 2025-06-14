import { memo } from "react";

// Type safety improvements
interface SpendingItemPaymentDetailsProps {
  amount: number;
  date: string;
  isSubscription?: boolean;
}

export const SpendingItemPaymentDetails = memo(function SpendingItemPaymentDetails({
  amount,
  date,
  isSubscription = false,
}: SpendingItemPaymentDetailsProps) {
  return (
    <div className="text-right">
      <div className="text-sm text-gray-400">{isSubscription ? "NEXT PAYMENT" : "PURCHASE DATE"}</div>
      <div className="text-2xl font-bold">${amount}</div>
      {isSubscription && (
        <div className="text-sm text-gray-400">
          <span>Due <span className="font-bold text-white ml-1">{date}</span></span>
        </div>
      )}
    </div>
  );
});