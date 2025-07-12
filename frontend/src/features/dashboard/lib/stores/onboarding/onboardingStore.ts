import { create } from 'zustand';

// Types from current type system
import type { OnboardingStorageType } from '@/types/domain/onboarding';
import type { CreatePhysicalLocationRequest } from '@/types/domain/physical-location';
import type { CreateDigitalLocationRequest } from '@/types/domain/digital-location';
import type { PhysicalLocationType, LocationIconBgColor } from '@/types/domain/location-types';

/**
 * --------------------------------------
 * Define partial "draft" types
 * --------------------------------------
 *
 * Since user is only filling out fields step by step, store
 * partial versions of the "Create" inputs. We omit server-generated fields
 * like IDs to keep the store minimal.
 *
 * Using Partial<> ensures we can do incremental updates safely.
*/

type PhysicalLocationDraft = Partial<CreatePhysicalLocationRequest> & {
  bgColor?: LocationIconBgColor;
};

type DigitalLocationDraft = Partial<CreateDigitalLocationRequest>;

interface OnboardingState {
  /** Tracks if the user has fully completed onboarding */
  onboardingComplete: boolean;

  /** Tracks whether user wants setup (true), no setup (false), or not chosen yet (null) */
  wantsSetup: boolean | null;

  /** Which storage type the user selected: PHYSICAL, DIGITAL, or BOTH */
  storageType: OnboardingStorageType | null;

  /** Temporarily store partial physical location data during onboarding */
  physicalDraft?: PhysicalLocationDraft;

  /** Temporarily store partial digital location data during onboarding */
  digitalDraft?: DigitalLocationDraft;

  /* Actions */
  setOnboardingComplete: (complete: boolean) => void;
  setWantsSetup: (wantsSetup: boolean) => void;
  setStorageType: (type: OnboardingStorageType) => void;
  updatePhysicalDraft: (draft: Partial<{
    name: string;
    locationType: PhysicalLocationType;
    mapCoordinates?: string;
    bgColor?: LocationIconBgColor;
  }>) => void;
  updateDigitalDraft: (draft: Partial<DigitalLocationDraft>) => void;
  resetOnboarding: () => void;
}

export const useOnboardingStore = create<OnboardingState>((set) => ({
  onboardingComplete: false,
  wantsSetup: null,
  storageType: null,
  physicalDraft: {},
  digitalDraft: {},

  setOnboardingComplete: (complete) =>
    set({ onboardingComplete: complete }),

  setWantsSetup: (wantsSetup) =>
    set({ wantsSetup }),

  setStorageType: (type) =>
    set({ storageType: type }),

  /** Merges partial updates into the existing physicalDraft */
  updatePhysicalDraft: (draft) =>
    set((state) => ({
      physicalDraft: {
        ...state.physicalDraft,
        ...draft,
      },
    })),

  /** Merges partial updates into the existing digitalDraft */
  updateDigitalDraft: (draft) =>
    set((state) => ({
      digitalDraft: {
        ...state.digitalDraft,
        ...draft,
      },
   })),

  /** Optional helper fn to clear/reset store if needed */
  resetOnboarding: () =>
    set(() => ({
      onboardingComplete: false,
      wantsSetup: null,
      storageType: null,
      physicalDraft: {},
      digitalDraft: {},
    })),
}));

export const useSetOnboardingStorageType = () => useOnboardingStore((state) => state.setStorageType);
export const useOnboardingPhysicalDraft = () => useOnboardingStore((state) => state.physicalDraft);
export const useUpdatePhysicalDraft = () => useOnboardingStore((state) => state.updatePhysicalDraft);
export const useOnboardingDigitalDraft = () => useOnboardingStore((state) => state.digitalDraft);
