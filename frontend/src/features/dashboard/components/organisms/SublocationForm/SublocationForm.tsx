import { useCallback, useState } from 'react';
import { useForm, FormProvider as HookFormProvider } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

// Shadcn UI Components
import { Button } from "@/shared/components/ui/button"
import { Input } from "@/shared/components/ui/input"
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
import { BookshelfIcon } from '@/shared/components/ui/CustomIcons/BookshelfIcon';
import { MediaConsoleIcon } from '@/shared/components/ui/CustomIcons/MediaConsoleIcon';
import { DrawerIcon } from '@/shared/components/ui/CustomIcons/DrawerIcon';
import { CabinetIcon } from '@/shared/components/ui/CustomIcons/CabinetIcon';
import { ClosetIcon } from '@/shared/components/ui/CustomIcons/ClosetIcon';
import { Package } from 'lucide-react';

// Utils
import { PhysicalLocationIcon } from '@/features/dashboard/lib/utils/getPhysicalLocationIcon';
import type { LocationIconBgColor } from '@/types/domain/location-types';

// Types
import type { LocationsBFFPhysicalLocationResponse } from '@/types/domain/physical-location';

// Constants
const SublocationType = {
  shelf: 'shelf',
  console: 'console',
  cabinet: 'cabinet',
  closet: 'closet',
  drawer: 'drawer',
  box: 'box',
} as const;

export const SublocationFormSchema = z.object({
  locationName: z.string().min(3, 'Location name must be at least 3 characters'),
  locationType: z.enum([SublocationType.shelf, SublocationType.console, SublocationType.cabinet, SublocationType.closet, SublocationType.drawer, SublocationType.box] as const),
});

type FormValues = z.infer<typeof SublocationFormSchema>;

interface SublocationData {
  id: string;
  name: string;
  locationType: string;
  createdAt?: Date;
  updatedAt?: Date;
}

interface SublocationFormProps {
  onSuccess?: (data: FormValues) => void;
  defaultValues?: FormValues;
  buttonText?: string;
  sublocationData?: SublocationData;
  isEditing?: boolean;
  onDelete?: (id: string) => void;
  onClose?: () => void;
  parentLocation: LocationsBFFPhysicalLocationResponse;
}

export function SublocationForm({
  buttonText = "Add Sublocation",
  defaultValues = {
    locationName: '',
    locationType: SublocationType.shelf,
  },
  sublocationData,
  isEditing = false,
  onDelete,
  onClose,
  onSuccess,
  parentLocation,
}: SublocationFormProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  // Create form default values based on whether we're editing or creating
  const formDefaultValues = isEditing && sublocationData ? {
    locationName: sublocationData.name,
    locationType: sublocationData.locationType as typeof SublocationType[keyof typeof SublocationType],
  } : defaultValues;

  const methods = useForm<FormValues>({
    resolver: zodResolver(SublocationFormSchema),
    defaultValues: formDefaultValues,
    mode: 'onTouched'
  });

  const { handleSubmit, control } = methods;

  const onSubmit = useCallback((data: FormValues) => {
    // TODO: Implement mutation logic
    console.log('Form submitted:', {
      ...data,
      parentLocationId: parentLocation.physicalLocationId,
      parentLocationName: parentLocation.name,
      parentLocationType: parentLocation.physicalLocationType,
      parentLocationBgColor: parentLocation.bgColor,
    });

    if (onClose) onClose();
    if (onSuccess) onSuccess(data);
  }, [parentLocation, onClose, onSuccess]);

  const handleDelete = useCallback((id: string) => {
    if (onDelete) onDelete(id);
    setDeleteDialogOpen(false);
  }, [onDelete]);

  return (
    <HookFormProvider {...methods}>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
        {/* Parent Location Info */}
        <div className="flex items-center gap-4 p-4 bg-slate-100 rounded-lg">
          <PhysicalLocationIcon
            type={parentLocation.physicalLocationType}
            bgColor={(parentLocation.bgColor as LocationIconBgColor) || 'red'}
          />
          <div>
            <p className="text-sm text-slate-500">Parent Location</p>
            <p className="font-medium">{parentLocation.name}</p>
          </div>
        </div>

        {/* Location Name */}
        <FormField
          control={control}
          name="locationName"
          render={({ field, fieldState: { error } }) => (
            <FormItem>
              <FormLabel>Sublocation Name</FormLabel>
              <FormControl>
                <Input
                  placeholder="Example: Study bookcase"
                  {...field}
                  aria-invalid={!!error}
                />
              </FormControl>
              <FormDescription>
                What shall we call this area?
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
              <FormLabel>Storage Unit Type</FormLabel>
              <Select
                onValueChange={field.onChange}
                defaultValue={field.value}
              >
                <FormControl>
                  <SelectTrigger aria-invalid={!!error}>
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
              {error && <FormMessage>{error.message}</FormMessage>}
            </FormItem>
          )}
        />

        {/* Form actions */}
        <div className="flex justify-between w-full mt-6">
          <Button
            type="submit"
            className={isEditing && onDelete ? "flex-1" : "w-full"}
          >
            {isEditing ? "Update Sublocation" : buttonText}
          </Button>

          {isEditing && onDelete && sublocationData && (
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
              >
                Cancel
              </Button>
              <Button
                variant="destructive"
                onClick={() => handleDelete(sublocationData.id)}
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