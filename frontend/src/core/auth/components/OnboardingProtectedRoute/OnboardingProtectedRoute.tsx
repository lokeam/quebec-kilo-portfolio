import { Navigate } from 'react-router-dom';
import { useAuthStatus } from '@/core/auth/hooks/useAuthStatus';
import { getOnboardingDebugState, logDebugInfo } from '@/core/utils/debug/onboardingDebug';
import { LoadingPage } from '@/shared/components/ui/loading/LoadingPage';

interface OnboardingProtectedRouteProps {
  children: React.ReactNode;
}

export default function OnboardingProtectedRoute({ children }: OnboardingProtectedRouteProps) {
  const { isLoading, isAuthenticated, hasCompletedOnboarding } = useAuthStatus();
  const debugState = getOnboardingDebugState();

  console.log(
    'OnboardingProtectedRoute:',
    'isAuthenticated =', isAuthenticated,
    'isLoading =', isLoading,
    'hasCompletedOnboarding =', hasCompletedOnboarding
  );

  // Show loading while auth status is being determined
  if (isLoading) {
    return (
      <div className="min-h-screen bg-background">
        <LoadingPage />
      </div>
    );
  }

  // Redirect to login if not authenticated
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // Debug: Force incomplete onboarding
  if (debugState.forceIncompleteOnboarding) {
    logDebugInfo('OnboardingProtectedRoute', 'Debug: Forcing incomplete onboarding - allowing access');
    return <>{children}</>;
  }

  // Debug: Bypass onboarding entirely
  if (debugState.bypassOnboarding) {
    logDebugInfo('OnboardingProtectedRoute', 'Debug: Bypassing onboarding - redirecting to home');
    return <Navigate to="/" replace />;
  }

  // Debug: Simulate profile error
  if (debugState.simulateProfileError) {
    logDebugInfo('OnboardingProtectedRoute', 'Debug: Simulating profile error');
    return <>{children}</>;
  }

  // Check if user has completed onboarding
  if (!hasCompletedOnboarding) {
    logDebugInfo('OnboardingProtectedRoute', 'User has not completed onboarding - allowing access');
    return <>{children}</>;
  }

  // User has completed onboarding - redirect to home
  logDebugInfo('OnboardingProtectedRoute', 'User has completed onboarding - redirecting to home');
  return (
    <>
      <LoadingPage />
      <Navigate to="/" replace />
    </>
  );
}