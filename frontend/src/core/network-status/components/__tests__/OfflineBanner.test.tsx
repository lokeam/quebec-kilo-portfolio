import { screen, act, waitFor } from '@testing-library/react';
import { OfflineBanner } from '@/core/network-status/components/OfflineBanner';
import { renderWithProvider } from '@/core/network-status/providers/__tests__/utils/renderUtils';

describe('OfflineBanner', () => {
  it('should not display when online', () => {
    renderWithProvider(<OfflineBanner />);
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it('should display when offline', async () => {
    renderWithProvider(<OfflineBanner />);

    act(() => {
      window.setNetworkStatus(false);
    });

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    });

    const alert = screen.getByRole('alert');
    expect(alert).toHaveTextContent('You are offline. Please check your internet connection.');
  });

  it('should not display when online', async () => {
    renderWithProvider(<OfflineBanner />);

    act(() => {
      window.setNetworkStatus(true);
    });

    await waitFor(() => {
      expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    });
  });
});