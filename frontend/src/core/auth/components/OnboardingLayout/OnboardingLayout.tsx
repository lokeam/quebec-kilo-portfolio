import { Outlet } from 'react-router-dom';
import { OfflineBanner } from '@/core/network-status/components/OfflineBanner';
import { OnboardingProgressBar } from '@/features/dashboard/components/organisms/OnboardingPage/OnboardingProgressBar';

/**
 * OnboardingLayout
 *
 * A simplified layout for onboarding pages that doesn't include
 * the full sidebar and navigation components. This provides a
 * clean, focused experience for new users completing onboarding.
 */
export default function OnboardingLayout() {
  return (
    <div className="min-h-screen bg-background">
      {/* Progress Bar */}
      <OnboardingProgressBar />

      {/* Network Status Provider - Offline Banner */}
      <OfflineBanner />

      {/* Main content area */}
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
}