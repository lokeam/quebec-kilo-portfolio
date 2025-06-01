import { useCallback } from 'react';

// Shadcn UI Components
import { Switch } from "@/shared/components/ui/switch"
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
import { Input } from "@/shared/components/ui/input";

// Custom Components
import { ResponsiveCombobox } from '@/shared/components/ui/ResponsiveCombobox/ResponsiveCombobox';
import { FormContainer } from '@/features/dashboard/components/templates/FormContainer';

// Hooks
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { useCreateDigitalLocation } from '@/core/api/queries/digitalLocation.queries';
//import { useCreateOnlineService } from '@/core/api/queries/useOnlineServiceMutations';


// Zod
import { z } from "zod";

// Icons
import { CalendarIcon } from 'lucide-react';

// Utils
import { format } from 'date-fns';
import { cn } from '@/shared/components/ui/utils';

// Constants
import { PAYMENT_METHODS } from '@/shared/constants/payment';

// Types
import type { SelectableItem } from '@/shared/components/ui/ResponsiveCombobox/ResponsiveCombobox';

// Utils, OnlineService
import { ServiceCombobox } from '@/features/dashboard/components/organisms/OnlineServicesPage/ServiceCombobox/ServiceCombobox';

export type FormValues = {
  name: string;
  isActive: boolean;
  url: string;
  billingCycle: string;
  costPerCycle: number;
  nextPaymentDate: Date;
  paymentMethod: string;
  isSubscriptionService: boolean;
};

const DEFAULT_FORM_VALUES: FormValues = {
  name: "",
  isActive: true,
  url: "",
  billingCycle: "",
  costPerCycle: 0,
  nextPaymentDate: new Date(),
  paymentMethod: "",
  isSubscriptionService: false
};

export const OnlineServiceFormSchema = z.discriminatedUnion('isSubscriptionService', [
  // Non-subscription service schema
  z.object({
    isSubscriptionService: z.literal(false),
    name: z.string().min(1, "Service name is required"),
    url: z.string(),
    paymentMethod: z.string().min(1, "Payment method is required"),
  }),
  // Subscription service schema
  z.object({
    isSubscriptionService: z.literal(true),
    name: z.string().min(1, "Service name is required"),
    url: z.string(),
    isActive: z.boolean(),
    billingCycle: z.string().min(1, "Billing cycle is required"),
    costPerCycle: z.number().min(0, "Cost must be greater than 0"),
    nextPaymentDate: z.date(),
    paymentMethod: z.string().min(1, "Payment method is required"),
  })
]);

interface OnlineServiceFormProps {
  onSuccess?: (data: FormValues) => void;
  onClose?: () => void;
  defaultValues?: Partial<FormValues>;
  buttonText?: string;
}

