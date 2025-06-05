import { useTheme } from '@/core/theme/hooks/useTheme';
import { LOCATION_ICON_COLORS, DEFAULT_COLORS } from '../constants/location-icon-colors';
import type { LocationIconBgColor } from '@/types/domain/location-types';

/**
 * Hook to get the appropriate background and text colors for a location icon based on the current theme
 *
 * @param color - The location icon background color name
 * @returns An object containing the text and background colors for the current theme
 *
 * @example
 * ```tsx
 * const { text, background } = useLocationBgColor('red');
 * return <div style={{ color: text, backgroundColor: background }}>Content</div>;
 * ```
 */
export function useLocationBgColor(color?: LocationIconBgColor) {
  const { isDarkMode } = useTheme();
  const theme = isDarkMode ? 'dark' : 'light';

  if (!color) {
    return DEFAULT_COLORS[theme];
  }

  return LOCATION_ICON_COLORS[color][theme];
}