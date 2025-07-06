import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { ThemeMode } from '@/core/theme/constants/themeConstants';
import { DEFAULT_THEME_MODE, SYSTEM_DARK_MODE_QUERY } from '@/core/theme/constants/themeConstants';
import { themeConfig } from '@/core/theme/theme.config';

interface ThemeState {
  mode: ThemeMode;
  isSystemPreference: boolean;
  actions: {
    changeTheme: (mode: ThemeMode) => void;
    enableSystemPreference: () => (() => void) | void;
    disableSystemPreference: () => void;
    updateDOM: (mode: ThemeMode) => void;
  }
};

export const useThemeStore = create<ThemeState>()(
  // Save theme choice to local storage
  persist(
    (set, get) => ({
      mode: DEFAULT_THEME_MODE,
      isSystemPreference: false,

      actions: {

        changeTheme: (mode: ThemeMode) => {
          if (mode === 'system') {
            set({
              mode: 'system',
              isSystemPreference: true,
            });

            // Get current system preference
            const darkModeMediaQuery = window.matchMedia(SYSTEM_DARK_MODE_QUERY);
            const isSystemSetToDarkMode = darkModeMediaQuery.matches;
            const currentTheme = isSystemSetToDarkMode ? 'dark' : 'light';

            get().actions.updateDOM(currentTheme);
          } else {
            set({
              mode,
              isSystemPreference: false,
            });
            get().actions.updateDOM(mode);
          }
        },

        enableSystemPreference: () => {
          const darkModeMediaQuery = window.matchMedia(SYSTEM_DARK_MODE_QUERY);
          const isSystemSetToDarkMode = darkModeMediaQuery.matches;

          set({
            mode: isSystemSetToDarkMode ? 'dark' : 'light',
            isSystemPreference: true,
          });

          const handler = (event: MediaQueryListEvent) => {
            if (get().isSystemPreference) {
              const newMode = event.matches ? 'dark' : 'light';
              set({ mode: newMode });
              get().actions.updateDOM(newMode);
            }
          };

          darkModeMediaQuery.addEventListener('change', handler);
          return () => darkModeMediaQuery.removeEventListener('change', handler);
        },

        disableSystemPreference: () => {
          set({ isSystemPreference: false });
        },

        updateDOM: (mode: ThemeMode) => {
          const windowDocumentElement = window.document.documentElement;

          // Update CSS classes
          windowDocumentElement.classList.remove('light', 'dark');
          windowDocumentElement.classList.add(mode);

          // Update CSS variables
          const config = themeConfig[mode];
          for (const key in config) {
            const value = config[key as keyof typeof config];
            windowDocumentElement.style.setProperty(`--${key}`, value);
          }
        }
      },
    }),
    {
      name: 'theme-storage',
      partialize: (state) => ({
        mode: state.mode,
        isSystemPreference: state.isSystemPreference,
      }),
    }
  )
);

// Selector hooks for better performance
export const useThemeMode = () => useThemeStore((state) => state.mode);
export const useThemeActions = () => useThemeStore((state) => state.actions);
export const useIsSystemPreference = () => useThemeStore((state) => state.isSystemPreference);