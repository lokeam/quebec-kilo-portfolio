import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { ThemeMode } from '../constants/themeConstants'
import { DEFAULT_THEME_MODE, SYSTEM_DARK_MODE_QUERY } from '../constants/themeConstants'
import { themeConfig } from '@/core/theme/theme.config';

interface ThemeState {
  mode: ThemeMode;
  isSystemPreference: boolean;
  actions: {
    toggleTheme: () => void;
    setTheme: (mode: ThemeMode) => void;
    enableSystemPreference: () => (() => void) | void; // Return cleanup function OR void
    disableSystemPreference: () => void;
    applyTheme: (mode: ThemeMode) => void;
  }
};

export const useThemeStore = create<ThemeState>()(
  persist(
    (set, get) => ({
      mode: DEFAULT_THEME_MODE,
      isSystemPreference: false,

      actions: {
        toggleTheme: () => {
          const currentMode = get().mode
          set({
            mode: currentMode === 'light' ? 'dark' : 'light',
            isSystemPreference: false,
          })
        },

        setTheme: (mode: ThemeMode) => {
          set({
            mode,
            isSystemPreference: false,
          })
        },

        enableSystemPreference: () => {
          const darkModeMediaQuery = window.matchMedia(SYSTEM_DARK_MODE_QUERY)
          const systemPrefersDark = darkModeMediaQuery.matches

          set({
            mode: systemPrefersDark ? 'dark' : 'light',
            isSystemPreference: true,
          });

          const handler = (event: MediaQueryListEvent) => {
            if (get().isSystemPreference) {
              set({ mode: event.matches ? 'dark' : 'light'});
            }
          };

          // Listen for system theme changes
          darkModeMediaQuery.addEventListener('change', handler);

          // Return cleanup function
          return () => darkModeMediaQuery.removeEventListener('change', handler);
        },

        disableSystemPreference: () => {
          set({ isSystemPreference: false })
        },

        applyTheme: (mode: ThemeMode) => {
          const root = window.document.documentElement;

          // Remove old theme
          root.classList.remove('light', 'dark');

          // Add new theme
          root.classList.add(mode);

          // Apply theme variables
          Object.entries(themeConfig[mode]).forEach(([key, value]) => {
            root.style.setProperty(`--${key}`, value)
          });
        }
      },
    }),
    {
      name: 'theme-storage',
      // Only persist the mode and isSystemPreference
      partialize: (state) => ({
        mode: state.mode,
        isSystemPreference: state.isSystemPreference,
      }),
    }
  )
);

// Decouple selector hooks for better performance
export const useThemeMode = () => useThemeStore((state) => state.mode);
export const useThemeActions = () => useThemeStore((state) => state.actions);
export const useIsSystemPreference = () => useThemeStore((state) => state.isSystemPreference);
