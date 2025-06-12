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
import type { LocationsBFFSublocationResponse, LocationsBFFPhysicalLocationResponse } from '@/types/domain/physical-location';
import type { LocationIconBgColor } from '@/types/domain/location-types';

// Utils
import { cn } from '@/shared/components/ui/utils';
import { PhysicalLocationIcon } from '@/features/dashboard/lib/utils/getPhysicalLocationIcon';
import { SublocationIcon } from '@/features/dashboard/lib/utils/getSublocationIcon';

// Components
import { SublocationDeleteWarning } from '../DeleteWarningContent/SublocationDeleteWarning';
import { PhysicalLocationDeleteWarning } from '../DeleteWarningContent/PhysicalLocationDeleteWarning';

// Type guard to determine if row is a physical location
const isPhysicalLocation = (
  row: LocationsBFFPhysicalLocationResponse | LocationsBFFSublocationResponse
): row is LocationsBFFPhysicalLocationResponse => {
  return 'physicalLocationId' in row;
};

interface PhysicalLocationsTableRowComponentProps {
  row: LocationsBFFPhysicalLocationResponse | LocationsBFFSublocationResponse;
  index: number;
  isSelected?: boolean;
  onSelectionChange?: (checked: boolean) => void;
  onEdit?: (location: LocationsBFFPhysicalLocationResponse | LocationsBFFSublocationResponse) => void;
  onDelete?: (id: string) => void;
  sublocations: LocationsBFFSublocationResponse[];
}

function PhysicalLocationsTableRowComponent({
  row,
  isSelected = false,
  onSelectionChange,
  onEdit,
  onDelete,
  sublocations
}: PhysicalLocationsTableRowComponentProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const handleEdit = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent row onClick from firing
    onEdit?.(row);
  }, [row, onEdit]);

  const handleDelete = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent row onClick from firing
    setDeleteDialogOpen(true);
    setDeleteError(null);
  }, []);

  const handleCheckboxChange = useCallback((checked: boolean) => {
    onSelectionChange?.(checked);
  }, [onSelectionChange]);

  const handleConfirmDelete = useCallback(() => {
    const id = isPhysicalLocation(row) ? row.physicalLocationId : row.sublocationId;
    if (!id || !onDelete) return;

    setIsDeleting(true);
    setDeleteError(null);

    try {
      onDelete(id);
    } catch (err) {
      setIsDeleting(false);
      setDeleteError("Something went wrong. We can't complete this operation now, please try again later.");
      console.error("Error deleting location:", err);
    }
  }, [row, onDelete]);

  const rowName = isPhysicalLocation(row) ? row.name : row.sublocationName;

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
            {isPhysicalLocation(row) ? (
              <>
                <SublocationIcon
                  type={row.physicalLocationType}
                  bgColor={(row.bgColor as LocationIconBgColor) || 'gray'}
                />
                <span className="text-gray-500">None - Please add sublocation</span>
              </>
            ) : (
              <>
                <SublocationIcon
                  type={row.sublocationType}
                  bgColor={row.parentLocationBgColor || 'gray'}
                />
                <span>{row.sublocationName}</span>
              </>
            )}
          </div>
        </TableCell>
        <TableCell>
          <div className="flex items-center gap-2">
            {isPhysicalLocation(row) ? (
              <>
                <PhysicalLocationIcon
                  type={row.physicalLocationType}
                  bgColor={(row.bgColor as LocationIconBgColor) || 'gray'}
                />
                <span>{row.name}</span>
              </>
            ) : (
              <>
                <PhysicalLocationIcon
                  type={row.parentLocationType}
                  bgColor={row.parentLocationBgColor || 'gray'}
                />
                <span>{row.parentLocationName}</span>
              </>
            )}
          </div>
        </TableCell>
        <TableCell>
          {isPhysicalLocation(row) ? (
            <span className="text-gray-500">No coordinates</span>
          ) : row.mapCoordinates ? (
            <a
              href={row.mapCoordinates.googleMapsLink}
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
          {isPhysicalLocation(row) ? '0 items' : `${row.storedItems} items`}
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
              onClick={handleEdit}
            >
              <IconEdit className="h-16 w-16" />
              <span className="sr-only">Edit {rowName}</span>
            </Button>
            <Button
              variant="destructive"
              size="sm"
              className="h-10 w-10 p-0"
              onClick={handleDelete}
            >
              <IconTrash className="h-16 w-16" />
              <span className="sr-only">Delete {rowName}</span>
            </Button>
          </div>
        </TableCell>
      </TableRow>

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete {rowName}</DialogTitle>
            <DialogDescription>
              {deleteError ? (
                <div className="text-red-500">
                  {deleteError}
                </div>
              ) : isPhysicalLocation(row) ? (
                <PhysicalLocationDeleteWarning
                  location={row}
                  associatedItems={sublocations}
                />
              ) : (
                <SublocationDeleteWarning location={row} />
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
      (isPhysicalLocation(prevProps.row) && isPhysicalLocation(nextProps.row)
        ? prevProps.row.physicalLocationId === nextProps.row.physicalLocationId
        : !isPhysicalLocation(prevProps.row) && !isPhysicalLocation(nextProps.row)
          ? prevProps.row.sublocationId === nextProps.row.sublocationId
          : false) &&
      prevProps.isSelected === nextProps.isSelected
    );
  }
);
