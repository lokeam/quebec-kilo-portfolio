import { memo, useCallback, useState } from 'react';

// Custom Components
import { Badge } from "@/shared/components/ui/badge";
import { DigitalLocationIcon } from '@/features/dashboard/lib/utils/getDigitalLocationIcon';

// Shadcn Components
import { Card } from "@/shared/components/ui/card"
import { Button } from '@/shared/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogTitle,
  DialogHeader,
  DialogDescription
} from '@/shared/components/ui/dialog';

// Icons
import { IconCalendarDollar, IconEdit, IconTrash } from '@tabler/icons-react';

// Utils
import {
  isPaidService,
  formatServicePrice,
  isRenewalMonth,
} from '@/features/dashboard/lib/utils/online-service-status';

// Types
import type { DigitalLocation } from '@/types/domain/digital-location';

// Constants
//import { SERVICE_STATUS_CODES } from '@/shared/constants/service.constants';
import { cn } from '@/shared/components/ui/utils';

interface SingleOnlineServiceCardProps {
  service: DigitalLocation;
  onDelete?: (id: string) => void;
  onEdit?: (service: DigitalLocation) => void;
  isWatchedByResizeObserver?: boolean;
}

export const SingleOnlineServiceCard = memo(({
  service,
  onDelete,
  onEdit,
  isWatchedByResizeObserver
}: SingleOnlineServiceCardProps) => {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const showRenewalBadge = service.isActive &&
    isPaidService(service) &&
    isRenewalMonth(service);

  const handleEditService = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent card onClick from firing
    onEdit?.(service);
  }, [service, onEdit]);

  const handleDeleteClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent card onClick from firing
    setDeleteDialogOpen(true);
    // Reset error state when opening the dialog
    setDeleteError(null);
  }, []);

  const handleConfirmDelete = useCallback(() => {
    if (!service.id || !onDelete) return;

    setIsDeleting(true);
    setDeleteError(null);

    try {
      // Call the actual delete function from props
      onDelete(service.id);

      // The dialog will be closed after successful deletion
      // We're not calling setIsDeleting(false) here because
      // we want the button to stay in loading state until the
      // dialog is closed after success, which happens through the mutation
      // Close the dialog after a short delay to allow the mutation to process
      // The mutation's onSuccess will handle the actual cleanup
      setTimeout(() => {
        setDeleteDialogOpen(false);
        setIsDeleting(false);
      }, 100);
    } catch (err) {
      // This catch block is for synchronous errors
      // Most errors will be caught by the mutation's onError
      setIsDeleting(false);
      setDeleteError("Something went wrong. We can't complete this operation now, please try again later.");
      console.error("Error deleting service:", err);
    }
  }, [service.id, onDelete]);

  return (
    <>
      <Card
        className={cn(
          "flex relative cursor-pointer group w-full min-h-[100px] max-h-[100px] p-4",
          "transition-all duration-200",
          "hover:ring-1 hover:ring-ring hover:ring-inset",
          "hover:shadow-lg",
          isWatchedByResizeObserver && 'data-card-sentinel'
        )}
      >
        <div className="flex items-center justify-between min-w-0 w-full">
          <div className="flex items-center gap-3 min-w-0">
            <div className="w-10 h-10 shrink-0 text-foreground flex items-center justify-center">
              <DigitalLocationIcon name={service.name} className="w-full h-full" />
            </div>
            <div className="flex flex-col">
              <span
                className="font-medium text-sm text-foreground truncate overflow-hidden"
                style={{
                  maxWidth: 'var(--label-max-width)',
                  display: 'block',
                }}
              >
                {service.name}
              </span>
              {isPaidService(service) && service.billingCycle && (
                <div className="flex flex-col">
                  <span className="text-xs text-muted-foreground capitalize">
                    {`${service.billingCycle} subscription`}
                  </span>
                </div>
              )}
            </div>
          </div>

          <div className="relative w-24">
            {/* Show on no hover */}
            <div className="absolute inset-0 flex justify-end items-center gap-1 text-sm shrink-0 transition-opacity duration-200 group-hover:opacity-0 group-hover:pointer-events-none">
              {!isPaidService(service) ? (
                  <span className="font-medium text-foreground">
                    --
                  </span>
                ) : (
                <>
                  <span className="font-medium text-foreground">
                    {formatServicePrice(service.monthlyCost)}
                  </span>
                  <span className="text-muted-foreground text-xs">/ 1 mo</span>
                </>
              )}
              {/* {service.isActive && (
                <Power className="h-5 w-5 ml-1 text-emerald-500" />
              )} */}
              {showRenewalBadge && (
                <Badge
                  variant="destructive"
                  className="ml-1 absolute top-3 right-2"
                >
                  <IconCalendarDollar className="h-5 w-5 ml-1" />
                  <span className="ml-1">Renews this month</span>
                </Badge>
              )}
            </div>

            {/* Show on hover */}
            <div className={cn(
              "flex items-center gap-2 transition-opacity duration-200 opacity-0 group-hover:opacity-100 invisible group-hover:visible",
            )}>
              <Button
                variant="outline"
                size="sm"
                className="h-10 w-10 p-0"
                onClick={handleEditService}
              >
                <IconEdit className="h-5 w-5" />
                <span className="sr-only">Edit {service.name}</span>
              </Button>
              <Button
                variant="outline"
                size="sm"
                className="h-10 w-10 p-0"
                onClick={handleDeleteClick}
              >
                <IconTrash className="h-5 w-5" />
                <span className="sr-only">Delete {service.name}</span>
              </Button>
            </div>
          </div>
        </div>
      </Card>

      {/* Delete confirmation dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete {service.name}</DialogTitle>
            <DialogDescription>
                              {deleteError ? (
                  <div className="text-destructive">
                    {deleteError}
                  </div>
              ) : (
                "Are you sure you want to delete this service? This action cannot be undone."
              )}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
              disabled={isDeleting && !deleteError}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleConfirmDelete}
              disabled={isDeleting && !deleteError}
            >
              {isDeleting && !deleteError ? (
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

SingleOnlineServiceCard.displayName = 'SingleOnlineServiceCard';