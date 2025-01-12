import { useEffect } from 'react';
import debounce from 'lodash/debounce';

/**
 * Configuration interface for the card label width behavior
 * @interface LabelWidthConfig
 */
interface LabelWidthConfig {
  /** CSS selector attribute to identify the target element */
  selectorAttribute: string;

  /** Breakpoint values in pixels for responsive behavior */
  breakpoints: {
    narrow: number;
    medium: number;
  };

  /** CSS width values to apply at different breakpoints */
  widths: {
    narrow: string;
    medium: string;
    wide: string;
  };
};


/**
 * A React hook that dynamically updates CSS custom properties for the SingleOnlineServiceCard
 * label widths based on the width of the sentinel element (ie: the first card in the list).
 *
 * @param config - Configuration object for controlling label width behavior
 * @param config.selectorAttribute - CSS selector to identify the target element
 * @param config.breakpoints - Pixel values defining width threshold
 * @param config.widths - CSS width values to apply at different breakpoints
 *
 * @example
 * ```tsx
 * useCardLabelWidth({
 *   selectorAttribute: '[data-card-container]',
 *   breakpoints: { narrow: 300, medium: 600 },
 *   widths: { narrow: '100%', medium: '75%', wide: '50%' }
 * });
 * ```
 *
 * @remarks
 * - Uses ResizeObserver to efficiently track element size changes
 * - Implements debouncing (16ms === 60fps) to optimize performance during resize
 * - Updates CSS custom properties: --label-max-width and --tw-text-opacity
 * - Automatically cleans up observers and event listeners on unmount
 *
 * @returns void
 */
export function useCardLabelWidth(config: LabelWidthConfig) {
  // Break specific config properties for optimal performance
  const { selectorAttribute, breakpoints, widths } = config;

  useEffect(() => {
    const updateLabelWidth = debounce(() => {
      const sentinelElement = document.querySelector(`${selectorAttribute}`);
      if (!sentinelElement) return;

      const cardWidth = sentinelElement.getBoundingClientRect().width;

      let labelWidth;

      if (cardWidth < breakpoints.narrow) {
        labelWidth = widths.narrow;
      } else if (cardWidth < breakpoints.medium) {
        labelWidth = widths.medium;
      } else {
        labelWidth = widths.wide;
      }

      document.documentElement.style.setProperty('--label-max-width', labelWidth);
      document.documentElement.style.setProperty('--tw-text-opacity', '1');
    }, 16); // 16ms is default debounce time and approx 60fps

    const resizeObserver = new ResizeObserver(updateLabelWidth);
    const sentinelElement = document.querySelector(`${selectorAttribute}`);

    if (sentinelElement) {
      resizeObserver.observe(sentinelElement);
    }

    return () => {
      resizeObserver.disconnect();
      updateLabelWidth.cancel();
    }
  }, [selectorAttribute, breakpoints, widths]);
};
