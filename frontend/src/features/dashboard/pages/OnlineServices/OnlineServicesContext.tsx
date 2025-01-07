import { createContext, useContext, ReactNode, useState } from 'react';

type ViewMode = 'grid' | 'list' | 'table';

interface OnlineServicesContextType {
  viewMode: ViewMode;
  setViewMode: (mode: ViewMode) => void;
}

const OnlineServicesContext = createContext<OnlineServicesContextType | undefined>(undefined);

interface OnlineServicesProviderProps {
  children: ReactNode;
}

export function OnlineServicesProvider({ children }: OnlineServicesProviderProps) {
  const [viewMode, setViewMode] = useState<ViewMode>(() => {
    const saved = localStorage.getItem('onlineServicesViewMode') as ViewMode;
    return saved || 'grid';
  });

  // Persist view mode to localStorage
  const handleSetViewMode = (mode: ViewMode) => {
    setViewMode(mode);
    localStorage.setItem('onlineServicesViewMode', mode);
  };

  return (
    <OnlineServicesContext.Provider
      value={{
        viewMode,
        setViewMode: handleSetViewMode
      }}
    >
      {children}
    </OnlineServicesContext.Provider>
  );
}

export function useOnlineServices() {
  const context = useContext(OnlineServicesContext);
  if (context === undefined) {
    throw new Error('useOnlineServices must be used within an OnlineServicesProvider');
  }
  return context;
}