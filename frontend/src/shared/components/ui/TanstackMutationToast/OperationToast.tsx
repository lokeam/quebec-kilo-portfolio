import type { ReactNode } from 'react';
import { toast as sonnerToast } from 'sonner';

// Icons
import { X } from '@/shared/components/ui/icons';

interface OperationToastProps {
  id: string | number;
  message: string;
  icon: ReactNode;
}

/**
 * Custom operation toast component with dismiss button on top right.
 * Mirrors the layout of IntroToast for consistency.
 */
export function OperationToast(props: OperationToastProps) {
  const { id, message, icon } = props;

  return (
    <div className="flex rounded-lg bg-white dark:bg-gray-800 shadow-lg ring-1 ring-black/5 dark:ring-white/10 w-full md:max-w-[420px] items-start p-4 gap-3">
      {/* Icon */}
      <div className="flex-shrink-0 mt-0.5">
        {icon}
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0">
        <div className="text-sm text-gray-900 dark:text-gray-100 leading-relaxed">
          {message}
        </div>
      </div>

      {/* Close Button */}
      <button
        onClick={() => sonnerToast.dismiss(id)}
        className="flex-shrink-0 rounded-md p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        aria-label="Dismiss toast"
      >
        <X className="h-4 w-4" />
      </button>
    </div>
  );
}