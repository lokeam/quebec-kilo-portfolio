import { Link } from 'react-router-dom'
import {
  ChevronsUpDown,
  LogOut,
  // Sparkles,
} from '@/shared/components/ui/icons'
import { IconSettings } from '@/shared/components/ui/icons'
import { Avatar, AvatarFallback, AvatarImage } from '@/shared/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/shared/components/ui/sidebar'

// Auth hook
import { useAuthContext } from '@/core/auth/context-provider/AuthContext'

export function NavUser() {
  const { user, logout } = useAuthContext();
  const { isMobile } = useSidebar();

  //console.log('NavUser Debug:', { user, isAuthenticated, isLoading });

  if (!user) {
    return <></>
  }

  // Safe access to user properties with fallbacks
  const userName = user.name || user.email || 'User';
  const userEmail = user.email || '';
  const userPicture = user.picture || '';
  const userInitial = userName.charAt(0).toUpperCase();

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu modal={false}>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size='lg'
              className='data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground'
            >
              <Avatar className='h-8 w-8 rounded-lg'>
                <AvatarImage src={userPicture} alt={userName} />

                {/* Safe fallback with user initial */}
                <AvatarFallback className='rounded-lg'>{userInitial}</AvatarFallback>

              </Avatar>
              <div className='grid flex-1 text-left text-sm leading-tight'>
                <span className='truncate font-semibold'>{userName}</span>
                <span className='truncate text-xs'>{userEmail}</span>
              </div>
              <ChevronsUpDown className='ml-auto size-4' />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className='w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg'
            side={isMobile ? 'bottom' : 'right'}
            align='end'
            sideOffset={4}
            avoidCollisions={false}
          >
            <DropdownMenuLabel className='p-0 font-normal'>
              <div className='flex items-center gap-2 px-1 py-1.5 text-left text-sm'>
                <Avatar className='h-8 w-8 rounded-lg'>
                  <AvatarImage src={userPicture} alt={userName} />

                  <AvatarFallback className='rounded-lg'>{userInitial}</AvatarFallback>
                </Avatar>
                <div className='grid flex-1 text-left text-sm leading-tight'>
                  <span className='truncate font-semibold'>{userName}</span>
                  <span className='truncate text-xs'>{userEmail}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            {/* <DropdownMenuGroup>
              <DropdownMenuItem>
                <Sparkles />
                Upgrade to Pro
              </DropdownMenuItem>
            </DropdownMenuGroup> */}
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              {/* <DropdownMenuItem asChild>
                <Link to='/settings'>
                  <CreditCard />
                  Billing
                </Link>
              </DropdownMenuItem> */}
              <DropdownMenuItem asChild>
                <Link to='/settings'>
                  <IconSettings />
                  Settings
                </Link>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={logout}>
              <LogOut />
              Sign out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}