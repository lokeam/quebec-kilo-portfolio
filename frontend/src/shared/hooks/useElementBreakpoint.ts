import { useEffect, useCallback } from 'react';
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

  const handleResize = useCallback((width: number) => {
    // Find the first matching breakpoint rule
    const matchingRule = [...breakpointRules]
      .reverse()
      .find(rule => width <= rule.breakpoint);

    // Apply matching rule value or default value
    onBreakpointChange(matchingRule ? matchingRule.value : defaultValue);
  }, [breakpointRules, defaultValue, onBreakpointChange]);

  useEffect(() => {
    const updateBreakpoint = debounce(() => {
      const element = document.querySelector(selector);
      if (!element) return;

      const width = element.getBoundingClientRect().width;
      handleResize(width);
    }, debounceDelay);

    const resizeObserver = new ResizeObserver(updateBreakpoint);
    const element = document.querySelector(selector);

    if (element) {
      resizeObserver.observe(element);
      // Initial check
      updateBreakpoint();
    }

    return () => {
      resizeObserver.disconnect();
      updateBreakpoint.cancel();
    };
  }, [selector, handleResize, debounceDelay]);
}