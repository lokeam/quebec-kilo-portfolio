export const FADE_IN_ANIMATION_SETTINGS = {
  initial: { opacity: 0 },
  animate: { opacity: 1 },
  exit: { opacity: 0 },
  transition: { duration: 0.2 },
};

export const STAGGER_CHILD_VARIANTS = {
  hidden: { opacity: 0, y: 20 },
  show: { opacity: 1, y: 0, transition: { duration: 0.4, type: "spring" } },
};

export const QKO_HEADERS = {
  headers: {
    "x-powered-by": "Q-KO.com - video game management platform",
  },
};

export const NAVIGATION_ROUTES = {
  ONBOARDING_INTRO: '/onboarding/intro',
  ONBOARDING_SELECT_STORAGE: '/onboarding/locations',
  ONBOARDING_SELECT_PHYSICAL_MAIN: '/onboarding/physical/location',
  ONBOARDING_SELECT_PHYSICAL_SUB: '/onboarding/physical/sublocation',
  ONBOARDING_SELECT_DIGITAL: '/onboarding/digital',
  ONBOARDING_COMPLETE: '/onboarding/complete',
} as const;
