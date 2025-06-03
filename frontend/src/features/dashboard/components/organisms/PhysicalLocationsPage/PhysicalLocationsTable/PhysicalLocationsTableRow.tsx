import { memo, useCallback, useState } from 'react';

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

// Icons
import { PaymentIcon } from 'react-svg-credit-card-payment-icons/dist';
import { IconEdit, IconTrash } from '@tabler/icons-react';

// Types
import type { DigitalLocation } from '@/types/domain/online-service';

// Hooks
import { useOnlineServicesToggleActive } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { isPaidService} from '@/features/dashboard/lib/utils/online-service-status';
import { useDeleteDigitalLocation } from '@/core/api/queries/digitalLocation.queries';

// Utils
import { cn } from '@/shared/components/ui/utils';
import { showToast } from '@/shared/components/ui/TanstackMutationToast/showToast';
import { DigitalLocationIcon } from '@/features/dashboard/lib/utils/getDigitalLocationIcon';

type PaymentMethodType = "Alipay" | "Amex" | "Code" | "CodeFront" | "Diners" | "Discover" | "Elo" | "Generic" | "Hiper" | "Hipercard" | "Jcb" | "Maestro" | "Mastercard" | "Mir" | "Paypal" | "Unionpay" | "Visa";

const isValidPaymentMethod = (method: string): method is PaymentMethodType => {
  const validMethods = [
    "Alipay", "Amex", "Code", "CodeFront", "Diners", "Discover", "Elo",
    "Generic", "Hiper", "Hipercard", "Jcb", "Maestro", "Mastercard",
    "Mir", "Paypal", "Unionpay", "Visa"
  ];
  return validMethods.includes(method);
};

interface PhysicalLocationsTableRowComponentProps {
  service: DigitalLocation;
  index: number
  isSelected?: boolean;
  onSelectionChange?: (checked: boolean) => void;
  onEdit?: (service: DigitalLocation) => void;
};

const createToggleActiveOnlineServiceToast = (label: string, isActive: boolean) => {
    showToast({
      message: `Recorded ${label} as ${isActive ? 'active' : 'inactive'}`,
      variant: 'success',
      duration: 2500,
    });
};

function PhysicalLocationsTableRowComponent({
  service,
  isSelected = false,
  onSelectionChange,
  onEdit
}: PhysicalLocationsTableRowComponentProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const deleteDigitalLocation = useDeleteDigitalLocation();

  const paymentDate = service.nextPaymentDate
    ? new Date(service.nextPaymentDate).toLocaleDateString('en-CA')
    : '--';

  // Handlers for online service activation
  const toggleActiveOnlineService = useOnlineServicesToggleActive();

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
  }, []);

  const handleConfirmDelete = useCallback(() => {
    deleteDigitalLocation.mutate(service.id, {
      onSuccess: () => {
        setDeleteDialogOpen(false);
      }
    });
  }, [service.id, deleteDigitalLocation]);

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
              <DigitalLocationIcon name={service.name} className="w-full h-full" />
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
          {isValidPaymentMethod(service.paymentMethod) ? (
            <PaymentIcon
              type={service.paymentMethod}
              format="flatRounded"
            />
          ) : (
            <PaymentIcon
              type="Generic"
              format="flatRounded"
            />
          )}
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

        {/* Edit and Delete Buttons */}
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
              variant="destructive"
              size="sm"
              className="h-10 w-10 p-0"
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
            <DialogTitle>Delete Physical Location</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {service.name}? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
              disabled={deleteDigitalLocation.isPending}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleConfirmDelete}
              disabled={deleteDigitalLocation.isPending}
            >
              {deleteDigitalLocation.isPending ? (
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
  PhysicalLocationsTableRowComponent,
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
