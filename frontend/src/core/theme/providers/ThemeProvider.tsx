import { useEffect, type ReactNode } from 'react'
import { useThemeStore } from '@/core/theme/stores/useThemeStore';

interface ThemeProviderProps {
  children: ReactNode;
  enableSystemPreference?: boolean;
};

export const ThemeProvider = ({
  children,
  enableSystemPreference = false
}: ThemeProviderProps) => {
  // Use selector pattern for better performance
  const actions = useThemeStore((state) => state.actions);
  const currentMode = useThemeStore((state) => state.mode);

  console.log('🎨 ThemeProvider rendering with mode:', currentMode);
  console.log('🔍 localStorage in ThemeProvider:', localStorage.getItem('qko-theme-storage'));

  // Apply theme to DOM when mode changes (including on initialization)
  useEffect(() => {
    console.log('🎨 Applying theme to DOM:', currentMode);
    actions.updateDOM(currentMode);
  }, [currentMode]); // actions is stable from Zustand, no need to include it

  // Handle system preference initialization
  useEffect(() => {
    if (enableSystemPreference) {
      const unsubscribe = actions.enableSystemPreference();

      return () => {
        if (typeof unsubscribe === 'function') {
          unsubscribe();
        }
        actions.disableSystemPreference();
      }
    }
  }, [enableSystemPreference, actions]);

  return (
    <>{children}</>
  );
}
