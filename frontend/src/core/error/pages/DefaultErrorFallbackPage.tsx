import { useNavigate } from 'react-router-dom';
import type { DefaultErrorFallbackProps } from '../types/error.types';
import { ErrorButton } from '../components/ErrorButton';
import { isApiError, getErrorMessage } from '../utils/error.utils';
import { ERROR_ROUTES } from '../constants/error.constants';

export const DefaultErrorFallbackPage = ({ error, resetErrorBoundary }: DefaultErrorFallbackProps) => {
  const navigate = useNavigate();

  // Redirect to login for unauthorized errors
  if(isApiError(error) && error.statusCode === 401) {
    navigate(ERROR_ROUTES.LOGIN, { replace: true });

    // Prevent flash of unstyled error content
    return null;
  }

  if (!isApiError(error)) {
    return (
      <div role="alert" data-testid="error-fallback">
        <h2>Something went wrong</h2>
        <ErrorButton
          variant="retry"
          onClick={resetErrorBoundary}
          label="Try Again"
        />
      </div>
    );
  }

  return (
    <div role="alert" data-testid="error-fallback">
      <h2>{getErrorMessage(error.statusCode)}</h2>
      <ErrorButton
        variant="home"
        onClick={() => {}}
        label="Return Home"
      />
    </div>
  );
};
