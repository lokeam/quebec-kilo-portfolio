import { useCallback, useEffect, useMemo } from 'react';
import debounce from 'lodash/debounce';

/**
 * Configuration for a single breakpoint rule
 */
export interface BreakpointRule<T> {
  /** Breakpoint width in pixels */
  breakpoint: number;
  /** Value to apply when width is below this breakpoint */
  value: T;
}

/**
 * Configuration for the element breakpoint behavior
 */
export interface ElementBreakpointConfig<T> {
  /** CSS selector to identify the target element */
  selector: string;
  /** Array of breakpoint rules, ordered from smallest to largest width */
  breakpointRules: BreakpointRule<T>[];
  /** Default value to use when width is above all breakpoints */
  defaultValue: T;
  /** Callback function to handle value changes */
  onBreakpointChange: (value: T) => void;
  /** Optional debounce delay in ms (default: 16) */
  debounceDelay?: number;
}

/**
 * A React hook that observes an element's width and applies different values based on breakpoints.
 *
 * @example
 * ```tsx
 * useElementBreakpoint({
 *   selector: '[data-card="wishlist"]',
 *   breakpointRules: [
 *     { breakpoint: 480, value: { showTags: false, showRating: false } },
 *     { breakpoint: 768, value: { showTags: true, showRating: false } }
 *   ],
 *   defaultValue: { showTags: true, showRating: true },
 *   onBreakpointChange: (value) => setVisibility(value)
 * });
 * ```
 *
 * @param config - Configuration object for the breakpoint behavior
 * @returns void
 */
export function useElementBreakpoint<T>(config: ElementBreakpointConfig<T>): void {
  const {
    selector,
    breakpointRules,
    defaultValue,
    onBreakpointChange,
    debounceDelay = 16
  } = config;

  // Create a stable reference to the resize handler
  const handleResize = useCallback((width: number) => {
    const matchingRule = [...breakpointRules]
      .reverse()
      .find(rule => width <= rule.breakpoint);
    onBreakpointChange(matchingRule ? matchingRule.value : defaultValue);
  }, [breakpointRules, defaultValue, onBreakpointChange]);

  // Create a stable reference to the debounce function
  const debouncedHandleResize = useMemo(() => {
    return debounce(handleResize, debounceDelay);
  }, [handleResize, debounceDelay]);

  // useEffect kicks everything off
  useEffect(() => {
    let resizeObserver: ResizeObserver | null = null;
    const targetEls = document.querySelectorAll(selector);

    try {
      // Handle initial size exactly once per element
      for (const el of targetEls) {
        const width = el.getBoundingClientRect().width;
        handleResize(width);
      }

      // Set up observer for future changes
      resizeObserver = new ResizeObserver(entries => {
        for (const entry of entries) {
          debouncedHandleResize(entry.contentRect.width);
        }
      });

      for (const el of targetEls) {
        resizeObserver?.observe(el);
      }
    } catch (error) {
      console.error('Error setting up ResizeObserver:', error);
    }

    // Cleanup observer on unmount
    return () => {
      debouncedHandleResize.cancel();
      resizeObserver?.disconnect();
    };
  },[selector, handleResize, debouncedHandleResize]);
}
