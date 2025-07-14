// Sentry error boundary
import * as Sentry from '@sentry/react';
import { useLocation } from 'react-router-dom';
import { categorizeError } from '@/core/error/categorization';

import { DefaultErrorFallbackPage } from '@/core/error/pages/DefaultErrorFallbackPage';

// Note: Removed Tanstack Query error boundary in favor of Sentry
// import { QueryErrorResetBoundary } from '@tanstack/react-query';


interface ErrorBoundaryProviderProps {
  children: React.ReactNode;
  FallbackComponent?: React.ComponentType<{
    error: Error;
    resetErrorBoundary: () => void;
  }>;
}

/**
 * Global error boundary provider that catches runtime errors and displays fallback UI
 * Based on React Error Boundary pattern: https://react.dev/reference/react/Component#catching-rendering-errors-with-an-error-boundary
 */
export const ErrorBoundaryProvider = ({
  children,
  FallbackComponent = DefaultErrorFallbackPage
}: ErrorBoundaryProviderProps) => {
  const location = useLocation();

  return (
    <Sentry.ErrorBoundary
      fallback={({ error, resetError }) => (
        <FallbackComponent
          error={error as Error}
          resetErrorBoundary={resetError}
        />
      )}
      beforeCapture={(scope, hint) => {
        // Categorize the error
        const error = (hint as { originalException?: Error })?.originalException;
        if (error) {
          const errorContext = categorizeError(error, {
            currentPage: location.pathname,
            timestamp: new Date().toISOString(),
          });

          // Add error categorization tags
          Object.entries(errorContext.tags).forEach(([key, value]) => {
            scope.setTag(key, value);
          });

          // Add error context
          scope.setContext('error_categorization', errorContext.context);
        }

        // Set error level based on priority
        scope.setLevel("error");

        // Add navigation context
        scope.setContext('navigation', {
          currentRoute: location.pathname,
          search: location.search,
          hash: location.hash,
        });

        // Add performance context if available
        if (window.performance) {
          const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
          if (navigation) {
            scope.setContext('performance', {
              pageLoadTime: navigation.loadEventEnd - navigation.loadEventStart,
              domContentLoaded: navigation.domContentLoadedEventEnd - navigation.domContentLoadedEventStart,
              firstPaint: performance.getEntriesByName('first-paint')[0]?.startTime,
              firstContentfulPaint: performance.getEntriesByName('first-contentful-paint')[0]?.startTime,
            });
          }
        }
      }}
    >
      {children}
    </Sentry.ErrorBoundary>
  );
};
