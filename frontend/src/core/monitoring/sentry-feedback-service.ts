import * as Sentry from '@sentry/react';

/**
 * Sentry Feedback Service
 *
 * Provides a simple, professional way for users to submit feedback
 * using Sentry's existing message capture capabilities.
 *
 * This approach:
 * - Uses Sentry's message capture for feedback
 * - Includes user context automatically
 * - Provides structured feedback data
 * - Integrates with existing error tracking
 * - No additional dependencies required
 */

export interface FeedbackData {
  type: 'bug' | 'feature' | 'general' | 'other';
  message: string;
  userExperience?: 'excellent' | 'good' | 'fair' | 'poor';
  currentPage?: string;
  additionalContext?: Record<string, unknown>;
}

/**
 * Submit user feedback to Sentry
 *
 * This creates a structured feedback message that appears in the Sentry dashboard.
 *
 */
export function submitFeedback(feedback: FeedbackData): void {
  // Create a structured feedback message
  const feedbackMessage = `[FEEDBACK] ${feedback.type.toUpperCase()}: ${feedback.message}`;

  // Add feedback context
  const feedbackContext = {
    feedback_type: feedback.type,
    user_experience: feedback.userExperience,
    current_page: feedback.currentPage || window.location.pathname,
    timestamp: new Date().toISOString(),
    ...feedback.additionalContext,
  };

  // Add tracking breadcrumb
  Sentry.addBreadcrumb({
    category: 'feedback',
    message: `User submitted ${feedback.type} feedback`,
    level: 'info',
    data: feedbackContext,
  });

  // Capture feedback as a msg in Sentry
  Sentry.captureMessage(feedbackMessage, {
    level: 'info', // Use info level for feedback (not errors)
    tags: {
      type: 'user_feedback',
      feedback_type: feedback.type,
      source: 'feedback_widget',
    },
    contexts: {
      feedback: feedbackContext,
    },
  });
}

/**
 * Submit bug report with error context
 *
 * Use this when users want to report a specific issue they encountered.
 */
export function submitBugReport(
  description: string,
  errorContext?: {
    errorMessage?: string;
    errorStack?: string;
    userSteps?: string;
  }
): void {
  const feedback: FeedbackData = {
    type: 'bug',
    message: description,
    currentPage: window.location.pathname,
    additionalContext: {
      error_message: errorContext?.errorMessage,
      error_stack: errorContext?.errorStack,
      user_steps: errorContext?.userSteps,
      browser_info: {
        userAgent: navigator.userAgent,
        url: window.location.href,
        timestamp: new Date().toISOString(),
      },
    },
  };

  submitFeedback(feedback);
}

/**
 * Submit feature request
 *
 * Use this when users want to request new features or improvements.
 */
export function submitFeatureRequest(
  description: string,
  priority?: 'low' | 'medium' | 'high'
): void {
  const feedback: FeedbackData = {
    type: 'feature',
    message: description,
    currentPage: window.location.pathname,
    additionalContext: {
      priority: priority || 'medium',
      feature_request: true,
    },
  };

  submitFeedback(feedback);
}

/**
 * Submit general feedback
 *
 * Use this for general comments, suggestions, or other feedback.
 */
export function submitGeneralFeedback(
  message: string,
  userExperience?: 'excellent' | 'good' | 'fair' | 'poor'
): void {
  const feedback: FeedbackData = {
    type: 'general',
    message,
    userExperience,
    currentPage: window.location.pathname,
  };

  submitFeedback(feedback);
}

/**
 * Get current page context for feedback
 */
export function getCurrentPageContext(): Record<string, unknown> {
  return {
    currentPage: window.location.pathname,
    currentUrl: window.location.href,
    timestamp: new Date().toISOString(),
    userAgent: navigator.userAgent,
  };
}