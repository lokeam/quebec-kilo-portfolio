import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { ThemeMode, UIMode } from '@/core/theme/constants/themeConstants';
import { DEFAULT_THEME_MODE, DEFAULT_UI_MODE, SYSTEM_DARK_MODE_QUERY } from '@/core/theme/constants/themeConstants';
import { themeConfig } from '@/core/theme/theme.config';

// Import restoration module to ensure it runs before store creation
import '@/core/theme/hooks/useThemeRestoration';

interface ThemeState {
  mode: ThemeMode;
  uiMode: UIMode;
  isSystemPreference: boolean;
  actions: {
    changeTheme: (mode: ThemeMode) => void;
    changeUIMode: (uiMode: UIMode) => void;
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
        uiMode: DEFAULT_UI_MODE,
        isSystemPreference: false,

      actions: {

        changeTheme: (mode: ThemeMode) => {
          console.log('🎨 Theme changing to:', mode);
          console.log('🔍 Current localStorage before change:', localStorage.getItem('qko-theme-storage'));

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

          // Check if data was saved after a short delay
          setTimeout(() => {
            console.log('🔍 localStorage after theme change:', localStorage.getItem('qko-theme-storage'));
          }, 100);
        },

        changeUIMode: (uiMode: UIMode) => {
          set({ uiMode });
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
      name: 'qko-theme-storage',
      partialize: (state) => ({
        mode: state.mode,
        uiMode: state.uiMode,
        isSystemPreference: state.isSystemPreference,
      }),
    }
  )
);

// Selector hooks for better performance
export const useThemeMode = () => useThemeStore((state) => state.mode);
export const useUIMode = () => useThemeStore((state) => state.uiMode);
export const useThemeActions = () => useThemeStore((state) => state.actions);
export const useIsSystemPreference = () => useThemeStore((state) => state.isSystemPreference);