import { useMemo } from 'react';
import { cn } from '@/shared/components/ui/utils';
import { useMediaQuery } from '@/shared/hooks/useMediaQuery';

/**
 * useDrawerPosition Hook
 *
 * A custom hook that manages responsive drawer positioning and styling for the DrawerContainer component.
 * This hook handles both mobile and desktop viewport sizes, providing appropriate positioning and animations.
 *
 * @example
 * ```tsx
 * const { positionStyles, currentPosition, isDesktop } = useDrawerPosition({
 *   drawerPosition: 'right',
 *   mobilePosition: 'bottom',
 *   desktopBreakpoint: '(min-width: 768px)'
 * });
 *
 * // Use with DrawerContainer
 * <DrawerContainer
 *   direction={currentPosition}
 *   className={positionStyles}
 * >
 *   {children}
 * </DrawerContainer>
 * ```
 *
 * @param {Object} props - Configuration options for the drawer
 * @param {DrawerPosition} [props.drawerPosition='right'] - Desktop drawer position ('right', 'bottom', 'left', 'top')
 * @param {'bottom' | 'top'} [props.mobilePosition='bottom'] - Mobile drawer position (limited to vertical options)
 * @param {string} [props.desktopBreakpoint='(min-width: 768px)'] - Media query for desktop breakpoint
 * @param {string} [props.className] - Additional CSS classes to merge with position styles
 *
 * @returns {Object} Hook return object
 * @returns {string} return.positionStyles - Combined Tailwind classes for positioning and animation
 * @returns {DrawerPosition} return.currentPosition - Current active position based on viewport
 * @returns {boolean} return.isDesktop - Current viewport matches desktop breakpoint
 *
 * @see DrawerContainer - The template component that consumes this hook
 *
 * @notes
 * - Position Styles:
 *   - Desktop: Uses percentage-based widths for side drawers
 *   - Mobile: Full-width bottom/top sheets
 *   - Includes proper z-indexing and background styles
 *
 * - Breakpoints:
 *   - md: 768px (35% width)
 *   - lg: 1024px (45% width)
 *   - xl: 1280px (35% width)
 *
 * - Accessibility:
 *   - Works with Radix UI's Dialog primitive
 *   - Maintains proper focus management
 *   - Supports screen readers through aria attributes
 */

type DrawerPosition = 'right' | 'bottom' | 'left' | 'top';

interface UseDrawerPositionProps {
  drawerPosition?: DrawerPosition;
  mobilePosition?: 'bottom' | 'top';
  desktopBreakpoint?: string;
  className?: string;
};

export function useDrawerPosition({
  drawerPosition = 'right',
  mobilePosition = 'bottom',
  desktopBreakpoint = '(min-width: 768px)',
  className,
}: UseDrawerPositionProps) {
  const isDesktop = useMediaQuery(desktopBreakpoint);

  const positionStyles = useMemo(() => {
    const position = isDesktop ? drawerPosition : mobilePosition;
    const baseStyles = 'fixed z-50 bg-background';

    const positions = {
      bottom: `
        bottom-0 left-0 right-0 mt-24 rounded-t-lg
        md:h-[85vh]
      `,
      right: `
        fixed top-0 bottom-0 right-0
        md:top-0 md:bottom-0 md:right-0 md:left-[30%] lg:left-[45%] xl:left-[65%]
        rounded-l-lg md:rounded-t-md
      `,
      left: `
        top-0 left-0 h-full
        md:left-0 md:right-[70%] lg:right-[55%] xl:right-[35%]
        rounded-r-lg
      `,
      top: `
        top-0 left-0 right-0 mb-24 rounded-b-lg
        md:h-[85vh]
      `,
    };

    return cn(baseStyles, positions[position], className);
  }, [isDesktop, drawerPosition, mobilePosition, className]);

  const currentPosition = isDesktop ? drawerPosition : mobilePosition;

  return {
    positionStyles,
    currentPosition,
    isDesktop,
  };
}
