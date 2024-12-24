import { describe, it, expect } from 'vitest'
import { render, screen, act } from '@testing-library/react'
import { useTheme } from '@mui/material'
import { ThemeProvider } from '@/core/theme/providers/ThemeProvider';
import { useThemeStore } from '@/core/theme/stores/useThemeStore'
import { SYSTEM_DARK_MODE_QUERY } from '@/core/theme/constants/themeConstants';
import { darkTheme } from '@/core/theme/darkTheme';
import { createTheme, Button } from '@mui/material'

// Test component that uses MUI theme
function ThemedButton() {
  const theme = useTheme()
  return (
    <Button
      data-testid="themed-button"
      sx={{ bgcolor: theme.palette.primary.main }}
    >
      Test Button
    </Button>
  )
}

describe('Theme System Integration', () => {
  beforeEach(() => {
    const { actions } = useThemeStore.getState()
    actions.setTheme('light')
    actions.disableSystemPreference()
  })

  it('should apply theme changes throughout the component tree', () => {
    render(
      <ThemeProvider>
        <ThemedButton />
      </ThemeProvider>
    )

    const { actions } = useThemeStore.getState()

    act(() => {
      actions.setTheme('dark')
    })

    const button = screen.getByTestId('themed-button')
    // Verify theme is actually applied to component
    expect(button).toHaveStyle({
      backgroundColor: 'rgb(96, 165, 250)' // dark theme primary color
    })
  })

  it('should sync system preference changes with UI', () => {
    const theme = createTheme(darkTheme);

    // Setup matchMedia mock to return dark mode
    window.matchMedia = vi.fn().mockImplementation(query => ({
      matches: query === SYSTEM_DARK_MODE_QUERY,
      media: query,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    }))

    render(
      <ThemeProvider enableSystemPreference>
        <ThemedButton />
      </ThemeProvider>
    )

    const button = screen.getByTestId('themed-button')
    expect(button).toHaveStyle({
      backgroundColor: theme.palette.primary.main
    })
  })
})