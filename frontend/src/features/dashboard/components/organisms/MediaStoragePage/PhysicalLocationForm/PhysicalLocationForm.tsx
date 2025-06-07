import { useCallback, useState } from 'react';
import { useForm, FormProvider as HookFormProvider } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useCreatePhysicalLocation } from '@/core/api/queries/physicalLocation.queries';
import type { CreatePhysicalLocationRequest } from '@/types/domain/physical-location';

// Shadcn UI Components
import { Button } from "@/shared/components/ui/button"
import { Input } from "@/shared/components/ui/input"
import { Switch } from "@/shared/components/ui/switch"
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

// Icons
import { IconCar } from '@tabler/icons-react';
import { House, Building, Building2, Warehouse } from 'lucide-react';

export const PhysicalLocationFormSchema = z.object({
  locationName: z.string().min(3, 'Location name must be at least 3 characters'),
  locationType: z.enum(['house', 'apartment', 'office', 'warehouse', 'vehicle'] as const),
  bgColor: z.enum(['red', 'green', 'blue', 'orange', 'gold', 'purple', 'brown', 'gray', 'pink'] as const),
  coordinates: z.object({
    enabled: z.boolean(),
    value: z.string()
  })
});

type FormValues = z.infer<typeof PhysicalLocationFormSchema>;

interface PhysicalLocationData {
  id: string;
  name: string;
  locationType: string;
  bgColor?: string;
  mapCoordinates?: string;
  createdAt?: Date;
  updatedAt?: Date;
}

interface PhysicalLocationFormProps {
  onSuccess?: (data: FormValues) => void;
  defaultValues?: FormValues;
  buttonText?: string;
  locationData?: PhysicalLocationData;
  isEditing?: boolean;
  onDelete?: (id: string) => void;
  onClose?: () => void;
}

export function PhysicalLocationForm({
  buttonText = "Add Location",
  defaultValues = {
    locationName: '',
    locationType: 'house' as const,
    bgColor: 'red' as const,
    coordinates: {
      enabled: false,
      value: ''
    }
  },
  locationData,
  isEditing = false,
  onDelete,
  onClose,
  onSuccess,
}: PhysicalLocationFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const methods = useForm<FormValues>({
    resolver: zodResolver(PhysicalLocationFormSchema),
    defaultValues,
    mode: 'onTouched'
  });

  const { handleSubmit, control } = methods;
  const createMutation = useCreatePhysicalLocation();

  const onSubmit = useCallback((data: FormValues) => {
    const locationPayload: CreatePhysicalLocationRequest = {
      name: data.locationName,
      locationType: data.locationType,
      type: data.locationType,
      bgColor: data.bgColor,
      mapCoordinates: data.coordinates.enabled ? data.coordinates.value : undefined
    };

    createMutation.mutate(locationPayload, {
      onSuccess: () => {
        if (onClose) onClose();
        if (onSuccess) onSuccess(data);
      }
    });
  }, [createMutation, onSuccess, onClose]);

  const handleDelete = useCallback((id: string) => {
    if (onDelete) onDelete(id);
    setDeleteDialogOpen(false);
  }, [onDelete]);

  return (
    <HookFormProvider {...methods}>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
        {/* Location Name */}
        <FormField
          control={control}
          name="locationName"
          render={({ field, fieldState: { error } }) => (
            <FormItem>
              <FormLabel>Location Name</FormLabel>
              <FormControl>
                <Input
                  placeholder="Enter a location name"
                  {...field}
                  aria-invalid={!!error}
                />
              </FormControl>
              <FormDescription>
                This is the name of the location where the media is stored.
              </FormDescription>
              {error && <FormMessage>{error.message}</FormMessage>}
            </FormItem>
          )}
        />

        {/* Location Type */}
        <FormField
          control={control}
          name="locationType"
          render={({ field, fieldState: { error } }) => (
            <FormItem>
              <FormLabel>Location Type</FormLabel>
              <Select
                onValueChange={field.onChange}
                defaultValue={field.value}
              >
                <FormControl>
                  <SelectTrigger aria-invalid={!!error}>
                    <SelectValue placeholder="Select a location type" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="house">
                    <div className="flex items-center gap-2">
                      <House size={20} color='#fff' className='mr-2'/>
                      <span>House</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="apartment">
                    <div className="flex items-center gap-2">
                      <Building size={20} color='#fff' className='mr-2'/>
                      <span>Apartment</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="office">
                    <div className="flex items-center gap-2">
                      <Building2 size={20} color='#fff' className='mr-2'/>
                      <span>Office</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="warehouse">
                    <div className="flex items-center gap-2">
                      <Warehouse size={20} color='#fff' className='mr-2'/>
                      <span>Warehouse</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="vehicle">
                    <div className="flex items-center gap-2">
                      <IconCar size={20} color='#fff' className='mr-2'/>
                      <span>Vehicle</span>
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
              <FormDescription>
                Think of this as the venue where the media is stored.
              </FormDescription>
              {error && <FormMessage>{error.message}</FormMessage>}
            </FormItem>
          )}
        />

        {/* Icon BG Color */}
        <FormField
          control={control}
          name="bgColor"
          render={({ field, fieldState: { error } }) => (
            <FormItem>
              <FormLabel>Icon Background Color</FormLabel>
              <Select
                onValueChange={field.onChange}
                defaultValue={field.value}
              >
                <FormControl>
                  <SelectTrigger aria-invalid={!!error}>
                    <SelectValue placeholder="Select a background color for your icon" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="red">Red</SelectItem>
                  <SelectItem value="blue">Blue</SelectItem>
                  <SelectItem value="green">Green</SelectItem>
                  <SelectItem value="gold">Gold</SelectItem>
                  <SelectItem value="purple">Purple</SelectItem>
                  <SelectItem value="orange">Orange</SelectItem>
                  <SelectItem value="brown">Brown</SelectItem>
                  <SelectItem value="gray">Gray</SelectItem>
                  <SelectItem value="pink">Pink</SelectItem>
                </SelectContent>
              </Select>
              <FormDescription>
                Customize the background color of your location icon.
              </FormDescription>
              {error && <FormMessage>{error.message}</FormMessage>}
            </FormItem>
          )}
        />

        {/* Coordinates */}
        <FormField
          control={control}
          name="coordinates"
          render={({ field, fieldState: { error } }) => (
            <FormItem>
              <FormLabel>Map Coordinates</FormLabel>
              <div className="flex items-center space-x-2">
                <Switch
                  checked={field.value.enabled}
                  onCheckedChange={(checked: boolean) => {
                    field.onChange({ ...field.value, enabled: checked });
                  }}
                />
                <FormLabel>Enable coordinates</FormLabel>
              </div>
              {field.value.enabled && (
                <FormControl>
                  <Input
                    placeholder="Enter coordinates (e.g., 45.5017,-73.5673)"
                    value={field.value.value}
                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                      field.onChange({ ...field.value, value: event.target.value });
                    }}
                  />
                </FormControl>
              )}
              <FormDescription>
                Optional: Add map coordinates for this location.
              </FormDescription>
              {error && <FormMessage>{error.message}</FormMessage>}
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
              ? "Creating..."
              : isEditing
                ? "Update Location"
                : buttonText}
          </Button>

          {isEditing && onDelete && locationData && (
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
      {isEditing && onDelete && locationData && (
        <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Confirm Deletion</DialogTitle>
              <DialogDescription>
                Are you sure you want to delete this location? This action cannot be undone.
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
                onClick={() => handleDelete(locationData.id)}
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
