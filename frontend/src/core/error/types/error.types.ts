export type ErrorSeverity = 'fatal' | 'error' | 'warning' | 'info' | 'debug';

export interface ErrorConfig {
  message?: string;
  severity?: ErrorSeverity;
  actionLabel?: string;
  onAction?: () => void;
};

export interface ErrorBoundaryProps {
  children: React.ReactNode;
  config?: ErrorConfig;
};
