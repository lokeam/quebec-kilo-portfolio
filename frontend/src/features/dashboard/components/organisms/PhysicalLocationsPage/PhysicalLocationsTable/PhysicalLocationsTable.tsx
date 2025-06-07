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
import type { LocationsBFFSublocationResponse } from '@/types/domain/physical-location'

interface PhysicalLocationsTableProps {
  sublocationRows: LocationsBFFSublocationResponse[]
  onEdit?: (sublocation: LocationsBFFSublocationResponse) => void
  onDelete?: (id: string) => void
}

export const PhysicalLocationsTable = memo(({
  sublocationRows,
  onEdit,
  onDelete
}: PhysicalLocationsTableProps) => {
  const [selectedRows, setSelectedRows] = useState<Set<string>>(new Set());
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const handleSelectionChange = (sublocationId: string, checked: boolean) => {
    setSelectedRows(prev => {
      const next = new Set(prev);
      if (checked) {
        next.add(sublocationId);
      } else {
        next.delete(sublocationId);
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
  const allSelected = sublocationRows.length > 0 && selectedRows.size === sublocationRows.length;

  console.log('Physical Locations Table, sublocationRows: ', sublocationRows);

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
                    setSelectedRows(new Set(sublocationRows.map(row => row.sublocationId)));
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
          {sublocationRows.map((sublocation, index) => (
            <PhysicalLocationsTableRow
              key={sublocation.sublocationId}
              sublocation={sublocation}
              index={index}
              isSelected={selectedRows.has(sublocation.sublocationId)}
              onSelectionChange={(checked) => handleSelectionChange(sublocation.sublocationId, checked)}
              onEdit={onEdit}
              onDelete={onDelete}
            />
          ))}
        </TableBody>
      </Table>

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Selected Sublocations</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {selectedRows.size} selected sublocation{selectedRows.size > 1 ? 's' : ''}? This action cannot be undone.
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
