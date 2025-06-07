import { useCallback, useEffect, useState } from 'react';

// Shadcn UI Components
import { Input } from "@/shared/components/ui/input"
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

// Hooks
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm, FormProvider as HookFormProvider } from "react-hook-form"

// Zod
import { z } from "zod"

// Icons
import { House, Building, Building2, Warehouse } from 'lucide-react';
import { IconCar } from '@tabler/icons-react';
import { Dialog, DialogContent, DialogFooter, DialogTitle, DialogHeader, DialogDescription } from '@/shared/components/ui/dialog';

// Utils
import { validateCoordinateFormat, parseCoordinates } from '@/features/dashboard/lib/utils/validateCoordinates';

export const PhysicalLocationFormSchema = z.object({
  locationName: z
    .string({
      required_error: "Please enter a location name",
    })
    .min(3, {
      message: "Location name must be at least 3 characters long",
    }),
  locationType: z
    .string({
      required_error: "Please select a location type",
    })
    .min(1, {
      message: "Please select a location type",
    }),
  bgColor: z
    .string({
      required_error: "Please select a background color",
    })
    .min(1, {
      message: "Please select a background color",
    }),
  coordinates: z.object({
    enabled: z.boolean().default(false),
    value: z.string().optional().superRefine((val, ctx) => {
      const parentEnabled = ctx.path && ctx.path.length > 1 ?
        ((ctx as z.RefinementCtx & { data?: { enabled: boolean } }).data?.enabled || false) : false;

      if (val === undefined && parentEnabled) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "Coordinates are required when enabled",
        });
        return;
      }

      if (val && parentEnabled) {
        const isValidFormat = validateCoordinateFormat(val);
        if (!isValidFormat) {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            message: "Coordinates must be in the format 'latitude,longitude' with valid decimal values",
          });
        }
      }
    }),
  }).default({ enabled: false, value: undefined }),
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
}

export function PhysicalLocationForm({
  buttonText = "Add Location",
  defaultValues = {
    locationName: '',
    locationType: '',
    bgColor: 'red',
    coordinates: {
      enabled: false,
      value: ''
    }
  },
  locationData,
  isEditing = false,
  onDelete,
  onSuccess,
}: PhysicalLocationFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const methods = useForm<FormValues>({
    resolver: zodResolver(PhysicalLocationFormSchema),
    defaultValues: isEditing && locationData
      ? {
          locationName: locationData.name || '',
          locationType: locationData.locationType || '',
          bgColor: locationData.bgColor || 'gray',
          coordinates: {
            enabled: !!locationData.mapCoordinates,
            value: locationData.mapCoordinates || ''
          }
        }
      : defaultValues,
    mode: 'onTouched',
  });

  useEffect(() => {
    if (isEditing && locationData) {
      methods.reset({
        locationName: locationData.name || '',
        locationType: locationData.locationType || '',
        bgColor: locationData.bgColor || 'gray',
        coordinates: {
          enabled: !!locationData.mapCoordinates,
          value: locationData.mapCoordinates || ''
        }
      });
    }
  }, [methods, isEditing, locationData]);

  const onSubmit = useCallback((data: FormValues) => {
    const locationPayload = {
      name: data.locationName,
      locationType: data.locationType,
      bgColor: data.bgColor,
      mapCoordinates: data.coordinates.enabled ? data.coordinates.value : ''
    };

    console.log('Form submitted with payload:', locationPayload);

    if (onSuccess) onSuccess(data);
  }, [onSuccess]);

  const handleDelete = useCallback((id: string) => {
    if (onDelete) onDelete(id);
    setDeleteDialogOpen(false);
  }, [onDelete]);

  return (
    <HookFormProvider {...methods}>
      <form onSubmit={methods.handleSubmit(onSubmit)} className="space-y-8">
        {/* Location Name */}
        <FormField
          control={methods.control}
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
          control={methods.control}
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
          control={methods.control}
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
          control={methods.control}
          name="coordinates"
          render={({ field, fieldState: { error } }) => (
            <FormItem className="space-y-4">
              <div className="flex flex-row items-center justify-between rounded-lg border p-4">
                <div className="space-y-0.5">
                  <FormLabel className="text-base">
                    Coordinates
                  </FormLabel>
                  <FormDescription>
                    Optionally add map coordinates for the location.
                  </FormDescription>
                </div>
                <FormControl>
                  <Switch
                    checked={field.value?.enabled}
                    onCheckedChange={(checked) => {
                      field.onChange({
                        enabled: checked,
                        value: checked ? field.value?.value : undefined
                      });
                    }}
                  />
                </FormControl>
              </div>

              {field.value?.enabled && (
                <FormItem>
                  <FormControl>
                    <Input
                      placeholder="Enter latitude and logitude coordinates or paste a GoogleMaps URL"
                      value={field.value.value ?? ''}
                      onChange={(event) => {
                        const input = event.target.value;
                        field.onChange({
                          enabled: true,
                          value: input || ''
                        });

                        const parsedCoordinates = parseCoordinates(input);
                        if (parsedCoordinates && parsedCoordinates !== input) {
                          setTimeout(() => {
                            field.onChange({
                              enabled: true,
                              value: parsedCoordinates
                            });
                          }, 300);
                        }
                      }}
                      aria-invalid={!!error}
                    />
                  </FormControl>
                  <FormDescription className="text-xs">
                    Format: latitude, longitude (e.g., 40.69007948941017, -74.04439419553563) or paste a Google Maps URL
                  </FormDescription>
                  {error && <FormMessage>{error.message}</FormMessage>}
                </FormItem>
              )}
            </FormItem>
          )}
        />

        {/* Form actions */}
        <div className="flex justify-between w-full mt-6">
          <Button
            type="submit"
            className={isEditing && onDelete ? "flex-1" : "w-full"}
          >
            {isEditing ? "Update Location" : buttonText}
          </Button>

          {isEditing && onDelete && locationData && (
            <Button
              type="button"
              variant="destructive"
              className="ml-2"
              onClick={() => setDeleteDialogOpen(true)}
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
