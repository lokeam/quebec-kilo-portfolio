import { describe, it, expect } from 'vitest';
import { render, screen, act } from '@testing-library/react'
import { Button } from '@/shared/components/ui/button';
import { ThemeProvider } from '@/core/theme/providers/ThemeProvider';
import { useThemeStore } from '@/core/theme/stores/useThemeStore';
import { SYSTEM_DARK_MODE_QUERY } from '@/core/theme/constants/themeConstants';

// Test component that displays current theme
function ThemedButton() {
  const theme = useThemeStore((state) => state.mode)
  return (
    <Button data-testid="themed-button">
      Current theme: {theme}
    </Button>
  );
};

describe('Theme System Integration', () => {
  beforeEach(() => {
    const { actions } = useThemeStore.getState();
    actions.setTheme('light');
    actions.disableSystemPreference();
    localStorage.clear();
    document.documentElement.className = '';
  });

  it('should apply theme changes throughout the component tree', () => {
    render(
      <ThemeProvider>
        <ThemedButton />
      </ThemeProvider>
    );

    const { actions } = useThemeStore.getState();

    act(() => {
      actions.setTheme('dark');
    })

    const button = screen.getByTestId('themed-button');
    expect(document.documentElement).toHaveClass('dark');
    expect(button).toHaveTextContent('Current theme: dark');
  });

  it('should sync system preference changes with UI', () => {
    // Setup matchMedia mock to return dark mode
    window.matchMedia = vi.fn().mockImplementation(query => ({
      matches: query === SYSTEM_DARK_MODE_QUERY,
      media: query,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    }));

    render(
      <ThemeProvider>
        <ThemedButton />
      </ThemeProvider>
    );

    // Enable system preference
    act(() => {
      const { actions } = useThemeStore.getState();
      actions.enableSystemPreference();
    });

    const button = screen.getByTestId('themed-button');
    expect(document.documentElement).toHaveClass('dark');
    expect(button).toHaveTextContent('Current theme: dark');
  });

  it('should persist theme preference', () => {
    const { unmount } = render(
      <ThemeProvider>
        <ThemedButton />
      </ThemeProvider>
    );

    act(() => {
      const { actions } = useThemeStore.getState();
      actions.setTheme('dark');
    })

    // Check immediate state
    expect(document.documentElement).toHaveClass('dark');
    expect(screen.getByTestId('themed-button')).toHaveTextContent('Current theme: dark');

    // Unmount then remount to verify persistence
    unmount();

    render(
      <ThemeProvider>
        <ThemedButton />
      </ThemeProvider>
    );

    // Verify that we're persisting state
    expect(document.documentElement).toHaveClass('dark')
    expect(screen.getByTestId('themed-button')).toHaveTextContent('Current theme: dark')
  })

  it('should apply correct CSS variables for each theme', () => {
    render(
      <ThemeProvider>
        <ThemedButton />
      </ThemeProvider>
    );

    // Check light theme variables
    expect(document.documentElement).toHaveStyle({
      '--background': '0 0% 100%',
      '--foreground': '222.2 47.4% 11.2%'
    });

    // Switch to dark theme
    act(() => {
      const { actions } = useThemeStore.getState()
      actions.setTheme('dark')
    });

    // Check dark theme variables
    expect(document.documentElement).toHaveStyle({
      '--background': '220 20% 12%',
      '--foreground': '220 10% 98%'
    });
  });
});
