// Theme mode type
export type ThemeMode = 'light' | 'dark' | 'system';

// Media query system theme preference
export const SYSTEM_DARK_MODE_QUERY = '(prefers-color-scheme: dark)';

// System theme preference
export const DEFAULT_THEME_MODE: ThemeMode = 'light';

// Transition duration for theme changes
export const THEME_TRANSITION_DURATION = 200;

// Theme color palette names for consistency
export const THEME_COLORS = {
  PRIMARY: 'primary',
  SECONDARY: 'secondary',
  ERROR: 'error',
  WARNING: 'warning',
  SUCCESS: 'success',
  BACKGROUND: 'background',
  TEXT: 'text',
} as const;

// Z-index values for consistent layering
export const Z_INDEX = {
  MODAL: 1300,
  SNACKBAR: 1400,
  TOOLTIP: 1500,
  DRAWER: 1200,
} as const;
