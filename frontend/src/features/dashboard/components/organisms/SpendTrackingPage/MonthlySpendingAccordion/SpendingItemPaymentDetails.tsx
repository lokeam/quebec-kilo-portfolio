import { memo } from "react";
import { formatCurrency } from "@/features/dashboard/lib/utils/formatCurrency";


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
      <div className="text-sm text-muted-foreground">{isSubscription ? "NEXT PAYMENT" : "AMOUNT"}</div>
      <div className="text-2xl font-bold">{formatCurrency(amount)}</div>
        <div className="text-sm text-muted-foreground">
          <span>{isSubscription ? "Due" : "Payment Date"}<span className="font-bold text-foreground ml-1">{date}</span></span>
        </div>
    </div>
  );
});