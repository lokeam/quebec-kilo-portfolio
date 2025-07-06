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
