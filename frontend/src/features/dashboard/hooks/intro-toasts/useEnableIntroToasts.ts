import { useAuth0 } from '@auth0/auth0-react';
import { useEffect, useState } from 'react';

// Queries
import { useUpdateUserMetadata } from '@/core/api/queries/user.queries';


// Extend Window interface for optimistic state
declare global {
  interface Window {
    __ONBOARDING_OPTIMISTIC_COMPLETE__?: boolean;
    __WANTS_INTRO_TOASTS__?: boolean;
  }
}

/*
  This hook is used to enable or disable intro toasts shown during onboarding.
  It uses optimistic updates to immediately reflect the user's choice,
  then updates the backend asynchronously.
 */
export const useEnableIntroToasts = () => {
  const { user } = useAuth0();
  const [wantsIntroToasts, setWantsIntroToasts] = useState<boolean | null>(null);
  const updateUserMetadataMutation = useUpdateUserMetadata();

  useEffect(() => {
    // Check for optimistic state first (set during onboarding)
    if (window.__WANTS_INTRO_TOASTS__ !== undefined) {
      // console.log('✅ Using optimistic wantsIntroToasts value:', window.__WANTS_INTRO_TOASTS__);
      setWantsIntroToasts(window.__WANTS_INTRO_TOASTS__);
      return;
    }

    // Check localStorage for optimistic state (persists across refreshes)
    const localStorageWantsIntroToasts = localStorage.getItem('__WANTS_INTRO_TOASTS__');
    if (localStorageWantsIntroToasts !== null) {
      const wantsToasts = localStorageWantsIntroToasts === 'true';
      // console.log('✅ Using localStorage optimistic wantsIntroToasts value:', wantsToasts);
      setWantsIntroToasts(wantsToasts);
      return;
    }

    // Only fallback to app_metadata if no optimistic state exists
    // console.log('🔄 Falling back to app_metadata for wantsIntroToasts');
    const fallbackValue = user?.app_metadata?.wantsIntroToasts === true;
    // console.log('🔄 Fallback value from app_metadata:', fallbackValue);
    setWantsIntroToasts(fallbackValue);
  }, [user]);

  const updatePreference = async (wantsToasts: boolean) => {
    if (!user) return;

    // Set optimistic state immediately
    window.__WANTS_INTRO_TOASTS__ = wantsToasts;
    localStorage.setItem('__WANTS_INTRO_TOASTS__', wantsToasts.toString());
    setWantsIntroToasts(wantsToasts);

    // console.log('📤 Updating wantsIntroToasts preference:', wantsToasts);

    try {
      await updateUserMetadataMutation.mutateAsync({
        wantsIntroToasts: wantsToasts
      });

      // console.log('✅ Successfully updated wantsIntroToasts in backend');
    } catch (error) {
      console.error('❌ Failed to update intro toasts preference:', error);
      // Keep the optimistic state even if backend fails
    }
  };

  // Return the resolved value (default to false if still loading to be safe)
  const resolvedValue = wantsIntroToasts !== null ? wantsIntroToasts : false;

  return { wantsIntroToasts: resolvedValue, updatePreference };
};