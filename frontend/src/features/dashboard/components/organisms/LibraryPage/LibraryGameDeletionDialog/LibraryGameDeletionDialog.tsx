import { useState, useMemo } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Checkbox } from '@/shared/components/ui/checkbox';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogTitle,
  DialogHeader,
  DialogDescription
} from '@/shared/components/ui/dialog';
import { BoxIcon } from '@/shared/components/ui/icons';
import { IconCloudDataConnection } from '@/shared/components/ui/icons';
import { useDeleteLibraryGameVersions } from '@/core/api/queries/gameLibrary.queries';
import type { LibraryGameItemRefactoredResponse } from "@/types/domain/library-types";

interface LibraryGameDeletionDialogProps {
  game: LibraryGameItemRefactoredResponse;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onDeleteSuccess: () => void;
}

interface DeletableVersion {
  id: string;
  type: 'physical' | 'digital';
  locationName: string;
  platformName: string;
  platformId: number;
  locationId: string;
  userGameId?: number;
}

export function LibraryGameDeletionDialog({
  game,
  open,
  onOpenChange,
  onDeleteSuccess
}: LibraryGameDeletionDialogProps) {
  const [selectedVersions, setSelectedVersions] = useState<Set<string>>(new Set());
  const [deleteAll, setDeleteAll] = useState(false);

  const deleteGameVersionsMutation = useDeleteLibraryGameVersions();

  // Parse game data into deletable versions
  const deletableVersions = useMemo((): DeletableVersion[] => {
    const versions: DeletableVersion[] = [];

    // Add physical versions
    game.physicalLocations.forEach(location => {
      location.gamePlatformVersions.forEach(platform => {
        versions.push({
          id: `physical-${location.sublocationId}-${platform.platformId}`,
          type: 'physical',
          locationName: `${location.parentLocationName} | ${location.sublocationName}`,
          platformName: platform.platformName,
          platformId: platform.platformId,
          locationId: location.sublocationId,
        });
      });
    });

    // Add digital versions
    game.digitalLocations.forEach(location => {
      location.gamePlatformVersions.forEach(platform => {
        versions.push({
          id: `digital-${location.digitalLocationId}-${platform.platformId}`,
          type: 'digital',
          locationName: location.digitalLocationName,
          platformName: platform.platformName,
          platformId: platform.platformId,
          locationId: location.digitalLocationId,
        });
      });
    });

    return versions;
  }, [game]);

  // Handle "Delete All" checkbox
  const handleDeleteAllChange = (checked: boolean) => {
    setDeleteAll(checked);
    if (checked) {
      setSelectedVersions(new Set());
    }
  };

  // Handle individual version checkbox
  const handleVersionChange = (versionId: string, checked: boolean) => {
    if (deleteAll) {
      setDeleteAll(false);
    }

    const newSelected = new Set(selectedVersions);
    if (checked) {
      newSelected.add(versionId);
    } else {
      newSelected.delete(versionId);
    }
    setSelectedVersions(newSelected);
  };

  // Handle form submission
  const handleSubmit = async () => {
    if (deleteGameVersionsMutation.isPending) return;

    try {
      // Prepare the request data
      const requestData = {
        gameId: game.id,
        deleteAll: deleteAll,
        versions: deleteAll ? [] : deletableVersions
          .filter(v => selectedVersions.has(v.id))
          .map(v => ({
            type: v.type,
            locationId: v.locationId,
            platformId: v.platformId,
          }))
      };

      await deleteGameVersionsMutation.mutateAsync(requestData);

      onDeleteSuccess();
      onOpenChange(false);

      // Reset state
      setSelectedVersions(new Set());
      setDeleteAll(false);
    } catch (error) {
      // Error handling is done in the mutation hook
      console.error('Error in handleSubmit:', error);
    }
  };

  // Determine which versions are selected
  const selectedVersionsList = deleteAll ? deletableVersions : deletableVersions.filter(v => selectedVersions.has(v.id));
  const hasSelections = deleteAll || selectedVersions.size > 0;
  const isSubmitting = deleteGameVersionsMutation.isPending;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Remove Game Versions</DialogTitle>
          <DialogDescription>
            Please select which platform versions you want to remove from your library.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Delete All Option */}
          <div className="flex items-center space-x-2 p-3 border rounded-lg">
            <Checkbox
              id="delete-all"
              checked={deleteAll}
              onCheckedChange={handleDeleteAllChange}
              disabled={isSubmitting}
            />
            <label
              htmlFor="delete-all"
              className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
            >
              Delete all versions of this game
            </label>
          </div>

          {/* Physical Versions */}
          {game.physicalLocations.length > 0 && (
            <div className="space-y-2">
              <h4 className="font-medium text-sm text-muted-foreground flex items-center gap-2">
                <BoxIcon className="h-4 w-4" />
                Physical Copies
              </h4>
              <div className="space-y-2">
                {deletableVersions
                  .filter(v => v.type === 'physical')
                  .map(version => (
                    <div key={version.id} className="flex items-center space-x-2 p-2 border rounded">
                      <Checkbox
                        id={version.id}
                        checked={deleteAll || selectedVersions.has(version.id)}
                        onCheckedChange={(checked) => handleVersionChange(version.id, checked as boolean)}
                        disabled={deleteAll || isSubmitting}
                      />
                      <label
                        htmlFor={version.id}
                        className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 flex-1"
                      >
                        <div className="font-medium">{version.locationName}</div>
                        <div className="text-muted-foreground">{version.platformName}</div>
                      </label>
                    </div>
                  ))}
              </div>
            </div>
          )}

          {/* Digital Versions */}
          {game.digitalLocations.length > 0 && (
            <div className="space-y-2">
              <h4 className="font-medium text-sm text-muted-foreground flex items-center gap-2">
                <IconCloudDataConnection className="h-4 w-4" />
                Digital Copies
              </h4>
              <div className="space-y-2">
                {deletableVersions
                  .filter(v => v.type === 'digital')
                  .map(version => (
                    <div key={version.id} className="flex items-center space-x-2 p-2 border rounded">
                      <Checkbox
                        id={version.id}
                        checked={deleteAll || selectedVersions.has(version.id)}
                        onCheckedChange={(checked) => handleVersionChange(version.id, checked as boolean)}
                        disabled={deleteAll || isSubmitting}
                      />
                      <label
                        htmlFor={version.id}
                        className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 flex-1"
                      >
                        <div className="font-medium">{version.locationName}</div>
                        <div className="text-muted-foreground">{version.platformName}</div>
                      </label>
                    </div>
                  ))}
              </div>
            </div>
          )}
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isSubmitting}
          >
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={handleSubmit}
            disabled={!hasSelections || isSubmitting}
          >
            {isSubmitting ? (
              <>
                <span className="animate-spin mr-2">⊚</span>
                Removing...
              </>
            ) : (
              `Remove ${deleteAll ? 'All' : selectedVersionsList.length} Version${selectedVersionsList.length === 1 ? '' : 's'}`
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}