import { render, act } from '@testing-library/react';
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
    rendered = renderWithProvider(<TestComponent />);
  });

  afterEach(() => {
    rendered.unmount();
  });

  it('handles network status changes correctly', async () => {
    const { getByTestId } = rendered;

    // Initial state
    expect(getByTestId('test-component')).toHaveTextContent('Network Status: online');

    // Offline
    await act(async () => {
      window.setNetworkStatus(false);
    });
    expect(getByTestId('test-component')).toHaveTextContent('Network Status: offline');

    // Back online
    await act(async () => {
      window.setNetworkStatus(true);
    });
    expect(getByTestId('test-component')).toHaveTextContent('Network Status: online');

    // Rapid changes
    await act(async () => {
      window.setNetworkStatus(false);
      window.setNetworkStatus(true);
      window.setNetworkStatus(false);
    });
    expect(getByTestId('test-component')).toHaveTextContent('Network Status: offline');
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