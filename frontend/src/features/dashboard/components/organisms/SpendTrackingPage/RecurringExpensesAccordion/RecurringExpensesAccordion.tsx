
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/shared/components/ui/accordion';
import { RecurringExpensesItem } from '@/features/dashboard/components/organisms/SpendTrackingPage/RecurringExpensesAccordion/RecurringExpensesItem';
import type { SpendTrackingService } from '@/features/dashboard/lib/types/service.types';

interface RecurringExpensesAccordionProps {
  thisMonth: SpendTrackingService[]
  future: SpendTrackingService[]
}

export function RecurringExpensesAccordion({ thisMonth, future }: RecurringExpensesAccordionProps) {

  return (
    <Accordion type="multiple" defaultValue={["this-month"]} className="col-span-full lg:col-span-2 flex flex-col h-full space-y-4 ">
      <AccordionItem value="this-month" className="border-none">
        <AccordionTrigger className="flex gap-2 text-xl font-semibold text-slate-200 hover:no-underline justify-start">
          {/* <ChevronDown className="h-6 w-6 shrink-0 text-slate-400 transition-transform duration-200" /> */}
          This month
        </AccordionTrigger>
        <AccordionContent className="pt-4">
          <div className="space-y-1">
            {thisMonth.map((item, index) => (
              <RecurringExpensesItem key={index} {...item} />
            ))}
          </div>
        </AccordionContent>
      </AccordionItem>

      <AccordionItem value="future" className="border-none">
        <AccordionTrigger className="flex gap-2 text-xl font-semibold text-slate-200 hover:no-underline justify-start">
          {/* <ChevronDown className="h-6 w-6 shrink-0 text-slate-400 transition-transform duration-200" /> */}
          Next month
        </AccordionTrigger>
        <AccordionContent className="pt-4">
          <div className="space-y-1">
            {future.map((item, index) => (
              <RecurringExpensesItem
                key={index}
                {...item}
                isWatchedByResizeObserver={true}
              />
            ))}
          </div>
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  )
}
