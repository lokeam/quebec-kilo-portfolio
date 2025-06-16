import { useState } from 'react';

// Components
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';
import { MonthlySpendingItemDetails } from '@/features/dashboard/components/organisms/SpendTrackingPage/MonthlySpendingAccordion/MonthlySpendingItemDetails';
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
import type { SpendItem, SingleYearlyTotalBFFResponse } from '@/types/domain/spend-tracking';

interface MonthlySpendingAccordionProps {
  thisMonth: SpendItem[];
  future: SpendItem[];
  oneTimeTotal: SingleYearlyTotalBFFResponse[];
}

export function MonthlySpendingAccordion({ thisMonth, future, oneTimeTotal }: MonthlySpendingAccordionProps) {
  const [selectedItem, setSelectedItem] = useState<SpendItem | null>(null);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  const handleItemClick = (item: SpendItem) => {
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
