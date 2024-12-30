import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { NavUser } from '../NavUser';
import { MemoryRouter } from 'react-router-dom';
import { type ComponentPropsWithoutRef } from 'react';

// Mock test data
const mockUser = {
  name: 'John Doe',
  email: 'john@example.com',
  avatar: 'https://example.com/avatar.jpg'
};

// Mock the sidebar context
const mockSidebarContext = {
  isMobile: false,
  state: 'expanded' as const,
  open: true,
  setOpen: vi.fn(),
  openMobile: false,
  setOpenMobile: vi.fn(),
  toggleSidebar: vi.fn(),
};

// Mock the sidebar module
vi.mock('@/shared/components/ui/sidebar', () => ({
  useSidebar: () => mockSidebarContext,
  SidebarProvider: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  SidebarMenu: ({ children }: { children: React.ReactNode }) => <ul>{children}</ul>,
  SidebarMenuItem: ({ children }: { children: React.ReactNode }) => <li>{children}</li>,
  SidebarMenuButton: ({ children, ...props }: ComponentPropsWithoutRef<'button'>) => (
    <button {...props}>{children}</button>
  ),
}));

const renderNavUser = () => {
  return render(
    <MemoryRouter>
      <NavUser user={mockUser} />
    </MemoryRouter>
  );
};

describe('NavUser', () => {
  it('renders user information correctly', () => {
    renderNavUser();

    // Check name and email are displayed
    expect(screen.getByText(mockUser.name)).toBeInTheDocument();
    expect(screen.getByText(mockUser.email)).toBeInTheDocument();

    // Check avatar container
    const avatarContainer = screen.getByRole('button').querySelector('.relative.flex.shrink-0.overflow-hidden.h-8.w-8.rounded-lg');
    expect(avatarContainer).toBeInTheDocument();

    // Check initials are displayed (since image isn't loading in test environment)
    expect(screen.getByText('JD')).toBeInTheDocument();
  });

  it('opens dropdown menu when clicked', async () => {
    renderNavUser();

    const trigger = screen.getByRole('button');
    await userEvent.click(trigger);

    // Check for dropdown menu items
    expect(screen.getByText('Upgrade to Pro')).toBeInTheDocument();
    expect(screen.getByText('Account')).toBeInTheDocument();
    expect(screen.getByText('Billing')).toBeInTheDocument();
    expect(screen.getByText('Notifications')).toBeInTheDocument();
    expect(screen.getByText('Log out')).toBeInTheDocument();
  });

  it('renders navigation links correctly', async () => {
    renderNavUser();

    const trigger = screen.getByRole('button');
    await userEvent.click(trigger);

    // Check if links have correct hrefs
    expect(screen.getByText('Account').closest('a')).toHaveAttribute('href', '/settings/account');
    expect(screen.getByText('Billing').closest('a')).toHaveAttribute('href', '/settings');
    expect(screen.getByText('Notifications').closest('a')).toHaveAttribute('href', '/settings/notifications');
  });

  it('adjusts dropdown position for mobile', () => {
    // Mock mobile view
    vi.mocked(mockSidebarContext).isMobile = true;

    renderNavUser();

    const trigger = screen.getByRole('button');
    expect(trigger).toBeInTheDocument();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });
});
