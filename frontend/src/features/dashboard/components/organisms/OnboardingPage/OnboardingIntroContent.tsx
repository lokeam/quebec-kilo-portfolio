

export const OnboardingIntroContent = () => (
  <div className="relative">
  {/* Step connector lines */}
  <div className="absolute top-5 left-[calc(25%-2rem)] right-[calc(25%-2rem)] h-[2px] bg-primary">
    <div className="absolute inset-0 bg-primary" />
  </div>

  {/* Steps */}
  <div className="relative grid grid-cols-3 gap-4">
    <div className="flex flex-col items-center gap-4">
      <div className="relative z-10 flex h-10 w-10 items-center justify-center rounded-full border-2 border-primary bg-background text-primary font-medium">
        1
      </div>
      <p className="text-center text-sm">
        Tell us a bit about your game library so we can customize Q-Ko to you
      </p>
    </div>

    <div className="flex flex-col items-center gap-4">
      <div className="relative z-10 flex h-10 w-10 items-center justify-center rounded-full border-2 border-primary bg-background text-primary font-medium">
        2
      </div>
      <p className="text-center text-sm">
        Based on your answers, we&apos;ll show you how to track everything you own
      </p>
    </div>

    <div className="flex flex-col items-center gap-4">
      <div className="relative z-10 flex h-10 w-10 items-center justify-center rounded-full border-2 border-primary bg-background text-primary font-medium">
        3
      </div>
      <p className="text-center text-sm">
        Choose what games or hardware you&apos;d like us to track deals for and we&apos;ll guide you the rest
        of the way
      </p>
    </div>
  </div>
  </div>
);

OnboardingIntroContent.displayName = 'OnboardingIntroContent';
