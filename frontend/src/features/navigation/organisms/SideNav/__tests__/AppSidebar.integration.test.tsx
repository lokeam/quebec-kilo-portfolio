import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { vi } from 'vitest';
import { AppSidebar } from '../AppSidebar';

/**
 * @note This component is primarily tested through:
 * - Unit tests of individual sidebar components
 * - E2E tests covering critical navigation paths
 *
 * Integration tests here cover basic rendering and critical paths only.
 * Full component testing is handled at the E2E level for better ROI.
 */

vi.mock('@/shared/components/ui/sidebar', () => ({
  Sidebar: vi.fn(({ children, ...props }) => (
    <div data-testid="sidebar" {...props}>
      {children}
    </div>
  )),
  SidebarContent: vi.fn(({ children }) => <div>{children}</div>),
  SidebarHeader: vi.fn(({ children }) => <div>{children}</div>),
  SidebarGroup: vi.fn(({ children }) => <div>{children}</div>),
  SidebarGroupLabel: vi.fn(({ children }) => <div>{children}</div>),
  SidebarMenu: vi.fn(({ children }) => <div>{children}</div>),
  SidebarMenuItem: vi.fn(({ children }) => <div>{children}</div>),
  SidebarMenuSub: vi.fn(({ children }) => <div>{children}</div>),
  SidebarMenuSubItem: vi.fn(({ children }) => <div>{children}</div>),
  SidebarMenuSubButton: vi.fn(({ children }) => <button>{children}</button>),
  SidebarMenuButton: vi.fn(({ children }) => <button>{children}</button>),
  SidebarFooter: vi.fn(({ children }) => <div>{children}</div>),
  SidebarRail: vi.fn(({ children }) => <div>{children}</div>),
  useSidebar: vi.fn().mockReturnValue({
    isMobile: false,
    isCollapsed: false,
    setIsCollapsed: vi.fn(),
    setIsMobile: vi.fn(),
  }),
}));

vi.mock('@/features/navigation/config/navigation', () => ({
  navigation: [
    {
      title: 'General',
      items: [
        {
          title: 'Dashboard',
          href: '/',
          icon: 'dashboard',
        },
      ],
    },
  ],
}));

describe('AppSidebar', () => {
  const renderAppSidebar = () => {
    render(
      <MemoryRouter>
        <AppSidebar />
      </MemoryRouter>
    );
  };

  it('renders basic navigation structure', () => {
    renderAppSidebar();
    expect(screen.getByTestId('sidebar')).toBeInTheDocument();
  });
});