/**
 * Simple test to verify our Sentry API monitoring service
 *
 * This file demonstrates how the service works and can be used
 * for manual testing or as documentation.
 */

import { logApiError, logAuthError, logApiSuccess, logSlowApiResponse, logApiBreadcrumb } from './sentry-api-monitor';

// Example usage - this shows how the service works
export function demonstrateApiMonitoring() {
  // Example 1: Log an API error
  const apiError = new Error('Failed to fetch users');
  logApiError(apiError, {
    url: '/api/users',
    method: 'GET',
    status: 500,
    errorMessage: 'Internal server error',
    responseData: { error: 'Database connection failed' }
  });

  // Example 2: Log an auth error
  const authError = new Error('Token refresh failed');
  logAuthError(authError, {
    originalUrl: '/api/protected-data',
    originalMethod: 'GET',
    originalStatus: 401,
    errorMessage: 'Failed to refresh authentication token'
  });

  // Example 3: Log a successful API call
  logApiSuccess({
    url: '/api/users',
    method: 'GET',
    status: 200,
    responseTime: 250
  });

  // Example 4: Log a slow API response
  logSlowApiResponse({
    url: '/api/heavy-operation',
    method: 'POST',
    responseTime: 6500, // 6.5 seconds
    threshold: 5000     // 5 second threshold
  });

  // Example 5: Add a breadcrumb for tracking what happened
  logApiBreadcrumb({
    url: '/api/upload',
    method: 'POST',
    message: 'File upload started',
    data: {
      fileName: 'document.pdf',
      fileSize: 2048576, // 2MB
      userId: 'user123'
    }
  });

  // Example 6: Add another breadcrumb showing progress
  logApiBreadcrumb({
    url: '/api/upload',
    method: 'POST',
    message: 'File upload 50% complete',
    data: {
      bytesUploaded: 1024288,
      totalBytes: 2048576
    }
  });
}

/**
 * Test that our functions can be called without errors
 */
export function testApiMonitoringFunctions() {
  try {
    // Test API error logging
    logApiError(new Error('Test error'), {
      url: '/test',
      method: 'GET',
      status: 404,
      errorMessage: 'Test error message'
    });

    // Test auth error logging
    logAuthError(new Error('Test auth error'), {
      originalUrl: '/test',
      originalMethod: 'GET',
      originalStatus: 401,
      errorMessage: 'Test auth error message'
    });

    // Test success logging
    logApiSuccess({
      url: '/test',
      method: 'GET',
      status: 200,
      responseTime: 100
    });

    // Test slow response logging
    logSlowApiResponse({
      url: '/test',
      method: 'GET',
      responseTime: 6000,
      threshold: 5000
    });

    // Test breadcrumb logging
    logApiBreadcrumb({
      url: '/test',
      method: 'POST',
      message: 'Test breadcrumb message',
      data: {
        testData: 'This is test data',
        timestamp: Date.now()
      }
    });

    console.log('✅ All Sentry API monitoring functions work correctly');
    return true;
  } catch (error) {
    console.error('❌ Error testing Sentry API monitoring:', error);
    return false;
  }
}