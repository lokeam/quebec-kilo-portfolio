import { type ReactNode, memo, useRef, useEffect } from 'react';

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

// Icons
import { Building2 } from 'lucide-react';
import { IconCloudDataConnection } from '@tabler/icons-react';

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
 *
 * 5. IMPORTANT USAGE PATTERNS:
 *    This component can be used in two distinct ways:
 *
 *    a) Auto-generated trigger buttons (simpler approach):
 *       - Provide either triggerAddLocation or triggerEditLocation prop
 *       - Don't provide open/onOpenChange props
 *       - The component will automatically create a button and handle drawer state
 *       - Example: <DrawerContainer triggerAddLocation="Add Item">...</DrawerContainer>
 *
 *    b) Controlled drawer with external buttons (more flexible):
 *       - Provide open and onOpenChange props to control drawer state externally
 *       - Create your own Button component with an onClick handler to set the open state to true
 *       - Don't use triggerAddLocation or triggerEditLocation with this approach
 *       - Example:
 *         const [isOpen, setIsOpen] = useState(false);
 *         return (
 *           <>
 *             <Button onClick={() => setIsOpen(true)}>Open Drawer</Button>
 *             <DrawerContainer open={isOpen} onOpenChange={setIsOpen}>...</DrawerContainer>
 *           </>
 *         )
 *
 *    ⚠️ WARNING: Do not mix these approaches! Using both triggerAddLocation/triggerEditLocation AND
 *    open/onOpenChange can lead to buttons not appearing or drawers not opening correctly.
 */

interface DrawerContainerProps {
// Trigger props
  triggerAddLocation?: ReactNode;
  triggerEditLocation?: ReactNode;
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
  triggerBtnIcon?: 'location' | 'digital';
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
      <DrawerTitle className="sr-only" data-testid="drawer-title">{title}</DrawerTitle>
      <DrawerDescription className="sr-only" data-testid="drawer-description">
        {description || title}
      </DrawerDescription>
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
  triggerAddLocation,
  triggerEditLocation,

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
  triggerBtnIcon,
  ...rest /* <-- This collects any additional props to pass to DrawerContent */
}: DrawerContainerProps) {
  const { positionStyles, currentPosition } = useDrawerPosition({
    drawerPosition,
    mobilePosition,
    desktopBreakpoint,
    className,
  });

  // Create ref for close button
  const closeButtonRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    if (open && closeButtonRef.current) {
      // Wait for drawer to open before focusing on close button
      setTimeout(() => {
        closeButtonRef.current?.focus();
      }, 100);
    }
  }, [open]);

  return (
    <Drawer
      open={open}
      onOpenChange={onOpenChange}
      direction={currentPosition}
      data-testid="drawer-container"
      modal={true} /* <-- This is the default value, but we're explicitly setting it here to handle focus management properly */
      aria-modal="true"
    >
      {/* Add Location Trigger */}
      {triggerAddLocation && (
        <DrawerTrigger asChild>
          <Button
            variant={triggerVariant}
            data-testid="drawer-trigger"
          >
            {triggerBtnIcon === 'location' && <Building2 />}
            {triggerBtnIcon === 'digital' && <IconCloudDataConnection />}
            {triggerAddLocation}
          </Button>
        </DrawerTrigger>
      )}

      {/* Edit Location Trigger */}
      {triggerEditLocation && (
        <DrawerTrigger asChild>
          <Button variant={triggerVariant} data-testid="drawer-trigger">
            {triggerEditLocation}
          </Button>
        </DrawerTrigger>
      )}

      {/* Use DrawerContent directly, passing through additional props */}
      <DrawerContent
        className={cn(positionStyles, className)}
        data-testid="drawer-content"
        aria-describedby="drawer-description"
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
