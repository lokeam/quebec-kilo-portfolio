import type { ApiError, HttpErrorCode } from '../types/error.types';
import { ERROR_MESSAGES } from '../constants/error.constants';

export function isApiError(error: unknown): error is ApiError {
  return (
    error !== null &&
    typeof error === 'object' &&
    'name' in error &&
    'message' in error &&
    'statusCode' in error &&
    (error as ApiError).name === 'ApiError' &&
    typeof (error as ApiError).statusCode === 'number'
  );
};

export function getErrorMessage(statusCode: HttpErrorCode): string {
  switch (statusCode) {
    case 401:
      return ERROR_MESSAGES.UNAUTHORIZED;
    case 403:
      return ERROR_MESSAGES.FORBIDDEN;
    case 404:
      return ERROR_MESSAGES.NOT_FOUND;
    case 500:
    case 503:
      return ERROR_MESSAGES.SERVER_ERROR;
    default:
      return ERROR_MESSAGES.UNKNOWN;
  }
};
