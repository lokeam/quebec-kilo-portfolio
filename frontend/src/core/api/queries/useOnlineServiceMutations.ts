import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import {
  createOnlineService,
  updateOnlineService,
  deleteOnlineService,
  type CreateOnlineServiceRequest
} from '@/core/api/services/onlineServices.service';

export function useCreateOnlineService() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (serviceData: CreateOnlineServiceRequest) => createOnlineService(serviceData),
    onSuccess: (data) => {
      // Invalidate services queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['services'] });

      // Show success toast
      toast.success(`${data.name} added to your services!`);
    },
    onError: (error) => {
      toast.error('Failed to add service');
      console.error('Add service error:', error);
    }
  });
}

export function useUpdateOnlineService() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (serviceData: CreateOnlineServiceRequest) => updateOnlineService(serviceData),
    onSuccess: (data) => {
      // Invalidate services queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['services'] });

      // Show success toast
      toast.success(`${data.name} updated successfully!`);
    },
    onError: (error) => {
      toast.error('Failed to update service');
      console.error('Update service error:', error);
    }
  });
}

export function useDeleteOnlineService() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (serviceId: string) => deleteOnlineService(serviceId),
    onSuccess: () => {
      // Invalidate services queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['services'] });

      // Show success toast
      toast.success('Service deleted successfully');
    },
    onError: (error) => {
      toast.error('Failed to delete service');
      console.error('Delete service error:', error);
    }
  });
}