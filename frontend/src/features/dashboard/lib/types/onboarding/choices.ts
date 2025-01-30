import type { OnboardingStep, WorkspaceChoice, StorageLocation } from '@/features/dashboard/lib/types/onboarding/base';

export interface PhysicalStorageData {
  storageLocations: StorageLocation[];
}

export interface DigitalStorageData {
  physicalLocationName: string;
  physicalLocationType: string;
  hasCoordinates?: boolean;
  coordinates?: string;
  subLocationName: string;
  subLocationType: string;
}

export interface LibraryItem {
  name: string;
  associatedWith: WorkspaceChoice;
}

export interface OnboardingState {
  currentStep: OnboardingStep;
  isComplete: boolean;
  storageChoice?: WorkspaceChoice;
  physicalStorageData?: PhysicalStorageData;
  digitalStorageData?: DigitalStorageData;
  libraryItems: LibraryItem[];
}
