import { useCallback, useEffect, useState } from 'react';

// Components
import { FormContainer } from '@/features/dashboard/components/templates/FormContainer';

// Shadcn UI Components
// import { Input } from "@/shared/components/ui/input"
// import { Switch } from "@/shared/components/ui/switch"

import { Button } from "@/shared/components/ui/button";

import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/components/ui/form"

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/components/ui/select"

import { Calendar } from '@/shared/components/ui/calendar';

import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '@/shared/components/ui/popover';

// Hooks
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { useCreateOnlineService, useUpdateOnlineService, useDeleteOnlineService } from '@/core/api/queries/useOnlineServiceMutations';

// Zod
import { z } from "zod"

// Icons
import { Dialog, DialogContent, DialogFooter, DialogTitle, DialogHeader, DialogDescription } from '@/shared/components/ui/dialog';
import { CalendarIcon } from 'lucide-react';

import { format } from 'date-fns';
import { cn } from '@/shared/components/ui/utils';
import { Input } from "@/shared/components/ui/input";
import { ResponsiveCombobox } from '@/shared/components/ui/ResponsiveCombobox/ResponsiveCombobox';
import { PAYMENT_METHODS } from '@/shared/constants/payment';

import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import type { SelectableItem } from '@/shared/components/ui/ResponsiveCombobox/ResponsiveCombobox';


// Utils
import { ServiceCombobox } from '../../ServiceCombobox/ServiceCombobox';

export type OnlineServiceFormData = {
  service: OnlineService | null;
  expenseType?: string;
  cost: number;
  billingPeriod?: string;
  nextPaymentDate?: Date;
  paymentMethod?: SelectableItem;
};

type FormValues = {
  service: OnlineService | null;
  expenseType?: string;
  cost: number;
  billingPeriod?: string;
  nextPaymentDate?: Date;
  paymentMethod?: SelectableItem;
};

export const OnlineServiceFormSchema = z.object({
  service: z.custom<OnlineService | null>().refine((service) => service !== null, "Please select a service"),
  expenseType: z.string().optional(),
  cost: z.number().default(0),
  billingPeriod: z.string().optional(),
  nextPaymentDate: z.date().optional(),
  paymentMethod: z.custom<SelectableItem>().optional()
});

const validateFormProgress = (formState: FormValues): boolean => {
  const {
    service,
    expenseType,
    cost,
    billingPeriod,
    nextPaymentDate,
    paymentMethod,
  } = formState;

  // First check if service exists
  if (!service) return false;

  // For non-subscription services, only require service selection
  if (!service.isSubscriptionService) {
    return true;
  }

  // For subscription services, validate required fields
  if (!expenseType) return false;
  if (typeof cost !== 'number' || cost <= 0) return false;
  if (!nextPaymentDate) return false;

  // Set billingPeriod to expenseType if not set
  const effectiveBillingPeriod = billingPeriod || expenseType;
  if (!effectiveBillingPeriod) return false;

  // Check payment method last
  if (!paymentMethod?.id) return false;

  return true;
}

interface OnlineServiceData {
  id: string;
  service: OnlineService | null;
  expenseType?: string;
  cost?: number;
  billingPeriod?: string;
  nextPaymentDate?: Date;
  paymentMethod?: SelectableItem;
  createdAt?: Date;
  updatedAt?: Date;
}

interface OnlineServiceFormProps {
  onSuccess?: (data: OnlineServiceFormData) => void;
  defaultValues?: Partial<FormValues>;
  buttonText?: string;
  serviceData?: OnlineServiceData;
  isEditing?: boolean;
  onDelete?: (id: string) => void;
}

