import { describe, it, expect, beforeEach } from 'vitest'
import { render, screen, act } from '@testing-library/react'
import { useTheme } from '@mui/material'
import { ThemeProvider } from '@/core/theme/providers/ThemeProvider'
import { useThemeStore } from '@/core/theme/stores/useThemeStore'

// Helper component to access theme
function ThemeTestConsumer() {
  const theme = useTheme()
  return (
    <div data-testid="theme-mode">{theme.palette.mode}</div>
  );
}

describe('ThemeProvider', () => {
  beforeEach(() => {
    // Reset store state
    const { actions } = useThemeStore.getState()
    actions.setTheme('light')
    actions.disableSystemPreference()
    localStorage.clear()
  })

  it('should provide theme based on store state', () => {
    render(
      <ThemeProvider>
        <ThemeTestConsumer />
      </ThemeProvider>
    )

    expect(screen.getByTestId('theme-mode')).toHaveTextContent('light')
  });

  it('should update theme when Zustand store state changes', () => {
    render(
      <ThemeProvider>
        <ThemeTestConsumer />
      </ThemeProvider>
    )

    act(() => {
      const { actions } = useThemeStore.getState()
      actions.setTheme('dark')
    })

    expect(screen.getByTestId('theme-mode')).toHaveTextContent('dark');
  });

  it('should enable system preference when prop is true', () => {
    render(
      <ThemeProvider enableSystemPreference>
        <ThemeTestConsumer />
      </ThemeProvider>
    )

    expect(useThemeStore.getState().isSystemPreference).toBe(true);
  });

  it('should disable system preference on unmount', () => {
    const { unmount } = render(
      <ThemeProvider enableSystemPreference>
        <ThemeTestConsumer />
      </ThemeProvider>
    )

    unmount();

    expect(useThemeStore.getState().isSystemPreference).toBe(false);
  })
})