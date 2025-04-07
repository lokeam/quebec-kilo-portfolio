// MediaPageLocationEditForm.tsx
import { useState, useCallback } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/shared/components/ui/card';
import { Input } from '@/shared/components/ui/input';
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Pencil, Trash2, MapPin } from 'lucide-react';

// Components
import { FormContainer } from '@/features/dashboard/components/templates/FormContainer';
import { getLocationTypeIcon } from '@/features/dashboard/lib/utils/getLocationIcon';

// Hooks
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';
import { useLocationUpdate, useLocationDelete } from '@/core/api/hooks/useLocationActions';

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/components/ui/select";

import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/components/ui/form";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/components/ui/dialog";

import { PhysicalLocationType } from '@/features/dashboard/lib/types/media-storage/constants';
import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';

// Form schema
const LocationEditFormSchema = z.object({
  name: z
    .string()
    .min(3, { message: "Location name must be at least 3 characters" }),
  locationType: z
    .nativeEnum(PhysicalLocationType, {
      required_error: "Please select a location type"
    }),
  mapCoordinates: z
    .string()
    .optional(),
});

type LocationFormValues = z.infer<typeof LocationEditFormSchema>;

interface MediaPageLocationEditFormProps {
  locationData: PhysicalLocation[];
  onSuccess: () => void;
}

export function MediaPageLocationEditForm({ locationData, onSuccess }: MediaPageLocationEditFormProps) {
  // State
  const [selectedLocation, setSelectedLocation] = useState<PhysicalLocation | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
  const [locationToDelete, setLocationToDelete] = useState<string | null>(null);

  // Hooks
  const domainMaps = useDomainMaps();

  // Form setup
  const form = useForm<LocationFormValues>({
    resolver: zodResolver(LocationEditFormSchema),
    defaultValues: {
      name: "",
      locationType: undefined,
      mapCoordinates: "",
    }
  });

  // Custom action hooks with success callbacks
  const { updateLocation, isUpdating } = useLocationUpdate(() => {
    setIsEditing(false);
    setSelectedLocation(null);
    onSuccess();
  });

  const { deleteLocation, isDeleting } = useLocationDelete(() => {
    setDeleteConfirmOpen(false);
    setLocationToDelete(null);
    onSuccess();
  });

  // Event handlers
  const handleEditClick = useCallback((location: PhysicalLocation) => {
    setSelectedLocation(location);
    form.reset({
      name: location.name,
      locationType: location.locationType,
      mapCoordinates: location.mapCoordinates,
    });
    setIsEditing(true);
  }, [form]);

  const handleDeleteClick = useCallback((locationId: string) => {
    setLocationToDelete(locationId);
    setDeleteConfirmOpen(true);
  }, []);

  const handleSubmit = useCallback((data: LocationFormValues) => {
    if (!selectedLocation) return;

    updateLocation({
      id: selectedLocation.id,
      name: data.name,
      locationType: data.locationType,
      mapCoordinates: data.mapCoordinates,
    });
  }, [selectedLocation, updateLocation]);

  const handleConfirmDelete = useCallback(() => {
    if (locationToDelete) {
      deleteLocation(locationToDelete);
    }
  }, [locationToDelete, deleteLocation]);

  const handleCancelEdit = useCallback(() => {
    setIsEditing(false);
  }, []);

  // UI helpers
  const getLocationTypeDisplay = useCallback((type: PhysicalLocationType) => {
    const typeMap = {
      [PhysicalLocationType.HOUSE]: 'House',
      [PhysicalLocationType.APARTMENT]: 'Apartment',
      [PhysicalLocationType.OFFICE]: 'Office',
      [PhysicalLocationType.WAREHOUSE]: 'Warehouse',
    };
    return typeMap[type] || type;
  }, []);

  const renderLocationIcon = useCallback((type: PhysicalLocationType) => {
    return getLocationTypeIcon(type, domainMaps);
  }, [domainMaps]);

  return (
    <div className="px-4">
      {isEditing && selectedLocation ? (
        // Edit form with Zod validation
        <FormContainer form={form} onSubmit={handleSubmit}>
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Location Name</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="locationType"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Location Type</FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select location type" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {Object.values(PhysicalLocationType).map((type) => (
                      <SelectItem key={type} value={type}>
                        <div className="flex items-center">
                          {renderLocationIcon(type as PhysicalLocationType)}
                          {getLocationTypeDisplay(type as PhysicalLocationType)}
                        </div>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="mapCoordinates"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Map Coordinates</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormDescription>
                  Optional location coordinates
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex justify-end space-x-2 pt-4">
            <Button
              variant="outline"
              type="button"
              onClick={handleCancelEdit}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isUpdating}
            >
              {isUpdating ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </FormContainer>
      ) : (
        // List of locations
        <div className="space-y-4 pb-4">
          {locationData.map((location) => (
            <Card key={location.id} className="overflow-hidden">
              <CardHeader className="pb-2">
                <div className="flex justify-between items-start pb-2">
                  <div>
                    <CardTitle>{location.name}</CardTitle>
                    <CardDescription className="flex items-center mt-2">
                      {renderLocationIcon(location.locationType)}
                      {getLocationTypeDisplay(location.locationType)}
                    </CardDescription>
                  </div>
                  <div className="flex space-x-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleEditClick(location)}
                    >
                      <Pencil className="h-4 w-4 mr-1" />
                      Edit
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleDeleteClick(location.id)}
                    >
                      <Trash2 className="h-4 w-4 mr-1" />
                      Delete
                    </Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent className="pb-3">
                {location.mapCoordinates && (
                  <div className="flex items-center text-sm text-muted-foreground">
                    <MapPin className="h-4 w-4 mr-1" />
                    {location.mapCoordinates}
                  </div>
                )}
              </CardContent>
              <CardFooter className="bg-muted/50 py-2 px-6 text-sm">
                <div className="flex justify-between w-full">
                  <span>Created: {location.createdAt.toLocaleDateString()}</span>
                  <span>Updated: {location.updatedAt.toLocaleDateString()}</span>
                </div>
              </CardFooter>
            </Card>
          ))}
        </div>
      )}

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
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
              onClick={() => setDeleteConfirmOpen(false)}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleConfirmDelete}
              disabled={isDeleting}
            >
              {isDeleting ? 'Deleting...' : 'Delete'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}