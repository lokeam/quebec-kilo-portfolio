/**
 * Debug Utilities for Development
 *
 * This file provides debugging tools for testing authentication, toasts, and optimistic state.
 * Available in development mode only.
 *
 * USAGE:
 * 1. Console Access: debugUtils.clearToastState()
 * 2. URL Parameters: http://localhost:3000/?debug=toasts
 * 3. Component Access: const debugUtils = useDebugUtils()
 *
 * @example
 * // Console access (pure functions only)
 * debugUtils.clearToastState()
 * debugUtils.testOptimisticIntroToasts()
 *
 * // Component access (all functions including hooks)
 * const debugUtils = useDebugUtils()
 * debugUtils.debugAuth0Claims()
 *
 * // URL parameter triggers
 * http://localhost:3000/?debug=toasts     // Clear toast state
 * http://localhost:3000/?debug=optimistic // Test optimistic toasts
 * http://localhost:3000/?debug=reset      // Reset optimistic state
 * http://localhost:3000/?debug=no-toasts  // Test no toasts
 * http://localhost:3000/?debug=claims     // Debug Auth0 claims
 * http://localhost:3000/?debug=state      // Debug current state
 */

import { useAuth0 } from '@auth0/auth0-react';
import { useEnableIntroToasts } from '@/features/dashboard/hooks/intro-toasts/useEnableIntroToasts';

// Define the debug utils type
interface DebugUtils {
  clearToastState: () => void;
  testOptimisticIntroToasts: () => void;
  resetOptimisticState: () => void;
  testNoToasts: () => void;
}

// Extend Window interface for optimistic state + console access
declare global {
  interface Window {
    __ONBOARDING_OPTIMISTIC_COMPLETE__?: boolean;
    __WANTS_INTRO_TOASTS__?: boolean;
    debugUtils?: DebugUtils;
  }
}

/**
 * Pure utility functions (no React hooks required)
 * These functions are available in both the browser console and via URL parameters
 */
export const debugUtils = {
  /**
   * Clears toast state from localStorage and reloads page
   * Useful for testing toast behavior after user has seen toasts
   *
   * @example
   * debugUtils.clearToastState()
   */
  clearToastState: () => {
    localStorage.removeItem('shownIntroToasts');
    console.log('�� Cleared shownIntroToasts from localStorage');
    window.location.reload();
  },

  /**
   * Tests optimistic intro toasts by setting optimistic state to true
   * Useful for testing toast behavior during onboarding
   *
   * @example
   * debugUtils.testOptimisticIntroToasts()
   */
  testOptimisticIntroToasts: () => {
    console.log('�� Current optimistic states:', {
      onboardingComplete: window.__ONBOARDING_OPTIMISTIC_COMPLETE__,
      wantsIntroToasts: window.__WANTS_INTRO_TOASTS__
    });

    window.__WANTS_INTRO_TOASTS__ = true;
    localStorage.setItem('__WANTS_INTRO_TOASTS__', 'true');
    console.log('✅ Set optimistic wantsIntroToasts to true');
    window.location.reload();
  },

  /**
   * Resets all optimistic state and localStorage
   * Useful for clearing all debug state and starting fresh
   *
   * @example
   * debugUtils.resetOptimisticState()
   */
  resetOptimisticState: () => {
    delete window.__ONBOARDING_OPTIMISTIC_COMPLETE__;
    delete window.__WANTS_INTRO_TOASTS__;
    localStorage.removeItem('__ONBOARDING_OPTIMISTIC_COMPLETE__');
    localStorage.removeItem('__WANTS_INTRO_TOASTS__');
    console.log('🧹 Reset optimistic state');
    window.location.reload();
  },

  /**
   * Tests "no toasts" scenario by setting optimistic state to false
   * Useful for testing behavior when user has disabled intro toasts
   *
   * @example
   * debugUtils.testNoToasts()
   */
  testNoToasts: () => {
    delete window.__ONBOARDING_OPTIMISTIC_COMPLETE__;
    delete window.__WANTS_INTRO_TOASTS__;
    localStorage.removeItem('__ONBOARDING_OPTIMISTIC_COMPLETE__');
    localStorage.removeItem('__WANTS_INTRO_TOASTS__');
    localStorage.removeItem('shownIntroToasts');

    window.__ONBOARDING_OPTIMISTIC_COMPLETE__ = true;
    window.__WANTS_INTRO_TOASTS__ = false;
    localStorage.setItem('__ONBOARDING_OPTIMISTIC_COMPLETE__', 'true');
    localStorage.setItem('__WANTS_INTRO_TOASTS__', 'false');
    window.location.reload();
  }
};

/**
 * Custom hook for debug functions that require React hooks
 *
 * This hook provides access to all debug utilities, including those that need
 * React hooks like useAuth0 and useEnableIntroToasts.
 *
 * @returns {Object} All debug utilities including hook-dependent functions
 *
 * @example
 * ```typescript
 * function MyShinyComponent() {
 *   const debugUtils = useDebugUtils();
 *
 *   return (
 *     <button onClick={debugUtils.debugAuth0Claims}>
 *       Debug Auth0 Claims
 *     </button>
 *   );
 * }
 * ```
 */
export const useDebugUtils = () => {
  /**
   * Debug Auth0 claims and user metadata
   * Logs ID token claims, user app_metadata, and current toast preferences
   *
   * @example
   * debugUtils.debugAuth0Claims()
   */
  const { user, getIdTokenClaims } = useAuth0();
  const { wantsIntroToasts } = useEnableIntroToasts();

  const debugAuth0Claims = async () => {
    try {
      const claims = await getIdTokenClaims();
      console.log('�� Debug Auth0 Claims:', claims);
      console.log('🔍 User app_metadata:', user?.app_metadata);
      console.log('�� Current wantsIntroToasts value:', wantsIntroToasts);
    } catch (error) {
      console.error('Failed to get claims:', error);
    }
  };

  /**
   * Debug current application state
   * Logs optimistic state, localStorage values, + user preferences
   *
   * @example
   * debugUtils.debugCurrentState()
   */
  const debugCurrentState = () => {
    console.log('🔍 Current Debug State:', {
      optimisticOnboardingComplete: window.__ONBOARDING_OPTIMISTIC_COMPLETE__,
      optimisticWantsIntroToasts: window.__WANTS_INTRO_TOASTS__,
      localStorageOnboardingComplete: localStorage.getItem('__ONBOARDING_OPTIMISTIC_COMPLETE__'),
      localStorageWantsIntroToasts: localStorage.getItem('__WANTS_INTRO_TOASTS__'),
      localStorageShownToasts: localStorage.getItem('shownIntroToasts'),
      userAppMetadata: user?.app_metadata,
      wantsIntroToasts: wantsIntroToasts
    });
  };

  return {
    ...debugUtils,
    debugAuth0Claims,
    debugCurrentState
  };
};

// Make pure functions available in browser console for development
if (process.env.NODE_ENV === 'development') {
  window.debugUtils = debugUtils;
}

// Enable URL parameter triggers
if (process.env.NODE_ENV === 'development') {
  const urlParams = new URLSearchParams(window.location.search);
  const debugMode = urlParams.get('debug');

  if (debugMode) {
    switch (debugMode) {
      case 'toasts':
        debugUtils.clearToastState();
        break;
      case 'optimistic':
        debugUtils.testOptimisticIntroToasts();
        break;
      case 'reset':
        debugUtils.resetOptimisticState();
        break;
      case 'no-toasts':
        debugUtils.testNoToasts();
        break;
    }
  }
}
