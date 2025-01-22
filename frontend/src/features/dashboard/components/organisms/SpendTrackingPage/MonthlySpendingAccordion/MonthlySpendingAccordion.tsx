// Shadcn UI Components
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger
} from '@/shared/components/ui/accordion';
import { Separator } from '@/shared/components/ui/separator';

// Components
import { MemoizedMonthlySpendingAccordionItem } from '@/features/dashboard/components/organisms/SpendTrackingPage/MonthlySpendingAccordion/MonthlySpendingAccordionItem';

// Types
import type { SpendTrackingService } from '@/features/dashboard/lib/types/service.types';
import { useState } from 'react';
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';
import { MonthlySpendingItemDetails } from './MonthlySpendingItemDetails';

interface MonthlySpendingAccordionProps {
  thisMonth: SpendTrackingService[];
  future: SpendTrackingService[];
}

export function MonthlySpendingAccordion({ thisMonth, future }: MonthlySpendingAccordionProps) {
  const [selectedItem, setSelectedItem] = useState<SpendTrackingService | null>(null);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  const handleItemClick = (item: SpendTrackingService) => {
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
                  {...item}
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
                  {...item}
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
          triggerText=""  // Empty because we're controlling open state externally
        >
          <MonthlySpendingItemDetails item={selectedItem} />
        </DrawerContainer>
      )}
    </>
  )
}
