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
import { IconEdit, IconTrash } from '@/shared/components/ui/icons';

// Types
import type { LocationsBFFSublocationResponse } from '@/types/domain/physical-location';

// Hooks
import { useDeleteSublocation } from '@/core/api/queries/physicalLocation.queries';

// Utils
import { cn } from '@/shared/components/ui/utils';
import { PhysicalLocationIcon } from '@/features/dashboard/lib/utils/getPhysicalLocationIcon';
import { SublocationIcon } from '@/features/dashboard/lib/utils/getSublocationIcon';

// Components
import { SublocationDeleteWarning } from '@/features/dashboard/components/organisms/PhysicalLocationsPage/DeleteWarningContent/SublocationDeleteWarning';

interface SingleSublocationCardProps {
  location: LocationsBFFSublocationResponse;
  onDelete?: (id: string) => void;
  onEdit?: (location: LocationsBFFSublocationResponse) => void;
  isWatchedByResizeObserver?: boolean;
}

export const SingleSublocationCard = memo(({
  location,
  onDelete,
  onEdit,
  isWatchedByResizeObserver
}: SingleSublocationCardProps) => {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const deleteMutation = useDeleteSublocation();

  const handleEditLocation = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent card onClick from firing
    onEdit?.(location);
  }, [location, onEdit]);

  const handleDeleteClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent card onClick from firing
    setDeleteDialogOpen(true);
    setDeleteError(null);
  }, []);

  const handleConfirmDelete = useCallback(() => {
    if (!location.sublocationId) return;

    deleteMutation.mutate([location.sublocationId], {
      onSuccess: () => {
        setDeleteDialogOpen(false);
        if (onDelete) onDelete(location.sublocationId);
      },
      onError: (error) => {
        setDeleteError(error instanceof Error ? error.message : 'Failed to delete sublocation');
      }
    });
  }, [location.sublocationId, deleteMutation, onDelete]);

  return (
    <>
      <Card
        className={cn(
          "flex flex-col relative cursor-pointer group w-full min-h-[180px] max-h-[180px] p-4",
          "transition-all duration-200",
          "hover:ring-1 hover:ring-ring hover:ring-inset",
          "hover:shadow-lg",
          isWatchedByResizeObserver && 'data-card-sentinel'
        )}
      >
        <div className="flex flex-col w-full">
          {/* First Row - Icons and Edit Section */}
          <div className="flex items-center justify-between w-full mb-8">
            {/* Left side - Sublocation Icon */}
            <div className="flex items-center">
              <SublocationIcon
                type={location.sublocationType}
                bgColor={location.parentLocationBgColor || 'gray'}
              />
            </div>

            {/* Right side - Physical Location Icon or Edit Section */}
            <div className="relative w-32">
              {/* Physical Location Icon (shown by default) */}
              <div className="flex items-center justify-end transition-opacity duration-200 group-hover:opacity-0">
                <PhysicalLocationIcon
                  type={location.parentLocationType}
                  bgColor={location.parentLocationBgColor || 'gray'}
                />
              </div>

              {/* Edit Section (shown on hover) */}
              <div className="absolute top-0 right-0 flex items-center gap-2 opacity-0 invisible transition-opacity duration-200 group-hover:opacity-100 group-hover:visible">
                <Button
                  variant="outline"
                  size="sm"
                  className="h-10 w-10 p-0"
                  onClick={handleEditLocation}
                >
                  <IconEdit className="h-5 w-5" />
                  <span className="sr-only">Edit {location.sublocationName}</span>
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  className="h-10 w-10 p-0"
                  onClick={handleDeleteClick}
                >
                  <IconTrash className="h-5 w-5" />
                  <span className="sr-only">Delete {location.sublocationName}</span>
                </Button>
              </div>
            </div>
          </div>

          {/* Second Row - Text Content */}
          <div className="flex flex-col">
            {/* Sublocation Name */}
            <div className="flex mb-2">
              <span
                className="font-bold text-lg text-foreground truncate overflow-hidden"
                style={{
                  maxWidth: 'var(--label-max-width)',
                  display: 'block',
                }}
              >
                {location.sublocationName}
              </span>
            </div>

            {/* Parent Location Name */}
            <div className="flex">
              <span
                className="font-medium text-sm text-muted-foreground truncate overflow-hidden"
                style={{
                  maxWidth: 'var(--label-max-width)',
                  display: 'block',
                }}
              >
                {`Parent Location: ${location.parentLocationName}`}
              </span>
            </div>
          </div>
        </div>
      </Card>

      {/* Delete confirmation dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete {location.sublocationName}</DialogTitle>
            <DialogDescription>
                              {deleteError ? (
                  <div className="text-destructive">
                    {deleteError}
                  </div>
              ) : (
                <SublocationDeleteWarning location={location} />
              )}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
              disabled={deleteMutation.isPending}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleConfirmDelete}
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? (
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
    </>
  );
});

SingleSublocationCard.displayName = 'SingleSublocationCard';
