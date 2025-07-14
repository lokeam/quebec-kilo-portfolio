import { AxiosError } from 'axios';
import { useNavigate } from 'react-router-dom';

// Components
import { ErrorPage } from '@/core/error/components/ErrorPage';

// Types
import type { DefaultErrorFallbackProps } from '../types/error.types';

// Utils
import { isApiError, getErrorMessage } from '@/core/error/utils/error.utils';

// Constants
import { ERROR_ROUTES } from '@/core/error/constants/error.constants';


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
    if (error.response?.status === 401) {
      const errorMessage = error.response?.data?.message || '';
      const isTokenRefreshIssue = errorMessage.includes('token') ||
                                 errorMessage.includes('expired') ||
                                 errorMessage.includes('refresh');
      if (isTokenRefreshIssue) {
        return (
          <ErrorPage
            variant="500"
            title="Session expired"
            subtext="Your session has expired. Please try again."
            buttonText="Retry"
            onButtonClick={resetErrorBoundary}
            role="alert"
            ariaLive="assertive"
          />
        );
      } else {
        navigate(ERROR_ROUTES.LOGIN, { replace: true });
        return (
          <ErrorPage
            variant="500"
            title="Redirecting to login..."
            role="alert"
            ariaLive="assertive"
          />
        );
      }
    }

    if (isApiError(error)) {
      return (
        <ErrorPage
          variant="500"
          title={getErrorMessage(error.statusCode)}
          buttonText="Try Again"
          onButtonClick={resetErrorBoundary}
          role="alert"
          ariaLive="assertive"
        />
      );
    }
  }

  // Runtime errors -> Generic message with retry
  return (
    <ErrorPage
      variant="500"
      title="Something went wrong"
      buttonText="Try Again"
      onButtonClick={resetErrorBoundary}
      role="alert"
      ariaLive="assertive"
    />
  );
};
