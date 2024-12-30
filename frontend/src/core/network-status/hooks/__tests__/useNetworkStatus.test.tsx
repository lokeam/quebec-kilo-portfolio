import { renderHook, act } from '@testing-library/react';
import { useNetworkStatus } from '../useNetworkStatus';
import { vi } from 'vitest';
import { NETWORK_STATUS_ERRORS } from '@/core/network-status/constants/errorConstants';
import { NetworkStatusProvider } from '../../providers/NetworkStatusProvider';

describe('useNetworkStatus', () => {
  it('returns current network status', () => {
    const { result } = renderHook(() => useNetworkStatus(), {
      wrapper: NetworkStatusProvider
    });

    expect(result.current.isOnline).toBe(true);
  });


  it('updates when network status changes', async () => {
    const { result } = renderHook(() => useNetworkStatus(), {
      wrapper: NetworkStatusProvider
    });

    await act(async () => {
      window.setNetworkStatus(false);
    });

    await act(async () => {
      window.setNetworkStatus(true);
    });

    expect(result.current.isOnline).toBe(true);
  });

  it('throws error when used outside provider', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    expect(() => {
      renderHook(() => useNetworkStatus());
    }).toThrowError(new Error(NETWORK_STATUS_ERRORS.PROVIDER_MISSING));

    consoleSpy.mockRestore();
  });
})