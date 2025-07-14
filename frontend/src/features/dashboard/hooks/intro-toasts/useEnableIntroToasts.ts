import { useAuthContext } from '@/core/auth/context-provider/AuthContext';
import { axiosInstance } from '@/core/api/client/axios-instance';

/*
  This hook is used to enable or disable intro toasts shown during onboarding.
  It is used to update the user's preference for intro toasts in the database.
 */
export const useEnableIntroToasts = () => {
  const { user } = useAuthContext();
  const wantsIntroToasts = user?.app_metadata?.wantsIntroToasts !== false;

  const updatePreference = async (wantsToasts: boolean) => {
    if (!user) return;
    try {
      await axiosInstance.patch('/v1/users/metadata', {
        wantsIntroToasts: wantsToasts
      });
    } catch (error) {
      console.error('Failed to update intro toasts preference:', error);
    }
  };

  return { wantsIntroToasts, updatePreference };
};