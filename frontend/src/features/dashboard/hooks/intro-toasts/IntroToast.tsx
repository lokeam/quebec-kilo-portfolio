import React from 'react';
import { toast as sonnerToast } from 'sonner';
import { X, Info } from '@/shared/components/ui/icons';

interface IntroToastProps {
  id: string | number;
  title?: string;
  description: string;
  icon?: React.ReactNode;
}

/**
 * Safely renders text with styled spans using a simple pattern matching approach.
 * Supports <bold>text</bold> and <highlight>text</highlight> patterns.
 */
function renderStyledText(text: string) {
  // Split by bold tags
  const parts = text.split(/(<bold>.*?<\/bold>)/g);

  return parts.map((part, index) => {
    const boldMatch = part.match(/<bold>(.*?)<\/bold>/);
    const highlightMatch = part.match(/<highlight>(.*?)<\/highlight>/);

    if (boldMatch) {
      return (
        <span key={index} className="font-bold text-blue-500">
          {boldMatch[1]}
        </span>
      );
    }

    if (highlightMatch) {
      return (
        <span key={index} className="font-bold text-green-500">
          {highlightMatch[1]}
        </span>
      );
    }

    return part;
  });
}

/**
 * Custom headless toast component for intro toasts.
 * Shows for 60 seconds and is dismissable.
 * Supports styled text using safe pattern matching.
 */
export function IntroToast(props: IntroToastProps) {
  const { id, title, description, icon = <Info className="h-5 w-5 text-blue-600" /> } = props;

  return (
    <div className="flex rounded-lg bg-white dark:bg-gray-800 shadow-lg ring-1 ring-black/5 dark:ring-white/10 w-full md:max-w-[420px] items-start p-4 gap-3">
      {/* Icon */}
      <div className="flex-shrink-0 mt-0.5">
        {icon}
      </div>

      {/* Content */}
      <div className="flex-1 min-w-0">
        {title && (
          <p className="text-sm font-medium text-gray-900 dark:text-gray-100 mb-1">
            {title}
          </p>
        )}
        <div className="text-sm text-gray-600 dark:text-gray-300 leading-relaxed">
          {renderStyledText(description)}
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

/**
 * Abstracted toast function for intro toasts.
 * Shows for 60 seconds and is dismissable.
 */
export function showIntroToast(toast: Omit<IntroToastProps, 'id'>) {
  return sonnerToast.custom((id) => (
    <IntroToast
      id={id}
      title={toast.title}
      description={toast.description}
      icon={toast.icon}
    />
  ), {
    duration: 60000, // 60 seconds
  });
}