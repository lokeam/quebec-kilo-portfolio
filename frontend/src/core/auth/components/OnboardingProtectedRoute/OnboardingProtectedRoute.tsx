import { Navigate } from 'react-router-dom';
import { useAuthContext } from '@/core/auth/context-provider/AuthContext';
import { useOnboardingStatus } from '@/core/auth/hooks/useOnboardingStatus';
import { getOnboardingDebugState, logDebugInfo } from '@/core/utils/debug/onboardingDebug';
import { LoadingPage } from '@/shared/components/ui/loading/LoadingPage';

interface OnboardingProtectedRouteProps {
  children: React.ReactNode;
}

export default function OnboardingProtectedRoute({ children }: OnboardingProtectedRouteProps) {
  const { isAuthenticated, isLoading: authLoading } = useAuthContext(); // Only get auth data
  const { hasCompletedOnboarding, isLoading: onboardingLoading } = useOnboardingStatus(); // Safe, no API calls
  const debugState = getOnboardingDebugState();

  console.log(
    'OnboardingProtectedRoute:',
    'isAuthenticated =', isAuthenticated,
    'authLoading =', authLoading,
    'hasCompletedOnboarding =', hasCompletedOnboarding
  );

  // Show loading while Auth0 is loading
  if (authLoading || onboardingLoading) {
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

  // Check if user has completed onboarding (based on Auth0 user data, no API calls)
  if (!hasCompletedOnboarding) {
    logDebugInfo('OnboardingProtectedRoute', 'User has not completed onboarding - allowing access');
    return <>{children}</>;
  }

  // User has completed onboarding - redirect to home
  logDebugInfo('OnboardingProtectedRoute', 'User has completed onboarding - redirecting to home');
  return <Navigate to="/" replace />;
}