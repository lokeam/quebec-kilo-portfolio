import { ErrorInfo } from 'react';

export type ErrorSeverity = 'fatal' | 'error' | 'warning' | 'info' | 'debug';

export interface ErrorConfig {
  message?: string;
  severity?: ErrorSeverity;
  actionLabel?: string;
  onAction?: () => void;
  onError?: (error: Error, errorInfo: React.ErrorInfo) => void;
  onReset?: () => void;
};

export interface ErrorBoundaryProps {
  children: React.ReactNode;
  config?: ErrorConfig;
};

export interface ErrorFallbackProps {
  error: Error;
  resetErrorBoundary: () => void;
}

export interface ErrorBoundaryProps {
  children: React.ReactNode;
  fallback?: React.ComponentType<ErrorFallbackProps>;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
  onReset?: () => void;
}