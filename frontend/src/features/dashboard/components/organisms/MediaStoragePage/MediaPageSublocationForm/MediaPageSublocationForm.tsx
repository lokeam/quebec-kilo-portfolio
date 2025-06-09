import { useEffect, useState } from 'react';
import { AxiosError } from 'axios';

// Components
import { FormContainer } from '@/features/dashboard/components/templates/FormContainer';

// Shadcn UI Components
import { Input } from '@/shared/components/ui/input';
import { Button } from '@/shared/components/ui/button';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/components/ui/form';

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/components/ui/dialog";

// Icons
import { BookshelfIcon } from '@/shared/components/ui/CustomIcons/BookshelfIcon';
import { MediaConsoleIcon } from '@/shared/components/ui/CustomIcons/MediaConsoleIcon';
import { DrawerIcon } from '@/shared/components/ui/CustomIcons/DrawerIcon';
import { CabinetIcon } from '@/shared/components/ui/CustomIcons/CabinetIcon';
import { ClosetIcon } from '@/shared/components/ui/CustomIcons/ClosetIcon';
import { Package } from 'lucide-react';

// Hooks
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { useLocationManager } from '@/core/api/hooks/useLocationManager';
import { toast } from 'sonner';

// Zod
import { z } from 'zod';

import { SublocationType } from '@/features/dashboard/lib/types/media-storage/constants';
import type { LocationPayload } from "@/core/api/queries/useLocationMutations";
import type { Sublocation } from "@/core/api/types/location";

export const SublocationFormSchema = z.object({
  locationName: z
    .string({
      required_error: "Please enter a location name",
    })
    .min(3, {
      message: "Location name must be at least 3 characters long",
    }),
  locationType: z
    .enum([SublocationType.shelf, SublocationType.console, SublocationType.cabinet, SublocationType.closet, SublocationType.drawer, SublocationType.box], {
      required_error: "Please select a location type",
    }).default({ enabled: false, value: undefined }),
});

interface LocationData {
  id: string;
  name: string;
  locationType: SublocationType;
  createdAt?: Date;
  updatedAt?: Date;
}

interface MediaPageSublocationFormProps {
  onSuccess?: (data: z.infer<typeof SublocationFormSchema>) => void;
  defaultValues?: z.infer<typeof SublocationFormSchema>;
  buttonText?: string;
  sublocationData?: LocationData; // For editing existing sublocation
  isEditing?: boolean;   // To indicate edit mode
  onDelete?: (id: string) => void;
  parentLocationId: string; // Required for creating/updating sublocations
}

export function MediaPageSublocationForm({
  onSuccess,
  defaultValues = {
    locationName: '',
    locationType: SublocationType.shelf,
  },
  buttonText = "Add Sublocation",
  sublocationData,
  isEditing = false,
  onDelete,
  parentLocationId,
}: MediaPageSublocationFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  // Setup location manager hook for sublocations
  const locationManager = useLocationManager({
    type: 'sublocation',
    onSuccess: () => {
      onSuccess?.(form.getValues());
    }
  });

  /* Specific form components creates their own useForm hook instances */
  const form = useForm<z.infer<typeof SublocationFormSchema>>({
    resolver: zodResolver(SublocationFormSchema),
    defaultValues: isEditing && sublocationData
      ? {
          locationName: sublocationData.name || "",
          locationType: sublocationData.locationType || SublocationType.shelf,
        }
      : defaultValues,
  });

  // If sublocationData changes and we're in edit mode, update form values
  useEffect(() => {
    if (isEditing && sublocationData) {
      form.reset({
        locationName: sublocationData.name || '',
        locationType: sublocationData.locationType || SublocationType.shelf,
      });
    }
  }, [form, isEditing, sublocationData]);

  const handleSubmit = async (values: z.infer<typeof SublocationFormSchema>) => {
    try {
      const payload: LocationPayload = {
        name: values.locationName,
        locationType: values.locationType,
        mapCoordinates: values.coordinates.enabled ? values.coordinates.value : undefined,
        physicalLocationId: parentLocationId,
      };

      if (isEditing && sublocationData) {
        await locationManager.update({
          ...payload,
          id: sublocationData.id,
        });
        toast.success('Sublocation updated successfully');
      } else {
        const sublocationPayload: Sublocation = {
          ...payload,
          physicalLocationId: parentLocationId,
        };
        await locationManager.create(sublocationPayload);
        toast.success('Sublocation created successfully');
      }
    } catch (error) {
      const axiosError = error as AxiosError<{ error: string }>;
      if (axiosError.response?.data?.error?.includes('already exists')) {
        toast.error('A sublocation with this name already exists in this physical location');
      } else {
        toast.error('Failed to save sublocation. Please try again.');
      }
    }
  };

  const handleDelete = (id: string) => {
    locationManager.delete(id);
    setDeleteDialogOpen(false);
    toast.success('Sublocation deleted successfully');
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
              <FormLabel>Sublocation Name</FormLabel>
              <FormControl>
                <Input placeholder="Example: Study bookcase" {...field} />
              </FormControl>
              <FormDescription>
                What shall we call this area?
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
              <FormLabel>Storage Unit Type</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Where are you keeping your media?" />
                  </SelectTrigger>
                </FormControl>

                <SelectContent>
                  <SelectItem value="shelf">
                    <div className="flex items-center gap-2">
                      <BookshelfIcon size={20} color='#fff' className='mr-2'/>
                      <span>Shelf / Shelving unit</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="console">
                    <div className="flex items-center gap-2">
                      <MediaConsoleIcon size={20} color='#fff' className='mr-2'/>
                      <span>Media console</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="cabinet">
                    <div className="flex items-center gap-2">
                      <CabinetIcon size={20} color='#fff' className='mr-2'/>
                      <span>Cabinet</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="closet">
                    <div className="flex items-center gap-2">
                      <ClosetIcon size={20} color='#fff' className='mr-2'/>
                      <span>Closet</span>
                    </div>
                    </SelectItem>
                  <SelectItem value="drawer">
                    <div className="flex items-center gap-2">
                      <DrawerIcon size={20} color='#fff' className='mr-2'/>
                      <span>Drawer</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="box">
                    <div className="flex items-center gap-2">
                      <Package size={20} color='#fff' className='mr-2'/>
                      <span>Storage container</span>
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
              <FormDescription>
                What kind of furniture or storage unit is this?
              </FormDescription>

              <FormMessage />
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
            isEditing ? "Update Sublocation" : buttonText
          )}
        </Button>

        {isEditing && onDelete && sublocationData && (
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

      {/* Delete Confirmation Dialog */}
      {isEditing && onDelete && sublocationData && (
        <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Confirm Deletion</DialogTitle>
              <DialogDescription>
                Are you sure you want to delete this sublocation? This action cannot be undone.
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
                onClick={() => handleDelete(sublocationData.id)}
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
