/**
 * User Service
 *
 * Provides functions for managing user operations through the backend API.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';
import type {
  RequestUserDeletionResponse,
  UserDeletionStatus,
} from '@/types/domain/user';


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
 * Requests deletion of the current user's account
 */
export const requestUserDeletion = (reason: string): Promise<RequestUserDeletionResponse> =>
  apiRequest('requestUserDeletion', () =>
    axiosInstance
      .post<UserDeletionResponseWrapper>(`${USER_DELETION_ENDPOINT}/request`, { reason })
      .then(response => response.data.data)
  );

/**
 * Cancels a pending user deletion request
 */
export const cancelUserDeletion = (): Promise<{ message: string }> =>
  apiRequest('cancelUserDeletion', () =>
    axiosInstance
      .post<UserDeletionResponseWrapper>(`${USER_DELETION_ENDPOINT}/cancel`)
      .then(response => response.data.data as { message: string })
  );

/**
 * Gets the current deletion status for the user
 */
export const getUserDeletionStatus = (): Promise<UserDeletionStatus> =>
  apiRequest('getUserDeletionStatus', () =>
    axiosInstance
      .get<UserDeletionStatusResponseWrapper>(`${USER_DELETION_ENDPOINT}/status`)
      .then(response => response.data.data)
  );
