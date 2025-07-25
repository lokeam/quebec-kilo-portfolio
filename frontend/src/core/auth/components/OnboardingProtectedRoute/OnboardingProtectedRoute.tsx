import { Navigate } from 'react-router-dom';
import { useAuthStatus } from '@/core/auth/hooks/useAuthStatus';
import { LoadingPage } from '@/shared/components/ui/loading/LoadingPage';

interface OnboardingProtectedRouteProps {
  children: React.ReactNode;
}

export default function OnboardingProtectedRoute({ children }: OnboardingProtectedRouteProps) {
  const { isLoading, isAuthenticated, hasCompletedOnboarding } = useAuthStatus();

  // console.log(
  //   'OnboardingProtectedRoute:',
  //   'isAuthenticated =', isAuthenticated,
  //   'isLoading =', isLoading,
  //   'hasCompletedOnboarding =', hasCompletedOnboarding
  // );

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

  // Check if user has completed onboarding
  if (!hasCompletedOnboarding) {
    return <>{children}</>;
  }

  // User has completed onboarding - redirect to home
  return (
    <>
      <LoadingPage />
      <Navigate to="/" replace />
    </>
  );
}