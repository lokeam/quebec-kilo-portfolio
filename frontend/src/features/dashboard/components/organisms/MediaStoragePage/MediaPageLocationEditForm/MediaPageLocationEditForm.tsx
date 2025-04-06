// MediaPageLocationEditForm.tsx
import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/shared/components/ui/card';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { Pencil, Trash2, MapPin } from 'lucide-react';

// Hooks
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';

// Components
import { getLocationTypeIcon } from '@/features/dashboard/lib/utils/getLocationIcon';

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/components/ui/select";

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

interface MediaPageLocationEditFormProps {
  locationData: PhysicalLocation[];
  onSuccess: () => void;
}

export function MediaPageLocationEditForm({ locationData, onSuccess }: MediaPageLocationEditFormProps) {
  const [selectedLocation, setSelectedLocation] = useState<PhysicalLocation | null>(null);
  const [editData, setEditData] = useState<Partial<PhysicalLocation>>({});
  const [isEditing, setIsEditing] = useState(false);
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
  const [locationToDelete, setLocationToDelete] = useState<string | null>(null);
  const domainMaps = useDomainMaps();

  const queryClient = useQueryClient();

  // Mutation for updating a location
  const updateLocationMutation = useMutation({
    mutationFn: async (data: Partial<PhysicalLocation>) => {
      // Simulate API call
      console.log('Updating location:', data);
      // In a real app, you would call your API here
      return new Promise(resolve => setTimeout(() => resolve(data), 1000));
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['mediaStorage'] });
      toast.success('Location updated successfully');
      setIsEditing(false);
      setSelectedLocation(null);
      onSuccess();
    },
    onError: (error) => {
      toast.error('Failed to update location: ' + error.message);
    },
  });

  // Mutation for deleting a location
  const deleteLocationMutation = useMutation({
    mutationFn: async (id: string) => {
      // Simulate API call
      console.log('Deleting location:', id);
      // In a real app, you would call your API here
      return new Promise(resolve => setTimeout(() => resolve({ id }), 1000));
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['mediaStorage'] });
      toast.success('Location deleted successfully');
      setDeleteConfirmOpen(false);
      setLocationToDelete(null);
      onSuccess();
    },
    onError: (error) => {
      toast.error('Failed to delete location: ' + error.message);
    },
  });

  // Handle edit button click
  const handleEditClick = (location: PhysicalLocation) => {
    setSelectedLocation(location);
    setEditData({
      name: location.name,
      locationType: location.locationType,
      mapCoordinates: location.mapCoordinates,
    });
    setIsEditing(true);
  };

  // Handle delete button click
  const handleDeleteClick = (locationId: string) => {
    setLocationToDelete(locationId);
    setDeleteConfirmOpen(true);
  };

  // Handle form submission
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedLocation) return;

    updateLocationMutation.mutate({
      id: selectedLocation.id,
      ...editData,
    });
  };

  // Get location type display name
  const getLocationTypeDisplay = (type: PhysicalLocationType) => {
    const typeMap = {
      [PhysicalLocationType.HOUSE]: 'House',
      [PhysicalLocationType.APARTMENT]: 'Apartment',
      [PhysicalLocationType.OFFICE]: 'Office',
      [PhysicalLocationType.WAREHOUSE]: 'Warehouse',
    };
    return typeMap[type] || type;
  };

  // Get location type icon
  const renderLocationIcon = (type: PhysicalLocationType) => {
    return getLocationTypeIcon(type, domainMaps);
  };

  return (
    <div className="px-4">

          {isEditing && selectedLocation ? (
            // Edit form
            <form onSubmit={handleSubmit} className="space-y-8 p-4">
              <div className="space-y-2">
                <Label htmlFor="name">Location Name</Label>
                <Input
                  id="name"
                  value={editData.name || ''}
                  onChange={(e) => setEditData({...editData, name: e.target.value})}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="locationType">Location Type</Label>
                <Select
                  value={editData.locationType}
                  onValueChange={(value) => setEditData({...editData, locationType: value as PhysicalLocationType})}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select location type" />
                  </SelectTrigger>
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
              </div>

              <div className="space-y-2">
                <Label htmlFor="mapCoordinates">Map Coordinates</Label>
                <Input
                  id="mapCoordinates"
                  value={editData.mapCoordinates || ''}
                  onChange={(e) => setEditData({...editData, mapCoordinates: e.target.value})}
                />
              </div>

              <div className="flex justify-end space-x-2 pt-4">
                <Button
                  variant="outline"
                  type="button"
                  onClick={() => setIsEditing(false)}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={updateLocationMutation.isPending}
                >
                  {updateLocationMutation.isPending ? 'Saving...' : 'Save Changes'}
                </Button>
              </div>
            </form>
          ) : (
            // List of locations
            <div className="space-y-4 pb-4">
              {locationData.map((location) => (
                <Card key={location.id} className="overflow-hidden">
                  <CardHeader className="pb-2">
                    <div className="flex justify-between items-start">
                      <div>
                        <CardTitle>{location.name}</CardTitle>
                        <CardDescription className="flex items-center mt-1">
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
              onClick={() => locationToDelete && deleteLocationMutation.mutate(locationToDelete)}
              disabled={deleteLocationMutation.isPending}
            >
              {deleteLocationMutation.isPending ? 'Deleting...' : 'Delete'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}