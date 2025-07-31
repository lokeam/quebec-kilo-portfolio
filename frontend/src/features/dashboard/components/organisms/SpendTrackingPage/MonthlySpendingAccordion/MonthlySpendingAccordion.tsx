import { useState } from 'react';

// Components
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer';
import { MonthlySpendingItemDetails } from '@/features/dashboard/components/organisms/SpendTrackingPage/MonthlySpendingAccordion/MonthlySpendingItemDetails';
import { MemoizedMonthlySpendingAccordionItem } from '@/features/dashboard/components/organisms/SpendTrackingPage/MonthlySpendingAccordion/MonthlySpendingAccordionItem';
import { SpendTrackingForm } from '@/features/dashboard/components/organisms/SpendTrackingPage/SpendTrackingForm/SpendTrackingForm';

// Shadcn UI Components
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger
} from '@/shared/components/ui/accordion';
import { Separator } from '@/shared/components/ui/separator';

// Utils
import { transformSpendingItemResponseToFormData } from '@/features/dashboard/lib/utils/spendTrackingTransformers';

// Types
import type { SpendItem, SingleYearlyTotalBFFResponse } from '@/types/domain/spend-tracking';

interface MonthlySpendingAccordionProps {
  thisMonth: SpendItem[];
  future: SpendItem[];
  oneTimeTotal: SingleYearlyTotalBFFResponse[];
}



export function MonthlySpendingAccordion({ thisMonth, future, oneTimeTotal }: MonthlySpendingAccordionProps) {
  /* Open / Close Details Drawer */
  const [selectedItem, setSelectedItem] = useState<SpendItem | null>(null);
  const [isDrawerOpen, setIsDrawerOpen] = useState<boolean>(false);

  /* Editing an item */
  const [isEditFormOpen, setIsEditFormOpen] = useState<boolean>(false);
  const [editingItem, setEditingItem] = useState<SpendItem | null>(null);

  const handleItemClick = (item: SpendItem) => {
    setSelectedItem(item);
    setIsDrawerOpen(true);
  };

  const handleEditClick = (item: SpendItem) => {
    setEditingItem(item);
    setIsEditFormOpen(true);
    setIsDrawerOpen(false); // Close the details drawer
  };

  /**
   * Closes the details drawer when an item is deleted
   */
  const handleDeleteClick = () => {
    setIsDrawerOpen(false);
    setSelectedItem(null);
  };

  // Helper to pick correct date field
  const getSortDate = (item: SpendItem): number => {
    return item.spendTransactionType === 'subscription'
    ? item.nextBillingDate ?? 0
    : item.purchaseDate ?? 0;
  }

  const sortedThisMonth = [...thisMonth].sort((a, b) => getSortDate(a) - getSortDate(b));
  const sortedFuture = [...future].sort((a, b) => getSortDate(a) - getSortDate(b));

  return (
    <>
      <Accordion type="multiple" defaultValue={["this-month"]} className="col-span-full flex flex-col h-full space-y-4">
        <AccordionItem value="this-month" className="border-none">
          <AccordionTrigger className="flex gap-2 text-xl font-semibold text-foreground hover:no-underline justify-start">
            This month
          </AccordionTrigger>
          <AccordionContent className="pt-4">
            <div className="space-y-1">
              {sortedThisMonth.map((item, index) => (
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
          <AccordionTrigger className="flex gap-2 text-xl font-semibold text-foreground hover:no-underline justify-start">
            Next month
          </AccordionTrigger>
          <AccordionContent className="pt-4">
            <div className="space-y-1">
              {sortedFuture.map((item, index) => (
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

      {editingItem && (
        <DrawerContainer
          open={isEditFormOpen}
          onOpenChange={setIsEditFormOpen}
          title={`Edit ${editingItem.title}`}
          description="Update your expense details"
        >
          <SpendTrackingForm
            isEditing={true}
            spendTrackingData={transformSpendingItemResponseToFormData(editingItem)}
            onClose={() => {
              setIsEditFormOpen(false);
              setEditingItem(null);
            }}
            onSuccess={() => {
              setIsEditFormOpen(false);
              setEditingItem(null);
              // Refresh data will be handled by cache invalidation
            }}
          />
        </DrawerContainer>
      )}

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
            onEdit={handleEditClick}
            onDelete={handleDeleteClick}
          />
        </DrawerContainer>
      )}
    </>
  );
}
