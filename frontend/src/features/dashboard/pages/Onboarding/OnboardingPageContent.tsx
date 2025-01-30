// Components
import { OnboardingWelcome } from '@/features/dashboard/components/organisms/OnboardingPage/OnboardingWelcome';

export function OnboardingPageContent() {
  return (
      <div className="mx-auto flex h-screen max-w-3xl flex-col items-center justify-center overflow-x-hidden">
        <OnboardingWelcome />
      </div>
  );
}
