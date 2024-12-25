import { useContext } from 'react';
import { NetworkStatusContext } from '@/core/network-status/providers/NetworkStatusProvider';
import { NETWORK_STATUS_ERRORS } from '@/core/network-status/constants/errorConstants';
import { NetworkStatusContextType } from '@/core/network-status/types/networkStatusType';

export const useNetworkStatus = (): NetworkStatusContextType => {
  const context = useContext(NetworkStatusContext);
  if (context === undefined) {
    throw new Error(NETWORK_STATUS_ERRORS.PROVIDER_MISSING);
  }
  return context;
};