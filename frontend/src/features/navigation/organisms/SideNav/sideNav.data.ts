import {
  IconHelp,
  IconLayoutDashboard,
  IconPackages,
  IconSettings,
  IconCloudDataConnection,
  IconShoppingCartHeart,
  IconBook,
} from '@tabler/icons-react'
import { Bell, Building2, LibraryBig, Gamepad2, Clapperboard, CircleDollarSign } from 'lucide-react'
import { type SidebarData } from '@/features/navigation/organisms/SideNav/sideNav.types'

export const sidebarData: SidebarData = {
  user: {
    name: 'John Doe',
    email: 'john.doe@example.com',
    avatar: '/avatars/shadcn.jpg',
  },
  domains: [
    {
      name: 'Games',
      logo: Gamepad2,
      description: 'Console and PC',
    },
    {
      name: 'Movies / Shows',
      logo: Clapperboard,
      description: 'Physical and Streaming',
    },
    {
      name: 'Books (beta)',
      logo: IconBook,
      description: 'Physical and Digital',
    },
  ],
  navGroups: [
    {
      title: 'General',
      items: [
        {
          title: 'Dashboard',
          url: '/',
          icon: IconLayoutDashboard,
        },
        {
          title: 'Notifications',
          url: '/notifications',
          icon: Bell,
        },
        // {
        //   title: 'Apps',
        //   url: '/apps',
        //   icon: IconPackages,
        // },
      ],
    },
    {
      title: 'Pages',
      items: [
        {
          title: 'Library',
          icon: LibraryBig,
          url: '/library',
        },
        {
          title: 'Spend Tracking',
          url: '/spend-tracking',
          badge: '5',
          icon: CircleDollarSign,
        },
        {
          title: 'Wishlist',
          url: '/wishlist',
          badge: '3',
          icon: IconShoppingCartHeart,
        },
        {
          title: 'Physical Locations',
          icon: Building2,
          url: '/physical-locations',
        },
        {
          title: 'Online Services',
          icon: IconCloudDataConnection,
          url: '/online-services',
        },
      ],
    },
    {
      title: 'Other',
      items: [
        {
          title: 'Settings',
          icon: IconSettings,
          url: '/settings',
        },
        {
          title: 'FAQs',
          url: '/support',
          icon: IconHelp,
        },
      ],
    },
  ],
}
