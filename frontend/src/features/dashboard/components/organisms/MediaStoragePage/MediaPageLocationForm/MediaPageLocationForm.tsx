import { useEffect, useState } from 'react';

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
import { toast } from "sonner"

// Zod
import { z } from "zod"

// Icons
import { House, Building, Building2, Warehouse } from 'lucide-react';
import { IconCar } from '@tabler/icons-react';
import { Dialog, DialogContent, DialogFooter, DialogTitle, DialogHeader, DialogDescription } from '@/shared/components/ui/dialog';

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
        if (val === undefined && (ctx as z.RefinementCtx & { parent: { enabled: boolean } }).parent.enabled) {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            message: "Coordinates are required when enabled",
          });
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

interface MediaPageLocationFormProps {
  onSuccess?: (data: z.infer<typeof PhysicalLocationFormSchema>) => void;
  defaultValues?: z.infer<typeof PhysicalLocationFormSchema>;
  buttonText?: string;
  locationData?: PhysicalLocationData; // NOTE: for editing existing information
  isEditing?: boolean;
  onDelete?: (id: string) => void;
}

export function MediaPageLocationForm({
  buttonText = "Add Location",
  onSuccess,
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
}: MediaPageLocationFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

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

  const handleSubmit = (data: z.infer<typeof PhysicalLocationFormSchema>) => {
    onSuccess?.(data);

    toast(isEditing ? "Location updated successfully" : "New location created successfully", {
      className: 'bg-green-500 text-white',
      duration: 2500,
    });

    // NOTE: for use in transforming data for API
    if (isEditing && locationData) {
      // Example: Call API directly if needed
      // updateLocation({
      //   id: locationData.id,
      //   name: data.locationName,
      //   locationType: data.locationType,
      //   mapCoordinates: data.coordinates.enabled ? data.coordinates.value : undefined,
      // });
    }
  };

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
                    placeholder="Enter coordiantes"
                    value={field.value.value ?? ''}
                    onChange={(event) => {
                      field.onChange({
                        enabled: true,
                        value: event.target.value || '' // Ensure value is never undefined
                      });
                    }}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          </FormItem>
        )}
      />

      {/* Form actions */}
      <div className="flex justify-between w-full mt-6">
        <Button type="submit" className={isEditing && onDelete ? "flex-1" : "w-full"}>
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
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={() => {
                onDelete(locationData.id);
                setDeleteDialogOpen(false);
              }}
            >
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    )}
  </>
  );
}
