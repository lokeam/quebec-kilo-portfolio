export type OnboardingStep =
  | 'WELCOME'
  | 'STORAGE'
  | 'PHYSICAL'
  | 'DIGITAL'
  | 'WISHLIST'
  | 'COMPLETE';

export type WorkspaceChoice =
  | 'PHYSICAL_STORAGE'
  | 'DIGITAL_STORAGE'
  | 'PHYSICAL_AND_DIGITAL_STORAGE';

export interface StorageLocation {
  name: string;
  classification: string;
  isOptional?: boolean;
};
