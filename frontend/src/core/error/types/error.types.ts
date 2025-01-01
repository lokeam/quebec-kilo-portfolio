export type HttpErrorCode = 400 | 401 | 403 | 404 | 500 | 503;

export interface ApiError extends Error {
  statusCode: HttpErrorCode;
  message: string;
}

export interface DefaultErrorFallbackProps {
  error: Error;
  resetErrorBoundary: () => void;
}

export interface ErrorButtonProps {
  onClick: () => void;
  variant: 'retry' | 'home' | 'back';
  label: string;
}
