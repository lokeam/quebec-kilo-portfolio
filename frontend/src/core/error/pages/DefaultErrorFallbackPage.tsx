import { useNavigate } from 'react-router-dom';
import type { DefaultErrorFallbackProps } from '../types/error.types';
import { ErrorButton } from '../components/ErrorButton';
import { isApiError, getErrorMessage } from '../utils/error.utils';
import { ERROR_ROUTES } from '../constants/error.constants';
import { AxiosError } from 'axios';

/**
 * Default error fallback page that handles both runtime and query errors
 *
 * Error handling hierarchy:
 * 1. Auth errors (401) -> Redirect to login
 * 2. Query/API errors -> Show error message with retry
 * 3. Runtime errors -> Show generic error with retry
 *
 * @see https://tanstack.com/query/latest/docs/react/reference/QueryErrorResetBoundary
 */
export const DefaultErrorFallbackPage = ({
  error,
  resetErrorBoundary
}: DefaultErrorFallbackProps) => {
  const navigate = useNavigate();

  // Handle Axios/Query errors
  if (error instanceof AxiosError) {
    // Auth errors -> Redirect to login
    if (error.response?.status === 401) {
      navigate(ERROR_ROUTES.LOGIN, { replace: true });
      return null;
    }

    // API errors -> Show message with retry
    if (isApiError(error)) {
      return (
        <div role="alert" data-testid="error-fallback">
          <h2>{getErrorMessage(error.statusCode)}</h2>
          <ErrorButton
            variant="retry"
            onClick={resetErrorBoundary}
            label="Try Again"
          />
          <ErrorButton
            variant="home"
            onClick={() => navigate('/')}
            label="Return Home"
          />
        </div>
      );
    }
  }

  // Runtime errors -> Generic message with retry
  return (
    <div role="alert" data-testid="error-fallback">
      <h2>Something went wrong</h2>
      <ErrorButton
        variant="retry"
        onClick={resetErrorBoundary}
        label="Try Again"
      />
      <ErrorButton
        variant="home"
        onClick={() => navigate('/')}
        label="Return Home"
      />
    </div>
  );
};
