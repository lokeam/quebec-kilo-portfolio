import { type ReactNode, memo } from 'react';

// Shadcn UI components
import { Button } from '@/shared/components/ui/button';
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from '@/shared/components/ui/drawer';

// Utils + Hooks
import { cn } from '@/shared/components/ui/utils';
import { useDrawerPosition } from '@/features/dashboard/lib/hooks/useDrawerPosition';

/**
 * DrawerContainer Component
 *
 * IMPORTANT IMPLEMENTATION NOTES:
 *
 * 1. Component Structure:
 *    - We use Shadcn UI's Drawer components directly (DrawerContent, DrawerHeader, etc.)
 *    - We only memoize our custom content components
 *
 * 2. Why This Approach Works:
 *    - Previous attempts to wrap Shadcn components with memo() failed because TypeScript
 *      couldn't reconcile all the internal React props these components expect
 *    - By using Shadcn components directly and only memoizing our content, we avoid
 *      TypeScript conflicts with component props
 *
 * 3. Key Patterns Used:
 *    - Component Composition: Using smaller, memoized components inside Shadcn containers
 *    - Props Spreading: Using ...rest to pass through additional props to DrawerContent
 *    - Explicit Test IDs: Each important element has a data-testid for testing
 *
 * 4. Performance Optimizations:
 *    - MemoizedDrawerHeaderContent: Prevents re-renders of the header content
 *    - MemoizedDrawerFooter: Prevents re-renders of the static footer
 */

interface DrawerContainerProps {
// Trigger props
  triggerText?: ReactNode;
  triggerVariant?: 'default' | 'outline' | 'secondary' | 'ghost';

  // Drawer content props
  title: string;
  description?: string;
  children: ReactNode;

  // Control props
  open: boolean;
  onOpenChange: (open: boolean) => void;

  // Optional customization
  className?: string;
  drawerPosition?: 'right' | 'bottom' | 'left' | 'top';
  mobilePosition?: 'bottom' | 'top';
  desktopBreakpoint?: string;
}

/**
 * Memoized component for drawer header content
 * Note: We memoize the content, not the DrawerHeader component itself
 */
const MemoizedDrawerHeaderContent = memo(function DrawerHeaderContent({
  title,
  description
}: {
  title: string;
  description?: string;
}) {
  return (
    <>
      <DrawerTitle data-testid="drawer-title">{title}</DrawerTitle>
      {description && <DrawerDescription>{description}</DrawerDescription>}
    </>
  );
});

/**
 * Memoized component for drawer footer
 * Contains the cancel button and its container
 */
const MemoizedDrawerFooter = memo(function DrawerFooter() {
  return (
    <DrawerClose asChild>
      <div className="px-4 pb-4 w-full">
      <Button
          variant="outline"
          className="w-full"
          data-testid="drawer-cancel-button"
        >
          Cancel
        </Button>
      </div>
    </DrawerClose>
  );
});

/**
 * Main DrawerContainer component
 *
 * Important: We use Shadcn's DrawerContent directly instead of trying to wrap it
 * This avoids TypeScript conflicts with component props
 */
export function DrawerContainer({
  triggerText,
  triggerVariant = 'default',
  title,
  description,
  children,
  open,
  onOpenChange,
  className,
  drawerPosition = 'right',
  mobilePosition = 'bottom',
  desktopBreakpoint = '(min-width: 768px)',
  ...rest /* <-- This collects any additional props to pass to DrawerContent */
}: DrawerContainerProps) {
  const { positionStyles, currentPosition } = useDrawerPosition({
    drawerPosition,
    mobilePosition,
    desktopBreakpoint,
    className,
  });

  return (
    <Drawer
      open={open}
      onOpenChange={onOpenChange}
      direction={currentPosition}
      data-testid="drawer-container"
      modal={true} /* <-- This is the default value, but we're explicitly setting it here to handle focus management properly */
      aria-modal="true"
    >
      {triggerText && (
        <DrawerTrigger asChild>
          <Button
            variant={triggerVariant}
            data-testid="drawer-trigger"
          >
            {triggerText}
          </Button>
        </DrawerTrigger>
      )}

      {/* Use DrawerContent directly, passing through additional props */}
      <DrawerContent
        className={cn(positionStyles, className)}
        data-testid="drawer-content"
        {...rest}
      >
        <DrawerHeader className="text-left" data-testid="drawer-header">
          <MemoizedDrawerHeaderContent title={title} description={description} />
        </DrawerHeader>

        <div className="px-4" data-testid="drawer-body">{children}</div>

        <MemoizedDrawerFooter />
      </DrawerContent>
    </Drawer>
  );
}
