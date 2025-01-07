import { useState } from 'react';

export type ViewMode = 'grid' | 'list' | 'table';

export function useViewMode() {
  const [viewMode, setViewMode] = useState<ViewMode>(() => {
    const saved = localStorage.getItem('onlineServicesViewMode') as ViewMode;
    return saved || 'grid';
  });

  const handleViewModeChange = (mode: ViewMode) => {
    setViewMode(mode);
    localStorage.setItem('onlineServicesViewMode', mode);
  };

  return { viewMode, setViewMode: handleViewModeChange };
}