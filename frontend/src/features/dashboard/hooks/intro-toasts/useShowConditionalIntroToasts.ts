import { useEffect } from 'react';
import { INTRO_TOASTS } from './intro-toasts-constants';
import { useEnableIntroToasts } from './useEnableIntroToasts';
import { showIntroToast } from './IntroToast';

/**
 * Shows an intro toast with optional conditional logic.
 * Only shows once per user (tracked in localStorage) and only if user has enabled intro toasts during onboarding.
 * Debug mode can be enabled via URL parameter ?debugToasts=true or localStorage debugIntroToasts=true
 *
 * @param toastId - ID of the toast to show
 * @param condition - Optional boolean condition that triggers the toast (defaults to true)
 *
 * @example
 * // Always show toast 1 (no condition)
 * useShowConditionalIntroToasts(1);
 *
 * // Show toast 3 when user has physical locations
 * useShowConditionalIntroToasts(3, hasPhysicalLocations);
 *
 * // Show toast 7 when user has games in library
 * useShowConditionalIntroToasts(7, hasGamesInLibrary);
 */
export function useShowConditionalIntroToasts(toastId: number, condition: boolean = true) {
  const { wantsIntroToasts } = useEnableIntroToasts();

  useEffect(() => {
    // Only show toast when condition is true
    if (!condition) return;

    // Check for debug mode
    const urlParams = new URLSearchParams(window.location.search);
    const debugMode = urlParams.get('debugToasts') === 'true' ||
                     localStorage.getItem('debugIntroToasts') === 'true';

    // Don't show toasts if user has disabled them (unless in debug mode)
    if (!wantsIntroToasts && !debugMode) return;

    const toast = INTRO_TOASTS.find(t => t.id === toastId);
    if (!toast) return;

    // In debug mode, always show the toast
    if (debugMode) {
      showIntroToast({
        title: toast.title,
        description: `${toast.message}`,
      });
      return;
    }

    // Check if toast has already been shown to this user
    const shownToasts = JSON.parse(localStorage.getItem('shownIntroToasts') || '[]');
    if (shownToasts.includes(toastId)) return;

    // Show the toast
    showIntroToast({
      title: toast.title,
      description: toast.message,
    });

    // Mark this toast as shown
    shownToasts.push(toastId);
    localStorage.setItem('shownIntroToasts', JSON.stringify(shownToasts));
  }, [condition, toastId, wantsIntroToasts]);
}