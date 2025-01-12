import { Dialog, DialogContent, DialogOverlay, DialogPortal, DialogTitle, DialogHeader, DialogFooter } from '@/shared/components/ui/dialog';
import { Input } from '@/shared/components/ui/input';
import { cn } from '@/shared/components/ui/utils';
import { VisuallyHidden } from '@radix-ui/react-visually-hidden';
import { type ReactNode } from 'react';

interface SearchDialogProps {
  /** Controls the open state of the dialog */
  open: boolean;

  /** Callback fired when the open state changes */
  onOpenChange: (open: boolean) => void;

  /** Current value of the search input */
  searchQuery: string;

  /** Callback fired when the search input changes */
  onSearchChange: (event: React.ChangeEvent<HTMLInputElement>) => void;

  /** Placeholder text for the search input */
  searchPlaceholder: string;

  /** Title of the dialog - shown in header or visually hidden for ARIA purposes*/
  dialogTitle: string;

  /** Content to render in the dialog body */
  children: ReactNode;

  /** Optional CSS classes to apply to the dialog */
  className?: string;

  /** Whether to visually hide the dialog header */
  hideHeader?: boolean;

  /** Optional footer content */
  footer?: ReactNode;

  /** Optional trigger element */
  trigger?: ReactNode;
}

/**
 * SearchDialog component - A reusable dialog for search functionality
 * Custom styles on the DialogContent component override Radix UI defaults
 *
 * @example
 * ```tsx
 * <SearchDialog
 *   open={isOpen}
 *   onOpenChange={setIsOpen}
 *   searchQuery={query}
 *   onSearchChange={handleChange}
 *   searchPlaceholder="Search..."
 *   dialogTitle="Search Items"
 * >
 *   {results.map(result => (
 *     <SearchResult key={result.id} {...result} />
 *   ))}
 * </SearchDialog>
 * ```
 *
 * @see {@link https://ui.shadcn.com/docs/components/dialog Dialog Component}
 */
export function SearchDialog({
  open,
  onOpenChange,
  searchQuery,
  onSearchChange,
  searchPlaceholder,
  dialogTitle,
  children,
  className,
  hideHeader = false,
  footer,
  trigger,
}: SearchDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      {trigger}

      <DialogPortal>
        <DialogOverlay className="fixed inset-0 bg-black/50 data-[state=open]:animate-fadeIn data-[state=closed]:animate-fadeOut" />

        <div className="fixed inset-0 overflow-hidden pt-20">
          <div className="flex items-start justify-center">
            <DialogContent
              className={cn(
                "fixed pt-10 px-4 left-[50%] top-[5%] z-50 w-[95vw] sm:w-[90vw] max-w-[940px] -translate-x-[50%] translate-y-0 bg-background rounded-lg shadow-lg",
                "transition-opacity duration-200",
                "opacity-0 data-[state=open]:opacity-100",
                className
              )}
            >
              <div className="flex flex-col h-[calc(100vh-120px)]">
                {hideHeader ? (
                  <VisuallyHidden>
                    <DialogTitle>{dialogTitle}</DialogTitle>
                  </VisuallyHidden>
                ) : (
                  <DialogHeader>
                    <DialogTitle>{dialogTitle}</DialogTitle>
                  </DialogHeader>
                )}

                <div className="shrink-0 p2 md:p-4">
                  <Input
                    placeholder={searchPlaceholder}
                    value={searchQuery}
                    onChange={onSearchChange}
                    className="w-full"
                  />
                </div>

                <div className="flex-1 overflow-y-auto scrollbar-thin scrollbar-thumb-gray-800 scrollbar-track-transparent">
                  <div className="py-2 space-y-2 md:p-4 md:space-y-4">
                    {children}
                  </div>
                </div>

                {footer && (
                  <DialogFooter>
                    {footer}
                  </DialogFooter>
                )}
              </div>
            </DialogContent>
          </div>
        </div>
      </DialogPortal>
    </Dialog>
  );
}
