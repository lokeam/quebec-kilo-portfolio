import { useAuth } from '@/core/auth/hooks/useAuth';
import { Navigate } from 'react-router-dom';
import { Loading } from '@/shared/components/ui/loading/Loading';

interface ProtectedRouteProps {
  children: JSX.Element;
};

function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading } = useAuth();

  // Show loading while Auth0 is initializing
  if (isLoading) {
    return <Loading />;
  }

  // Redirect to login if not authenticated
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // User is authenticated, show the protected content
  return children;
}

export default ProtectedRoute;
