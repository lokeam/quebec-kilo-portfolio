import * as React from 'react'
import {
  SidebarMenu,
  SidebarMenuItem,
} from '@/shared/components/ui/sidebar'
import QkoLogo from '@/shared/components/ui/LogoMap/brand/qko_logo'

type QKOLoaderProps = {
  text?: string
  className?: string
  logoClassName?: string
}

export default function QKOLoader({
  logoClassName
}: QKOLoaderProps) {
  return (
    <div className='qko-loader flex aspect-square size-8 items-center justify-center rounded-lg text-sidebar-primary-foreground'>
      <QkoLogo className={`logo-svg size-8 ${logoClassName || ''}`} />
      <span className="sr-only">Loading...</span>
    </div>
  )
}