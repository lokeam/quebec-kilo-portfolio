declare global {
  interface Window {
    setNetworkStatus: (status: boolean) => void;
  }
}

export {};