import { ReactNode, useEffect, useMemo } from 'react'
import { ThemeProvider as MUIThemeProvider, createTheme } from '@mui/material'
import { useThemeStore } from '@/core/theme/stores/useThemeStore';
import { lightTheme } from '@/core/theme/lightTheme';
import { darkTheme } from '@/core/theme/darkTheme';
import type { Theme } from '@mui/material';

interface ThemeProviderProps {
  children: ReactNode
  enableSystemPreference?: boolean
};

export const ThemeProvider = ({
  children,
  enableSystemPreference = false
}: ThemeProviderProps) => {
  // Use selector pattern for better performance
  const mode = useThemeStore((state) => state.mode);
  const actions = useThemeStore((state) => state.actions);

  // Memoize theme creation to prevent unnecessary recalculations
  const theme: Theme = useMemo(
    () => createTheme(mode === 'light' ? lightTheme : darkTheme),
    [mode]
  );

  // Handle system preference initialization
  useEffect(() => {
    if (enableSystemPreference) {
      const unsubscribe = actions.enableSystemPreference();

      // Cleanup system preference listeners when component unmounts
      // or when enableSystemPreference changes
      return () => {
        if (typeof unsubscribe === 'function') {
          unsubscribe();
        }
        actions.disableSystemPreference();
      }
    }
  }, [enableSystemPreference, actions]);

  return (
    <MUIThemeProvider theme={theme}>
      {children}
    </MUIThemeProvider>
  )
}