export function OnlineServiceForm({
  buttonText = "Add Service",
  defaultValues = {
    service: null,
    cost: 0
  },
  serviceData,
  isEditing = false,
  onDelete,
}: OnlineServiceFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  // Replace locationManager with mutation hooks using correct property names
  const createMutation = useCreateOnlineService();
  const updateMutation = useUpdateOnlineService();
  const deleteMutation = useDeleteOnlineService();

  const isLoading = createMutation.isPending || updateMutation.isPending || deleteMutation.isPending;

  /* Specific form components creates their own useForm hook instances */
  const form = useForm<FormValues>({
    resolver: zodResolver(OnlineServiceFormSchema),
    defaultValues: isEditing && serviceData
      ? {
        service: serviceData.service,
        expenseType: serviceData.expenseType,
        cost: serviceData.cost || 0,
        billingPeriod: serviceData.billingPeriod,
        nextPaymentDate: serviceData.nextPaymentDate,
        paymentMethod: serviceData.paymentMethod,
      }
      : defaultValues
  });

  const { watch, formState: { isValid, errors } } = form;

  const selectedService = watch('service');
  console.log('Selected Service:', selectedService);
  const isFreePlan = !selectedService?.isSubscriptionService;
  console.log('Is Free Plan:', isFreePlan, 'isSubscriptionService:', selectedService?.isSubscriptionService);

  // If service data changes AND we are editing, update form values
  useEffect(() => {
    if (isEditing && serviceData) {
      form.reset({
        service: serviceData.service,
        expenseType: serviceData.expenseType,
        cost: serviceData.cost || 0,
        billingPeriod: serviceData.billingPeriod,
        nextPaymentDate: serviceData.nextPaymentDate,
        paymentMethod: serviceData.paymentMethod,
      });
    }
  }, [form, isEditing, serviceData]);

  const isFormValid = isValid && validateFormProgress({
    service: watch('service'),
    expenseType: watch('expenseType'),
    cost: watch('cost'),
    billingPeriod: watch('billingPeriod'),
    nextPaymentDate: watch('nextPaymentDate'),
    paymentMethod: watch('paymentMethod'),
  });

  const handleSubmit = useCallback((data: z.infer<typeof OnlineServiceFormSchema>) => {
    const servicePayload = {
      id: isEditing && serviceData ? serviceData.id : undefined,
      name: data.service?.name || 'Unknown Service',
      parentId: null,
      type: 'digital' as const,
      parentLocationId: 'root',
      metadata: {
        service: data.service,
        expenseType: data.expenseType,
        cost: data.cost,
        billingPeriod: data.billingPeriod,
        nextPaymentDate: data.nextPaymentDate,
        paymentMethod: data.paymentMethod,
      }
    };

    if (isEditing && serviceData) {
      updateMutation.mutate(servicePayload);
    } else {
      createMutation.mutate(servicePayload);
    }
  }, [isEditing, serviceData, createMutation, updateMutation]);

  const handleDelete = useCallback((id: string) => {
    deleteMutation.mutate(id);
    setDeleteDialogOpen(false);
  }, [deleteMutation]);

  return (
    <>
      <FormContainer form={form} onSubmit={handleSubmit}>
        {/* Service Selection */}
        <FormField
          control={form.control}
          name="service"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Service name <span className="text-red-500">*</span></FormLabel>
              <ServiceCombobox
                onServiceSelect={(service: OnlineService) => {
                  field.onChange(service);
                  form.trigger('service');
                }}
                placeholder="Search for a service..."
                emptyMessage="No services found."
              />
              <FormMessage>{errors.service?.message}</FormMessage>
            </FormItem>
          )}
        />

        {/* Expense Section - Only show for paid services */}
        {!isFreePlan && (
          <FormField
            control={form.control}
            name="expenseType"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Billing cycle <span className="text-red-500">*</span></FormLabel>
                <Select onValueChange={field.onChange} value={field.value}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="1 month">Monthly</SelectItem>
                    <SelectItem value="3 months">Quarterly</SelectItem>
                    <SelectItem value="6 months">Bi-Annually</SelectItem>
                    <SelectItem value="1 year">Annually</SelectItem>
                  </SelectContent>
                </Select>
                <FormDescription>
                  How often do you pay to use this service?
                </FormDescription>
                <FormMessage>{errors.expenseType?.message}</FormMessage>
              </FormItem>
            )}
          />
        )}

        {/* Cost Section - Only show for paid services */}
        {!isFreePlan && (
          <FormField
            control={form.control}
            name="cost"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Cost <span className="text-red-500">*</span></FormLabel>
                <div className="flex">
                  <div className="flex-none flex items-center px-3 border border-r-0 rounded-l-md bg-transparent">
                    <span className="text-sm text-gray-500">$</span>
                  </div>
                  <FormControl>
                    <Input
                      type="number"
                      step="0.01"
                      min="0"
                      className="rounded-l-none"
                      value={field.value || ''}
                      onChange={(event) => {
                        const value = event.target.value;
                        field.onChange(value === '' ? 0 : parseFloat(value));
                      }}
                      onBlur={field.onBlur}
                    />
                  </FormControl>
                </div>
                <FormDescription>
                  Approximately how much is each payment?
                </FormDescription>
                <FormMessage>{errors.cost?.message}</FormMessage>
              </FormItem>
            )}
          />
        )}

        {/* Billing Section - Only show for paid services */}
        {!isFreePlan && (
          <FormField
            control={form.control}
            name="nextPaymentDate"
            render={({ field }) => (
              <FormItem className="flex flex-col">
                <FormLabel>Next payment date <span className="text-red-500">*</span></FormLabel>
                <Popover>
                  <PopoverTrigger asChild>
                    <FormControl>
                      <Button
                        variant="outline"
                        className={cn(
                          "w-full pl-3 text-left font-normal",
                          !field.value && "text-muted-foreground"
                        )}
                      >
                        {field.value ? (
                          format(field.value, "PPP")
                        ) : (
                          <span>Pick a date</span>
                        )}
                        <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                      </Button>
                    </FormControl>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                      mode="single"
                      selected={field.value}
                      onSelect={field.onChange}
                      disabled={(date) =>
                        date < new Date()
                      }
                      initialFocus
                    />
                  </PopoverContent>
                </Popover>
                <FormDescription>
                  When is your next payment due?
                </FormDescription>
                <FormMessage>{errors.nextPaymentDate?.message}</FormMessage>
              </FormItem>
            )}
          />
        )}

        {/* Payment Method Section - Only show for paid services */}
        {!isFreePlan && (
          <FormField
            control={form.control}
            name="paymentMethod"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Payment method <span className="text-red-500">*</span></FormLabel>
                <ResponsiveCombobox
                  onSelect={(method: SelectableItem) => {
                    field.onChange(method);
                    form.trigger('paymentMethod');
                  }}
                  items={Object.values(PAYMENT_METHODS)}
                  placeholder="Select a Payment Method"
                  emptyMessage="No payment methods found."
                />
                <FormDescription>
                  How do you pay for this service?
                </FormDescription>
                <FormMessage>{errors.paymentMethod?.message?.toString()}</FormMessage>
              </FormItem>
            )}
          />
        )}

        {/* Submit Button */}
        <div className="flex justify-between w-full mt-6">
          <Button
            type="submit"
            className={isEditing && onDelete ? "flex-1" : "w-full"}
            disabled={!isFormValid || isLoading}
          >
            {(createMutation.isPending || updateMutation.isPending) ? (
              <>
                <span className="animate-spin mr-2">⊚</span>
                {isEditing ? "Updating..." : "Creating..."}
              </>
            ) : (
              isEditing ? "Update Service" : buttonText
            )}
          </Button>

          {isEditing && onDelete && serviceData && (
            <Button
              type="button"
              variant="destructive"
              className="ml-2"
              onClick={() => setDeleteDialogOpen(true)}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? "Deleting..." : "Delete"}
            </Button>
          )}
        </div>
      </FormContainer>

      {/* Delete Dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirm Deletion</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this service? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
              disabled={deleteMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={() => handleDelete(serviceData?.id || '')}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? (
                <>
                  <span className="animate-spin mr-2">⊚</span>
                  Deleting...
                </>
              ) : (
                "Delete"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
