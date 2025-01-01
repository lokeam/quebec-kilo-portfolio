import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import { DefaultErrorFallbackPage } from '@/core/error/pages/DefaultErrorFallbackPage';
import { ERROR_MESSAGES } from '@/core/error/constants/error.constants';
import type { ApiError, HttpErrorCode } from '@/core/error/types/error.types';
import { isApiError } from '../../utils/error.utils';

// Test factory for creating typed API errors
const createApiError = (statusCode: HttpErrorCode, message: string): ApiError => {
  const error = {
    message,
    statusCode,
    name: 'ApiError'
  };
  return error;
};

// Mock react-router-dom
vi.mock('react-router-dom', () => ({
  useNavigate: () => vi.fn(),
}));

// Mock ErrorButton component
vi.mock('../../components/ErrorButton', () => ({
  ErrorButton: ({
    label,
    onClick,
    variant
  }: {
    label: string;
    onClick: () => void;
    variant: string;
  }) => (
    <button onClick={onClick} data-testid={`error-button-${variant}`}>
      {label}
    </button>
  )
}));

describe('DefaultErrorFallback', () => {
  const resetErrorBoundary = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows generic error for non-API errors', () => {
    render(
      <DefaultErrorFallbackPage
        error={new Error('Test error')}
        resetErrorBoundary={resetErrorBoundary}
      />
    );

    expect(screen.getByText('Something went wrong')).toBeInTheDocument();
    expect(screen.getByTestId('error-button-retry')).toBeInTheDocument();
  });

  it('shows appropriate message for API errors', () => {
    const apiError = createApiError(404, 'Not Found');

    console.log('API Error:', apiError);
    console.log('isApiError result:', isApiError(apiError));

    render(
      <DefaultErrorFallbackPage
        error={apiError}
        resetErrorBoundary={resetErrorBoundary}
      />
    );

    expect(screen.getByText(ERROR_MESSAGES.NOT_FOUND)).toBeInTheDocument();
    expect(screen.getByTestId('error-button-home')).toBeInTheDocument();
  });

  it('calls resetErrorBoundary when retry button is clicked', () => {
    render(
      <DefaultErrorFallbackPage
        error={new Error('Test error')}
        resetErrorBoundary={resetErrorBoundary}
      />
    );

    fireEvent.click(screen.getByTestId('error-button-retry'));
    expect(resetErrorBoundary).toHaveBeenCalledTimes(1);
  });
});