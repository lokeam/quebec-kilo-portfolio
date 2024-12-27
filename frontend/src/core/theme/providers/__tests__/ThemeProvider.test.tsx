import { describe, it, expect, beforeEach } from 'vitest'
import { render, screen, act } from '@testing-library/react'
import { ThemeProvider } from '@/core/theme/providers/ThemeProvider'
import { useThemeStore } from '@/core/theme/stores/useThemeStore'
import { mockThemeConfig, getThemeVariables } from '@/core/theme/__tests__/theme.test-utils';

vi.mock('@core/theme/theme.config', () => ({
  themeConfig: mockThemeConfig,
}));

// Helper component to access theme
function ThemeTestConsumer() {
  const mode = useThemeStore((state) => state.mode);
  return (
    <div data-testid="theme-mode">{ mode }</div>
  );
}

describe('ThemeProvider', () => {
  beforeEach(() => {
    // Reset store state
    const { actions } = useThemeStore.getState();
    actions.setTheme('light');
    actions.disableSystemPreference();
    localStorage.clear();
    document.documentElement.className = '';
    document.documentElement.style.cssText = '';
  });

  it('should provide theme based on store state', () => {
    render(
      <ThemeProvider>
        <ThemeTestConsumer />
      </ThemeProvider>
    )

    expect(screen.getByTestId('theme-mode')).toHaveTextContent('light')
    expect(document.documentElement).toHaveClass('light');

    const themeVariables = getThemeVariables();
    expect(themeVariables.background).toBe(mockThemeConfig.light.background);
  });

  it('should update theme when Zustand store state changes', () => {
    render(
      <ThemeProvider>
        <ThemeTestConsumer />
      </ThemeProvider>
    )

    act(() => {
      const { actions } = useThemeStore.getState();
      actions.setTheme('dark');
    })

    expect(screen.getByTestId('theme-mode')).toHaveTextContent('dark');
    expect(document.documentElement).toHaveClass('dark');

    const themeVariables = getThemeVariables();
    expect(themeVariables.background).toBe(mockThemeConfig.dark.background);
  });

  it('should enable system preference when prop is true', () => {

    // Mock system dark mode preference
    Object.defineProperty(window, 'matchMedia', {
      value: vi.fn().mockImplementation(query => ({
        matches: true,
        media: query,
        addEventListener: vi.fn(),
        removeEventListener: vi.fn()
      }))
    });

    render(
      <ThemeProvider enableSystemPreference>
        <ThemeTestConsumer />
      </ThemeProvider>
    )

    expect(useThemeStore.getState().isSystemPreference).toBe(true);
    expect(screen.getByTestId('theme-mode')).toHaveTextContent('dark');
  });

  it('should disable system preference on unmount', () => {
    const { unmount } = render(
      <ThemeProvider enableSystemPreference>
        <ThemeTestConsumer />
      </ThemeProvider>
    )

    unmount();

    expect(useThemeStore.getState().isSystemPreference).toBe(false);
  });

  it('should apply CSS theme variables correctly', () => {
    render(
      <ThemeProvider>
        <ThemeTestConsumer />
      </ThemeProvider>
    );

    const themeVariables = getThemeVariables();
    Object.entries(mockThemeConfig.light).forEach(([key, value]) => {
      expect(themeVariables[key as keyof typeof themeVariables]).toBe(value);
    });
  });
});
