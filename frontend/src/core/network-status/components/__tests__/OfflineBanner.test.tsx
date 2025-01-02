import { screen, act, waitFor } from '@testing-library/react';
import { OfflineBanner } from '@/core/network-status/components/OfflineBanner';
import { renderWithProvider } from '@/core/network-status/providers/__tests__/utils/renderUtils';

describe('OfflineBanner', () => {
  beforeEach(() => {
    // Reset navigator.onLine before each test
    Object.defineProperty(navigator, 'onLine', {
      configurable: true,
      value: true
    });
  });

  it('should not display when online', () => {
    renderWithProvider(<OfflineBanner />);
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it('should display when offline', async () => {
    renderWithProvider(<OfflineBanner />);

    await act(async () => {
      Object.defineProperty(navigator, 'onLine', {
        configurable: true,
        value: false
      });
      window.dispatchEvent(new Event('offline'));
    });

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    }, { timeout: 1000 });

    // Check for heading and text content
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('You are offline');
    expect(screen.getByText('Please check your internet connection.')).toBeInTheDocument();

    // Verify icon presence
    const alert = screen.getByRole('alert');
    expect(alert.querySelector('svg')).toBeInTheDocument();

    // Verify banner styling
    expect(alert).toHaveClass('offline_banner', 'bg-charcoal');
  });

  it('should not display when online', async () => {
    renderWithProvider(<OfflineBanner />);

    // First make it offline
    await act(async () => {
      Object.defineProperty(navigator, 'onLine', {
        configurable: true,
        value: false
      });
      window.dispatchEvent(new Event('offline'));
    });

    // Then back online
    await act(async () => {
      Object.defineProperty(navigator, 'onLine', {
        configurable: true,
        value: true
      });
      window.dispatchEvent(new Event('online'));
    });

    await waitFor(() => {
      expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    });
  });
});