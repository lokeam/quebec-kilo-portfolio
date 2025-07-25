import { useAuthStatus } from '@/core/auth/hooks/useAuthStatus';
import { Navigate } from 'react-router-dom';
import { LoadingPage } from '@/shared/components/ui/loading/LoadingPage';

interface ProtectedRouteProps {
  children: JSX.Element;
};

function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isLoading, isAuthenticated, needsOnboarding } = useAuthStatus();

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

  // Redirect if needs onboarding
  if (needsOnboarding) {
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
