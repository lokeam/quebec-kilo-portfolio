import { useCallback, useState } from 'react';
import { useForm, FormProvider as HookFormProvider } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useCreateSpendItem } from '@/core/api/queries/spendTracking.queries';
import { useGetAllDigitalLocations } from '@/core/api/queries/digitalLocation.queries';
import type { CreateOneTimePurchaseRequest } from '@/types/domain/spend-tracking';

// Shadcn UI Components
import { Button } from "@/shared/components/ui/button"
import { Calendar } from '@/shared/components/ui/calendar';
import { Input } from "@/shared/components/ui/input"
import { Checkbox } from "@/shared/components/ui/checkbox"
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
import { Dialog, DialogContent, DialogFooter, DialogTitle, DialogHeader, DialogDescription } from '@/shared/components/ui/dialog';
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '@/shared/components/ui/popover';

// Utils
import { cn } from '@/shared/components/ui/utils';
import { format } from 'date-fns';

// Icons
import { IconCloudDataConnection, IconDeviceGamepad2, IconPackage, IconDisc, IconCpu2, IconSparkles } from '@tabler/icons-react';
import { CalendarIcon } from 'lucide-react';

// Constants
import { PAYMENT_METHODS } from '@/shared/constants/payment';

// Custom Components
import { ResponsiveCombobox } from '@/shared/components/ui/ResponsiveCombobox/ResponsiveCombobox';
//import { FormContainer } from '@/features/dashboard/components/templates/FormContainer';

// Types
import type { SelectableItem } from '@/shared/components/ui/ResponsiveCombobox/ResponsiveCombobox';

export const SpendTrackingFormSchema = z.object({
  title: z.string().min(3, 'Title must be at least 3 characters'),
  spending_category_id: z.number().min(1, 'Category is required'),
  amount: z.number().min(0.01, "Cost must be greater than 0"),
  payment_method: z.string().min(1, "Payment method is required"),
  purchase_date: z.date(),
  digital_location_id: z.string().optional(),
  is_wishlisted: z.boolean().default(false),
});

type FormValues = z.infer<typeof SpendTrackingFormSchema>;

interface SpendTrackingData {
  id: string;
  title: string;
  spending_category_id: number;
  payment_method: string;
  amount: number;
  purchase_date?: Date;
  digital_location_id?: string;
  is_wishlisted?: boolean;
}

interface SpendTrackingFormProps {
  onSuccess?: (data: FormValues) => void;
  defaultValues?: FormValues;
  buttonText?: string;
  spendTrackingData?: SpendTrackingData;
  isEditing?: boolean;
  onDelete?: (id: string) => void;
  onClose?: () => void;
}

// Spending category options with IDs matching the database
const SPENDING_CATEGORIES = [
  { id: 1, name: 'Hardware', icon: IconCpu2 },
  { id: 2, name: 'DLC', icon: IconDeviceGamepad2 },
  { id: 3, name: 'In-Game Purchase', icon: IconSparkles },
  { id: 4, name: 'Physical Game', icon: IconDisc },
  { id: 5, name: 'Digital Game', icon: IconCloudDataConnection },
  { id: 6, name: 'Misc', icon: IconPackage },
];

