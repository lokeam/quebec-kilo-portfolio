import * as React from 'react'
import { ChevronsUpDown } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu/DropdownMenu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/shared/components/ui/sidebar'

export function DomainSwitcher({
  domains,
}: {
  domains: {
    name: string
    logo: React.ElementType
    description: string
  }[]
}) {
  const { isMobile } = useSidebar()
  const [activeDomain, setActiveDomain] = React.useState(domains[0])

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size='lg'
              className='data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground'
            >
              {/* SVG Domain Logo */}
              <div className='flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground'>
                <activeDomain.logo className='size-4' />
              </div>

              {/* Domain Name and description */}
              <div className='grid flex-1 text-left text-sm leading-tight'>
                <span className='truncate font-semibold'>
                  {activeDomain.name}
                </span>
                <span className='truncate text-xs'>{activeDomain.description}</span>
              </div>
              <ChevronsUpDown className='ml-auto' />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className='w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg'
            align='start'
            side={isMobile ? 'bottom' : 'right'}
            sideOffset={4}
          >
            <DropdownMenuLabel className='text-xs text-muted-foreground'>
              Domains
            </DropdownMenuLabel>
            {domains.map((Domain, index) => (
              <DropdownMenuItem
                key={Domain.name}
                onClick={() => setActiveDomain(Domain)}
                className='gap-2 p-2'
              >
                <div className='flex size-6 items-center justify-center rounded-sm border'>
                  <Domain.logo className='size-4 shrink-0' />
                </div>
                {Domain.name}
                <DropdownMenuShortcut>⌘{index + 1}</DropdownMenuShortcut>
              </DropdownMenuItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}