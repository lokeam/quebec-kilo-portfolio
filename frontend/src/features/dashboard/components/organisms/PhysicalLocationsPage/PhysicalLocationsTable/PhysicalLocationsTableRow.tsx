import { memo, useCallback, useState } from 'react';

// Shadcn Components
import { Button } from '@/shared/components/ui/button';
import { Checkbox } from "@/shared/components/ui/checkbox";
import { TableCell, TableRow } from "@/shared/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogTitle,
  DialogHeader,
  DialogDescription
} from '@/shared/components/ui/dialog';

// Icons
import { IconEdit, IconTrash } from '@tabler/icons-react';

// Types
import type { LocationsBFFSublocationResponse } from '@/types/domain/physical-location';

// Utils
import { cn } from '@/shared/components/ui/utils';
import { PhysicalLocationIcon } from '@/features/dashboard/lib/utils/getPhysicalLocationIcon';
import { SublocationIcon } from '@/features/dashboard/lib/utils/getSublocationIcon';

interface PhysicalLocationsTableRowComponentProps {
  sublocation: LocationsBFFSublocationResponse;
  index: number;
  isSelected?: boolean;
  onSelectionChange?: (checked: boolean) => void;
  onEdit?: (sublocation: LocationsBFFSublocationResponse) => void;
  onDelete?: (id: string) => void;
}

function PhysicalLocationsTableRowComponent({
  sublocation,
  isSelected = false,
  onSelectionChange,
  onEdit,
  onDelete
}: PhysicalLocationsTableRowComponentProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const handleEditService = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent row onClick from firing
    onEdit?.(sublocation);
  }, [sublocation, onEdit]);

  const handleDeleteService = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent row onClick from firing
    setDeleteDialogOpen(true);
    setDeleteError(null);
  }, []);

  const handleCheckboxChange = useCallback((checked: boolean) => {
    onSelectionChange?.(checked);
  }, [onSelectionChange]);

  const handleConfirmDelete = useCallback(() => {
    if (!sublocation.sublocationId || !onDelete) return;

    setIsDeleting(true);
    setDeleteError(null);

    try {
      onDelete(sublocation.sublocationId);
    } catch (err) {
      setIsDeleting(false);
      setDeleteError("Something went wrong. We can't complete this operation now, please try again later.");
      console.error("Error deleting sublocation:", err);
    }
  }, [sublocation.sublocationId, onDelete]);

  return (
    <>
      <TableRow
        className={cn(
          "h-[80px] relative group hover:border-",
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
          <div className="flex items-center gap-2">
            <SublocationIcon type={sublocation.sublocationType} bgColor={sublocation.parentLocationBgColor || 'gray'} />
            <span>{sublocation.sublocationName}</span>
          </div>
        </TableCell>
        <TableCell>
          <div className="flex items-center gap-2">
            <PhysicalLocationIcon type={sublocation.parentLocationType} bgColor={sublocation.parentLocationBgColor || 'gray'} />
            <span>{sublocation.parentLocationName}</span>
          </div>
        </TableCell>
        <TableCell>
          {sublocation.mapCoordinates ? (
            <a
              href={sublocation.mapCoordinates.googleMapsLink}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-500 hover:underline"
            >
              View on Google Maps
            </a>
          ) : (
            <span className="text-gray-500">No coordinates</span>
          )}
        </TableCell>
        <TableCell>
          {sublocation.storedItems} items
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
              <span className="sr-only">Edit {sublocation.sublocationName}</span>
            </Button>
            <Button
              variant="destructive"
              size="sm"
              className="h-10 w-10 p-0"
              onClick={handleDeleteService}
            >
              <IconTrash className="h-16 w-16" />
              <span className="sr-only">Delete {sublocation.sublocationName}</span>
            </Button>
          </div>
        </TableCell>
      </TableRow>

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete {sublocation.sublocationName}</DialogTitle>
            <DialogDescription asChild>
              {deleteError ? (
                <div className="text-red-500">
                  {deleteError}
                </div>
              ) : (
                <div className="space-y-4">
                  <p>Are you sure you want to delete this sublocation?</p>
                  <p>This action cannot be undone.</p>
                </div>
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

export const PhysicalLocationsTableRow = memo(
  PhysicalLocationsTableRowComponent,
  (prevProps, nextProps) => {
    return (
      prevProps.index === nextProps.index &&
      prevProps.sublocation.sublocationId === nextProps.sublocation.sublocationId &&
      prevProps.sublocation.sublocationName === nextProps.sublocation.sublocationName &&
      prevProps.sublocation.storedItems === nextProps.sublocation.storedItems &&
      prevProps.isSelected === nextProps.isSelected
    );
  }
);
