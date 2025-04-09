import { useState, useCallback } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/shared/components/ui/card';
import { Pencil, Trash2, MapPin } from 'lucide-react';

// Components
import { PhysicalLocationFormSingle } from '../PhysicalLocationFormSingle/PhysicalLocationFormSingle';
import { getLocationTypeIcon } from '@/features/dashboard/lib/utils/getLocationIcon';

// Hooks
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';
import { useLocationDelete } from '@/core/api/hooks/useLocationManager';

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

interface PhysicalLocationDrawerListProps {
  locationData: PhysicalLocation[];
  onSuccess: () => void;
}

export function PhysicalLocationDrawerList({ locationData, onSuccess }: PhysicalLocationDrawerListProps) {
  // State
  const [selectedLocation, setSelectedLocation] = useState<PhysicalLocation | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
  const [locationToDelete, setLocationToDelete] = useState<string | null>(null);

  // Hooks
  const domainMaps = useDomainMaps();

  // Custom action hooks with success callbacks
  const { deleteLocation, isDeleting } = useLocationDelete(() => {
    setDeleteConfirmOpen(false);
    setLocationToDelete(null);
    onSuccess();
  });

  // Event handlers
  const handleEditClick = useCallback((location: PhysicalLocation) => {
    setSelectedLocation(location);
    setIsEditing(true);
  }, []);

  const handleDeleteClick = useCallback((locationId: string) => {
    setLocationToDelete(locationId);
    setDeleteConfirmOpen(true);
  }, []);

  const handleEditSuccess = useCallback(() => {
    setIsEditing(false);
    setSelectedLocation(null);
    onSuccess();
  }, [onSuccess]);

  const handleConfirmDelete = useCallback(() => {
    if (locationToDelete) {
      deleteLocation(locationToDelete);
    }
  }, [locationToDelete, deleteLocation]);

  const handleCancelEdit = useCallback(() => {
    setIsEditing(false);
    setSelectedLocation(null);
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
    <div className="">
      {isEditing && selectedLocation ? (
        // Use PhysicalLocationFormSingle for editing
        <div>
          <Button
            variant="ghost"
            onClick={handleCancelEdit}
            className="mb-4"
          >
            ← Back to locations
          </Button>

          <PhysicalLocationFormSingle
            locationData={selectedLocation}
            isEditing={true}
            onSuccess={handleEditSuccess}
            onDelete={(id) => deleteLocation(id)}
          />
        </div>
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
