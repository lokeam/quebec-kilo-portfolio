/**
 * User Service
 *
 * Provides functions for managing user operations through the backend API.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';

// Types
import type {
  RequestUserDeletionResponse,
  UserDeletionStatus,
  UserProfile,
  CreateUserProfileRequest,
  UpdateUserProfileRequest,
  UserProfileResponseWrapper,
} from '@/types/domain/user';

// Constants
import { API_BASE_PATH } from '@/core/api/config';


// Response wrapper for user deletion request
interface UserDeletionResponseWrapper {
  success: boolean;
  userID: string;
  data: RequestUserDeletionResponse;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

// Response wrapper for user deletion status
interface UserDeletionStatusResponseWrapper {
  success: boolean;
  userID: string;
  data: UserDeletionStatus;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

const USER_DELETION_ENDPOINT = '/v1/users/deletion';

/**
 * Gets the user's current profile
 *
*/
export const getUserProfile = (): Promise<UserProfile> =>
  apiRequest('getUserProfile', () =>
    axiosInstance
      .get<UserProfileResponseWrapper>(`${API_BASE_PATH}/users/profile`)
      .then(response => response.data.data)
  );

/**
 * Creates a new user profile
 */
export const createUserProfile = (data: CreateUserProfileRequest): Promise<UserProfile> =>
  apiRequest('createUserProfile', () =>
    axiosInstance
      .post<UserProfileResponseWrapper>(`${API_BASE_PATH}/users`, data)
      .then(response => response.data.data)
  );

/**
 * Updates the user's profile
 */
export const updateUserProfile = (data: UpdateUserProfileRequest): Promise<UserProfile> =>
  apiRequest('updateUserProfile', () =>
    axiosInstance
      .put<UserProfileResponseWrapper>(`${API_BASE_PATH}/users/profile`, data)
      .then(response => response.data.data)
  );


/**
 * Requests deletion of the current user's account
 */
export const requestUserDeletion = (reason: string): Promise<RequestUserDeletionResponse> =>
  apiRequest('requestUserDeletion', () =>
    axiosInstance
      .post<UserDeletionResponseWrapper>(`${API_BASE_PATH}${USER_DELETION_ENDPOINT}/request`, { reason })
      .then(response => response.data.data)
  );

/**
 * Cancels a pending user deletion request
 */
export const cancelUserDeletion = (): Promise<{ message: string }> =>
  apiRequest('cancelUserDeletion', () =>
    axiosInstance
      .post<UserDeletionResponseWrapper>(`${API_BASE_PATH}${USER_DELETION_ENDPOINT}/cancel`)
      .then(response => response.data.data as { message: string })
  );

/**
 * Gets the current deletion status for the user
 */
export const getUserDeletionStatus = (): Promise<UserDeletionStatus> =>
  apiRequest('getUserDeletionStatus', () =>
    axiosInstance
      .get<UserDeletionStatusResponseWrapper>(`${API_BASE_PATH}${USER_DELETION_ENDPOINT}/status`)
      .then(response => response.data.data)
  );

/**
 * Updates user metadata in Auth0
 */
export const updateUserMetadata = (metadata: Record<string, unknown>): Promise<void> =>
  apiRequest('updateUserMetadata', () =>
    axiosInstance
      .patch(`${API_BASE_PATH}/users/metadata`, metadata)
      .then(() => {}) // Return void - no data needed
  );