export function OnlineServiceForm({
  buttonText = "Add Service",
  defaultValues = DEFAULT_FORM_VALUES,
  onClose,
}: OnlineServiceFormProps) {
  const createMutation = useCreateDigitalLocation();

  const isLoading = createMutation.isPending;

  const form = useForm<FormValues>({
    resolver: zodResolver(OnlineServiceFormSchema),
    defaultValues
  });

  const { watch, formState: { isValid, errors } } = form;
  const isSubscriptionService = watch('isSubscriptionService');

  const handleSubmit = useCallback((data: FormValues) => {
    const servicePayload = {
      name: data.name,
      isActive: data.isSubscriptionService ? data.isActive : true,
      url: data.url,
      isSubscription: data.isSubscriptionService,
      payment_method: data.paymentMethod,
      ...(data.isSubscriptionService
        ? {
            subscription: {
              billing_cycle: data.billingCycle,
              cost_per_cycle: data.costPerCycle,
              next_payment_date: data.nextPaymentDate.toISOString(),
              payment_method: data.paymentMethod
            }
          }
        : {}
      )
    };

    createMutation.mutate(servicePayload, {
      onSuccess: () => {
        if (onClose) onClose();
      }
    });
  }, [createMutation, onClose]);

  return (
    <FormContainer form={form} onSubmit={handleSubmit}>
      <div className="max-h-[calc(100vh-12rem)] overflow-y-auto pr-4 -mr-4 [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-track]:bg-transparent [&::-webkit-scrollbar-thumb]:bg-[#086fe8] [&::-webkit-scrollbar-thumb]:rounded-full space-y-5">
        {/* Service Name */}
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Service name <span className="text-red-500">*</span></FormLabel>
              <ServiceCombobox
                onServiceSelect={(service) => {
                  field.onChange(service.name);
                  form.setValue('url', service.url);
                  form.setValue('isSubscriptionService', service.isSubscriptionService);
                  form.trigger('name');
                }}
                placeholder="Search for a service..."
                emptyMessage="No services found."
              />
              <FormMessage>{errors.name?.message}</FormMessage>
            </FormItem>
          )}
        />

        {/* Payment Method is always required */}
        <FormField
              control={form.control}
              name="paymentMethod"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Payment method <span className="text-red-500">*</span></FormLabel>
                  <FormDescription>
                    How do you make payments for this service?
                  </FormDescription>
                  <ResponsiveCombobox
                    onSelect={(method: SelectableItem) => {
                      field.onChange(method.id);
                      form.trigger('paymentMethod');
                    }}
                    items={Object.values(PAYMENT_METHODS)}
                    placeholder="Select a Payment Method"
                    emptyMessage="No payment methods found."
                  />
                  <FormMessage>{errors.paymentMethod?.message}</FormMessage>
                </FormItem>
              )}
            />

        {/* Subscription Fields - Only show for subscription services */}
        {isSubscriptionService && (
          <>
            {/* Active Status */}
            <FormField
              control={form.control}
              name="isActive"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
                  <div className="space-y-0.5">
                    <FormLabel>Active Status</FormLabel>
                    <FormDescription>
                      Is this service currently active?
                    </FormDescription>
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            {/* Billing Cycle */}
            <FormField
              control={form.control}
              name="billingCycle"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Billing cycle <span className="text-red-500">*</span></FormLabel>
                  <FormDescription>
                    What is your subscription plan?
                  </FormDescription>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="1 month">1 month</SelectItem>
                      <SelectItem value="3 month">3 month</SelectItem>
                      <SelectItem value="6 month">6 month</SelectItem>
                      <SelectItem value="12 month">12 month</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage>{errors.billingCycle?.message}</FormMessage>
                </FormItem>
              )}
            />

            {/* Cost */}
            <FormField
              control={form.control}
              name="costPerCycle"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Cost per cycle <span className="text-red-500">*</span></FormLabel>
                  <FormDescription>
                    How much do you pay per cycle?
                  </FormDescription>
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
                  <FormMessage>{errors.costPerCycle?.message}</FormMessage>
                </FormItem>
              )}
            />

            {/* Next Payment Date */}
            <FormField
              control={form.control}
              name="nextPaymentDate"
              render={({ field }) => (
                <FormItem className="flex flex-col">
                  <FormLabel>Next payment date <span className="text-red-500">*</span></FormLabel>
                  <FormDescription>
                    When is your next payment due?
                  </FormDescription>
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
                  <FormMessage>{errors.nextPaymentDate?.message}</FormMessage>
                </FormItem>
              )}
            />


          </>
        )}
      </div>

      {/* Submit Button - Always visible at the bottom */}
      <div className="flex justify-between w-full mt-6 sticky bottom-0 bg-background pt-4 border-t">
        <Button
          type="submit"
          className="w-full"
          disabled={!isValid || isLoading}
        >
          {isLoading ? (
            <>
              <span className="animate-spin mr-2">⊚</span>
              Creating...
            </>
          ) : (
            buttonText
          )}
        </Button>
      </div>
    </FormContainer>
  );
}
