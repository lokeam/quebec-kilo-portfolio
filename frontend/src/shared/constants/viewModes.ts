/**
 * Possible view modes for different features.
 * This provides type safety and runtime access to valid modes.
 */
export const viewModes = ['grid', 'list', 'table'] as const;
export type ViewMode = typeof viewModes[number];

/**
 * Simple function to validate view modes at runtime
 * This helps us safely handle values from localStorage
 */
export function isValidViewMode(value: unknown): value is ViewMode {
  return typeof value === 'string' && viewModes.includes(value as ViewMode);
}

/**
 * A simple function to manage localStorage for view modes
 * This handles persistence with proper type checking
 */
export function getStoredViewMode<T extends ViewMode>(
  key: string,
  defaultMode: T,
  allowedModes: readonly T[]
): T {
  const stored = localStorage.getItem(key);

  /* Check if stored viewMode is valid for this page */
  if (stored && allowedModes.includes(stored as T)) {
    return stored as T;
  }

  return defaultMode;
}

/**
 * Define which view modes are available for each feature
 * This creates a clear mapping of features to their allowed modes
 */
export const featureViewModes = {
  onlineServices: {
    allowed: ['grid', 'list', 'table'] as const,
    default: 'grid' as const,
    storageKey: 'onlineServices-view-mode'
  },
  library: {
    allowed: ['grid', 'list'] as const,
    default: 'grid' as const,
    storageKey: 'library-view-mode'
  }
} as const;
