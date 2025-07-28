import type { ReactNode } from 'react';
import { toast as sonnerToast } from 'sonner';

// components
import { OperationToast } from './OperationToast';

// constants
import { TOAST_DURATIONS } from '@/shared/constants/toast.config';

// icons
import { AlertTriangle, Ban, CheckCircle, Info } from '@/shared/components/ui/icons';

// type
import type { ToastConfig, ToastVariant } from '@/types/toast';

interface TanstackMutationToastProps {
  message: string;
  config?: ToastConfig
  variant?: ToastVariant;
  duration?: number;
}

const TOAST_ICONS: Record<ToastVariant, ReactNode> = {
  success: <CheckCircle    className="w-5 h-5 text-green-600" />,
  error:   <Ban            className="w-5 h-5 text-red-600" />,
  info:    <Info           className="w-5 h-5 text-blue-600" />,
  warning: <AlertTriangle  className="w-5 h-5 text-yellow-600" />,
};

export function showToast({ message, variant = 'success', duration = TOAST_DURATIONS.STANDARD
}: TanstackMutationToastProps) {
  return sonnerToast.custom((t) => (
    <OperationToast
      id={t}
      message={message}
      icon={TOAST_ICONS[variant]}
    />
  ), {
    duration,
  });
}
