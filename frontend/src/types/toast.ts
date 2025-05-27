export interface ToastConfig {
  /**
   * Duration in milliseconds
   */
  duration?: number;

  /**
   * Optional CSS class name for styling
   */
  className?: string;

  /**
   * Optional description text shown below the main message
   */
  description?: string;

  /**
   * Optional position of the toast
   */
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right';

  /**
   * Optional variant of the toast, used to determine the display icon
   */
  variant?: ToastVariant;
}

export type ToastVariant = 'success' | 'error' | 'info' | 'warning';