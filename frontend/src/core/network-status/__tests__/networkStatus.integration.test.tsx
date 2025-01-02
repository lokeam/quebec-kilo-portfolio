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
    // Reset navigator.onLine before each test
    Object.defineProperty(navigator, 'onLine', {
      configurable: true,
      value: true
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it('should handle real browser online/offline events', async () => {
    renderApp();

    // Initially no banner
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();

    // Simulate real browser offline event
    await act(async () => {
      Object.defineProperty(navigator, 'onLine', {
        configurable: true,
        value: false
      });
      window.dispatchEvent(new Event('offline'));
    });

    // Wait for the banner to appear
    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    }, { timeout: 1000 }); // Account for 500ms delay

    const banner = screen.getByRole('alert');

    // Check for heading and text content
    expect(screen.getByRole('heading', { level: 2 })).toHaveTextContent('You are offline');
    expect(screen.getByText('Please check your internet connection.')).toBeInTheDocument();

    // Verify icon presence
    expect(banner.querySelector('svg')).toBeInTheDocument();

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
    await act(async () => {
      Object.defineProperty(navigator, 'onLine', {
        configurable: true,
        value: true
      });
      window.dispatchEvent(new Event('online'));
    });

    // Wait for banner to disappear
    await waitFor(() => {
      expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    });

    // Verify no layout shift occurred
    expect(appContent.clientHeight).toBe(initialHeight);
  });

  it('should handle rapid network status changes', async () => {
    renderApp();

    await act(async () => {
      // Simulate network flapping
      Object.defineProperty(navigator, 'onLine', { value: false });
      window.dispatchEvent(new Event('offline'));

      Object.defineProperty(navigator, 'onLine', { value: true });
      window.dispatchEvent(new Event('online'));

      Object.defineProperty(navigator, 'onLine', { value: false });
      window.dispatchEvent(new Event('offline'));
    });

    // Wait for the final offline state to be reflected
    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    }, { timeout: 1000 });
  });

  it('should cleanup event listeners', () => {
    const removeEventListenerSpy = vi.spyOn(window, 'removeEventListener');
    const { unmount } = renderApp();

    unmount();

    expect(removeEventListenerSpy).toHaveBeenCalledWith('online', expect.any(Function));
    expect(removeEventListenerSpy).toHaveBeenCalledWith('offline', expect.any(Function));
  });
});