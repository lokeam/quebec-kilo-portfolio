import React, { useState, useMemo } from 'react'

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

// Icons
import { IconTrash } from '@tabler/icons-react'

// Types
import type { PhysicalLocation } from '@/types/domain/physical-location'
import type { SublocationRowData } from './PhysicalLocationsTableRow'

interface PhysicalLocationsTableProps {
  services: PhysicalLocation[]
  onEdit?: (sublocation: SublocationRowData) => void
}

export function PhysicalLocationsTable({ services, onEdit }: PhysicalLocationsTableProps) {
  const [selectedRows, setSelectedRows] = useState<string[]>([]);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  // Flatten the physical locations into sublocation rows
  const sublocationRows = useMemo(() => {
    const rows = services.flatMap(location => {
      return (location.sublocations || []).map(sublocation => {
        const row = {
          sublocationId: sublocation.id,
          sublocationName: sublocation.name,
          sublocationType: sublocation.type,
          parentLocationId: location.id,
          parentLocationName: location.name,
          parentLocationType: location.locationType,
          mapCoordinates: location.mapCoordinates?.googleMapsLink || '',
          bgColor: location.bgColor,
          storedItems: sublocation.metadata?.notes ? parseInt(sublocation.metadata.notes) : 0
        };

        return row;
      });
    });

    return rows;
  }, [services]);

  // Handle select all checkbox
  const handleSelectAll = (checked: boolean) => {
    if (checked) {
      setSelectedRows(sublocationRows.map(row => row.sublocationId));
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
    // TODO: Implement delete functionality
    setSelectedRows([]);
    setDeleteDialogOpen(false);
  };

  // Calculate if all rows are selected
  const allSelected = sublocationRows.length > 0 && selectedRows.length === sublocationRows.length;

  console.log('Physical Locations Table, sublocationRows: ', sublocationRows);

  return (
    <div className="w-full">
      <Table className="rounded-xl border">
        <TableHeader>
          <TableRow>
            <TableHead>
              <Checkbox
                checked={allSelected}
                onCheckedChange={handleSelectAll}
                aria-label="Select all"
              />
            </TableHead>
            <TableHead>Sublocation Name</TableHead>
            <TableHead>Parent Location</TableHead>
            <TableHead>Google Maps Link</TableHead>
            <TableHead>Number of Games Stored</TableHead>
            <TableHead className="w-[100px]">
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
        </TableHeader>
        <TableBody>
          {sublocationRows.map((sublocation, index) => (
            <PhysicalLocationsTableRow
              key={sublocation.sublocationId}
              sublocation={sublocation}
              index={index}
              isSelected={selectedRows.includes(sublocation.sublocationId)}
              onSelectionChange={(checked) => handleRowSelection(sublocation.sublocationId, checked)}
              onEdit={onEdit}
            />
          ))}
        </TableBody>
      </Table>

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Selected Sublocations</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {selectedRows.length} selected sublocation{selectedRows.length > 1 ? 's' : ''}? This action cannot be undone.
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
}
