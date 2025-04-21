import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import {
  createOnlineService,
  updateOnlineService,
  deleteOnlineService,
  type CreateOnlineServiceRequest
} from '@/core/api/services/onlineServices.service';

interface MutationOptions {
  onSuccessCallback?: () => void;
}

export function useCreateOnlineService(options?: MutationOptions) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (serviceData: CreateOnlineServiceRequest) => createOnlineService(serviceData),
    onSuccess: (_, variables) => {
      // Hopefully contains variables from teh original serviceData
      const serviceData = variables;

      // Invalidate services queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['services'] });

      // Invalidate digitalLocations query in order to update the UI
      queryClient.invalidateQueries({ queryKey: [`digitalLocations`] });

      // Show success toast
      const serviceName = serviceData.name || "New service";

      toast.success(`${serviceName} added to your services!`, {
        description: "Your new service has been successfully created.",
        duration: 15000
      });

      // Call the optional callback if provided
      if (options?.onSuccessCallback) {
        options.onSuccessCallback();
      }
    },
    onError: (error) => {
      toast.error('Failed to add service');
      console.error('Add service error:', error);
    }
  });
}

export function useUpdateOnlineService(options?: MutationOptions) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (serviceData: CreateOnlineServiceRequest) => updateOnlineService(serviceData),
    onSuccess: (data) => {
      // Invalidate services queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['services'] });

      // Invalidate digitalLocations query in order to update the UI
      queryClient.invalidateQueries({ queryKey: [`digitalLocations`] });

      // Show success toast
      toast.success(`${data.name} updated successfully!`, {
        description: "Your service has been successfully updated.",
        duration: 15000
      });

      // Call the optional callback if provided
      if (options?.onSuccessCallback) {
        options.onSuccessCallback();
      }
    },
    onError: (error: Error) => {
      toast.error('Failed to update service');
      console.error('Update service error:', error);
    }
  });
}

export function useDeleteOnlineService(options?: MutationOptions) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (serviceId: string) => deleteOnlineService(serviceId),
    onSuccess: () => {
      // Invalidate services queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['services'] });

      // Invalidate digitalLocations query in order to update the UI
      queryClient.invalidateQueries({ queryKey: [`digitalLocations`] });

      // Show success toast
      toast.success('Service deleted successfully', {
        description: "Your service has been successfully deleted.",
        duration: 15000
      });

      // Call the optional callback if provided
      if (options?.onSuccessCallback) {
        options.onSuccessCallback();
      }
    },
    onError: (error) => {
      toast.error('Failed to delete service');
      console.error('Delete service error:', error);
    }
  });
}