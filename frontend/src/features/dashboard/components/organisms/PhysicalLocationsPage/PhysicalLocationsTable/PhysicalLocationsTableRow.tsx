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
import type { SublocationRowData } from '@/core/api/adapters/analytics.adapter';

// Utils
import { cn } from '@/shared/components/ui/utils';
import { PhysicalLocationIcon } from '@/features/dashboard/lib/utils/getPhysicalLocationIcon';
import { SublocationIcon } from '@/features/dashboard/lib/utils/getSublocationIcon';

interface PhysicalLocationsTableRowComponentProps {
  sublocation: SublocationRowData;
  index: number;
  isSelected?: boolean;
  onSelectionChange?: (checked: boolean) => void;
  onEdit?: (sublocation: SublocationRowData) => void;
}

function PhysicalLocationsTableRowComponent({
  sublocation,
  isSelected = false,
  onSelectionChange,
  onEdit
}: PhysicalLocationsTableRowComponentProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const handleEditService = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent row onClick from firing
    onEdit?.(sublocation);
  }, [sublocation, onEdit]);

  const handleDeleteService = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent row onClick from firing
    setDeleteDialogOpen(true);
  }, []);

  const handleCheckboxChange = useCallback((checked: boolean) => {
    onSelectionChange?.(checked);
  }, [onSelectionChange]);

  const handleConfirmDelete = useCallback(() => {
    // TODO: Implement delete functionality
    setDeleteDialogOpen(false);
  }, []);

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
            <SublocationIcon type={sublocation.sublocationType} bgColor={sublocation.bgColor} />
            <span>{sublocation.sublocationName}</span>
          </div>
        </TableCell>
        <TableCell>
          <div className="flex items-center gap-2">
            <PhysicalLocationIcon type={sublocation.parentLocationType} bgColor={sublocation.bgColor} />
            <span>{sublocation.parentLocationName}</span>
          </div>
        </TableCell>
        <TableCell>
          {sublocation.mapCoordinates ? (
            <a
              href={sublocation.mapCoordinates}
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
            <DialogTitle>Delete Sublocation</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {sublocation.sublocationName}? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleConfirmDelete}
            >
              Delete
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
