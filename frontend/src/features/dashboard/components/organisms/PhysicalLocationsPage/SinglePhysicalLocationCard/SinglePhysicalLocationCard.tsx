import { memo, useCallback, useState } from 'react';

// Shadcn Components
import { Card } from "@/shared/components/ui/card"
import { Button } from '@/shared/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogTitle,
  DialogHeader,
  DialogDescription
} from '@/shared/components/ui/dialog';

// Icons
import { IconEdit, IconTrash } from '@tabler/icons-react';
import { PackagePlus } from 'lucide-react';

// Types
import type { LocationsBFFPhysicalLocationResponse, LocationsBFFSublocationResponse } from '@/types/domain/physical-location';
import type { LocationIconBgColor } from '@/types/domain/location-types';

// Utils
import { cn } from '@/shared/components/ui/utils';
import { PhysicalLocationIcon } from '@/features/dashboard/lib/utils/getPhysicalLocationIcon';

// Components
import { PhysicalLocationDeleteWarning } from '../DeleteWarningContent/PhysicalLocationDeleteWarning';

interface SinglePhysicalLocationCardProps {
  location: LocationsBFFPhysicalLocationResponse;
  sublocations: LocationsBFFSublocationResponse[];
  onDelete?: (id: string) => void;
  onEdit?: (location: LocationsBFFPhysicalLocationResponse) => void;
  onAddSublocation?: (location: LocationsBFFPhysicalLocationResponse) => void;
  isWatchedByResizeObserver?: boolean;
  isSelectionMode?: boolean;
  onSelect?: (location: LocationsBFFPhysicalLocationResponse) => void;
  isSelected?: boolean;
}

export const SinglePhysicalLocationCard = memo(({
  location,
  sublocations,
  onDelete,
  onEdit,
  onAddSublocation,
  isWatchedByResizeObserver,
  isSelectionMode = false,
  onSelect,
  isSelected = false
}: SinglePhysicalLocationCardProps) => {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const handleEditLocation = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onEdit?.(location);
  }, [location, onEdit]);

  const handleAddSublocation = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onAddSublocation?.(location);
  }, [location, onAddSublocation]);

  const handleDeleteClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    setDeleteDialogOpen(true);
    setDeleteError(null);
  }, []);

  const handleConfirmDelete = useCallback(() => {
    if (!location.physicalLocationId || !onDelete) return;

    setDeleteError(null);

    try {
      onDelete(location.physicalLocationId);
    } catch (err) {
      setDeleteError("Something went wrong. We can't complete this operation now, please try again later.");
      console.error("Error deleting location:", err);
    }
  }, [location.physicalLocationId, onDelete]);

  const handleCardClick = useCallback((e: React.MouseEvent) => {
    if (isSelectionMode && onSelect) {
      e.stopPropagation();
      onSelect(location);
    }
  }, [isSelectionMode, onSelect, location]);

  return (
    <>
      <Card
        className={cn(
          "flex flex-col relative cursor-pointer group w-full min-h-[100px] max-h-[100px] p-4",
          "transition-all duration-200",
          isSelectionMode
            ? "hover:ring-2 hover:ring-blue-500"
            : "hover:ring-1 hover:ring-ring hover:ring-inset",
          isSelected && "ring-2 ring-blue-500",
          isWatchedByResizeObserver && 'data-card-sentinel'
        )}
        onClick={handleCardClick}
        role={isSelectionMode ? "button" : undefined}
        aria-selected={isSelectionMode ? isSelected : undefined}
        aria-label={isSelectionMode ? `Select ${location.name} as parent location` : undefined}
      >
        <div className="flex flex-col w-full">
          <div className="flex items-center justify-between w-full mb-8">
            <div className="flex items-center">
              <span
                className="font-bold text-lg text-foreground truncate overflow-hidden"
                style={{
                  maxWidth: 'var(--label-max-width)',
                  display: 'block',
                }}
              >
                {location.name}
              </span>
            </div>

            <div className="relative w-32">
              <div className="flex items-center justify-end transition-opacity duration-200 group-hover:opacity-0">
                <PhysicalLocationIcon
                  type={location.physicalLocationType}
                  bgColor={(location.bgColor as LocationIconBgColor) || 'red'}
                />
              </div>

              {!isSelectionMode && (
                <div className="absolute top-0 right-0 flex items-center gap-2 opacity-0 invisible transition-opacity duration-200 group-hover:opacity-100 group-hover:visible">
                  <Button
                    variant="outline"
                    size="sm"
                    className="h-10 w-10 p-0"
                    onClick={handleAddSublocation}
                  >
                    <PackagePlus className="h-5 w-5" />
                    <span className="sr-only">Add sublocation to {location.name}</span>
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    className="h-10 w-10 p-0"
                    onClick={handleEditLocation}
                  >
                    <IconEdit className="h-5 w-5" />
                    <span className="sr-only">Edit {location.name}</span>
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    className="h-10 w-10 p-0"
                    onClick={handleDeleteClick}
                  >
                    <IconTrash className="h-5 w-5" />
                    <span className="sr-only">Delete {location.name}</span>
                  </Button>
                </div>
              )}
            </div>
          </div>
        </div>
      </Card>

      {!isSelectionMode && (
        <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Delete {location.name}</DialogTitle>
              <DialogDescription>
                {deleteError ? (
                  <div className="text-red-500">
                    {deleteError}
                  </div>
                ) : (
                  <PhysicalLocationDeleteWarning
                    location={location}
                    associatedItems={sublocations}
                  />
                )}
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setDeleteDialogOpen(false)}
                disabled={deleteError !== null}
              >
                Cancel
              </Button>
              <Button
                variant="destructive"
                onClick={handleConfirmDelete}
                disabled={deleteError !== null}
              >
                Delete
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </>
  );
});

SinglePhysicalLocationCard.displayName = 'SinglePhysicalLocationCard';
