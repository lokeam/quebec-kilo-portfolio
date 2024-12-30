import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi, describe, it, expect, beforeEach } from 'vitest';
import { NavGroup } from '../NavGroup';
import type { NavGroup as INavGroup, NavLink, NavCollapsible } from '../sideNav.types';
import { SidebarProvider } from '@/shared/components/ui/sidebar';
import * as sidebarModule from '@/shared/components/ui/sidebar';
import { MemoryRouter } from 'react-router-dom';
import { Home, Settings, Users } from 'lucide-react';

// Mock the sidebar context
const mockSidebarContext = {
  state: 'expanded' as const,
  isMobile: false,
  open: false,
  setOpen: vi.fn(),
  openMobile: false,
  setOpenMobile: vi.fn(),
  setCollapsed: vi.fn(),
  toggleSidebar: vi.fn(),
};

// Mock test data with proper typing
const mockNavGroup: INavGroup = {
  title: 'Main Navigation',
  items: [
    {
      title: 'Home',
      url: '/',
      icon: Home,
      badge: '3',
    } as NavLink,
    {
      title: 'Settings',
      icon: Settings,
      items: [
        {
          title: 'Users',
          url: '/settings/users',
          icon: Users,
          badge: 'New',
        } as NavLink,
        {
          title: 'Preferences',
          url: '/settings/preferences',
        } as NavLink,
      ],
    } as NavCollapsible,
  ],
};

// Mock react-router-dom
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    BrowserRouter: ({ children }: { children: React.ReactNode }) => children,
    useLocation: vi.fn().mockReturnValue({ pathname: '/settings', search: '' })
  };
});

const renderWithProviders = (ui: React.ReactNode, { route = '/', sidebarState = 'expanded' } = {}) => {
  return render(
    <MemoryRouter initialEntries={[route]}>
      <SidebarProvider defaultOpen={sidebarState === 'expanded'}> {/* Changed to use defaultOpen */}
        {ui}
      </SidebarProvider>
    </MemoryRouter>
  );
};

describe('NavGroup', () => {

  beforeEach(() => {
    vi.spyOn(sidebarModule, 'useSidebar').mockImplementation(() => mockSidebarContext);
  });

  describe('Basic rendering', () => {
    it('renders group title', () => {
      renderWithProviders(<NavGroup {...mockNavGroup} />);
      expect(screen.getByText('Main Navigation')).toBeInTheDocument();
    });

    it('renders all top-level items', () => {
      renderWithProviders(<NavGroup {...mockNavGroup} />);
      expect(screen.getByText('Home')).toBeInTheDocument();
      expect(screen.getByText('Settings')).toBeInTheDocument();
    });

    it('renders badges when provided', () => {
      renderWithProviders(<NavGroup {...mockNavGroup} />);
      expect(screen.getByText('3')).toBeInTheDocument();
    });
  });

  describe('Navigation links', () => {
    it('renders simple links correctly', () => {
      renderWithProviders(<NavGroup {...mockNavGroup} />);
      const homeLink = screen.getByRole('link', { name: /home/i });
      expect(homeLink).toHaveAttribute('href', '/');
    });

    it('closes mobile sidebar when link is clicked', async () => {
      renderWithProviders(<NavGroup {...mockNavGroup} />);
      const homeLink = screen.getByRole('link', { name: /home/i });
      await userEvent.click(homeLink);
      expect(mockSidebarContext.setOpenMobile).toHaveBeenCalledWith(false);
    });
  });

  describe('Collapsible menus', () => {
    it('renders collapsible menu when item has sub-items and sidebar is expanded', async () => {
      renderWithProviders(<NavGroup {...mockNavGroup} />);
      const settingsButton = screen.getByRole('button', { name: /settings/i });
      await userEvent.click(settingsButton);

      expect(screen.getByText('Users')).toBeInTheDocument();
      expect(screen.getByText('Preferences')).toBeInTheDocument();
    });

    it('shows sub-items with correct URLs', async () => {
      renderWithProviders(<NavGroup {...mockNavGroup} />);
      const settingsButton = screen.getByRole('button', { name: /settings/i });
      await userEvent.click(settingsButton);

      const usersLink = screen.getByRole('link', { name: /users/i });
      expect(usersLink).toHaveAttribute('href', '/settings/users');
    });
  });

  describe('Collapsed state', () => {
    it('renders dropdown menu when sidebar is collapsed', async () => {

      vi.spyOn(sidebarModule, 'useSidebar').mockImplementation(() => ({
        ...mockSidebarContext,
        state: 'collapsed'
      }));

      renderWithProviders(<NavGroup {...mockNavGroup} />, {
        route: '/settings/users',
        sidebarState: 'collapsed'
      });

      // Act - Click the Settings button to open dropdown
      const settingsButton = screen.getByRole('button', { name: /settings/i });
      await userEvent.click(settingsButton);

      // Assert - Look for the sub-items directly
      expect(screen.getByText('Users')).toBeInTheDocument();
      expect(screen.getByText('Preferences')).toBeInTheDocument();
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });
});