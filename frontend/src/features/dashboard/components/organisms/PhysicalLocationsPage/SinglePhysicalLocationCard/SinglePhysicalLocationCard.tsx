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

// Types
import type { LocationsBFFPhysicalLocationResponse, LocationsBFFSublocationResponse } from '@/types/domain/physical-location';

// Utils
import { cn } from '@/shared/components/ui/utils';
import { PhysicalLocationIcon } from '@/features/dashboard/lib/utils/getPhysicalLocationIcon';

interface SinglePhysicalLocationCardProps {
  location: LocationsBFFPhysicalLocationResponse;
  sublocations: LocationsBFFSublocationResponse[];
  onDelete?: (id: string) => void;
  onEdit?: (location: LocationsBFFPhysicalLocationResponse) => void;
  isWatchedByResizeObserver?: boolean;
}

export const SinglePhysicalLocationCard = memo(({
  location,
  sublocations,
  onDelete,
  onEdit,
  isWatchedByResizeObserver
}: SinglePhysicalLocationCardProps) => {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const associatedSublocations = sublocations.filter(
    sublocation => sublocation.parentLocationId === location.physicalLocationID
  );

  const handleEditLocation = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onEdit?.(location);
  }, [location, onEdit]);

  const handleDeleteClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    setDeleteDialogOpen(true);
    setDeleteError(null);
  }, []);

  const handleConfirmDelete = useCallback(() => {
    if (!location.physicalLocationID || !onDelete) return;

    setIsDeleting(true);
    setDeleteError(null);

    try {
      onDelete(location.physicalLocationID);
    } catch (err) {
      setIsDeleting(false);
      setDeleteError("Something went wrong. We can't complete this operation now, please try again later.");
      console.error("Error deleting location:", err);
    }
  }, [location.physicalLocationID, onDelete]);

  return (
    <>
      <Card
        className={cn(
          "flex flex-col relative cursor-pointer group w-full min-h-[100px] max-h-[100px] p-4 bg-gradient-to-b from-slate-900 to-slate-950 border-slate-800",
          "transition-all duration-200",
          "hover:ring-1 hover:ring-white/20 hover:ring-inset",
          "hover:shadow-[0_0_4px_0_rgba(95,99,104,0.6),0_0_6px_2px_rgba(95,99,104,0.6)]",
          isWatchedByResizeObserver && 'data-card-sentinel'
        )}
      >
        <div className="flex flex-col w-full">
          <div className="flex items-center justify-between w-full mb-8">
            <div className="flex items-center">
              <span
                className="font-bold text-lg text-white truncate overflow-hidden"
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
                  bgColor={location.bgColor}
                />
              </div>

              <div className="absolute top-0 right-0 flex items-center gap-2 opacity-0 invisible transition-opacity duration-200 group-hover:opacity-100 group-hover:visible">
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
                  className="h-10 w-10 p-0 text-red-500 hover:text-red-600 hover:bg-red-100"
                  onClick={handleDeleteClick}
                >
                  <IconTrash className="h-5 w-5" />
                  <span className="sr-only">Delete {location.name}</span>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </Card>

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete {location.name}</DialogTitle>
            <DialogDescription asChild>
              {deleteError ? (
                <div className="text-red-500">
                  {deleteError}
                </div>
              ) : (
                <div className="space-y-4">
                  <p>Are you sure you want to delete this physical location?</p>
                  {associatedSublocations.length > 0 && (
                    <>
                      <p>You will also delete all associated sublocations:</p>
                      <ul className="list-disc pl-4 space-y-1">
                        {associatedSublocations.map((sublocation) => (
                          <li key={sublocation.sublocationId}>
                            {sublocation.sublocationName}
                          </li>
                        ))}
                      </ul>
                    </>
                  )}
                  <p>This action cannot be undone.</p>
                </div>
              )}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
              disabled={isDeleting && !deleteError}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleConfirmDelete}
              disabled={isDeleting && !deleteError}
            >
              {isDeleting && !deleteError ? (
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

SinglePhysicalLocationCard.displayName = 'SinglePhysicalLocationCard';
