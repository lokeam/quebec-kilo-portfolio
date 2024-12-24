import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ErrorFallbackPage } from '@/core/error/components/ErrorFallbackPage';
import '@testing-library/jest-dom';

describe('ErrorFallbackPage', () => {
  // Setup
  const defaultProps = {
    error: new Error('Test error'),
    resetErrorBoundary: vi.fn(),
  };

  // Rendering tests
  describe('rendering', () => {
    it('renders default error state without config', () => {
      render(<ErrorFallbackPage {...defaultProps} />);

      expect(screen.getByText('Oops!')).toBeInTheDocument();
      expect(screen.getByText('Something went wrong')).toBeInTheDocument();
      expect(screen.getByText('Test error')).toBeInTheDocument();
    });

    it ('renders fatal error state when specified', () => {
      render(
      <ErrorFallbackPage
        {...defaultProps}
        config={{ severity: 'fatal' }}
      />
    );

    expect(screen.getByText('Critical Error')).toBeInTheDocument();
    });

    it('displays custom error message when provided', () => {
      const customMessage = 'Custom error message';
      render(
        <ErrorFallbackPage
          {...defaultProps}
          config={{ message: customMessage }}
        />
      );

      expect(screen.getByText(customMessage)).toBeInTheDocument();
    });

    it('shows actual error message when available', () => {
      const errorMessage = 'Specific error details';
      render(
        <ErrorFallbackPage
          {...defaultProps}
          error={new Error(errorMessage)}
        />
      );

      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });
  });

  // Interaction tests
  describe('interactions', () => {
    it('calls resetErrorBoundary when default button is clicked', () => {
      render(<ErrorFallbackPage {...defaultProps} />);

      fireEvent.click(screen.getByRole('button'));
      expect(defaultProps.resetErrorBoundary).toHaveBeenCalledTimes(1);
    });

    it('calls custom action when provided', () => {
      const customAction = vi.fn();
      render(
        <ErrorFallbackPage
          {...defaultProps}
          config={{ onAction: customAction }}
        />
      );

      fireEvent.click(screen.getByRole('button'));
      expect(customAction).toHaveBeenCalledTimes(1);
      expect(defaultProps.resetErrorBoundary).not.toHaveBeenCalled();
    });

    it('displays custom action label when provided', () => {
      const actionLabel = 'Custom Action';
      render(
        <ErrorFallbackPage
          {...defaultProps}
          config={{ actionLabel }}
        />
      );

      expect(screen.getByText(actionLabel)).toBeInTheDocument();
    })
  })

  // Logging tests
  describe('error logging', () => {
    it('logs error to the console', () => {
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      const error = new Error('Test error');

      render(<ErrorFallbackPage {...defaultProps} error={error} />);

      expect(consoleSpy).toHaveBeenCalledWith('[ErrorBoundary] Error:', error);
      consoleSpy.mockRestore();
    });
  });

  // Accessibility tests
  describe('accessibility', () => {
    it('maintains proper heading heirarchy', () => {
      render(<ErrorFallbackPage {...defaultProps} /> );
      const headings = screen.getAllByRole('heading');
      expect(headings[0]).toHaveAttribute('class', expect.stringContaining('h2'));
      expect(headings[1]).toHaveAttribute('class', expect.stringContaining('h5'));
    })

    it('ensures error button is keyboard accessible with Enter key press', () => {
      render(<ErrorFallbackPage {...defaultProps} />);

      const button = screen.getByRole('button');
      button.focus();
      expect(button).toHaveFocus();

      // Changed from keyPress to keyDown
      fireEvent.keyDown(button, { key: 'Enter', code: 'Enter' });
      fireEvent.keyUp(button, { key: 'Enter', code: 'Enter' });
      expect(defaultProps.resetErrorBoundary).toHaveBeenCalled();
    });

    // Additional keyboard tests
    it('ensures error button is keyboard accessible with Space key press', () => {
      render(<ErrorFallbackPage {...defaultProps} />);

      const button = screen.getByRole('button');
      button.focus();
      expect(button).toHaveFocus();

      // Material-UI requires both keyDown and keyUp for Space
      fireEvent.keyDown(button, { key: ' ', code: 'Space' });
      fireEvent.keyUp(button, { key: ' ', code: 'Space' });
      expect(defaultProps.resetErrorBoundary).toHaveBeenCalled();
    });
  });

  // Edge cases
  describe('edge cases', () => {
    it('handles undefined error message', () => {
      render(
        <ErrorFallbackPage
          {...defaultProps}
          error={new Error()}
        />
      );

      expect(screen.getByText('We apologize for the inconvenience. Please try again.')).toBeInTheDocument();
    });

    it('handles missing config values', () => {
      render(
        <ErrorFallbackPage
          {...defaultProps}
          config={{ message: undefined, actionLabel: undefined }}
        />
      );

      expect(screen.getByText('Something went wrong')).toBeInTheDocument();
    });
  });
});
