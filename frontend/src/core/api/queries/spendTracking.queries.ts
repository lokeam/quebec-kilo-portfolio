// Tanstack query
import { useMutation, useQueryClient } from '@tanstack/react-query';

// Base query hook
import { useAPIQuery } from '@/core/api/queries/useAPIQuery';

// Adapters
import { spendTrackingAdapter } from '@/core/api/adapters/spendTracking.adapter';

// Service Layer methods
import {
  getSpendTrackingItemById,
  createOneTimePurchase,
  updateSpendTrackingItem,
  deleteSpendTrackingItem,
  getSpendTrackingPageBFFResponse,
} from '@/core/api/services/spendTracking.service';

// Utils
import { showToast } from '@/shared/components/ui/TanstackMutationToast/showToast';

// Types
import type {
  CreateOneTimePurchaseRequest,
  SpendItem,
  SingleYearlyTotalBFFResponse,
  SpendTrackingBFFResponse,
  MediaCategory,
  TransactionType,
  SpendingItemBFFResponse,
} from '@/types/domain/spend-tracking';

// Constants
import { TOAST_SUCCESS_MESSAGES } from '@/shared/constants/toast.success.messages';
import { TOAST_ERROR_MESSAGES } from '@/shared/constants/toast.error.messages';
import { TOAST_DURATIONS } from '@/shared/constants/toast.config';
import { dashboardKeys } from './dashboard.queries';

export const spendTrackingKeys = {
  all: ['spend-tracking'] as const,
  lists: () => [...spendTrackingKeys.all, 'list'] as const,
  list: (filters: string) => [...spendTrackingKeys.lists(), { filters }] as const,
  details: () => [...spendTrackingKeys.all, 'detail'] as const,
  detail: (id: string) => [...spendTrackingKeys.details(), id] as const,
};

interface SpendTrackingPageData {
  currentMonthItems: SpendItem[];
  nextMonthItems: SpendItem[];
  yearlyTotals: SingleYearlyTotalBFFResponse[];
}

/**
 * Hook to fetch all spend tracking data for the BFF page
 */
export const useGetSpendTrackingPageBFFResponse = () => {
  return useAPIQuery<SpendTrackingBFFResponse>({
    queryKey: spendTrackingKeys.lists(),
    queryFn: async () => {
      try {
        const response = await getSpendTrackingPageBFFResponse();

        return spendTrackingAdapter.transformSpendTrackingResponse(response);
      } catch(error) {
        console.error('[DEBUG] useGetSpendTrackingPageBFFResponse: Error fetching data:', error);
        throw error;
      }
    },
    staleTime: 0, // Consider data stale immediately
    refetchOnMount: true, // Refetch when component mounts
    refetchOnWindowFocus: true, // Refetch when window regains focus
  });
};

/**
 * Hook to fetch a single spend item
 */
export const useGetSingleSpendItem = (id: string) => {
  return useAPIQuery<SpendItem>({
    queryKey: spendTrackingKeys.detail(id),
    queryFn: async () => {
      const item = await getSpendTrackingItemById(id);
      return item;
    }
  });
};

/**
 * Hook to create a spend item
 */
export const useCreateSpendItem = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateOneTimePurchaseRequest) => createOneTimePurchase(data), // ← CHANGE THIS LINE
    onSuccess: (data) => {
      console.log('�� DEBUG: useCreateSpendItem onSuccess:', {
        data,
        queryKey: spendTrackingKeys.lists(),
        currentCache: queryClient.getQueryData(spendTrackingKeys.lists())
      });
      queryClient.invalidateQueries({ queryKey: spendTrackingKeys.lists() });
      queryClient.invalidateQueries({ queryKey: dashboardKeys.bff() });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.SPEND_TRACKING.ADD_ITEM,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onMutate: async (newItem) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({ queryKey: spendTrackingKeys.lists() });

      // Snapshot the previous value
      const previousItems = queryClient.getQueryData<SpendTrackingPageData>(spendTrackingKeys.lists());

      // Create a temporary ID for the new item
      const tempItem: SpendingItemBFFResponse = {
        id: 'temp-' + Date.now(),
        title: newItem.title,
        amount: newItem.amount,
        spendTransactionType: 'one_time_purchase' as TransactionType,
        mediaType: 'misc' as MediaCategory, // Default to misc for new items
        paymentMethod: newItem.payment_method,
        isActive: true,
        createdAt: Date.now(),
        updatedAt: Date.now(),
        purchaseDate: new Date(newItem.purchase_date).getTime(),
        isDigital: newItem.is_digital || false,
        isWishlisted: newItem.is_wishlisted || false
      };

      // Optimistically update the cache
      queryClient.setQueryData<SpendTrackingPageData>(spendTrackingKeys.lists(), (old) => {
        if (!old) return { currentMonthItems: [tempItem], nextMonthItems: [], yearlyTotals: [] };
        return {
          ...old,
          currentMonthItems: [...old.currentMonthItems, tempItem]
        };
      });

      return { previousItems };
    },
    onError: (error) => {
      console.log('❔🔎 useCreateSpendItem query onError, error - ', error);
      showToast({
        message: TOAST_ERROR_MESSAGES.SPEND_TRACKING.ADD_ITEM.DEFAULT,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onSettled: () => {
      // Always refetch after error or success
      queryClient.invalidateQueries({ queryKey: spendTrackingKeys.lists() });
    }
  });
};

/**
 * Hook to update a spend item
 */
export const useUpdateSpendItem = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<SpendItem> }) =>
      updateSpendTrackingItem(id, data),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: spendTrackingKeys.lists() });
      queryClient.invalidateQueries({ queryKey: spendTrackingKeys.detail(String(data.id)) });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.SPEND_TRACKING.UPDATE_ITEM,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error) => {
      console.log('❔🔎 useUpdateSpendItem query onError, error - ', error);
      showToast({
        message: TOAST_ERROR_MESSAGES.SPEND_TRACKING.UPDATE_ITEM.DEFAULT,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    }
  });
};

/**
 * Hook to delete a spend item
 */
export const useDeleteSpendItem = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteSpendTrackingItem,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: spendTrackingKeys.lists() });
      queryClient.invalidateQueries({ queryKey: spendTrackingKeys.detail(id) });

      showToast({
        message: TOAST_SUCCESS_MESSAGES.SPEND_TRACKING.DELETE_ITEM,
        variant: 'success',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    },
    onError: (error) => {
      console.log('❔🔎 useDeleteSpendItem query onError, error - ', error);
      showToast({
        message: TOAST_ERROR_MESSAGES.SPEND_TRACKING.DELETE_ITEM.DEFAULT,
        variant: 'error',
        duration: TOAST_DURATIONS.EXTENDED,
      });
    }
  });
};