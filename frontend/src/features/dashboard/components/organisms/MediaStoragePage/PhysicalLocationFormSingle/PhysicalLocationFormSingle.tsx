import { useCallback, useEffect, useState } from 'react';

// Components
import { FormContainer } from '@/features/dashboard/components/templates/FormContainer';

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
import { useForm } from "react-hook-form"

// Zod
import { z } from "zod"

// Icons
import { House, Building, Building2, Warehouse } from 'lucide-react';
import { IconCar } from '@tabler/icons-react';
import { Dialog, DialogContent, DialogFooter, DialogTitle, DialogHeader, DialogDescription } from '@/shared/components/ui/dialog';
import { useLocationManager } from '@/core/api/hooks/useLocationManager';

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
    }),
    coordinates: z.object({
      enabled: z.boolean().default(false),
      value: z.string().optional().superRefine((val, ctx) => {
        // Check if ctx.path exists and has parent context
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
          // Validate coordinate format
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

interface PhysicalLocationData {
  id: string;
  name: string;
  locationType: string;
  bgColor?: string;
  mapCoordinates?: string;
  createdAt?: Date;
  updatedAt?: Date;
}

interface PhysicalLocationFormSingleProps {
  onSuccess?: (data: z.infer<typeof PhysicalLocationFormSchema>) => void;
  defaultValues?: z.infer<typeof PhysicalLocationFormSchema>;
  buttonText?: string;
  locationData?: PhysicalLocationData; // NOTE: for editing existing information
  isEditing?: boolean;
  onDelete?: (id: string) => void;
}

export function PhysicalLocationFormSingle({
  buttonText = "Add Location",
defaultValues = {
    locationName: '',
    locationType: '',
    coordinates: {
      enabled: false,
      value: '' // Ensure value is never undefined
    }
  },
  locationData,
  isEditing = false,
  onDelete,
  onSuccess,
}: PhysicalLocationFormSingleProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  // Setup location manager hook with proper type
  const locationManager = useLocationManager({
    type: 'physical'
  });

  /* Specific form components creates their own useForm hook instances */
  const form = useForm<z.infer<typeof PhysicalLocationFormSchema>>({
    resolver: zodResolver(PhysicalLocationFormSchema),
    defaultValues: isEditing && locationData
      ? {
        locationName: locationData.name || '',
        locationType: locationData.locationType || '',
        coordinates: {
          enabled: !!locationData.mapCoordinates,
          value: locationData.mapCoordinates || ''
        }
      }
      : defaultValues
  });

  // If location data changes AND we are editing, update form values
  useEffect(() => {
    if (isEditing && locationData) {
      form.reset({
        locationName: locationData.name || '',
        locationType: locationData.locationType || '',
        coordinates: {
          enabled: !!locationData.mapCoordinates,
          value: locationData.mapCoordinates || ''
        }
      });
    }
  }, [form, isEditing, locationData]);

  const handleSubmit = useCallback((data: z.infer<typeof PhysicalLocationFormSchema>) => {
    // Transform form data to match API expectations
    const locationPayload = {
      id: isEditing && locationData ? locationData.id : undefined,
      name: data.locationName,
      locationType: data.locationType,
      mapCoordinates: data.coordinates.enabled ? data.coordinates.value : undefined,
    };

    if (isEditing && locationData) {
      locationManager.update(locationPayload);
    } else {
      locationManager.create(locationPayload);
    }

    if (onSuccess) onSuccess(data);
  }, [isEditing, locationData, locationManager, onSuccess]);

  const handleDelete = useCallback((id: string) => {
    locationManager.delete(id);
    setDeleteDialogOpen(false);
  }, [locationManager]);

  return (
  <>
    <FormContainer form={form} onSubmit={handleSubmit}>
      {/* Location Name */}
      <FormField
        control={form.control}
        name="locationName"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Location Name</FormLabel>
            <FormControl>
              <Input placeholder="Enter a location name" {...field} />
            </FormControl>
            <FormDescription>
              This is the name of the location where the media is stored.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      {/* Location Type */}
      <FormField
        control={form.control}
        name="locationType"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Location Type</FormLabel>
            <Select onValueChange={field.onChange} defaultValue={field.value}>
              <FormControl>
                <SelectTrigger>
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
                <SelectItem value="commercialStorage">
                  <div className="flex items-center gap-2">
                    <Warehouse size={20} color='#fff' className='mr-2'/>
                    <span>Commercial Storage</span>
                  </div>
                </SelectItem>
                <SelectItem value="vehicle">
                  <div className="flex items-center gap-2">
                    <IconCar size={25} color='#fff' className='mr-2'/>
                    <span>Vehicle</span>
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
            <FormDescription>
              Think of this as the venue where the media is stored.
            </FormDescription>

            <FormMessage />
          </FormItem>
        )}
      />

      {/* Coordinates */}
      <FormField
        control={form.control}
        name="coordinates"
        render={({ field }) => (
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

                      // Update with user input for immediate feedback
                      field.onChange({
                        enabled: true,
                        value: event.target.value || '' // Ensure value is never undefined
                      });

                      // Attempt to parse coordinates or Google Maps URL
                      const parsedCoordinates = parseCoordinates(input);

                      // If parsing successful, update with formatted coordinates
                      if (parsedCoordinates && parsedCoordinates !== input) {
                        // Use timeout to allow original input to register and give ui feedback that something happened
                        setTimeout(() => {
                          field.onChange({
                            enabled: true,
                            value: parsedCoordinates
                          });
                        }, 300);
                      }
                    }}
                  />
                </FormControl>
                <FormDescription className="text-xs">
                  Format: latitude, longitude (e.g., 40.69007948941017, -74.04439419553563) or paste a Google Maps URL
                </FormDescription>
                <FormMessage />
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
          disabled={locationManager.isCreating || locationManager.isUpdating}
        >
          {(locationManager.isCreating || locationManager.isUpdating) ? (
            <>
              <span className="animate-spin mr-2">⊚</span>
              {isEditing ? "Updating..." : "Creating..."}
            </>
          ) : (
            isEditing ? "Update Location" : buttonText
          )}
        </Button>

        {isEditing && onDelete && locationData && (
          <Button
            type="button"
            variant="destructive"
            className="ml-2"
            onClick={() => setDeleteDialogOpen(true)}
            disabled={locationManager.isDeleting}
          >
            {locationManager.isDeleting ? "Deleting..." : "Delete"}
          </Button>
        )}
      </div>
    </FormContainer>

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
              disabled={locationManager.isDeleting}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={() => handleDelete(locationData.id)}
              disabled={locationManager.isDeleting}
            >
              {locationManager.isDeleting ? (
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
    )}
  </>
  );
}
