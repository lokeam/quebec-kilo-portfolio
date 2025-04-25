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
import { toast } from 'sonner';

// Icons
import { Monitor } from 'lucide-react';
import { PaymentIcon } from 'react-svg-credit-card-payment-icons/dist';
import { IconEdit, IconTrash } from '@tabler/icons-react';

// Types
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';

// Hooks
import { useOnlineServicesToggleActive, useOnlineServicesIsActive } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { formatCurrency, isServiceFree } from '@/features/dashboard/lib/utils/online-service-status';

// Utils
import { cn } from '@/shared/components/ui/utils';
import { validatePaymentMethod } from '@/shared/constants/payment';

// Helper function to get the correct cost based on the billing cycle
const getCycleBasedCost = (service: OnlineService): string => {
  // Early return for free services
  if (!service.billing || isServiceFree(service)) {
    return "FREE";
  }

  const { cycle, fees } = service.billing;

  // Safely extract numeric values from string costs
  const parseCost = (costString: string | undefined): number => {
    if (!costString) return 0;
    const numericValue = parseFloat(costString.replace(/[^0-9.]/g, ''));
    return isNaN(numericValue) ? 0 : numericValue;
  };

  // Use a simple mapping for direct returns
  const costMap: Record<string, string> = {
    "1 year": fees.annual,
    "3 months": fees.quarterly,
    "1 month": fees.monthly
  };

  // Handle the special calculated case for bi-annual
  if (cycle === "6 months") {
    const monthlyValue = parseCost(fees.monthly);
    return monthlyValue > 0
      ? formatCurrency((monthlyValue * 6).toFixed(2))
      : fees.monthly; // Fallback to the raw value if parsing failed
  }

  // Return from mapping or fallback to monthly as default
  return costMap[cycle] || fees.monthly;
};

interface OnlineServicesTableRowProps {
  service: OnlineService;
  index: number
  isSelected?: boolean;
  onSelectionChange?: (checked: boolean) => void;
  onDelete?: (id: string) => void;
  onEdit?: (service: OnlineService) => void;
};

const createToggleActiveOnlineServiceToast = (label: string, isActive: boolean) => {
  toast(`Recorded ${label} as ${isActive ? 'active' : 'inactive'}`, {
    className: 'bg-green-500 text-white',
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

  const hasValidLogo = Boolean(service.logo);
  const paymentDate = service.billing?.renewalDate ? `${service.billing.renewalDate.month} ${service.billing.renewalDate.day}` : '--';
  const isFree = isServiceFree({ billing: service.billing } as OnlineService);

  // Handlers for online service activation
  const toggleActiveOnlineService = useOnlineServicesToggleActive();
  const isActive = useOnlineServicesIsActive(service.name);

  const handleToggleActiveOnlineService = useCallback((isChecked: boolean) => {
    toggleActiveOnlineService(service.name, isChecked);
    createToggleActiveOnlineServiceToast(service.label, isChecked);
  }, [service.name, service.label, toggleActiveOnlineService]);

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
            {hasValidLogo ? (
                <SVGLogo
                  domain="games"
                  name={service.logo as LogoName<'games'>}
                  className="h-8 w-8"
                />
              ) : (
                <Monitor className="h-8 w-8 text-slate-500" />
              )}
            </div>
            <div className="flex flex-col">
              <span className="font-medium">{service.label}</span>
              <span className="text-sm text-muted-foreground">
                {service.tier?.currentTier || "Standard subscription"}
              </span>
            </div>
          </div>
        </TableCell>
        <TableCell>
          <Switch
            checked={isActive}
            onCheckedChange={handleToggleActiveOnlineService}
          />
        </TableCell>
        <TableCell>
          {
            isFree ? (
              <span>--</span>
            ) : (
              <span>{service.billing?.cycle}</span>
            )
          }
        </TableCell>
        <TableCell>
          {
            isFree ? (
              <span>--</span>
            ) : (
              <span>{getCycleBasedCost(service)}</span>
            )
          }
        </TableCell>
        <TableCell>
          <PaymentIcon
            type={validatePaymentMethod(
              typeof service.billing?.paymentMethod === 'string'
                ? service.billing.paymentMethod
                : service.billing?.paymentMethod?.id
            )}
            format="flatRounded"
          />
        </TableCell>
        <TableCell>
          {
            isFree ? (
              <span>--</span>
            ) : (
              <span>{paymentDate}</span>
            )
          }
        </TableCell>
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
              <span className="sr-only">Edit {service.label}</span>
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="h-10 w-10 p-0 text-red-500 hover:text-red-600 hover:bg-red-100"
              onClick={handleDeleteService}
            >
              <IconTrash className="h-16 w-16" />
              <span className="sr-only">Delete {service.label}</span>
            </Button>
          </div>
        </TableCell>
      </TableRow>

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete {service.label}</DialogTitle>
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
      prevProps.service.label === nextProps.service.label &&
      prevProps.service.logo === nextProps.service.logo &&
      prevProps.service.tier?.currentTier === nextProps.service.tier?.currentTier &&
      prevProps.service.billing?.cycle === nextProps.service.billing?.cycle &&
      prevProps.service.billing?.fees.monthly === nextProps.service.billing?.fees.monthly &&
      prevProps.service.billing?.paymentMethod === nextProps.service.billing?.paymentMethod &&
      prevProps.service.billing?.paymentMethod === nextProps.service.billing?.paymentMethod &&
      prevProps.service.status === nextProps.service.status &&
      prevProps.isSelected === nextProps.isSelected
    );
  }
);
