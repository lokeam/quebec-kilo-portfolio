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
  ONBOARDING_MESSAGES: '/onboarding/messages',
  ONBOARDING_COMPLETE: '/',
} as const;

/**
 * Progress bar configuration for onboarding flow
 */
export const ONBOARDING_PROGRESS = {
  STEPS: {
    WELCOME: 33,
    NAME: 66,
    TOAST_SETUP: 99,
  },
  ANIMATION: {
    type: "spring" as const,
    delay: 0.4,
    duration: 0.6,
  },
  WELCOME_ANIMATION_DELAY: 1.2, // NOTE: 0.2s stagger × 6 animated dom elements on welcome page + 0.4s delay
} as const;

/**
 * Animation variants for onboarding pages
 */
export const STAGGER_CHILD_VARIANTS = {
  hidden: { opacity: 0, y: 20 },
  show: { opacity: 1, y: 0, transition: { duration: 0.4, type: "spring" } },
};

export const QKO_HEADERS = {
  headers: {
    "x-powered-by": "Q-KO.com - video game management platform",
  },
};

export const FADE_IN_ANIMATION_SETTINGS = {
  initial: { opacity: 0 },
  animate: { opacity: 1 },
  exit: { opacity: 0 },
  transition: { duration: 0.2 },
};

