import { Suspense } from 'react';

// Components
import { ErrorBoundary } from 'react-error-boundary';
import { OnlineServicesPageErrorFallback } from '@/features/dashboard/pages/OnlineServices/OnlineServicesPageErrorFallback';
import { HomePageSkeleton } from '@/features/dashboard/pages/HomePage/HomePageSkeleton'
//import { OnlineServiceForm } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServiceForm/OnlineServiceForm';
import { ServiceCombobox } from '@/features/dashboard/components/organisms/ServiceCombobox/ServiceCombobox';
import { StateControlledButton } from '@/shared/components/ui/StateControlledButton/StateControlledButton';
import { DebouncedResponsiveCombobox } from '@/shared/components/ui/DebouncedResponsiveCombobox/DebouncedResponsiveCombobox';

// ShadCN UI Components
import { Button } from '@/shared/components/ui/button';
import {
  Form,
  FormControl,
//  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/components/ui/form"
import { Input } from '@/shared/components/ui/input';
//import { Label } from '@/shared/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/components/ui/select';
import { Calendar } from '@/shared/components/ui/calendar';
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '@/shared/components/ui/popover';

// Hooks + Utils
import { useNavigate } from 'react-router-dom';
import { cn } from '@/shared/components/ui/utils';
import { format } from 'date-fns';
import { useForm } from "react-hook-form"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"

// Icons
import { CalendarIcon } from 'lucide-react';

// Constants
import { NAVIGATION_ROUTES } from '@/features/dashboard/lib/types/onboarding/constants';
import { PAYMENT_METHODS } from '@/shared/constants/payment';
import { BILLING_CYCLES } from '@/features/dashboard/lib/types/spend-tracking/constants';

// Types
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
//import type { BillingDetails } from '@/features/dashboard/lib/types/online-services/services';
//import type { ServiceType } from '@/shared/constants/service.constants';
import type { SelectableItem } from '@/shared/components/ui/DebouncedResponsiveCombobox/DebouncedResponsiveCombobox';


// Tightly Coupled Helper fn:
interface FormState {
  selectedService: any;
  expenseType: string;
  cost?: number;
  billingPeriod?: string;
  nextPayment?: Date;
  paymentMethod?: SelectableItem;
}

const validateFormProgress = (formState: FormState):boolean => {
  const {
    selectedService,
    expenseType,
    cost,
    billingPeriod,
    nextPayment,
    paymentMethod,
  } = formState;

  // First check if service exists
  if (!selectedService) return false;

  // For free services, only require service selection
  if (selectedService?.billing?.cycle === 'NA') {
    return true;
  }

  // For paid services, validate required fields
  if (!expenseType) return false;
  if (typeof cost !== 'number' || cost <= 0) return false;
  if (!nextPayment) return false;

  // Set billingPeriod to expenseType if not set
  const effectiveBillingPeriod = billingPeriod || expenseType;
  if (!effectiveBillingPeriod) return false;

  // Check payment method last
  if (!paymentMethod?.id) return false;

  return true;
}

const formSchema = z.object({
  service: z.custom<OnlineService | null>()
    .nullable()
    .refine((service) => service !== null, {
      message: "Please select a service"
    }),
  expenseType: z.string().optional().refine(
    (val) => !val || ['1 month', '3 months', '6 months', '1 year'].includes(val),
    "Please select a valid expense type"
  ),
  cost: z.number()
    .min(0, "Cost must be at least 0")
    .default(0),
  billingPeriod: z.string().optional(),
  nextPaymentDate: z.date()
    .optional()
    .refine((date) => {
      if (!date) return true;
      return date >= new Date();
    }, "Payment date must be in the future"),
  paymentMethod: z.custom<SelectableItem>()
    .optional()
    .refine((method) => !method || (method.id && method.displayName), {
      message: "Please select a valid payment method"
    })
});

interface FormValues {
  service: OnlineService | null;
  expenseType: string | undefined;
  cost: number;
  billingPeriod?: string;
  nextPaymentDate?: Date;
  paymentMethod?: SelectableItem;
}


export default function OnboardingPageDigital() {

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      service: null,
      expenseType: undefined,
      cost: 0, // Note: This needs to be initialized to 0 to avoid validation errors
      billingPeriod: undefined,
      nextPaymentDate: undefined,
      paymentMethod: undefined,
    }
  });
  const navigate = useNavigate();
  const { watch, formState: {
    isValid,
    errors,
  }} = form;
  const selectedService = watch('service');
  const expenseType = watch('expenseType');
  const cost = watch('cost');
  //const billingPeriod = watch('billingPeriod');

  //const [date, setDate] = useState<Date>();
  const isFreePlan =
    !selectedService?.billing?.cycle
    || selectedService?.billing?.cycle === BILLING_CYCLES.NA;

  const shouldShowExpenseSection = !!selectedService;
  const shouldShowCostSection = isFreePlan ? false : !!expenseType;
  const shouldShowBillingSection = isFreePlan ? false : !!cost;


  const isFormValid = isValid && validateFormProgress({
    selectedService: watch('service'),
    expenseType: watch('expenseType') || '',
    cost: watch('cost'),
    billingPeriod: watch('billingPeriod'),
    nextPayment: watch('nextPaymentDate'),
    paymentMethod: watch('paymentMethod') ,
  });

  const onSubmit = (data: FormValues) => {
    console.log('Form data:', data);
    navigate(NAVIGATION_ROUTES.ONBOARDING_COMPLETE);
  }

  console.log('Form validation:', {
    isValid,
    formProgress: validateFormProgress({
      selectedService: watch('service'),
      expenseType: watch('expenseType') || '',
      cost: watch('cost'),
      billingPeriod: watch('billingPeriod'),
      nextPayment: watch('nextPaymentDate'),
      paymentMethod: watch('paymentMethod')
    }),
    values: form.getValues()
  });

  return (
    <ErrorBoundary
      FallbackComponent={OnlineServicesPageErrorFallback}
      resetKeys={[location.pathname]}
    >
      <Suspense fallback={<HomePageSkeleton />}>
        <div className="mx-auto flex h-screen max-w-3xl flex-col items-center justify-center overflow-x-hidden">
          <h1 className="text-3xl font-bold mb-3">Now let's add a cloud service</h1>
          <p className="text-lg mb-2">(e.g. "Steam", "Apple Arcade", "Xbox Game Pass", etc.)</p>
          <p className="text-sm mb-4 text-muted-foreground">Don't worry, you can always edit this later.</p>

          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmit)}
              className="shrink-0 p2 md:p-4 space-y-6"
            >
              {/* Service Section */}
              <FormField
                control={form.control}
                name="service"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>What is the name of this service? <span className="text-red-500">*</span></FormLabel>
                    <ServiceCombobox
                      onServiceSelect={(service: OnlineService) => {
                        field.onChange(service);
                        form.trigger('service'); // Trigger custom validation when valid
                      }}
                      placeholder="Search for a service..."
                      emptyMessage="No services found."
                    />
                    <FormMessage>{errors.service?.message}</FormMessage>
                  </FormItem>
                )}
              />

              {/* Expense Section */}
              <div className={cn(
                "transition-all duration-700",
                shouldShowExpenseSection ? "fade-in" : "opacity-0 hidden",
                isFreePlan && "hidden"
              )}>
                <FormField
                  control={form.control}
                  name="expenseType"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>How often do you pay to use this service? <span className="text-red-500">*</span></FormLabel>
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
                    </FormItem>
                  )}
                />
              </div>

              {/* Cost Section */}
              <div className={cn(
                  "space-y-2 transition-all duration-700",
                  shouldShowCostSection ? "fade-in" : "opacity-0 hidden"
              )}>
                <FormField
                  control={form.control}
                  name="cost"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Approximately how much is each payment?</FormLabel>
                      <div className="flex">
                        <div className="flex-none flex items-center px-3 border border-r-0 rounded-l-md bg-gray-50">
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
                    </FormItem>
                  )}
                />
              </div>

              {/* Billing Section */}
              <div className={cn(
                "space-y-6 transition-all duration-700",
                shouldShowBillingSection ? "fade-in" : "opacity-0 hidden"
              )}>
                <FormField
                  control={form.control}
                  name="nextPaymentDate"
                  render={({ field }) => (
                    <FormItem className="flex flex-col">
                      <FormLabel>Around what date is your next payment due?</FormLabel>
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
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="paymentMethod"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>How do you pay for this service?</FormLabel>
                      <DebouncedResponsiveCombobox
                        onSelect={(method: SelectableItem) => {
                          field.onChange(method);
                          form.trigger('paymentMethod');
                        }}
                        items={Object.values(PAYMENT_METHODS)}
                        placeholder="Select a Payment Method"
                        emptyMessage="No payment methods found."
                      />
                      <FormMessage>{errors.paymentMethod?.message?.toString()}</FormMessage>
                    </FormItem>
                  )}
                />
              </div>

              <StateControlledButton
                type="submit"
                isDisabled={!isFormValid}
                className={`flex flex-row justify-between items-center mr-1 w-full hover:text-white transition duration-500 ease-in-out disabled:pointer-events-none disabled:border-gray-700/60 disabled:text-gray-600 dark:disabled:text-gray-400`}
              >
                Submit
              </StateControlledButton>

            </form>
          </Form>
        </div>
      </Suspense>
    </ErrorBoundary>
  );
}
