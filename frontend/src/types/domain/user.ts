
export interface RequestUserDeletionRequest {
  reason: string;
};

export interface RequestUserDeletionResponse {
  message: string;
  gracePeriodEnd?: string; // ISO date string - axios will convert from snake_case
};

export interface UserDeletionStatus {
  isActive: boolean;
  isDeleted: boolean;
  isDeletionRequested: boolean;
  isInGracePeriod: boolean;
  gracePeriodEnd?: string;
  deletionReason?: string;
};
