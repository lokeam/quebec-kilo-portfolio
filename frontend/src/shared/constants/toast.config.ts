// constants
import type { ToastVariant } from "@/types/toast";

export const TOAST_DURATIONS = {
  STANDARD: 60000, // 60 seconds
  EXTENDED: 90000, // 90 seconds
} as const;

export const TOAST_CLASSNAMES = {
  SUCCESS: 'success-toast',
  ERROR: 'error-toast',
  WARNING: 'warning-toast',
  INFO: 'info-toast'
} as const;

export const TOAST_POSITION = 'top-left' as const;

export const DEFAULT_TOAST_CONFIG = {
  duration: TOAST_DURATIONS.STANDARD,
  className: TOAST_CLASSNAMES.SUCCESS,
  position: TOAST_POSITION,
  variant: 'success',
} as const;

export const TOAST_BG_CLASSES: Record<ToastVariant, string> = {
  success: 'border-2 border-solid border-green-600',
  error:   'border-2 border-solid border-red-600',
  info:    'border-2 border-solid border-blue-600',
  warning: 'border-2 border-solid border-yellow-600',
}