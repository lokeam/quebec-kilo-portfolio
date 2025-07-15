/**
 * Text Wrapping Utility
 *
 * Makes potentially long blocks of text look better by automatically adjusting how it wraps across lines.
 * Instead of ugly ragged edges, this creates nice, even text blocks.
 *
 * What it does:
 * - Makes text wrap more evenly (similar justified text)
 * - Works automatically when your page resizes
 * - Uses modern browser features when available
 * - Cleans up after itself to prevent memory issues
 *
 * Ratio explained:
 * - 1.0 = Text is as compact as possible (narrowest width)
 * - 0.9 = Mostly compact with a little breathing room
 * - 0.5 = Halfway between compact and spread out
 * - 0.0 = Text spreads out to fill the full width
 *
 * @example
 * // Apply general text balancing to all headings
 * <script>
 *   import { initTextBalancing } from '../utils/text-balance.js';
 *   document.addEventListener('DOMContentLoaded', initTextBalancing);
 * </script>
 *
 * @example
 * // Apply text balancing to specific text
 * <script>
 *   import { balanceElement } from '../utils/text-balance.js';
 *   balanceElement('.hero-title', 0.9); // 90% = very compact, 10% = spread out
 * </script>
 *
 * @example
 * // Fix one specific piece of text
 * <script>
 *   import { balanceText } from '../utils/text-balance.js';
 *   const element = document.querySelector('.my-heading');
 *   balanceText(element, 0.8); // 80% = compact, 20% = spread out
 * </script>
 *
 * @example
 * // Memory management to clean up when we're done (prevents memory leaks)
 * <script>
 *   import { cleanupTextBalancing } from '../utils/text-balance.js';
 *   cleanupTextBalancing();
 * </script>
 */

const SYMBOL_KEY = '__wrap_b';
const SYMBOL_NATIVE_KEY = '__wrap_n';

/**
 * Check if browser supports native text-balancing
 * @returns {boolean} True if browser supports CSS text-wrap: balance
 */
const isTextWrapBalanceSupported = () => {
  return typeof CSS !== 'undefined' && CSS.supports('text-wrap', 'balance');
};

/**
 * Binary search algorithm for text balancing
 * @param {HTMLElement} element - The element to balance
 * @param {number} ratio - Balance ratio (0-1, default: 1)
 */
const relayout = (element: HTMLElement, ratio: number = 1) => {
  const wrapper = element;
  const container = wrapper.parentElement;

  if (!container) return;

  const updateWidth = (width: number) => {
    wrapper.style.maxWidth = width + 'px';
  };

  // Reset wrapper width
  wrapper.style.maxWidth = '';

  // Get initial container size
  const width = container.clientWidth;
  const height = container.clientHeight;

  if (!width) return;

  // Binary search to find optimal width
  let lower = width / 2 - 0.25;
  let upper = width + 0.5;
  let middle: number;

  // Ensure we don't search widths lower than when text overflows
  updateWidth(lower);
  lower = Math.max(wrapper.scrollWidth, lower);

  while (lower + 1 < upper) {
    middle = Math.round((lower + upper) / 2);
    updateWidth(middle);

    if (container.clientHeight === height) {
      upper = middle;
    } else {
      lower = middle;
    }
  }

  // Update wrapper width with ratio
  updateWidth(upper * ratio + width * (1 - ratio));
};

/**
 * Initialize global relayout function
 * Sets up browser support detection and global relayout function
 */
const initGlobalRelayout = () => {
  if (typeof window === 'undefined') return;

  // Set native support flag
  (window as any)[SYMBOL_NATIVE_KEY] = isTextWrapBalanceSupported() ? 1 : 2;

  // Set global relayout function
  (window as any)[SYMBOL_KEY] = relayout;
};

/**
 * Make text in an element look better by adjusting how it wraps
 * @param {HTMLElement} element - The element containing the text
 * @param {number} ratio - Controls text width: 1 = compact, 0 = spread out (default: 1)
 *
 * @example
 * const element = document.querySelector('.my-heading');
 * balanceText(element, 0.8); // 80% compact, 20% spread out
 */
export function balanceText(element: HTMLElement, ratio: number = 1) {
  // Skip if browser supports native text-balancing
  if (typeof window !== 'undefined' && (window as any)[SYMBOL_NATIVE_KEY] === 1) {
    element.style.textWrap = 'balance';
    return;
  }

  // Initialize global function if needed
  if (typeof window !== 'undefined' && !(window as any)[SYMBOL_KEY]) {
    initGlobalRelayout();
  }

  // Apply balancing
  if (typeof window !== 'undefined') {
    (window as any)[SYMBOL_KEY](element, ratio);
  }

  // Set up ResizeObserver for responsive behavior
  if (typeof ResizeObserver !== 'undefined') {
    const container = element.parentElement;
    if (container) {
      const observer = new ResizeObserver(() => {
        if (typeof window !== 'undefined') {
          (window as any)[SYMBOL_KEY](element, ratio);
        }
      });
      observer.observe(container);

      // Store observer for cleanup
      (element as any).__wrap_o = observer;
    }
  }
}

/**
 * Make all headings on the page look better automatically
 * Finds all h1-h6 elements and improves how their text wraps
 *
 * @example
 * // In Astro component
 * <script>
 *   import { initTextBalancing } from '../utils/text-balance.js';
 *   document.addEventListener('DOMContentLoaded', initTextBalancing);
 * </script>
 */
export function initTextBalancing() {
  const headings = document.querySelectorAll('h1, h2, h3, h4, h5, h6');
  headings.forEach((element) => balanceText(element as HTMLElement));
}

/**
 * Make text in multiple elements look better
 * @param {string} selector - CSS selector to find elements (like '.hero-title')
 * @param {number} ratio - Controls text width: 1 = compact, 0 = spread out (default: 1)
 *
 * @example
 * // Make all elements with class 'hero-title' look better
 * balanceElement('.hero-title', 0.9); // 90% compact, 10% spread out
 *
 * @example
 * // Make all headings in a specific container look better
 * balanceElement('.content h1, .content h2', 0.8); // 80% compact, 20% spread out
 */
export function balanceElement(selector: string, ratio: number = 1) {
  const elements = document.querySelectorAll(selector);
  elements.forEach((element) => balanceText(element as HTMLElement, ratio));
}

/**
 * Clean up when you're done (prevents memory leaks)
 * Call this when you're done with the page or component
 *
 * @example
 * // Clean up when page is about to unload
 * window.addEventListener('beforeunload', cleanupTextBalancing);
 *
 * @example
 * // Clean up in component
 * <script>
 *   import { cleanupTextBalancing } from '../utils/text-balance.js';
 *   cleanupTextBalancing();
 * </script>
 */
export function cleanupTextBalancing() {
  const elements = document.querySelectorAll('[data-br]');
  elements.forEach((element) => {
    const observer = (element as any).__wrap_o;
    if (observer) {
      observer.disconnect();
      delete (element as any).__wrap_o;
    }
  });
}