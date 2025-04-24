import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import {
  createOnlineService,
  updateOnlineService,
  deleteOnlineService,
  type CreateOnlineServiceRequest
} from '@/core/api/services/onlineServices.service';
import { AxiosError } from 'axios';

interface MutationOptions {
  onSuccessCallback?: () => void;
}

export function useCreateOnlineService(options?: MutationOptions) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (serviceData: CreateOnlineServiceRequest) => createOnlineService(serviceData),
    onSuccess: (_, variables) => {
      // Access serviceData from variables parameter
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
    onSuccess: (_, variables) => {
      // Access original serviceData from variables
      const serviceData = variables;

      // Invalidate services queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['services'] });

      // Invalidate digitalLocations query in order to update the UI
      queryClient.invalidateQueries({ queryKey: [`digitalLocations`] });

      // Show success toast
      toast.success(`${serviceData.name} updated successfully!`, {
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
    onError: (error: AxiosError) => {
      // Determine the type of error and show appropriate message
      let errorMessage = 'Failed to delete service';
      let errorDescription = "Something went wrong. We can't complete this operation now, please try again later.";

      // Handle specific error types
      if (error.response) {
        // Server returned an error response (4xx, 5xx)
        const status = error.response.status;

        if (status === 401 || status === 403) {
          errorMessage = 'Permission denied';
          errorDescription = "You don't have permission to delete this service.";
        } else if (status === 404) {
          errorMessage = 'Service not found';
          errorDescription = "The service you're trying to delete doesn't exist or was already deleted.";
        } else if (status >= 500) {
          errorMessage = 'Server error';
          errorDescription = "The server encountered an error. Please try again later.";
        }
      } else if (error.request) {
        // Request was made but no response received (network error)
        errorMessage = 'Network error';
        errorDescription = "Couldn't connect to the server. Please check your internet connection.";
      }

      // Show error toast with appropriate message
      toast.error(errorMessage, {
        description: errorDescription,
        duration: 15000
      });

      console.error('Delete service error:', error);
    }
  });
}