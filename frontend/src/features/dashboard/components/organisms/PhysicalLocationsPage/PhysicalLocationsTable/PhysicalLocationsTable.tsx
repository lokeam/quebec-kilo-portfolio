import { memo, useState } from 'react'

// ShadCN Components
import { Button } from '@/shared/components/ui/button'
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/components/ui/table"
import { Checkbox } from "@/shared/components/ui/checkbox"

// Custom Components
import { PhysicalLocationsTableRow } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/PhysicalLocationsTable/PhysicalLocationsTableRow'

import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogTitle,
  DialogHeader,
  DialogDescription
} from '@/shared/components/ui/dialog';

// Types
import type { LocationsBFFSublocationResponse, LocationsBFFPhysicalLocationResponse } from '@/types/domain/physical-location'

interface PhysicalLocationsTableProps {
  sublocationRows: LocationsBFFSublocationResponse[]
  physicalLocationRows: LocationsBFFPhysicalLocationResponse[]
  onEdit?: (location: LocationsBFFPhysicalLocationResponse | LocationsBFFSublocationResponse) => void
  onDelete?: (id: string) => void
}

export const PhysicalLocationsTable = memo(({
  sublocationRows,
  physicalLocationRows,
  onEdit,
  onDelete
}: PhysicalLocationsTableProps) => {
  const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set());
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const handleSelectionChange = (id: string, checked: boolean) => {
    setSelectedRows(prev => {
      const next = new Set(prev);
      if (checked) {
        next.add(id);
      } else {
        next.delete(id);
      }
      return next;
    });
  };

  const handleDeleteSelectedRows = () => {
    if (selectedRows.size === 0) return;
    setDeleteDialogOpen(true);
  };

  const handleConfirmDelete = () => {
    // TODO: Implement delete functionality
    setSelectedRows(new Set());
    setDeleteDialogOpen(false);
  };

  // Calculate if all rows are selected
  const allSelected = (sublocationRows.length + physicalLocationRows.length) > 0 &&
    selectedRows.size === (sublocationRows.length + physicalLocationRows.length);

  return (
    <div className="w-full">
      <Table className="rounded-xl border">
        <TableHeader>
          <TableRow>
            <TableHead className="w-[50px]">
              <Checkbox
                checked={allSelected}
                onCheckedChange={(checked) => {
                  if (checked) {
                    const allIds = [
                      ...sublocationRows.map(row => row.sublocationId),
                      ...physicalLocationRows.map(row => row.physicalLocationId)
                    ];
                    setSelectedRows(new Set(allIds));
                  } else {
                    setSelectedRows(new Set());
                  }
                }}
                aria-label="Select all"
              />
            </TableHead>
            <TableHead>Sublocation</TableHead>
            <TableHead>Parent Location</TableHead>
            <TableHead>Map Coordinates</TableHead>
            <TableHead>Stored Items</TableHead>
            <TableHead className="w-[100px]">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {/* Render sublocation rows */}
          {sublocationRows.map((sublocation, index) => (
            <PhysicalLocationsTableRow
              key={sublocation.sublocationId}
              row={sublocation}
              index={index}
              isSelected={selectedRows.has(sublocation.sublocationId)}
              onSelectionChange={(checked) => handleSelectionChange(sublocation.sublocationId, checked)}
              onEdit={onEdit}
              onDelete={onDelete}
            />
          ))}

          {/* Render physical location rows */}
          {physicalLocationRows.map((location, index) => (
            <PhysicalLocationsTableRow
              key={location.physicalLocationId}
              row={location}
              index={sublocationRows.length + index}
              isSelected={selectedRows.has(location.physicalLocationId)}
              onSelectionChange={(checked) => handleSelectionChange(location.physicalLocationId, checked)}
              onEdit={onEdit}
              onDelete={onDelete}
            />
          ))}
        </TableBody>
      </Table>

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Selected Locations</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {selectedRows.size} selected location{selectedRows.size > 1 ? 's' : ''}? This action cannot be undone.
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
    </div>
  );
});
