import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { DomainSwitcher } from '../DomainSwitcher';
import { vi, describe, it, expect, beforeEach } from 'vitest'; // Or use jest if you're using Jest
import * as sidebarModule from '@/shared/components/ui/sidebar';

type SidebarMenuButtonProps = {
  children: React.ReactNode;
  className?: string;
  size?: 'sm' | 'lg';
  onClick?: () => void;
};

// Mock sidebar context with all required properties
const mockSidebarContext = {
  state: 'expanded' as const,
  isMobile: false,
  open: false,
  setOpen: vi.fn(),
  openMobile: false,
  setOpenMobile: vi.fn(),
  setCollapsed: vi.fn(),
};

// Sidebar mock
vi.mock('@/shared/components/ui/sidebar', () => ({
  ...vi.importActual('@/shared/components/ui/sidebar'),
  useSidebar: () => mockSidebarContext,
  SidebarMenu: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  SidebarMenuItem: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  SidebarMenuButton: ({ children, ...props }: SidebarMenuButtonProps) => (
    <button role="button" {...props}>{children}</button>
  ),
}));

// Dropdown menu mock with all required components
vi.mock('@/shared/components/ui/dropdown-menu/DropdownMenu', () => ({
  DropdownMenu: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="dropdown-menu">{children}</div>
  ),
  DropdownMenuTrigger: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="dropdown-trigger">{children}</div>
  ),
  DropdownMenuContent: ({
    children,
    side,
    align,
    sideOffset,
    className
  }: {
    children: React.ReactNode;
    side?: string;
    align?: string;
    sideOffset?: number;
    className?: string;
  }) => (
    <div
      data-testid="dropdown-content"
      data-side={side}
      data-align={align}
      data-offset={sideOffset}
      className={className}
    >
      {children}
    </div>
  ),
  DropdownMenuItem: ({
    children,
    onClick,
    className
  }: {
    children: React.ReactNode;
    onClick?: () => void;
    className?: string;
  }) => (
    <div
      role="menuitem"
      data-testid="dropdown-item"
      onClick={onClick}
      className={className}
    >
      {children}
    </div>
  ),
  DropdownMenuLabel: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="dropdown-label">{children}</div>
  ),
  DropdownMenuShortcut: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="dropdown-shortcut">{children}</div>
  ),
}));

// Mock test data
const mockDomains = [
  {
    name: 'Domain 1',
    description: 'First domain description',
    logo: () => <svg data-testid="domain1-logo" />,
  },
  {
    name: 'Domain 2',
    description: 'Second domain description',
    logo: () => <svg data-testid="domain2-logo" />,
  },
];

describe('DomainSwitcher', () => {
  it('renders with initial domain selected', () => {
    render(<DomainSwitcher domains={mockDomains} />);

    const domainTitle = screen.getByText('Domain 1', {
      selector: 'span.truncate.font-semibold'
    });
    expect(domainTitle).toBeInTheDocument();
    expect(screen.getByText('First domain description')).toBeInTheDocument();

    // Be more specific about which logo we're looking for
    const activeLogo = screen.getAllByTestId('domain1-logo')[0]; // Get the first instance
    expect(activeLogo).toBeInTheDocument();
  });

  it('opens dropdown menu when clicked', async () => {
    render(<DomainSwitcher domains={mockDomains} />);

    const trigger = screen.getByRole('button');
    await userEvent.click(trigger);

    expect(screen.getByText('Domains')).toBeInTheDocument();
    const menuItems = screen.getAllByRole('menuitem');
    expect(menuItems).toHaveLength(mockDomains.length);
  });

  it('switches active domain when new domain is selected', async () => {
    render(<DomainSwitcher domains={mockDomains} />);

    // Open dropdown
    const trigger = screen.getByRole('button');
    await userEvent.click(trigger);

    // Click second domain
    const menuItems = screen.getAllByRole('menuitem');
    await userEvent.click(menuItems[1]);

    // Verify domain switch
    const activeDomainTitle = screen.getByText('Domain 2', {
      selector: 'span.truncate.font-semibold'
    });
    expect(activeDomainTitle).toBeInTheDocument();
    expect(screen.getByText('Second domain description')).toBeInTheDocument();

    // Be more specific about which logo we're looking for
    const activeLogo = screen.getAllByTestId('domain2-logo')[0];
    expect(activeLogo).toBeInTheDocument();
  });

  it('displays keyboard shortcuts for domains', async () => {
    render(<DomainSwitcher domains={mockDomains} />);

    const trigger = screen.getByRole('button');
    await userEvent.click(trigger);

    expect(screen.getByText('⌘1')).toBeInTheDocument();
    expect(screen.getByText('⌘2')).toBeInTheDocument();
  });
});

describe('Mobile behavior for DomainSwitcher', () => {
  beforeEach(() => {
    vi.spyOn(sidebarModule, 'useSidebar').mockImplementation(() => ({
      ...mockSidebarContext,
      isMobile: true,
      toggleSidebar: vi.fn(),
    }));
  });

  it('positions dropdown on bottom for mobile view', async () => {
    render(<DomainSwitcher domains={mockDomains} />);

    // Open the dropdown
    const trigger = screen.getByRole('button');
    await userEvent.click(trigger);

    // Get the dropdown content
    const dropdownContent = screen.getByTestId('dropdown-content');

    // Verify the positioning props
    expect(dropdownContent).toHaveAttribute('data-side', 'bottom');
    expect(dropdownContent).toHaveAttribute('data-align', 'start');
    expect(dropdownContent).toHaveAttribute('data-offset', '4');
  });

  it('positions dropdown on right for desktop view', async () => {
    // Override the mock to test desktop view
    vi.spyOn(sidebarModule, 'useSidebar').mockImplementation(() => ({
      ...mockSidebarContext,
      isMobile: false,
      toggleSidebar: vi.fn(),
    }));

    render(<DomainSwitcher domains={mockDomains} />);

    // Open the dropdown
    const trigger = screen.getByRole('button');
    await userEvent.click(trigger);

    // Get the dropdown content
    const dropdownContent = screen.getByTestId('dropdown-content');

    // Verify the positioning props
    expect(dropdownContent).toHaveAttribute('data-side', 'right');
    expect(dropdownContent).toHaveAttribute('data-align', 'start');
    expect(dropdownContent).toHaveAttribute('data-offset', '4');
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });
});
