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

// Icons
import { IconTrash } from '@tabler/icons-react'

// Types
import type { DigitalLocation } from '@/types/domain/online-service'

interface OnlineServicesTableProps {
  services: DigitalLocation[]
  onDelete?: (id: string) => void
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

export function OnlineServicesTable({ services, onDelete, onEdit }: OnlineServicesTableProps) {
  const [selectedRows, setSelectedRows] = useState<string[]>([]);

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
    // TODO: Wire up delete logic
    console.log('Deleting rows: ', selectedRows);
    setSelectedRows([]);
  }

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
              onDelete={onDelete}
              onEdit={onEdit}
            />
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
