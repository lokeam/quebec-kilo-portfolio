import { Toaster } from 'sonner';

// hooks
import { useThemeStore } from '@/core/theme/stores/useThemeStore';

// constants
import { TOAST_BG_CLASSES, TOAST_POSITION } from '@/shared/constants/toast.config';

export function TanstackMutationToast() {
  const { mode } = useThemeStore();

  return (
    <Toaster
      theme={mode}
      position={TOAST_POSITION}
      toastOptions={{
        classNames: {
          toast: "border-2 border-solid border-gray-200, rounded-lg shadow-lg",
          success: TOAST_BG_CLASSES.success,
          error: TOAST_BG_CLASSES.error,
          warning: TOAST_BG_CLASSES.warning,
          info: TOAST_BG_CLASSES.info
        }
      }}
    />
  )
}
