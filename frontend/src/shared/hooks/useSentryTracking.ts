import * as Sentry from '@sentry/react';
import { useCallback } from 'react';
import { useLocation } from 'react-router-dom';

/**
 * Custom hook for tracking user actions in Sentry
 *
 * Features:
 * - Track specific user interactions (clicks, form submissions, etc.)
 * - Add context about what the user was doing
 * - Monitor feature usage patterns
 * - Provide debugging context for errors
 */
export function useSentryTracking() {
  const location = useLocation();

  /**
   * Track a user action with context
   */
  const trackAction = useCallback((
    action: string,
    context?: Record<string, unknown>
  ) => {
    Sentry.addBreadcrumb({
      category: 'user_action',
      message: `User performed: ${action}`,
      level: 'info',
      data: {
        action,
        currentPage: location.pathname,
        timestamp: new Date().toISOString(),
        ...context,
      },
    });
  }, [location.pathname]);

  /**
   * Track form submissions with form data context (legacy version)
   */
  const trackFormSubmissionLegacy = useCallback((
    formName: string,
    formData?: Record<string, unknown>,
    success = true
  ) => {
    Sentry.addBreadcrumb({
      category: 'form',
      message: `Form ${success ? 'submitted' : 'failed'}: ${formName}`,
      level: success ? 'info' : 'error',
      data: {
        formName,
        success,
        currentPage: location.pathname,
        formData: formData ? JSON.stringify(formData) : undefined,
        timestamp: new Date().toISOString(),
      },
    });
  }, [location.pathname]);

  /**
   * Track search queries with search term
   */
  const trackSearch = useCallback((
    searchTerm: string,
    searchType = 'general',
    resultsCount?: number
  ) => {
    Sentry.addBreadcrumb({
      category: 'search',
      message: `Search performed: ${searchType}`,
      level: 'info',
      data: {
        searchTerm,
        searchType,
        resultsCount,
        currentPage: location.pathname,
        timestamp: new Date().toISOString(),
      },
    });
  }, [location.pathname]);

  /**
   * Track navigation between pages
   */
  const trackNavigation = useCallback((
    fromPage: string,
    toPage: string,
    navigationType = 'click'
  ) => {
    Sentry.addBreadcrumb({
      category: 'navigation',
      message: `Navigation: ${fromPage} → ${toPage}`,
      level: 'info',
      data: {
        fromPage,
        toPage,
        navigationType,
        timestamp: new Date().toISOString(),
      },
    });
  }, []);

  /**
   * Track feature usage
   */
  const trackFeatureUsage = useCallback((
    feature: string,
    action: string,
    context?: Record<string, unknown>
  ) => {
    Sentry.addBreadcrumb({
      category: 'feature_usage',
      message: `Feature used: ${feature} - ${action}`,
      level: 'info',
      data: {
        feature,
        action,
        currentPage: location.pathname,
        timestamp: new Date().toISOString(),
        ...context,
      },
    });
  }, [location.pathname]);

  /**
   * Track form submissions with performance monitoring
   */
  const trackFormSubmission = useCallback((
    formName: string,
    formData?: Record<string, unknown>,
    success = true,
    duration?: number
  ) => {
    // Track form submission performance
    if (duration && duration > 1000) {
      Sentry.addBreadcrumb({
        category: 'performance',
        message: `Slow form submission: ${formName} took ${duration}ms`,
        level: 'warning',
        data: {
          formName,
          duration,
          threshold: 1000,
        },
      });
    }

    Sentry.addBreadcrumb({
      category: 'form',
      message: `Form ${success ? 'submitted' : 'failed'}: ${formName}`,
      level: success ? 'info' : 'error',
      data: {
        formName,
        success,
        currentPage: location.pathname,
        formData: formData ? JSON.stringify(formData) : undefined,
        duration,
        timestamp: new Date().toISOString(),
      },
    });
  }, [location.pathname]);

  /**
   * Track data fetching performance
   */
  const trackDataFetching = useCallback((
    operation: string,
    duration: number,
    success = true,
    context?: Record<string, unknown>
  ) => {
    // Track slow data fetching
    if (duration > 1000) {
      Sentry.addBreadcrumb({
        category: 'performance',
        message: `Slow data fetch: ${operation} took ${duration}ms`,
        level: 'warning',
        data: {
          operation,
          duration,
          threshold: 1000,
          ...context,
        },
      });
    }

    Sentry.addBreadcrumb({
      category: 'data_fetching',
      message: `Data ${success ? 'fetched' : 'failed'}: ${operation}`,
      level: success ? 'info' : 'error',
      data: {
        operation,
        success,
        duration,
        currentPage: location.pathname,
        timestamp: new Date().toISOString(),
        ...context,
      },
    });
  }, [location.pathname]);

  /**
   * Track user interactions with performance monitoring
   */
  const trackUserInteraction = useCallback((
    interaction: string,
    duration?: number,
    context?: Record<string, unknown>
  ) => {
    // Track slow interactions
    if (duration && duration > 500) {
      Sentry.addBreadcrumb({
        category: 'performance',
        message: `Slow interaction: ${interaction} took ${duration}ms`,
        level: 'warning',
        data: {
          interaction,
          duration,
          threshold: 500,
          ...context,
        },
      });
    }

    Sentry.addBreadcrumb({
      category: 'user_interaction',
      message: `User interaction: ${interaction}`,
      level: 'info',
      data: {
        interaction,
        duration,
        currentPage: location.pathname,
        timestamp: new Date().toISOString(),
        ...context,
      },
    });
  }, [location.pathname]);

  /**
   * Track errors with user context
   */
  const trackError = useCallback((
    error: Error,
    context?: Record<string, unknown>
  ) => {
    Sentry.setContext('user_action', {
      currentPage: location.pathname,
      timestamp: new Date().toISOString(),
      ...context,
    });

    Sentry.captureException(error, {
      tags: {
        type: 'user_action_error',
        page: location.pathname,
      },
    });
  }, [location.pathname]);

  return {
    trackAction,
    trackFormSubmission,
    trackFormSubmissionLegacy,
    trackDataFetching,
    trackUserInteraction,
    trackSearch,
    trackNavigation,
    trackFeatureUsage,
    trackError,
  };
}