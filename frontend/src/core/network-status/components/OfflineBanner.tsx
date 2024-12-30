
import { useNetworkStatus } from '@/core/network-status/hooks/useNetworkStatus';
import { WifiOff } from 'lucide-react';


export const OfflineBanner = () => {
  const { isOnline } = useNetworkStatus();

  if (isOnline) return null;

  return (
      <div
        className="offline_banner bg-charcoal relative box-border flex flex-col min-h-min items-center content-center h-full py-4 px-2"
        role="alert"
        style={{
          position: 'fixed',
          top: '0px',
          left: '0px',
          right: '0px',
          width: '100%',
          zIndex: 1000
        }}
      >
        <div className="flex flex-col w-full items-center content-center">
          <div className="w-64 flex flex-row">
            <div className="relative inline-flex items-center justify-center overflow-hidden bg-transparent h-9 w-9 rounded-full">
              <WifiOff color="#fff" size={30} />
          </div>
            <div className="flex flex-col items-center justify-center text-lg font-bold text-center align-middle ml-3">You are offline. Please check your internet connection.</div>
          </div>
        </div>
      </div>
  );
};
