import type { OnboardingStep } from '@/features/dashboard/lib/types/onboarding/base';
import type { OnboardingState } from '@/features/dashboard/lib/types/onboarding/choices';

/**
 * Determines whether the user can proceed to the next step based on their
 * current progress and data completeness.
 */
export const selectCanProceed = (state: OnboardingState): boolean => {
  switch (state.currentStep) {
    case 'WELCOME':
      // User can always proceed from welcome screen
      return true;

    case 'STORAGE':
      // Can only proceed if they've made a workspace choice
      return !!state.storageChoice;

    case 'PHYSICAL':
      // For physical storage, check if storage locations are defined
      if (
        state.storageChoice === 'PHYSICAL_STORAGE' ||
        state.storageChoice === 'PHYSICAL_AND_DIGITAL_STORAGE'
      ) {
        return !!state.physicalStorageData?.storageLocations?.length;
      }
      // Skip validation if this step isn't required
      return true;

    case 'DIGITAL':
      // For digital storage, check if digital storage data is complete
      if (
        state.storageChoice === 'DIGITAL_STORAGE' ||
        state.storageChoice === 'PHYSICAL_AND_DIGITAL_STORAGE'
      ) {
        return !!state.digitalStorageData?.storageLocations?.length;
      }
      // Skip validation if this step isn't required
      return true;

    case 'WISHLIST':
      // Require at least one wishlist item before completion
      return !!state.wishlistItems?.length;

    case 'COMPLETE':
      return false;

    default:
      return false;
  }
};

/**
 * Determines the next step in the onboarding flow based on the current
 * state and workspace choice. Returns null if there is no next step.
 */
export const selectNextStep = (state: OnboardingState): OnboardingStep | null => {
  switch (state.currentStep) {
    case 'WELCOME':
      // After welcome, always go to storage choice
      return 'STORAGE';

    case 'STORAGE':
      // After storage choice, route based on workspace selection
      switch (state.storageChoice) {
        case 'PHYSICAL_STORAGE':
          return 'PHYSICAL';
        case 'DIGITAL_STORAGE':
          return 'DIGITAL';
        case 'PHYSICAL_AND_DIGITAL_STORAGE':
          return 'PHYSICAL';
        default:
          return null;
      }

    case 'PHYSICAL':
      // After physical storage setup
      if (state.storageChoice === 'PHYSICAL_AND_DIGITAL_STORAGE') {
        // If they chose both, go to digital next
        return 'DIGITAL';
      }
      // Otherwise, proceed to wishlist
      return 'WISHLIST';

    case 'DIGITAL':
      // After digital storage, always go to wishlist
      return 'WISHLIST';

    case 'WISHLIST':
      // After wishlist, complete the onboarding
      return 'COMPLETE';

    case 'COMPLETE':
      // No next step after completion
      return null;

    default:
      return null;
  }
};

/**
 * Determines whether a step should be accessible based on the current
 * onboarding progress. Used to prevent users from skipping steps.
 */
export const selectIsStepAccessible = (
  state: OnboardingState,
  step: OnboardingStep
): boolean => {
  // Helper function to get numerical value for step ordering
  const getStepOrder = (step: OnboardingStep): number => {
    const stepOrder: Record<OnboardingStep, number> = {
      WELCOME: 0,
      STORAGE: 1,
      PHYSICAL: 2,
      DIGITAL: 3,
      WISHLIST: 4,
      COMPLETE: 5
    };
    return stepOrder[step];
  };

  // Get current step order
  const currentOrder = getStepOrder(state.currentStep);
  const targetOrder = getStepOrder(step);

  // Can't access steps beyond current progress
  if (targetOrder > currentOrder) {
    return false;
  }

  // Special handling for PHYSICAL and DIGITAL steps
  if (step === 'PHYSICAL') {
    return state.storageChoice === 'PHYSICAL_STORAGE' ||
           state.storageChoice === 'PHYSICAL_AND_DIGITAL_STORAGE';
  }

  if (step === 'DIGITAL') {
    return state.storageChoice === 'DIGITAL_STORAGE' ||
           state.storageChoice === 'PHYSICAL_AND_DIGITAL_STORAGE';
  }

  // All other previous steps are accessible
  return true;
};

/**
 * Returns the appropriate path for a given step, used for navigation.
 */
export const selectStepPath = (step: OnboardingStep): string => {
  const pathMap: Record<OnboardingStep, string> = {
    WELCOME: '/onboarding/welcome',
    STORAGE: '/onboarding/storage',
    PHYSICAL: '/onboarding/physical',
    DIGITAL: '/onboarding/digital',
    WISHLIST: '/onboarding/wishlist',
    COMPLETE: '/'
  };

  return pathMap[step];
};
