
import { useNetworkStatus } from '@/core/network-status/hooks/useNetworkStatus';
import { WifiOff } from 'lucide-react';
import { IconCloudOff } from '@tabler/icons-react';


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
          <div className="w-auto flex flex-row">
            <div className="flex flex-col items-center text-lg font-bold text-center align-middle">
              <h2 className="text-xl mb-2">You are offline</h2>
              <p className="text-md mb-2">Please check your internet connection.</p>
              <IconCloudOff color="#fff" size={42} />
            </div>
          </div>
        </div>
      </div>
  );
};
