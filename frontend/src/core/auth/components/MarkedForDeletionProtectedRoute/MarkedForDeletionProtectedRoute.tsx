import { useAuthContext } from '@/core/auth/context-provider/AuthContext';
import { Navigate } from 'react-router-dom';
import { LoadingPage } from '@/shared/components/ui/loading/LoadingPage';
import { useGetUserDeletionStatus } from '@/core/api/queries/user.queries';

interface DeletionRouteProps {
  children: JSX.Element;
}

export default function MarkedForDeletionProtectedRoute({ children }: DeletionRouteProps) {
  const { isAuthenticated, isLoading: authLoading } = useAuthContext();
  const { data: deletionStatus, isLoading: deletionLoading } = useGetUserDeletionStatus();

  // Show loading while Auth0 is doing its thing
  if (authLoading || deletionLoading) {
    return (
      <div className="min-h-screen bg-background">
        <LoadingPage />
      </div>
    );
  }

  // Redirect to login if not authenticated
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  // Redirect to home if account is not marked for deletion
  if (!deletionStatus?.isDeletionRequested) {
    return <Navigate to="/" />;
  }

  // User is actually authenticated AND account is marked for deletion
  return children;
}