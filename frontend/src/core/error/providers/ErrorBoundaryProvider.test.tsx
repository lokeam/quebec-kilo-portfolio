import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ErrorBoundaryProvider } from '../providers/ErrorBoundaryProvider';

describe('ErrorBoundaryProvider', () => {
  // Test component that can trigger errors
  const ThrowError = ({ shouldThrow = false}) => {
    if (shouldThrow) {
      throw new Error('Test error');
    }

    return <div>Normal render</div>
  };

  // Test setup
  const defaultConfig = {
    message: 'Custom Error Message',
    severity: 'fatal' as const,
    actionLabel: 'Custom Action'
  };

  describe('normal rendering', () => {
    it('renders children when no error occurs', () => {
      render(
        <ErrorBoundaryProvider>
          <div>Test Content</div>
        </ErrorBoundaryProvider>
      );

      expect(screen.getByText('Test Content')).toBeInTheDocument();
    });

    it('accepts and passes config to error page', () => {
      render(
        <ErrorBoundaryProvider config={defaultConfig}>
          <ThrowError shouldThrow />
        </ErrorBoundaryProvider>
      );

      expect(screen.getByText('Custom Error Message')).toBeInTheDocument();
      expect(screen.getByText('Custom Action')).toBeInTheDocument();
    });
  });

  describe('error handling', () => {
    beforeEach(() => {
      // Silence console.error for expected errors
      vi.spyOn(console, 'error').mockImplementation(() => {});
    });

    afterEach(() => {
      vi.restoreAllMocks();
    });

    it('catches and handles runtime errors', () => {
      render(
        <ErrorBoundaryProvider>
          <ThrowError shouldThrow />
        </ErrorBoundaryProvider>
      );

      expect(screen.getByText(/something went wrong/i)).toBeInTheDocument();
    });

    it('catches errors in nested components', () => {
      render(
        <ErrorBoundaryProvider>
          <div>
            <div>
              <ThrowError shouldThrow />
            </div>
          </div>
        </ErrorBoundaryProvider>
      );

      expect(screen.getByText(/something went wrong/i)).toBeInTheDocument();
    });

    it('logs errors to console', () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      render(
        <ErrorBoundaryProvider>
          <ThrowError shouldThrow />
        </ErrorBoundaryProvider>
      );

      // Only check for our specific error log
      expect(consoleSpy).toHaveBeenCalledWith(
        '[ErrorBoundary]:',
        expect.any(Error)
      );

      consoleSpy.mockRestore();
    });
  });

  describe('error recovery', () => {
    it('handles reset functionality', () => {
      const onReset = vi.fn();
      const onAction = vi.fn();

      render(
        <ErrorBoundaryProvider
        config={{
          onReset,
          onAction,
          actionLabel: 'Try Again'  // Ensure button is properly labeled
        }}
        >
          <ThrowError shouldThrow />
        </ErrorBoundaryProvider>
      );

      // Verify error state
      expect(screen.getByText('Something went wrong')).toBeInTheDocument();

      // Click reset button using default label
      fireEvent.click(screen.getByRole('button', { name: /try again/i }));

      // Verify action was called (this should happen first)
      expect(onAction).toHaveBeenCalled();

      // onReset is called via resetErrorBoundary
      expect(onReset).toHaveBeenCalled();
    });

    it('displays error message', () => {
      render(
        <ErrorBoundaryProvider>
          <ThrowError shouldThrow />
        </ErrorBoundaryProvider>
      );

      expect(screen.getByText('Test error')).toBeInTheDocument();
      expect(screen.getByRole('button')).toHaveTextContent('Try Again');
    });
  });

  describe('configuration', () => {
    it('accepts custom error handlers', () => {
      const customErrorHandler = vi.fn();

      render(
        <ErrorBoundaryProvider config={{ onError: customErrorHandler }}>
          <ThrowError shouldThrow />
        </ErrorBoundaryProvider>
      );

      expect(customErrorHandler).toHaveBeenCalled();
    })
  });

  describe('edge cases', () => {
    it('handles undefined children', () => {
      render(
        <ErrorBoundaryProvider>
          {undefined}
        </ErrorBoundaryProvider>
      )

      expect(true).toBeTruthy();
    });

    it('handles errors during error boundary reset', () => {
      const { rerender } = render(
        <ErrorBoundaryProvider>
          <ThrowError shouldThrow />
        </ErrorBoundaryProvider>
      );

      fireEvent.click(screen.getByRole('button'));

      rerender(
        <ErrorBoundaryProvider>
          <ThrowError shouldThrow />
        </ErrorBoundaryProvider>
      );

      expect(screen.getByText(/something went wrong/i)).toBeInTheDocument();
    });
  });
})