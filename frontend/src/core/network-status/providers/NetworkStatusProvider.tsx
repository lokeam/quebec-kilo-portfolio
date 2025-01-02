import {
  createContext,
  useState,
  useEffect,
  useContext,
  useCallback,
  type ReactNode,
} from 'react';

import { NETWORK_STATUS_ERRORS } from '../constants/errorConstants';
import { type NetworkStatusContextType } from '../types/networkStatusType';

// Context for managing network connectivity state across app. Undefined until Provider is mounted.
export const NetworkStatusContext = createContext<NetworkStatusContextType | undefined>(undefined);

/*
  Monitors and broadcasts network connectivity status.
  Must be placed high in component tree to provide status to all components.
*/
export function NetworkStatusProvider({ children}: { children: ReactNode }) {
  // Init state w/ current network status
  const [isOnline, setIsOnline] = useState<boolean>(navigator.onLine);
  const [connectionQuality, setConnectionQuality] = useState<string>('unknown');

  // Use useCallback to memoize the handler
  const updateOnlineStatus = useCallback((status: boolean) => {
    // Double check actual status against navigator.onLine
    if (status === navigator.onLine) {
      setIsOnline(status);
    }
  }, []);

  useEffect(() => {
    let timeoutId: number;

    const handleOnline = () => {
      // Clear any pending offline updates
      window.clearTimeout(timeoutId);
      updateOnlineStatus(true);
    };

    const handleOffline = () => {
      // Add slight delay to offline state to prevent flashing
      timeoutId = window.setTimeout(() => {
        updateOnlineStatus(false);
      }, 500); // Based on Chrome's internal delay for offline status
    };

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
      window.clearTimeout(timeoutId);
    };
  }, [updateOnlineStatus]);

  useEffect(() => {
    // Monitor connection quality if available
    if ('connection' in navigator) {
      const connection = (navigator as any).connection;
      const updateConnectionQuality = () => {
        setConnectionQuality(connection.effectiveType);
      };

      connection.addEventListener('change', updateConnectionQuality);
      return () => connection.removeEventListener('change', updateConnectionQuality);
    }
  }, []);

  return (
    <NetworkStatusContext.Provider value={{
      isOnline,
      connectionQuality
    }}>
      { children }
    </NetworkStatusContext.Provider>
  );
}

/*
  Custom hook to access network status from any component.
  Must be used within a NetworkStatusProvider
*/
export const useNetworkStatus = () => {
  const context = useContext(NetworkStatusContext);
  if (context === undefined) {
    throw new Error(NETWORK_STATUS_ERRORS.PROVIDER_MISSING);
  }

  return context;
}