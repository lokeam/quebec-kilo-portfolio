import {
  createContext,
  useState,
  useEffect,
  useContext,
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

  useEffect(() => {
    const handleOnline = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    // Initial status check
    setIsOnline(navigator.onLine);

    // Cleanup fn to remove event listeners on unmount
    return  () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    }
  }, []);

  return (
    <NetworkStatusContext.Provider value= {{ isOnline}} >
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