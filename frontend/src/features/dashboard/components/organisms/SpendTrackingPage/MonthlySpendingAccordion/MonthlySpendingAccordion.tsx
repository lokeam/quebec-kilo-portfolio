import { useState } from 'react';

// Components
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';
import { MonthlySpendingItemDetails } from './MonthlySpendingItemDetails';
import { MemoizedMonthlySpendingAccordionItem } from '@/features/dashboard/components/organisms/SpendTrackingPage/MonthlySpendingAccordion/MonthlySpendingAccordionItem';

// Shadcn UI Components
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger
} from '@/shared/components/ui/accordion';
import { Separator } from '@/shared/components/ui/separator';

// Types
import type { SubscriptionSpend } from '@/features/dashboard/lib/types/spend-tracking/subscription';
import type { OneTimeSpend } from '@/features/dashboard/lib/types/spend-tracking/purchases';
import type { YearlySpending } from '@/features/dashboard/lib/types/spend-tracking/base';


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
          title={""}
          description={""}
          triggerTextAdd=""  // Empty because we're controlling open state externally
        >
          <MonthlySpendingItemDetails
            item={selectedItem}
            oneTimeTotal={oneTimeTotal}
          />
        </DrawerContainer>
      )}
    </>
  )
}
