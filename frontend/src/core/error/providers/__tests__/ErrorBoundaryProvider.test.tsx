import { render, screen } from '@testing-library/react';
import { ErrorBoundaryProvider } from '@/core/error/providers/ErrorBoundaryProvider';
import { vi } from 'vitest';

// Mock monitoring service
vi.mock('@/services/monitoring', () => ({
  sendToSentry: vi.fn(),
}));

// Mock router since DefaultErrorFallbackPage uses it
vi.mock('react-router-dom', () => ({
  useNavigate: () => vi.fn(),
}));

const ThrowError = () => {
  throw new Error('Test error');
};

describe('ErrorBoundaryProvider', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders fallback when child throws', () => {
    // Suppress console.error for expected test error
    const consoleSpy = vi.spyOn(console, 'error');
    consoleSpy.mockImplementation(() => {});

    render(
      <ErrorBoundaryProvider>
        <ThrowError />
      </ErrorBoundaryProvider>
    );

    expect(screen.getByTestId('error-fallback')).toBeInTheDocument();
    consoleSpy.mockRestore();
  });

  it('renders children when no error occurs', () => {
    render(
      <ErrorBoundaryProvider>
        <div data-testid="child">Content</div>
      </ErrorBoundaryProvider>
    );

    expect(screen.getByTestId('child')).toBeInTheDocument();
  });
});
