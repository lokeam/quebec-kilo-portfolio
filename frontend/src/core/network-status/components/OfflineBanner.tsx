import { Alert } from '@mui/material';
import { useNetworkStatus } from '@/core/network-status/hooks/useNetworkStatus';

export const OfflineBanner = () => {
  const { isOnline } = useNetworkStatus();

  if (isOnline) {
    return null;
  }

  return (
    <Alert
      severity="error"
      sx={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        zIndex: 9999,
        width: '100%',
        textAlign: 'center',
      }}
    >
      You are offline. Please check your internet connection.
    </Alert>
  )


}