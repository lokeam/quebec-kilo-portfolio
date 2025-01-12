import { Card } from "@/shared/components/ui/card";
import { XCircle } from 'lucide-react';
import type { OnlineService } from "@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata"

interface RecurringExpensesProps {
  services: OnlineService[]
}

export function RecurringExpensesCard({ services }: RecurringExpensesProps) {
  const calculateYearlyTotal = (services: OnlineService[]) => {
    return services.reduce((total, service) => {
      if (service.price === 'FREE' || service.price === 'NA') return total

      const price = parseFloat(service.price.replace('$', ''))
      if (isNaN(price)) return total

      // Convert all prices to yearly
      const yearlyPrice = service.billingCycle.includes('mo')
        ? price * 12
        : service.billingCycle.includes('yr')
          ? price
          : price

      return total + yearlyPrice
    }, 0)
  }

  const yearlyTotal = calculateYearlyTotal(services)

  return (
    <Card className="p-6 bg-black">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold">Recurring expenses</h2>
        <XCircle className="w-5 h-5 text-gray-400" />
      </div>
      <div className="flex items-baseline gap-1 overflow-hidden">
        <span className="text-3xl font-bold">$</span>
        <span className="text-4xl font-bold">{yearlyTotal.toFixed(2)}</span>
        <span className="text-gray-500 ml-2">USD</span>
        <span className="text-gray-500 ml-1">Yearly</span>
      </div>
    </Card>
  );
}

