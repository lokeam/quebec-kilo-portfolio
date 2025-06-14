import { useState } from 'react';

// Components
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';
import { MonthlySpendingItemDetails } from './MonthlySpendingItemDetails';
import { MemoizedMonthlySpendingAccordionItem } from './MonthlySpendingAccordionItem';

// Shadcn UI Components
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger
} from '@/shared/components/ui/accordion';
import { Separator } from '@/shared/components/ui/separator';

// Local Type Definitions
interface BaseSpendItem {
  id: string;
  title: string;
  amount: number;
  spendTransactionType: 'subscription' | 'one-time';
  paymentMethod: string;
  mediaType: string;
  serviceName?: {
    id: string;
    displayName: string;
  };
  createdAt: number;
  updatedAt: number;
  isActive: boolean;
}

interface SubscriptionSpend extends BaseSpendItem {
  spendTransactionType: 'subscription';
  billingCycle: string;
  nextBillingDate: number;
  yearlySpending: Array<{
    year: number;
    amount: number;
  }>;
}

interface OneTimeSpend extends BaseSpendItem {
  spendTransactionType: 'one-time';
  isDigital: boolean;
  isWishlisted: boolean;
  purchaseDate: number;
}

interface YearlySpending {
  year: number;
  amount: number;
}

interface MonthlySpendingAccordionProps {
  thisMonth: (SubscriptionSpend | OneTimeSpend)[];
  future: (SubscriptionSpend | OneTimeSpend)[];
  oneTimeTotal: YearlySpending[];
}

export function MonthlySpendingAccordion({ thisMonth, future, oneTimeTotal }: MonthlySpendingAccordionProps) {
  const [selectedItem, setSelectedItem] = useState<SubscriptionSpend | OneTimeSpend | null>(null);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  const handleItemClick = (item: SubscriptionSpend | OneTimeSpend) => {
    setSelectedItem(item);
    setIsDrawerOpen(true);
  };

  return (
    <>
      <Accordion type="multiple" defaultValue={["this-month"]} className="col-span-full flex flex-col h-full space-y-4">
        <AccordionItem value="this-month" className="border-none">
          <AccordionTrigger className="flex gap-2 text-xl font-semibold text-slate-200 hover:no-underline justify-start">
            This month
          </AccordionTrigger>
          <AccordionContent className="pt-4">
            <div className="space-y-1">
              {thisMonth.map((item, index) => (
                <MemoizedMonthlySpendingAccordionItem
                  key={index}
                  item={item}
                  onClick={() => handleItemClick(item)}
                />
              ))}
            </div>
          </AccordionContent>
        </AccordionItem>

        <Separator />

        <AccordionItem value="future" className="border-none">
          <AccordionTrigger className="flex gap-2 text-xl font-semibold text-slate-200 hover:no-underline justify-start">
            Next month
          </AccordionTrigger>
          <AccordionContent className="pt-4">
            <div className="space-y-1">
              {future.map((item, index) => (
                <MemoizedMonthlySpendingAccordionItem
                  key={index}
                  item={item}
                  onClick={() => handleItemClick(item)}
                />
              ))}
            </div>
          </AccordionContent>
        </AccordionItem>
      </Accordion>

      {selectedItem && (
        <DrawerContainer
          open={isDrawerOpen}
          onOpenChange={setIsDrawerOpen}
          title={selectedItem.title}
          description=""
        >
          <MonthlySpendingItemDetails
            item={selectedItem}
            oneTimeTotal={oneTimeTotal}
          />
        </DrawerContainer>
      )}
    </>
  );
}
