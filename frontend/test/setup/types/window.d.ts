// Window interface definition for matchMedia mock testing
export interface Window {
  matchMedia: (query: string) => MediaQueryList;
}