import React from 'react'
import { TableCell, TableRow } from "@/shared/components/ui/table"
import { Checkbox } from "@/shared/components/ui/checkbox"
import { Switch } from "@/shared/components/ui/switch"
import { Monitor } from 'lucide-react'

interface OnlineService {
  name: string
  label: string
  logo: string
  tier: string
  billingCycle: string
  currency: string
  price: string
}

interface OnlineServicesTableRowProps {
  service: OnlineService
  index: number
}

const getStatusStyle = (index: number) => {
  const styles = [
    'bg-red-100 text-red-800',
    'bg-green-100 text-green-800',
    'bg-yellow-100 text-yellow-800'
  ]
  return styles[index % 3]
}

const getStatusText = (index: number) => {
  const statuses = ['Inactive', 'Publish', 'Scheduled']
  return statuses[index % 3]
}

const generateSKU = (name: string) => {
  return name.toUpperCase().slice(0, 3) + Math.floor(Math.random() * 90000 + 10000)
}

export function OnlineServicesTableRow({ service, index }: OnlineServicesTableRowProps) {
  return (
    <TableRow className="h-[72px]">
      <TableCell>
        <Checkbox />
      </TableCell>
      <TableCell>
        <div className="flex items-center gap-3">
          <div className="h-12 w-12 rounded-lg bg-slate-100 flex items-center justify-center">
            <img
              src={`/placeholder.svg?height=48&width=48`}
              alt={service.label}
              className="h-8 w-8"
            />
          </div>
          <div className="flex flex-col">
            <span className="font-medium">{service.label}</span>
            <span className="text-sm text-muted-foreground">
              {service.tier || "Standard subscription"}
            </span>
          </div>
        </div>
      </TableCell>
      <TableCell>
        <div className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-full bg-purple-100 flex items-center justify-center">
            <Monitor className="h-4 w-4 text-purple-500" />
          </div>
          <span>Gaming Services</span>
        </div>
      </TableCell>
      <TableCell>
        <Switch />
      </TableCell>
      <TableCell>{generateSKU(service.name)}</TableCell>
      <TableCell>{service.price}</TableCell>
      <TableCell>{Math.floor(Math.random() * 900 + 100)}</TableCell>
      <TableCell>
        <div className={`
          inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
          ${getStatusStyle(index)}
        `}>
          {getStatusText(index)}
        </div>
      </TableCell>
    </TableRow>
  )
}

