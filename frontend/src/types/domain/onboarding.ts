/**
 * Onboarding Types
 *
 * Types for the onboarding flow and storage selection
 */

/**
 * Storage type options for onboarding
 */
export const OnboardingStorageType = {
  PHYSICAL: 'physical',
  DIGITAL: 'digital',
  BOTH: 'both'
} as const;

export type OnboardingStorageType = typeof OnboardingStorageType[keyof typeof OnboardingStorageType];

/**
 * Navigation routes for onboarding flow
 */
export const NAVIGATION_ROUTES = {
  ONBOARDING_WELCOME: '/onboarding/welcome',
  ONBOARDING_NAME: '/onboarding/name',
  ONBOARDING_INTRO: '/onboarding/intro',
  ONBOARDING_COMPLETE: '/',
} as const;

/**
 * Animation variants for onboarding pages
 */
export const STAGGER_CHILD_VARIANTS = {
  hidden: { opacity: 0, y: 20 },
  show: { opacity: 1, y: 0, transition: { duration: 0.4, type: "spring" } },
};