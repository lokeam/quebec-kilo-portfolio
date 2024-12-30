interface ProtectedRouteProps {
  children: JSX.Element;
};

function ProtectedRoute({ children }: ProtectedRouteProps) {
  // TODO: Add Auth0 middleware

  // TODO: Add Loading component

  return children;
}

export default ProtectedRoute;
