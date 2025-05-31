import { memo, useCallback, useState } from 'react';

// Custom Components
import SVGLogo from "@/shared/components/ui/LogoMap/LogoMap";
import type { LogoName } from "@/shared/components/ui/LogoMap/LogoMap";

// Shadcn Components
import { Button } from '@/shared/components/ui/button';
import { TableCell, TableRow } from "@/shared/components/ui/table";
import { Checkbox } from "@/shared/components/ui/checkbox";
import { Switch } from "@/shared/components/ui/switch";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogTitle,
  DialogHeader,
  DialogDescription
} from '@/shared/components/ui/dialog';
//import { toast } from 'sonner';

// Icons
import { Monitor } from 'lucide-react';
import { PaymentIcon } from 'react-svg-credit-card-payment-icons/dist';
import { IconEdit, IconTrash } from '@tabler/icons-react';

// Types
//import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import type { DigitalLocation } from '@/types/domain/online-service';

// Hooks
import { useOnlineServicesToggleActive, useOnlineServicesIsActive } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { formatCurrency, isPaidService} from '@/features/dashboard/lib/utils/online-service-status';

// Utils
import { cn } from '@/shared/components/ui/utils';
import { validatePaymentMethod } from '@/shared/constants/payment';
import { showToast } from '@/shared/components/ui/TanstackMutationToast/showToast';
import { DigitalLocationIcon } from '@/features/dashboard/lib/utils/getDigitalLocationIcon';


// Helper function to get the correct cost based on the billing cycle
// const getCycleBasedCost = (service: DigitalLocation): string => {
//   // Early return for free services
//   if (!isPaidService(service)) return "FREE";

//   const { billingCycle, costPerCycle } = service;

//   // Safely extract numeric values from string costs
//   const parseCost = (costPerCycle: string | number): number => {
//     if (!costPerCycle) return 0;

//     const numericValue = parseFloat(costPerCycle.toString().replace(/[^0-9.]/g, ''));
//     return isNaN(numericValue) ? 0 : numericValue;
//   };

//   // Use a simple mapping for direct returns
//   const costMap: Record<string, string> = {
//     "1 year": fees.annual,
//     "3 months": fees.quarterly,
//     "1 month": fees.monthly
//   };

//   // Handle the special calculated case for bi-annual
//   if (billingCycle === "6 months") {
//     const monthlyValue = parseCost(fees.monthly);
//     return monthlyValue > 0
//       ? formatCurrency((monthlyValue * 6).toFixed(2))
//       : fees.monthly; // Fallback to the raw value if parsing failed
//   }

//   // Return from mapping or fallback to monthly as default
//   return costMap[cycle] || fees.monthly;
// };

interface OnlineServicesTableRowProps {
  service: DigitalLocation;
  index: number
  isSelected?: boolean;
  onSelectionChange?: (checked: boolean) => void;
  onDelete?: (id: string) => void;
  onEdit?: (service: DigitalLocation) => void;
};

const createToggleActiveOnlineServiceToast = (label: string, isActive: boolean) => {
    showToast({
      message: `Recorded ${label} as ${isActive ? 'active' : 'inactive'}`,
      variant: 'success',
      duration: 2500,
    });
};

