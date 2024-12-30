import { render, screen, act, waitFor } from '@testing-library/react';
import { NetworkStatusProvider } from '../providers/NetworkStatusProvider';
import { OfflineBanner } from '@/core/network-status/components/OfflineBanner';
import { vi } from 'vitest';

describe('Network Status Integration', () => {
  const renderApp = () => {
    return render(
      <NetworkStatusProvider>
        <div data-testid="app-content" style={{ height: '100vh' }}>
          <OfflineBanner />
          <div>Main Content</div>
        </div>
      </NetworkStatusProvider>
    );
  };

  beforeEach(() => {
    // Reset any previous network status
    window.setNetworkStatus(true);
  });

  it('should handle real browser online/offline events', async () => {
    renderApp();

    // Initially no banner
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();

    // Simulate real browser offline event
    act(() => {
      window.dispatchEvent(new Event('offline'));
    });

    // Wait for the banner to appear
    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    });

    const banner = screen.getByRole('alert');
    expect(banner).toHaveTextContent('You are offline. Please check your internet connection.');

    // Verify banner positioning and layout
    expect(banner).toHaveStyle({
      position: 'fixed',
      top: '0px',
      left: '0px',
      right: '0px',
      width: '100%',
      zIndex: '1000',
    });

    // Ensure no content shift when banner appears
    const appContent = screen.getByTestId('app-content');
    const initialHeight = appContent.clientHeight;

    // Simulate real browser online event
    act(() => {
      window.dispatchEvent(new Event('online'));
    });

    // Wait for banner to disappear
    await waitFor(() => {
      expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    });

    // Verify no layout shift occurred
    expect(appContent.clientHeight).toBe(initialHeight);
  });

  it('should propagate state changes through component tree', async () => {
    const { container } = renderApp();

    // Take initial snapshot
    expect(container).toMatchSnapshot();

    // Go offline
    await act(async () => {
      window.setNetworkStatus(false);
    });

    // Take offline snapshot
    expect(container).toMatchSnapshot();

    // Go back online
    await act(async () => {
      window.setNetworkStatus(true);
    });

    // Take final snapshot
    expect(container).toMatchSnapshot();
  });

  it('should handle rapid network status changes', async () => {
    renderApp();

    await act(async () => {
      // Simulate network flapping
      window.dispatchEvent(new Event('offline'));
      window.dispatchEvent(new Event('online'));
      window.dispatchEvent(new Event('offline'));
    });

    // Should end up showing banner in offline state
    expect(screen.getByRole('alert')).toBeInTheDocument();
  });

  it('should cleanup event listeners', () => {
    const removeEventListenerSpy = vi.spyOn(window, 'removeEventListener');
    const { unmount } = renderApp();

    unmount();

    expect(removeEventListenerSpy).toHaveBeenCalledWith('online', expect.any(Function));
    expect(removeEventListenerSpy).toHaveBeenCalledWith('offline', expect.any(Function));
  });
});