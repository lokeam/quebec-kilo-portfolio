"use client"

import React from 'react'
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/components/ui/table"
import { Checkbox } from "@/shared/components/ui/checkbox"
import { OnlineServicesTableRow } from './OnlineServicesTableRow'

interface OnlineService {
  name: string
  label: string
  logo: string
  tier: string
  billingCycle: string
  currency: string
  price: string
}

interface OnlineServicesTableProps {
  services: OnlineService[]
}

const TableHeaderRow: React.FC = () => (
  <TableRow>
    <TableHead className="w-[50px]">
      <Checkbox />
    </TableHead>
    <TableHead>Service</TableHead>
    <TableHead>Category</TableHead>
    <TableHead>Active</TableHead>
    <TableHead>Tier</TableHead>
    <TableHead>Billing Cycle</TableHead>
    <TableHead>Amount</TableHead>
    <TableHead>Payment Method</TableHead>
  </TableRow>
)

export function OnlineServicesTable({ services }: OnlineServicesTableProps) {
  return (
    <div className="w-full">
      <Table>
        <TableHeader>
          <TableHeaderRow />
        </TableHeader>
        <TableBody>
          {services.map((service, index) => (
            <OnlineServicesTableRow key={service.name} service={service} index={index} />
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
