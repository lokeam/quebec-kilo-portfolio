/**
 * Test file for Sentry Feedback Service
 *
 * This file demonstrates how the feedback service works and can be used
 * for manual testing or as documentation.
 */

import { submitFeedback, submitBugReport, submitFeatureRequest, submitGeneralFeedback } from './sentry-feedback-service';

// Example usage - this shows how the service works
export function testFeedbackService() {
  console.log('🧪 Testing Sentry Feedback Service...');

  // Test 1: Submit general feedback
  console.log('📝 Testing general feedback submission...');
  submitGeneralFeedback(
    'The dashboard is really helpful! I love how easy it is to track my spending.',
    'excellent'
  );

  // Test 2: Submit bug report
  console.log('🐛 Testing bug report submission...');
  submitBugReport(
    'The search function is not working properly on mobile devices',
    {
      errorMessage: 'Search results not loading',
      userSteps: '1. Open app on mobile\n2. Try to search for a game\n3. Results don\'t appear',
    }
  );

  // Test 3: Submit feature request
  console.log('💡 Testing feature request submission...');
  submitFeatureRequest(
    'It would be great to have a dark mode option for the app',
    'medium'
  );

  // Test 4: Submit custom feedback
  console.log('📋 Testing custom feedback submission...');
  submitFeedback({
    type: 'other',
    message: 'Just wanted to say thanks for building such a useful app!',
    currentPage: '/dashboard',
    additionalContext: {
      userAgent: navigator.userAgent,
      screenResolution: `${screen.width}x${screen.height}`,
    },
  });

  console.log('✅ Feedback service tests completed!');
  console.log('📊 Check your Sentry dashboard for the feedback messages.');
}

// Export for manual testing
export { testFeedbackService as demonstrateFeedbackService };