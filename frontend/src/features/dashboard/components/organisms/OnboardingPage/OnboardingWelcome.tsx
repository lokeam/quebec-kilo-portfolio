

import { Button, } from '@/shared/components/ui/button';

import { motion } from 'framer-motion';
import { STAGGER_CHILD_VARIANTS } from '@/types/domain/onboarding';
import { useNavigate } from 'react-router-dom';
import { NAVIGATION_ROUTES } from '@/types/domain/onboarding';

export function OnboardingWelcome() {
  const navigate = useNavigate();

  return (
    <motion.div
      className="z-10"
      exit={{ opacity: 0, scale: 0.95 }}
      transition={{ duration: 0.3, type: "spring" }}
    >
      <motion.div
        variants={{
          show: {
            transition: {
              staggerChildren: 0.2,
            },
          },
        }}
        initial="hidden"
        animate="show"
        className="mx-5 flex flex-col items-center space-y-10 text-center sm:mx-auto"
      >
        <motion.h1
          className="font-display text-4xl font-bold text-foreground transition-colors sm:text-5xl"
          variants={STAGGER_CHILD_VARIANTS}
        >
          Welcome to{" "}
          <span className="font-bold tracking-tighter">QKO</span> (<span className="text-2xl text-primary">Beta</span>)
        </motion.h1>
        <motion.p
          className="max-w-md text-accent-foreground/80 transition-colors sm:text-lg"
          variants={STAGGER_CHILD_VARIANTS}
        >
          Your personal game management assistant that helps:
        </motion.p>
        <motion.p
          className="max-w-md text-accent-foreground/80 transition-colors sm:text-lg"
          variants={STAGGER_CHILD_VARIANTS}
        >
          Keep track of your all your physical and digital games.
        </motion.p>
        <motion.p
          className="max-w-md text-accent-foreground/80 transition-colors sm:text-lg"
          variants={STAGGER_CHILD_VARIANTS}
        >
          Track recurring online service fees and one-time purchases.
        </motion.p>
        <motion.p
          className="max-w-md text-accent-foreground/80 transition-colors sm:text-lg"
          variants={STAGGER_CHILD_VARIANTS}
        >
          This app is still in <span className="font-bold text-primary">early access</span> release!
        </motion.p>
        <motion.p
          className="max-w-md text-accent-foreground/80 transition-colors sm:text-lg"
          variants={STAGGER_CHILD_VARIANTS}
        >
          Please report any bugs or issues you find!
        </motion.p>
        <motion.p
          className="max-w-md text-accent-foreground/80 transition-colors sm:text-lg"
          variants={STAGGER_CHILD_VARIANTS}
        >
          Thank you for your help!
        </motion.p>
        <motion.div
          variants={STAGGER_CHILD_VARIANTS}
          // className="rounded  px-10 py-2 font-medium transition-colors text-gray-900 bg-gray-100 hover:text-gray-100 hover:bg-gray-500"
        >
          <Button
            className="px-10 text-base font-medium"
            onClick={() => navigate(NAVIGATION_ROUTES.ONBOARDING_NAME)}
          >
            Get Started
          </Button>
        </motion.div>
      </motion.div>
    </motion.div>
  );
}