import { useAuthContext } from '@/core/auth/context-provider/AuthContext';
import { useOnboardingStatus } from '@/core/auth/hooks/useOnboardingStatus';
import { Navigate } from 'react-router-dom';
import { LoadingPage } from '@/shared/components/ui/loading/LoadingPage';
import { getOnboardingDebugState, logDebugInfo } from '@/core/utils/debug/onboardingDebug';

interface ProtectedRouteProps {
  children: JSX.Element;
};

function ProtectedRoute({ children }: ProtectedRouteProps) {
  const {
    isAuthenticated,
    isLoading,
  } = useAuthContext();

  const { hasCompletedOnboarding, isLoading: onboardingLoading } = useOnboardingStatus();

  const debugState = getOnboardingDebugState();

  // Show loading while Auth0 is initializing OR while onboarding status is being determined
  if (isLoading || onboardingLoading) {
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

  // Debug bypass: Skip onboarding check entirely
  if (debugState.bypassOnboarding) {
    logDebugInfo('ProtectedRoute', 'Debug bypass enabled - allowing access to protected routes');
    return children;
  }

  // CRITICAL FIX: Only proceed if we have a definitive onboarding status
  // This prevents the brief flash where hasCompletedOnboarding is false but we haven't determined it yet
  if (hasCompletedOnboarding === undefined) {
    // Still determining onboarding status, show loading
    return (
      <div className="min-h-screen bg-background">
        <LoadingPage />
      </div>
    );
  }

  // Check if user needs to complete onboarding
  // This applies to:
  // 1. New users who haven't completed onboarding yet
  // 2. Existing users who somehow don't have firstName/lastName
  if (!hasCompletedOnboarding) {
    logDebugInfo('ProtectedRoute', 'User needs onboarding - redirecting to onboarding');
    return (
      <>
        <LoadingPage />
        <Navigate to="/onboarding/welcome" replace />
      </>
    );
  }

  // User is authenticated and has completed onboarding, show the protected content
  return children;
}

export default ProtectedRoute;
