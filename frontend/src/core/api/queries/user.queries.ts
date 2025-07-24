/**
 * User Query Hooks
 *
 * Provides React Query hooks for managing user operations.
 */

import { AxiosError } from 'axios';

// React Query
import { useMutation, useQueryClient } from '@tanstack/react-query';

// Base Query Hook
import { useAPIQuery } from '@/core/api/queries/useAPIQuery';

// Service Layer methods
import {
  requestUserDeletion,
  cancelUserDeletion,
  getUserDeletionStatus,
  createUserProfile,
  updateUserProfile,
  getUserProfile,
  updateUserMetadata,
} from '@/core/api/services/user.service';

// Utils
import { showToast } from '@/shared/components/ui/TanstackMutationToast/showToast';
import { logger } from '@/core/utils/logger/logger';

// Types
import type { ApiError } from '@/types/api/response';
import type {
  UserDeletionStatus,
  UserProfile,
  CreateUserProfileRequest,
  UpdateUserProfileRequest,
} from '@/types/domain/user';

// Constants
import { TOAST_SUCCESS_MESSAGES } from '@/shared/constants/toast.success.messages';
import { TOAST_DURATIONS } from '@/shared/constants/toast.config';
import { TOAST_ERROR_MESSAGES } from '@/shared/constants/toast.error.messages';


/**
 * Query key factory for user queries
 */
export const userKeys = {
  all: ['user'] as const,
  profile: () => [...userKeys.all, 'profile'] as const,
  deletion: () => [...userKeys.all, 'deletion'] as const,
  deletionStatus: () => [...userKeys.deletion(), 'status'] as const,
};

/**
 * Get user profile
 *
*/
export const useGetUserProfile = () => {
  return useAPIQuery<UserProfile>({
    queryKey: userKeys.profile(),
    queryFn: async () => {
      const profile = await getUserProfile();

      if (!profile) {
        throw new Error('No user profile received from server');
      }

      return profile;
    },
    staleTime: 30000, // 30 seconds - profile doesn't change frequently
  })
}

/**
 * Hook to create user profile
 */
export const useCreateUserProfile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateUserProfileRequest) => createUserProfile(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.profile() });
      queryClient.invalidateQueries({ queryKey: userKeys.all });
    },
  });
};

/**
 * Hook to update user profile
 */
export const useUpdateUserProfile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateUserProfileRequest) => updateUserProfile(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.profile() });
      queryClient.invalidateQueries({ queryKey: userKeys.all });
    },
  });
};

/**
 * Hook to request user account deletion
 */
export const useRequestUserDeletion = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (reason: string) => requestUserDeletion(reason),
    onSuccess: (data) => {
      // Log the deletion request
      logger.debug('User deletion requested successfully', {
        message: data.message,
        gracePeriodEnd: data.gracePeriodEnd
      });

      // Invalidate user-related queries
      queryClient.invalidateQueries({ queryKey: userKeys.all });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.USER.DELETION_REQUESTED,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error: AxiosError<ApiError>) => {
      const errorMessages = TOAST_ERROR_MESSAGES.USER.DELETION_REQUEST as {
        DEFAULT: string;
        PERMISSION: string;
        SERVER: string;
      };

      let errorMessage = errorMessages.DEFAULT;

      if (error.response) {
        const status = error.response.status;
        if (status === 401 || status === 403) {
          errorMessage = errorMessages.PERMISSION;
        } else if (status >= 500) {
          errorMessage = errorMessages.SERVER;
        }
      }

      logger.error('Failed to request user deletion', { error });
      showToast({
        message: errorMessage,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    }
  });
};

/**
 * Hook to cancel user account deletion
 */
export const useCancelUserDeletion = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => cancelUserDeletion(),
    onSuccess: (data) => {
      // Log the cancellation
      logger.debug('User deletion cancelled successfully', {
        message: data.message
      });

      // Invalidate user-related queries
      queryClient.invalidateQueries({ queryKey: userKeys.all });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.USER.DELETION_CANCELLED,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error: AxiosError<ApiError>) => {
      const errorMessages = TOAST_ERROR_MESSAGES.USER.DELETION_CANCEL as {
        DEFAULT: string;
        PERMISSION: string;
        SERVER: string;
      };

      let errorMessage = errorMessages.DEFAULT;

      if (error.response) {
        const status = error.response.status;
        if (status === 401 || status === 403) {
          errorMessage = errorMessages.PERMISSION;
        } else if (status >= 500) {
          errorMessage = errorMessages.SERVER;
        }
      }

      logger.error('Failed to cancel user deletion', { error });
      showToast({
        message: errorMessage,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    }
  });
};

/**
 * Hook to get user deletion status
 */
export const useGetUserDeletionStatus = () => {
  return useAPIQuery<UserDeletionStatus>({
    queryKey: userKeys.deletionStatus(),
    queryFn: async () => {
      const status = await getUserDeletionStatus();
      return status;
    },
    staleTime: 30000, // 30 seconds - status doesn't change frequently
    refetchOnMount: true,
  });
};

/**
 * Hook to update user metadata
 */
export const useUpdateUserMetadata = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (metadata: Record<string, unknown>) => updateUserMetadata(metadata),
    onSuccess: () => {
      // Invalidate user profile queries to refresh data
      queryClient.invalidateQueries({ queryKey: userKeys.profile() });
      queryClient.invalidateQueries({ queryKey: userKeys.all });
    },
  });
};
