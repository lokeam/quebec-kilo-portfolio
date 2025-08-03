import { useLocation } from 'react-router-dom';
import { useEffect, useState } from 'react';

import { motion } from 'framer-motion';

// Constants
import { NAVIGATION_ROUTES, ONBOARDING_PROGRESS } from '@/types/domain/onboarding';

export function OnboardingProgressBar() {
  const location = useLocation();
  const [shouldAnimate, setShouldAnimate] = useState(false);

  // Determine current progress based on route
  const getCurrentProgress = () => {
    switch (location.pathname) {
      case NAVIGATION_ROUTES.ONBOARDING_WELCOME:
        return ONBOARDING_PROGRESS.STEPS.WELCOME;
      case NAVIGATION_ROUTES.ONBOARDING_NAME:
        return ONBOARDING_PROGRESS.STEPS.NAME;
      case NAVIGATION_ROUTES.ONBOARDING_MESSAGES:
        return ONBOARDING_PROGRESS.STEPS.TOAST_SETUP;
      default:
        return 0;
    }
  };

  const currentProgress = getCurrentProgress();

  // NOTE: Need to handle animation timing based on current page
  useEffect(() => {
    setShouldAnimate(false);

    const timer = setTimeout(() => {
      setShouldAnimate(true);
    }, location.pathname === NAVIGATION_ROUTES.ONBOARDING_WELCOME ? ONBOARDING_PROGRESS.WELCOME_ANIMATION_DELAY * 1000 : 100);

    return () => clearTimeout(timer);
  }, [location.pathname]);

  return (
    <div className="fixed top-0 left-0 right-0 z-50 h-4 bg-background/80 backdrop-blur-sm">
      <motion.div
        className="h-full bg-primary"
        initial={{ scaleX: 0 }}
        animate={shouldAnimate ? { scaleX: currentProgress / 100 } : { scaleX: 0 }}
        transition={ONBOARDING_PROGRESS.ANIMATION}
        style={{ transformOrigin: 'left' }}
      />
    </div>
  );
}