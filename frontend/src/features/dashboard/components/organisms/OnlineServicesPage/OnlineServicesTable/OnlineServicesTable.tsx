import React, { useState } from 'react'

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
import { OnlineServicesTableRow } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServicesTable/OnlineServicesTableRow'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogTitle,
  DialogHeader,
  DialogDescription
} from '@/shared/components/ui/dialog';

// Query hooks
import { useDeleteDigitalLocation } from '@/core/api/queries/digitalLocation.queries'

// Icons
import { IconTrash } from '@tabler/icons-react'

// Types
import type { DigitalLocation } from '@/types/domain/online-service'

interface OnlineServicesTableProps {
  services: DigitalLocation[]
  onEdit?: (service: DigitalLocation) => void
}

const TableHeaderRow: React.FC = () => (
  <TableRow>
    <TableHead className="w-[50px]" />
    <TableHead>Service</TableHead>
    <TableHead>Active</TableHead>
    <TableHead>Billing Cycle</TableHead>
    <TableHead>Amount</TableHead>
    <TableHead>Payment Method</TableHead>
    <TableHead>Renewal Date</TableHead>
  </TableRow>
)

export function OnlineServicesTable({ services, onEdit }: OnlineServicesTableProps) {
  const [selectedRows, setSelectedRows] = useState<string[]>([]);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const deleteDigitalLocation = useDeleteDigitalLocation();


  // Handle select all checkbox
  const handleSelectAll = (checked: boolean) => {
    if (checked) {
      setSelectedRows(services.map(service => service.id));
    } else {
      setSelectedRows([]);
    }
  }

  // Handle individual row selection
  const handleRowSelection = (id: string, checked: boolean) => {
    if (checked) {
      setSelectedRows(prev => [...prev, id]);
    } else {
      setSelectedRows(prev => prev.filter(rowId => rowId !== id));
    }
  }

  const handleDeleteSelectedRows = () => {
    if (selectedRows.length === 0) return;
    setDeleteDialogOpen(true);
  };

  const handleConfirmDelete = () => {
    deleteDigitalLocation.mutate(selectedRows, {
      onSuccess: () => {
        setSelectedRows([]);
        setDeleteDialogOpen(false);
      }
    });
  };

  // Calculate if all rows are selected
  const allSelected = services.length > 0 && selectedRows.length === services.length;

  return (
    <div className="w-full">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[50px]">
              <Checkbox
                checked={allSelected}
                onCheckedChange={handleSelectAll}
                aria-label="Select all"
              />
            </TableHead>
            <TableHead colSpan={6}>
              {selectedRows.length > 0 && (
                <Button
                  variant="destructive"
                  size="sm"
                  onClick={handleDeleteSelectedRows}
                  className="flex items-center gap-2"
                >
                  <IconTrash size={16} />
                  Delete Selected ({selectedRows.length})
                </Button>
              )}
            </TableHead>
          </TableRow>
          <TableHeaderRow />
        </TableHeader>
        <TableBody>
          {services.map((service, index) => (
            <OnlineServicesTableRow
              key={service.name}
              service={service}
              index={index}
              isSelected={selectedRows.includes(service.id)}
              onSelectionChange={(checked) => handleRowSelection(service.id, checked)}
              onEdit={onEdit}
            />
          ))}
        </TableBody>
      </Table>

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Selected Services</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {selectedRows.length} selected service{selectedRows.length > 1 ? 's' : ''}? This action cannot be undone.
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
    </div>
  );
}
