import { useThemeStore } from '../stores/useThemeStore';

/**
 * Hook to access theme information
 * @returns Object containing theme information and actions
 */
export function useTheme() {
  const mode = useThemeStore((state) => state.mode);

  return {
    isDarkMode: mode === 'dark',
    mode,
  };
}