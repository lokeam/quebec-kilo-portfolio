import * as React from 'react'
import {
  SidebarMenu,
  SidebarMenuItem,
} from '@/shared/components/ui/sidebar'
import QkoLogo from '@/shared/components/ui/LogoMap/brand/qko_logo'

type BrandLogoProps = {
  text?: string
  className?: string
  logoClassName?: string
}

export function BrandLogo({
  className,
  logoClassName
}: BrandLogoProps) {
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <div
          className={`peer/menu-button flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left outline-none ring-sidebar-ring transition-[width,height,padding] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 group-has-[[data-sidebar=menu-action]]/menu-item:pr-8 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[active=true]:bg-sidebar-accent data-[active=true]:font-medium data-[active=true]:text-sidebar-accent-foreground data-[state=open]:hover:bg-sidebar-accent data-[state=open]:hover:text-sidebar-accent-foreground group-data-[collapsible=icon]:!size-8 [&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 h-12 text-sm group-data-[collapsible=icon]:!p-0 ${className || ''}`}
          data-sidebar="menu-button"
          data-size="lg"
        >
          {/* Logo Container */}
          <div className='flex aspect-square size-8 items-center justify-center rounded-lg text-sidebar-primary-foreground'>
            <QkoLogo className={`size-8 ${logoClassName || ''}`} />
          </div>

          {/* Text */}
          {/* <span className='truncate font-semibold'>
            {text}
          </span> */}
        </div>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}