import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '@/shared/components/ui/sidebar'
import { NavGroup } from '@/features/navigation/organisms/SideNav/NavGroup'
import { NavUser } from '@/features/navigation/organisms/SideNav/NavUser'
import { TeamSwitcher } from '@/features/navigation/organisms/SideNav/TeamSwitcher'
import { sidebarData } from '@/features/navigation/organisms/SideNav/sideNav.data'

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible='icon' variant='floating' {...props}>
      <SidebarHeader>
        <TeamSwitcher teams={sidebarData.teams} />
      </SidebarHeader>
      <SidebarContent>
        {sidebarData.navGroups.map((props) => (
          <NavGroup key={props.title} {...props} />
        ))}
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={sidebarData.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}