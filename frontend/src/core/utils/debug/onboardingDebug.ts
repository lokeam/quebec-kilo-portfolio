/**
 * Onboarding Debug Utilities
 *
 * Provides debug toggles and utilities for developing onboarding flow
 * These should be removed in production
 */

interface OnboardingDebugState {
  /** Bypass onboarding flow entirely */
  bypassOnboarding: boolean;

  /** Force user to appear as new user */
  forceNewUser: boolean;

  /** Force user to appear as having incomplete onboarding */
  forceIncompleteOnboarding: boolean;

  /** Force user to appear as having completed onboarding */
  forceCompletedOnboarding: boolean;

  /** Simulate profile loading errors */
  simulateProfileError: boolean;

  /** Show debug info in console */
  showDebugInfo: boolean;
}

// Extend Window interface for debug state
declare global {
  interface Window {
    onboardingDebug: OnboardingDebugState;
  }
}

// Debug state - can be modified via browser console
window.onboardingDebug = {
  bypassOnboarding: false,
  forceNewUser: false,
  forceIncompleteOnboarding: false,
  forceCompletedOnboarding: false,
  simulateProfileError: false,
  showDebugInfo: true,
};

/**
 * Get current debug state
 */
export const getOnboardingDebugState = (): OnboardingDebugState => {
  return window.onboardingDebug || {};
};

/**
 * Check if debug mode is enabled
 */
export const isDebugMode = (): boolean => {
  const debugState = getOnboardingDebugState();
  return Object.values(debugState).some(Boolean);
};

/**
 * Get debug info for console logging
 */
export const getDebugInfo = () => {
  const debugState = getOnboardingDebugState();
  return {
    debugMode: isDebugMode(),
    debugState,
    instructions: `
🎛️ Onboarding Debug Controls:
• window.onboardingDebug.bypassOnboarding = true/false
• window.onboardingDebug.forceNewUser = true/false
• window.onboardingDebug.forceIncompleteOnboarding = true/false
• window.onboardingDebug.forceCompletedOnboarding = true/false
• window.onboardingDebug.simulateProfileError = true/false
• window.onboardingDebug.showDebugInfo = true/false
    `
  };
};

/**
 * Log debug info if enabled
 */
export const logDebugInfo = (context: string, data?: unknown) => {
  const debugState = getOnboardingDebugState();
  if (debugState.showDebugInfo) {
    console.log(`🔧 [ONBOARDING DEBUG] ${context}`, data);
  }
};