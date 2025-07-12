
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

// User Profile Types
export interface UserProfile {
  id: string;
  email: string;
  firstName?: string;
  lastName?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateUserProfileRequest {
  email: string;        // From Auth0 user.email
  auth0UserID: string;  // From Auth0 user.sub
  firstName: string;    // User provides during onboarding
  lastName: string;     // User provides during onboarding
}

export interface UpdateUserProfileRequest {
  firstName?: string;    // User provides during onboarding
  lastName?: string;     // User provides during onboarding
}

// Response wrapper for user profile operations
export interface UserProfileResponseWrapper {
  success: boolean;
  userID: string;
  data: UserProfile;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}
