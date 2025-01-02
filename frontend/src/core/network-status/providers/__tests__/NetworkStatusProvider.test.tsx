import { render, act, waitFor } from '@testing-library/react';
import { NetworkStatusProvider, useNetworkStatus } from '@/core/network-status/providers/NetworkStatusProvider';
import { vi } from 'vitest';
import { NETWORK_STATUS_ERRORS } from '@/core/network-status/constants/errorConstants';
import { renderWithProvider } from './utils/renderUtils';

const TestComponent = () => {
  const { isOnline } = useNetworkStatus();
  return (
    <div data-testid="test-component">
      Network Status: {isOnline ? 'online' : 'offline'}
    </div>
  );
};

describe('NetworkStatusProvider', () => {
  let rendered: ReturnType<typeof renderWithProvider>;

  beforeEach(() => {
    // Reset navigator.onLine before each test
    Object.defineProperty(navigator, 'onLine', {
      configurable: true,
      value: true
    });
    rendered = renderWithProvider(<TestComponent />);
  });

  afterEach(() => {
    rendered.unmount();
    vi.clearAllMocks();
  });

  it('handles network status changes correctly', async () => {
    const { getByTestId } = rendered;

    // Initial state
    expect(getByTestId('test-component')).toHaveTextContent('Network Status: online');

    // Simulate offline
    await act(async () => {
      Object.defineProperty(navigator, 'onLine', {
        configurable: true,
        value: false
      });
      window.dispatchEvent(new Event('offline'));
    });

    // Wait for the delayed status update
    await waitFor(() => {
      expect(getByTestId('test-component')).toHaveTextContent('Network Status: offline');
    }, { timeout: 1000 }); // Increase timeout to account for the 500ms delay

    // Simulate back online
    await act(async () => {
      Object.defineProperty(navigator, 'onLine', {
        configurable: true,
        value: true
      });
      window.dispatchEvent(new Event('online'));
    });

    await waitFor(() => {
      expect(getByTestId('test-component')).toHaveTextContent('Network Status: online');
    });

    // Test rapid changes
    await act(async () => {
      // Offline
      Object.defineProperty(navigator, 'onLine', { value: false });
      window.dispatchEvent(new Event('offline'));
      // Immediately online
      Object.defineProperty(navigator, 'onLine', { value: true });
      window.dispatchEvent(new Event('online'));
      // Finally offline
      Object.defineProperty(navigator, 'onLine', { value: false });
      window.dispatchEvent(new Event('offline'));
    });

    await waitFor(() => {
      expect(getByTestId('test-component')).toHaveTextContent('Network Status: offline');
    });
  });

  it('handles provider nesting and error cases', () => {
    const NestedTestComponent = () => {
      const { isOnline } = useNetworkStatus();
      return (
        <div data-testid="nested-test-component">
          Network Status: {isOnline ? 'online' : 'offline'}
        </div>
      );
    };

    const { getByTestId } = render(
      <NetworkStatusProvider>
        <NetworkStatusProvider>
          <NestedTestComponent />
        </NetworkStatusProvider>
      </NetworkStatusProvider>
    );
    expect(getByTestId('nested-test-component')).toHaveTextContent('Network Status: online');

    // Test error when used outside provider
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    expect(() => {
      render(<TestComponent />);
    }).toThrow(NETWORK_STATUS_ERRORS.PROVIDER_MISSING);
    consoleSpy.mockClear();
  });

  it('handles cleanup and event listeners', () => {
    const removeEventListenerSpy = vi.spyOn(window, 'removeEventListener');
    const { unmount } = renderWithProvider(<TestComponent />);

    unmount();

    expect(removeEventListenerSpy).toHaveBeenCalledWith('online', expect.any(Function));
    expect(removeEventListenerSpy).toHaveBeenCalledWith('offline', expect.any(Function));
  });
});