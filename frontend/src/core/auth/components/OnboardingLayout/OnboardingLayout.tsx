import { Outlet } from 'react-router-dom';
import { OfflineBanner } from '@/core/network-status/components/OfflineBanner';

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
      {/* Network Status Provider - Offline Banner */}
      <OfflineBanner />

      {/* Main content area */}
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
}