export function SpendTrackingForm({
  buttonText = "Add Expense",
  defaultValues = {
    title: '',
    spending_category_id: 5, // Default to Digital Game
    amount: 0,
    payment_method: '',
    purchase_date: new Date(),
    is_wishlisted: false,
  },
  spendTrackingData,
  isEditing = false,
  onDelete,
  onClose,
  onSuccess,
}: SpendTrackingFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  // Create form default values based on whether we're editing or creating
  const formDefaultValues = isEditing && spendTrackingData ? {
    title: spendTrackingData.title,
    spending_category_id: spendTrackingData.spending_category_id,
    amount: spendTrackingData.amount,
    payment_method: spendTrackingData.payment_method,
    purchase_date: spendTrackingData.purchase_date || new Date(),
    digital_location_id: spendTrackingData.digital_location_id,
    is_wishlisted: spendTrackingData.is_wishlisted || false,
  } : defaultValues;

  const form = useForm<FormValues>({
    resolver: zodResolver(SpendTrackingFormSchema),
    defaultValues: formDefaultValues,
    mode: 'onTouched'
  });

  const { handleSubmit, control, watch, formState: { errors } } = form;
  const createMutation = useCreateSpendItem();

  // Get digital locations for dropdown
  const { data: digitalLocations = [] } = useGetAllDigitalLocations();

  // Watch spending category to determine if digital
  const selectedCategory = watch('spending_category_id');
  const isDigitalCategory = selectedCategory === 5; // Digital Game

  const onSubmit = useCallback((data: FormValues) => {
    const payload: CreateOneTimePurchaseRequest = {
      title: data.title,
      spending_category_id: data.spending_category_id,
      amount: data.amount,
      payment_method: data.payment_method,
      purchase_date: data.purchase_date.toISOString(),
      digital_location_id: data.digital_location_id,
      is_wishlisted: data.is_wishlisted,
      is_digital: isDigitalCategory || !!data.digital_location_id,
    };

    createMutation.mutate(payload, {
      onSuccess: () => {
        if (onClose) onClose();
        if (onSuccess) onSuccess(data);
      }
    });
  }, [createMutation, isDigitalCategory, onSuccess, onClose]);

  const handleDelete = useCallback((id: string) => {
    if (onDelete) onDelete(id);
    setDeleteDialogOpen(false);
  }, [onDelete]);

  return (
    <HookFormProvider {...form}>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
        {/* Title */}
        <FormField
          control={control}
          name="title"
          render={({ field, fieldState: { error } }) => (
            <FormItem>
              <FormLabel>Title</FormLabel>
              <FormControl>
                <Input
                  placeholder="What did you buy?"
                  {...field}
                  aria-invalid={!!error}
                />
              </FormControl>
              <FormDescription>
                This is the name of the one-time purchase.
              </FormDescription>
              {error && <FormMessage>{error.message}</FormMessage>}
            </FormItem>
          )}
        />

        {/* Spending Category */}
        <FormField
          control={control}
          name="spending_category_id"
          render={({ field, fieldState: { error } }) => (
            <FormItem>
              <FormLabel>Category</FormLabel>
              <Select
                onValueChange={(value) => field.onChange(parseInt(value))}
                value={field.value.toString()}
              >
                <FormControl>
                  <SelectTrigger aria-invalid={!!error}>
                    <SelectValue placeholder="Select a category" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {SPENDING_CATEGORIES.map((category) => {
                    const IconComponent = category.icon;
                    return (
                      <SelectItem key={category.id} value={category.id.toString()}>
                        <div className="flex items-center gap-2">
                          <IconComponent size={20} color='#fff' className='mr-2'/>
                          <span>{category.name}</span>
                        </div>
                      </SelectItem>
                    );
                  })}
                </SelectContent>
              </Select>
              <FormDescription>
                What type of purchase is this?
              </FormDescription>
              {error && <FormMessage>{error.message}</FormMessage>}
            </FormItem>
          )}
        />

        {/* Digital Location (only show for digital categories) */}
        {isDigitalCategory && (
          <FormField
            control={control}
            name="digital_location_id"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Digital Platform (Optional)</FormLabel>
                <Select
                  onValueChange={field.onChange}
                  value={field.value || undefined}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select a digital platform" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {digitalLocations.map((location) => (
                      <SelectItem key={location.id} value={location.id}>
                        {location.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>
                  Which digital platform did you purchase this on?
                </FormDescription>
              </FormItem>
            )}
          />
        )}

        {/* Payment Method */}
        <FormField
          control={form.control}
          name="payment_method"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Payment method <span className="text-red-500">*</span></FormLabel>
              <FormDescription>
                How did you pay for this purchase?
              </FormDescription>
              <ResponsiveCombobox
                onSelect={(method: SelectableItem) => {
                  field.onChange(method.id.toLowerCase());
                  form.trigger('payment_method');
                }}
                items={Object.values(PAYMENT_METHODS)}
                placeholder="Select a Payment Method"
                emptyMessage="No payment methods found."
                initialValue={field.value}
              />
              <FormMessage>{errors.payment_method?.message}</FormMessage>
            </FormItem>
          )}
        />

        {/* Amount */}
        <FormField
          control={form.control}
          name="amount"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Cost <span className="text-red-500">*</span></FormLabel>
              <FormDescription>
                How much did you spend?
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
              <FormMessage>{errors.amount?.message}</FormMessage>
            </FormItem>
          )}
        />

        {/* Purchase Date */}
        <FormField
          control={form.control}
          name="purchase_date"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel>Purchase date <span className="text-red-500">*</span></FormLabel>
              <FormDescription>
                When did you make this purchase?
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
                      date > new Date()
                    }
                    initialFocus
                  />
                </PopoverContent>
              </Popover>
              <FormMessage>{errors.purchase_date?.message}</FormMessage>
            </FormItem>
          )}
        />

        {/* Wishlist Checkbox */}
        <FormField
          control={form.control}
          name="is_wishlisted"
          render={({ field }) => (
            <FormItem className="flex flex-row items-start space-x-3 space-y-0">
              <FormControl>
                <Checkbox
                  checked={field.value}
                  onCheckedChange={field.onChange}
                />
              </FormControl>
              <div className="space-y-1 leading-none">
                <FormLabel>Add to wishlist</FormLabel>
                <FormDescription>
                  Mark this as a wishlist item for future reference.
                </FormDescription>
              </div>
            </FormItem>
          )}
        />

        {/* Form actions */}
        <div className="flex justify-between w-full mt-6">
          <Button
            type="submit"
            className={isEditing && onDelete ? "flex-1" : "w-full"}
            disabled={createMutation.isPending}
          >
            {createMutation.isPending
              ? isEditing ? "Updating..." : "Creating..."
              : isEditing
                ? "Update Purchase"
                : buttonText}
          </Button>

          {isEditing && onDelete && spendTrackingData && (
            <Button
              type="button"
              variant="destructive"
              className="ml-2"
              onClick={() => setDeleteDialogOpen(true)}
              disabled={createMutation.isPending}
            >
              Delete
            </Button>
          )}
        </div>
      </form>

      {/* Delete confirmation dialog */}
      {isEditing && onDelete && spendTrackingData && (
        <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Confirm Deletion</DialogTitle>
              <DialogDescription>
                Are you sure you want to delete this purchase? This action cannot be undone.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setDeleteDialogOpen(false)}
              >
                Cancel
              </Button>
              <Button
                variant="destructive"
                onClick={() => handleDelete(spendTrackingData.id)}
              >
                Delete
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </HookFormProvider>
  );
}