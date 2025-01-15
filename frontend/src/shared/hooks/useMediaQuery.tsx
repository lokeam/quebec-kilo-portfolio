import { useEffect, useState } from 'react';

/**
 * useMediaQuery Hook
 *
 * A custom hook that provides reactive media query matching functionality.
 * Listens to viewport changes and returns whether or not the current viewport matches the provided media query.
 *
 * @example
 * ```tsx
 * // Basic usage
 * const isDesktop = useMediaQuery('(min-width: 768px)');
 *
 * // With tailwind breakpoints
 * const isLarge = useMediaQuery('(min-width: 1024px)');
 *
 * // Complex queries
 * const isPrintAndLarge = useMediaQuery('print and (min-width: 768px)');
 * ```
 *
 * @param {string} query - CSS media query string
 * @returns {boolean} Whether the viewport matches the media query
 *
 * @notes
 * - Uses the Window.matchMedia() API
 * - Automatically cleans up event listeners on unmount
 * - Updates reactively when viewport changes
 * - SSR safe (defaults to false)
 *
 * @bestPractices
 * - Memoize the query string if it's dynamically generated
 * - Use standard CSS media query syntax
 * - Consider using Tailwind's breakpoint conventions:
 *   - sm: 640px
 *   - md: 768px
 *   - lg: 1024px
 *   - xl: 1280px
 *   - 2xl: 1536px
 *
 * @example Advanced usage with dynamic queries
 * ```tsx
 * const breakpoint = useMemo(() => `(min-width: ${size}px)`, [size]);
 * const matches = useMediaQuery(breakpoint);
 * ```
 *
 * @see useDrawerPosition - Example of hook being used for responsive drawer positioning
 * @see https://developer.mozilla.org/en-US/docs/Web/API/Window/matchMedia
 */

export function useMediaQuery(query: string) {
  const [matches, setMatches] = useState(false);

  useEffect(() => {
    const media = window.matchMedia(query);
    setMatches(media.matches);

    const listener = (event: MediaQueryListEvent) => setMatches(event.matches);
    media.addEventListener('change', listener);

    return () => media.removeEventListener('change', listener);
  }, [query]);

  return matches;
}
