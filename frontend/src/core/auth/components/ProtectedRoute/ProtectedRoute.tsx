import { useAuth } from '@/core/auth/hooks/useAuth';
import { Navigate } from 'react-router-dom';
import { Loading } from '@/shared/components/ui/loading/Loading';
import { getOnboardingDebugState, logDebugInfo } from '@/core/utils/debug/onboardingDebug';

interface ProtectedRouteProps {
  children: JSX.Element;
};

function ProtectedRoute({ children }: ProtectedRouteProps) {
  const {
    isAuthenticated,
    isLoading,
    hasCompletedOnboarding,
    isProfileLoading
  } = useAuth();

  const debugState = getOnboardingDebugState();

  // Show loading while Auth0 is initializing or profile is loading
  if (isLoading || isProfileLoading) {
    return <Loading />;
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

  // Check if user needs to complete onboarding
  // This applies to:
  // 1. New users who haven't completed onboarding yet
  // 2. Existing users who somehow don't have firstName/lastName
  if (!hasCompletedOnboarding) {
    logDebugInfo('ProtectedRoute', 'User needs onboarding - redirecting to onboarding');
    return <Navigate to="/onboarding/welcome" replace />;
  }

  // User is authenticated and has completed onboarding, show the protected content
  return children;
}

export default ProtectedRoute;
