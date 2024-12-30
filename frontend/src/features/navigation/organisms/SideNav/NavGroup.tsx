import { type ReactNode } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { ChevronRight } from 'lucide-react'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/shared/components/ui/collapsible'
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  useSidebar,
} from '@/shared/components/ui/sidebar'
import { Badge } from '@/shared/components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu/DropdownMenu'
import type { NavCollapsible, NavItem, NavLink } from './sideNav.types'
import type { NavGroup } from './sideNav.types'

export function NavGroup({ title, items }: NavGroup) {
  const { state } = useSidebar();
  const location = useLocation();
  const href = location.pathname + location.search;
  return (
    <SidebarGroup>
      <SidebarGroupLabel>{title}</SidebarGroupLabel>
      <SidebarMenu>
        {items.map((item) => {
          const key = `${item.title}-${item.url}`

          if (!item.items) {
            return (
              <SidebarMenuLink key={key} item={item} href={href} />
            );
          }

          if (state === 'collapsed') {
            return (
              <SidebarMenuCollapsedDropdown key={key} item={item} href={href} />
            );
          }

          return (
            <SidebarMenuCollapsible key={key} item={item} href={href} />
          );
        })}
      </SidebarMenu>
    </SidebarGroup>
  )
}

const NavBadge = ({ children }: { children: ReactNode }) => (
  <Badge className='text-xs rounded-full px-1 py-0'>{children}</Badge>
);

const SidebarMenuLink = ({
  item,
  href
}: {
  item: NavLink;
  href: string;
}) => {
  const { setOpenMobile } = useSidebar();
  return (
    <SidebarMenuItem>
      {/* Sidebar Menu Button */}
      <SidebarMenuButton
        asChild
        isActive={checkIsActive(href, item)}
        tooltip={item.title}
      >
        {/* Sidebar Link */}
        <Link to={item.url} onClick={() => setOpenMobile(false)}>
          {/* Sidebar Icon */}
          {item.icon && <item.icon />}

          {/* Sidebar Title */}
          <span>{item.title}</span>

          {/* Sidebar Badge */}
          {item.badge && <NavBadge>{item.badge}</NavBadge>}
        </Link>
      </SidebarMenuButton>
    </SidebarMenuItem>
  );
}

const SidebarMenuCollapsible = ({
  item,
  href,
}: {
  item: NavCollapsible
  href: string
}) => {
  const { setOpenMobile } = useSidebar()
  return (
    <Collapsible
      asChild
      defaultOpen={checkIsActive(href, item, true)}
      className='group/collapsible'
    >
      {/* Sidebar Menu Item */}
      <SidebarMenuItem>
        {/* Collapsible Trigger */}
        <CollapsibleTrigger asChild>

          <SidebarMenuButton tooltip={item.title}>
            {/* Sidebar Icon */}
            {item.icon && <item.icon />}

            {/* Sidebar Title */}
            <span>{item.title}</span>

            {/* Sidebar Badge */}
            {item.badge && <NavBadge>{item.badge}</NavBadge>}

            {/* Sidebar Chevron */}
            <ChevronRight className='ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90' />
          </SidebarMenuButton>
        </CollapsibleTrigger>

        {/* Collapsible Content */}
        <CollapsibleContent className='CollapsibleContent'>
          <SidebarMenuSub data-test-id="menu-sub">
            {item.items.map((subItem) => (
              <SidebarMenuSubItem key={subItem.title}>
                <SidebarMenuSubButton
                  asChild
                  isActive={checkIsActive(href, subItem)}
                >
                  <Link to={subItem.url} onClick={() => setOpenMobile(false)}>
                    {subItem.icon && <subItem.icon />}
                    <span>{subItem.title}</span>
                    {subItem.badge && <NavBadge>{subItem.badge}</NavBadge>}
                  </Link>
                </SidebarMenuSubButton>
              </SidebarMenuSubItem>
            ))}
          </SidebarMenuSub>
        </CollapsibleContent>
      </SidebarMenuItem>
    </Collapsible>
  )
}

const SidebarMenuCollapsedDropdown = ({
  item,
  href,
}: {
  item: NavCollapsible
  href: string
}) => {
  return (
    <SidebarMenuItem>
      <DropdownMenu>
        {/* Sidebar Trigger */}
        <DropdownMenuTrigger asChild>


          <SidebarMenuButton
            tooltip={item.title}
            isActive={checkIsActive(href, item)}
          >
            {/* Sidebar Icon */}
            {item.icon && <item.icon />}

            {/* Sidebar Title */}
            <span>{item.title}</span>

            {/* Sidebar Badge */}
            {item.badge && <NavBadge>{item.badge}</NavBadge>}

            {/* Sidebar Chevron */}
            <ChevronRight className='ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90' />
          </SidebarMenuButton>

        </DropdownMenuTrigger>


        {/* Dropdown Menu Content */}
        <DropdownMenuContent side='right' align='start' sideOffset={4}>

        {/* Dropdown Menu Label */}
          <DropdownMenuLabel>
            {item.title} {item.badge ? `(${item.badge})` : ''}
          </DropdownMenuLabel>

          {/* Dropdown Menu Separator */}
          <DropdownMenuSeparator />

          {/* Dropdown Menu Items */}
          {item.items.map((sub) => (
            <DropdownMenuItem key={`${sub.title}-${sub.url}`} asChild>

              {/* Dropdown Menu Item Link */}
              <Link
                to={sub.url}
                className={`${checkIsActive(href, sub) ? 'bg-secondary' : ''}`}
              >
                {/* Dropdown Menu Item Icon */}
                {sub.icon && <sub.icon />}

                {/* Dropdown Menu Item Title */}
                <span className='max-w-52 text-wrap'>{sub.title}</span>
                {sub.badge && (
                  <span className='ml-auto text-xs'>{sub.badge}</span>
                )}
              </Link>

            </DropdownMenuItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  )
}

// TODO: Decouple this function from NavGroup to Routing layer
function checkIsActive(href: string, item: NavItem, mainNav = false) {
  const currentPath = typeof href === 'string' ? href: String(href);
  const itemURL = typeof item.url === 'string' ? item.url : String(item.url);

  return (
    currentPath === item.url || // endpoint?search-parameter
    currentPath.split('?')[0] === item.url  || // endpoint
    !!item?.items?.filter((index) => index.url === currentPath).length || // if child nav is active
    (mainNav &&
      currentPath.split('/')[1] !== '' &&
      currentPath.split('/')[1] === itemURL.split('/')[1])
  );
}
