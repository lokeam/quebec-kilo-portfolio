import { memo, useState, useCallback } from 'react';

// Components
import { MemoizedDashboardBadge } from '@/features/dashboard/components/molecules/DashboardBadge/DashboardBadge';
import { SpendingItemYearGrid } from './SpendingItemYearGrid';
import { SpendingItemPaymentDetails } from './SpendingItemPaymentDetails';

// Shadcn UI Components
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader } from '@/shared/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogTitle,
  DialogHeader,
  DialogDescription
} from '@/shared/components/ui/dialog';

// Hooks + Utils
import { useFormattedDate } from '@/features/dashboard/lib/hooks/useFormattedDate';
import { useDisplayAnnualSpendingData } from '@/features/dashboard/lib/hooks/useDisplayAnnualSpendingData';
import { normalizeOneTimePurchaseMediaType } from '@/features/dashboard/lib/utils/normalizeOneTimePurchaseMediaType';
import { formatCurrency } from '@/features/dashboard/lib/utils/formatCurrency';
import { useDeleteSpendItems } from '@/core/api/queries/spendTracking.queries';

// Types
import type { SpendingItemBFFResponse, SingleYearlyTotalBFFResponse } from '@/types/domain/spend-tracking';
import { MediaCategory } from '@/types/domain/spend-tracking';
import { TransactionType } from '@/types/domain/spend-tracking';

// Icons
import { PaymentIcon } from 'react-svg-credit-card-payment-icons/dist';
import { DigitalLocationIcon } from '@/features/dashboard/lib/utils/getDigitalLocationIcon';
import { MediaIcon } from '@/features/dashboard/lib/utils/getMediaIcon';
import { IconTrash } from '@tabler/icons-react';

/**
 * Props for the MonthlySpendingItemDetails component
 */
interface MonthlySpendingItemDetailsProps {
  item: SpendingItemBFFResponse;
  oneTimeTotal: SingleYearlyTotalBFFResponse[];
  onEdit?: (item: SpendingItemBFFResponse) => void;
  onDelete?: () => void;
}

type PaymentMethodType = "Alipay" | "Amex" | "Code" | "CodeFront" | "Diners" | "Discover" | "Elo" | "Generic" | "Hiper" | "Hipercard" | "Jcb" | "Maestro" | "Mastercard" | "Mir" | "Paypal" | "Unionpay" | "Visa";

/**
 * Displays detailed information about a spending item in a drawer
 *
 * Features:
 * - Shows item details, payment info, and yearly spending data
 * - Provides edit and delete functionality for one-time purchases
 * - Handles optimistic updates and error states
 * - Integrates with the spend tracking mutation system
 */
export const MonthlySpendingItemDetails = memo(function MonthlySpendingItemDetails({
  item,
  oneTimeTotal,
  onEdit,
  onDelete,
}: MonthlySpendingItemDetailsProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<boolean>(false);

  // Display Hooks
  const { spendingData, title, isSubscription } = useDisplayAnnualSpendingData(item, oneTimeTotal);

  // Mutation Hooks
  const deleteSpendItems = useDeleteSpendItems();


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

  /**
   * Opens the delete confirmation dialog
   */
  const handleDelete = useCallback((event: React.MouseEvent) => {
    event.stopPropagation();
    setDeleteDialogOpen(true);
  }, []);

  /**
   * Initiates the deletion process
   */
  const handleConfirmDelete = useCallback(() => {
    const itemIdString = item.id.toString();
    deleteSpendItems.mutate([itemIdString], {
      onSuccess: () => {
        setDeleteDialogOpen(false);
        onDelete?.();
      }
    });
  }, [deleteSpendItems, item.id, onDelete]);

  return (
    <>
      <Card className="mb-4">
        <CardHeader className="space-y-1.5">
          <div className="flex flex-row items-center justify-between space-y-4">
            <div className="flex flex-col">
              <div className="flex flex-row gap-4">
                <MemoizedDashboardBadge
                  variant="outline"
                  className="bg-accent text-accent-foreground border-border w-auto"
                  data-testid="media-type-badge"
                >
                  {normalizeOneTimePurchaseMediaType(item.mediaType)}
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
                  className="bg-accent text-accent-foreground border-border"
                >
                  {item.billingCycle}
                </MemoizedDashboardBadge>
                <MemoizedDashboardBadge
                  variant="outline"
                  className="bg-accent text-accent-foreground border-border"
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
              <div className="w-12 h-12 bg-muted rounded-lg flex items-center justify-center">
                <PaymentIcon
                  type={item.paymentMethod as PaymentMethodType}
                  format="flatRounded"
                />
              </div>
              <div>
                <div className="font-semibold">{item.paymentMethod}</div>
                <div className="text-sm text-muted-foreground">{formatCurrency(item.amount)}</div>
              </div>
            </div>
          </div>

          {/* Edit / Delete Buttons */}
          {item.spendTransactionType !== TransactionType.SUBSCRIPTION && (
            <div>
              <div className="flex-col w-full">
                {onEdit && (
                  <Button
                    variant="outline"
                    onClick={() => onEdit(item)}
                    className="w-full mb-4"
                  >
                    Edit Expense
                  </Button>
                )}
                {onDelete && (
                  <Button
                    variant="destructive"
                    onClick={handleDelete}
                    disabled={deleteSpendItems.isPending}
                    className="w-full"
                  >
                    {deleteSpendItems.isPending ? (
                      <>
                        <span className="animate-spin mr-2">⊚</span>
                        Deleting...
                      </>
                    ) : (
                      <>
                        <IconTrash className="h-4 w-4 mr-2" />
                        Delete Expense
                      </>
                    )}
                  </Button>
                )}
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Delete Dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Expense</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete "{item.title}"? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
              disabled={deleteSpendItems.isPending}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleConfirmDelete}
              disabled={deleteSpendItems.isPending}
            >
              {deleteSpendItems.isPending ? (
                <>
                  <span className="animate-spin mr-2">⊚</span>
                  Deleting...
                </>
              ) : (
                "Delete"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
});