function OnlineServicesTableRowComponent({
  service,
  isSelected = false,
  onSelectionChange,
  onDelete,
  onEdit
}: OnlineServicesTableRowProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const paymentDate = service.nextPaymentDate
    ? new Date(service.nextPaymentDate).toLocaleDateString('en-CA') // This will format as YYYY-MM-DD
    : '--';

  // Handlers for online service activation
  const toggleActiveOnlineService = useOnlineServicesToggleActive();
  const isActive = useOnlineServicesIsActive(service.name);

  const handleToggleActiveOnlineService = useCallback((isChecked: boolean) => {
    toggleActiveOnlineService(service.id, isChecked);
    createToggleActiveOnlineServiceToast(service.name, isChecked);
  }, [service.id, service.name, toggleActiveOnlineService]);

  const handleCheckboxChange = useCallback((checked: boolean) => {
    onSelectionChange?.(checked);
  }, [onSelectionChange]);

  const handleEditService = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent row onClick from firing
    onEdit?.(service);
  }, [service, onEdit]);

  const handleDeleteService = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent row onClick from firing
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
      // Note: We're not calling setIsDeleting(false) here because
      // we want the button to stay in loading state until the
      // dialog is closed after success, which happens through the mutation
    } catch (err) {
      // This catch block is for synchronous errors
      // Most errors will be caught by the mutation's onError
      setIsDeleting(false);
      setDeleteError("Something went wrong. We can't complete this operation now, please try again later.");
      console.error("Error deleting service:", err);
    }
  }, [service.id, onDelete]);

  console.log('OnlineServicesTableRow', service);

  return (
    <>
      <TableRow
        className={cn(
          "h-[72px] relative group hover:border-",
          isSelected && "bg-muted/50",
          "transition-all duration-200",
          "hover:ring-1 hover:ring-white/20 hover:ring-inset",
          "hover:shadow-[0_0_4px_0_rgba(95,99,104,0.6),0_0_6px_2px_rgba(95,99,104,0.6)]",
        )}
      >
        <TableCell>
          <Checkbox
            checked={isSelected}
            onCheckedChange={handleCheckboxChange}
          />
        </TableCell>
        <TableCell>
          <div className="flex items-center gap-3">
            <div className="h-12 w-12 rounded-lg bg-black flex items-center justify-center">
              <DigitalLocationIcon name={service.logo} className="w-full h-full" />
            </div>
            <div className="flex flex-col">
              <span className="font-medium">{service.name}</span>
              <span className="text-sm text-muted-foreground">{service.locationType}</span>
            </div>
          </div>
        </TableCell>
        <TableCell>
          <Switch
            checked={service.isActive}
            onCheckedChange={handleToggleActiveOnlineService}
          />
        </TableCell>
        <TableCell>
          {
            isPaidService(service) ? (
              <span>{service.billingCycle}</span>
            ) : (
              <span>--</span>
            )
          }
        </TableCell>
        <TableCell>
          {
            isPaidService(service) ? (
              <span>{service.costPerCycle.toString()}</span>
            ) : (
              <span>--</span>
            )
          }
        </TableCell>
        <TableCell>
          <PaymentIcon
            type={validatePaymentMethod(
              typeof service.paymentMethod === 'string'
                ? service.paymentMethod
                : service.paymentMethod
            )}
            format="flatRounded"
          />
        </TableCell>
        <TableCell>
          {
            isPaidService(service) ? (
              <span>{paymentDate}</span>
            ) : (
              <span>--</span>
            )
          }
        </TableCell>
        {/* Edit + Delete buttons shown on hover */}
        <TableCell>
          <div className={cn(
            "flex items-center gap-2 transition-opacity duration-200",
            "opacity-0 group-hover:opacity-100",
            isSelected && "opacity-100"
          )}>
            <Button
              variant="outline"
              size="sm"
              className="h-10 w-10 p-0"
              onClick={handleEditService}
            >
              <IconEdit className="h-16 w-16" />
              <span className="sr-only">Edit {service.name}</span>
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="h-10 w-10 p-0 text-red-500 hover:text-red-600 hover:bg-red-100"
              onClick={handleDeleteService}
            >
              <IconTrash className="h-16 w-16" />
              <span className="sr-only">Delete {service.name}</span>
            </Button>
          </div>
        </TableCell>
      </TableRow>

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete {service.name}</DialogTitle>
            <DialogDescription>
              {deleteError ? (
                <div className="text-red-500">
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
}

export const OnlineServicesTableRow = memo(
  OnlineServicesTableRowComponent,
  (prevProps, nextProps) => {
    return (
      prevProps.index === nextProps.index &&
      prevProps.service.name === nextProps.service.name &&
      prevProps.service.billingCycle === nextProps.service.billingCycle &&
      prevProps.service.monthlyCost === nextProps.service.monthlyCost &&
      prevProps.service.paymentMethod === nextProps.service.paymentMethod &&
      prevProps.service.nextPaymentDate === nextProps.service.nextPaymentDate &&
      prevProps.service.isActive === nextProps.service.isActive &&
      prevProps.isSelected === nextProps.isSelected
    );
  }
);
