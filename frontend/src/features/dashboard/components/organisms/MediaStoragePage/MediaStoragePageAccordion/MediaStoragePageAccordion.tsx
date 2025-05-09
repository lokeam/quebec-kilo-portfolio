import { useState } from 'react';

// Shadcn UI Components
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/shared/components/ui/accordion';
import { Avatar } from '@/shared/components/ui/avatar';
import { Skeleton } from "@/shared/components/ui/skeleton";

// Components
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer' ;
import {
  MediaStoragePageAccordionCard,
} from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/MediaStoragePageAccordionCard';
import { MediaPageSublocationForm } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaPageSublocationForm/MediaPageSublocationForm';
import { Button } from '@/shared/components/ui/button';
import { PencilIcon, Trash2Icon } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/components/ui/dialog";

// Hooks
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';
import { getLocationIcon } from '@/features/dashboard/lib/utils/getLocationIcon';
import { useLocationManager } from '@/core/api/hooks/useLocationManager';

// Types
import type { PhysicalLocation } from '@/types/domain/physical-location';
import type { DigitalLocation } from '@/types/domain/digital-location';
import type { LocationCardData } from '@/types/domain/location-card';
import type { MediaStorageMetadata } from '@/types/api/storage';
import type { GameItem } from '@/types/domain/game-item';

// Guards
import { isPhysicalLocation } from '@/features/dashboard/lib/types/media-storage/guards';
import { toast } from 'sonner';

interface MediaStoragePageAccordionProps {
  locationData: PhysicalLocation[] | DigitalLocation[];
  title: string;
  meta: MediaStorageMetadata;
  type: 'physical' | 'digital';
  isLoading?: boolean;
}

export function MediaStoragePageAccordion({
  locationData,
  title,
  meta,
  type,
  isLoading = false
}: MediaStoragePageAccordionProps) {
  const [activeCard, setActiveCard] = useState<LocationCardData | null>(null);
  const [openAddDrawer, setOpenAddDrawer] = useState<boolean>(false);
  const [openEditDrawer, setOpenEditDrawer] = useState<boolean>(false);
  const [openDeleteDialog, setOpenDeleteDialog] = useState<boolean>(false);
  const [activeLocationIndex, setActiveLocationIndex] = useState<number | null>(null);
  const [activeCardLocationIndex, setActiveCardLocationIndex] = useState<number | null>(null);
  const domainMaps = useDomainMaps();
  const locationManager = useLocationManager({ type: 'sublocation' });

  // Debug outputs
  console.log("MediaStoragePageAccordion rendered:", {
    type,
    title,
    locationCount: locationData?.length || 0,
    hasData: !!locationData && locationData.length > 0
  });

  if (locationData && locationData.length > 0) {
    console.log("First location data:", {
      name: locationData[0]?.name,
      isPhysical: isPhysicalLocation(locationData[0])
    });
  }

  const handleSetActive = (card: LocationCardData, locationIndex: number) => {
    setActiveCard(card);
    setActiveCardLocationIndex(locationIndex);
    console.log("Active card:", card, "from location index:", locationIndex);
  }

  // Function to handle opening the edit drawer
  const handleEditClick = () => {
    // Only allow editing if the active card belongs to the currently open location
    if (activeCard && activeCardLocationIndex === activeLocationIndex) {
      setOpenEditDrawer(true);
    } else {
      // Show a toast notification if no card is selected or if it's from a different location
      toast.error("Please select a sublocation from this location to edit");
    }
  };

  // Function to handle deleting a sublocation
  const handleDeleteClick = () => {
    // Only allow deleting if the active card belongs to the currently open location
    if (activeCard && isActiveCardFromCurrentLocation) {
      setOpenDeleteDialog(true);
    } else {
      // Show a toast notification if no card is selected or if it's from a different location
      toast.error("Please select a sublocation from this location to delete");
    }
  };

  const handleConfirmDelete = () => {
    if (activeCard && activeLocationIndex !== null) {
      locationManager.delete(activeCard.id);
      setOpenDeleteDialog(false);
      setActiveCard(null);
      setActiveCardLocationIndex(null);
      toast.success(`Sublocation ${activeCard.name} deleted successfully`);
    }
  };

  // Auto-select first sublocation when accordion item is opened
  const handleAccordionValueChange = (value: string) => {
    if (value) {
      // Extract the index from the value (e.g., "item-1" -> 0)
      const index = parseInt(value.split('-')[1]) - 1;
      setActiveLocationIndex(index);

      // Get the location at this index
      const location = locationData[index];

      // If it's a physical location with exactly one sublocation, auto-select it
      if (isPhysicalLocation(location) && location.sublocations?.length === 1) {
        const sublocation = location.sublocations[0];
        const cardData: LocationCardData = {
          id: sublocation.id,
          name: sublocation.name,
          description: sublocation.description,
          locationType: sublocation.type,
          bgColor: sublocation.metadata?.bgColor,
          items: sublocation.items as GameItem[],
          sublocations: [],
          mapCoordinates: sublocation.metadata?.notes,
          createdAt: sublocation.createdAt,
          updatedAt: sublocation.updatedAt
        };
        setActiveCard(cardData);
        setActiveCardLocationIndex(index);
      } else {
        // If the active card is from a different location, clear it
        if (activeCardLocationIndex !== index) {
          setActiveCard(null);
          setActiveCardLocationIndex(null);
        }
      }
    } else {
      // Reset when all accordions are closed
      setActiveLocationIndex(null);
    }
  };

  // Check if the active card belongs to the current location
  const isActiveCardFromCurrentLocation = activeCardLocationIndex === activeLocationIndex;

  // Get relevant counts based on type
  const itemCount = type === 'physical' ? meta.counts.items.physical : meta.counts.items.digital;
  const locationCount = type === 'physical' ? meta.counts.locations.physical : meta.counts.locations.digital;

  // Helper function to get total items for a location
  const getLocationItemCount = (location: PhysicalLocation | DigitalLocation) => {
    if (isPhysicalLocation(location)) {
      // For physical locations, count the number of sublocations
      return location.sublocations?.length || 0;
    } else {
      // For digital locations, count the number of items
      return (location as DigitalLocation).items?.length || 0;
    }
  }

  // Handle add sublocation button click
  const handleAddSublocationClick = () => {
    console.log("Add Sublocation button clicked");
    setOpenAddDrawer(true);
  };

  return (
    <div className="flex flex-col gap-4 border rounded-md p-4 mb-10">
      <DrawerContainer
        open={openAddDrawer}
        onOpenChange={setOpenAddDrawer}
        title="Add Sublocation"
        description="Where in your physical location do you keep your games?"
      >
        <MediaPageSublocationForm
          parentLocationId={activeLocationIndex !== null ? locationData[activeLocationIndex].id : ''}
          onSuccess={() => setOpenAddDrawer(false)}
        />
      </DrawerContainer>

      <div className="space-y-2">
        <h2 className='text-2xl font-bold tracking-tight'>{title}</h2>
        <div className="text-sm text-muted-foreground">
          <p>{itemCount} items</p>
          <p>{locationCount} locations</p>
        </div>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          <Skeleton className="h-12 w-full" />
          <Skeleton className="h-12 w-full" />
          <Skeleton className="h-12 w-full" />
        </div>
      ) : (
      <Accordion
        type="single"
        collapsible
        className="w-full"
        onValueChange={handleAccordionValueChange}
      >
        {Array.isArray(locationData) ? locationData.map((location, index) => (
          <AccordionItem key={`item-${index + 1}`} value={`item-${index + 1}`}>
            <AccordionTrigger>
              <div className="flex items-center space-x-3">
                <Avatar className="h-8 w-8 bg-muted flex items-center justify-center">
                  <div className="flex items-center justify-center w-full h-full">
                    {getLocationIcon(location, type, domainMaps)}
                  </div>
                </Avatar>
                <div className="flex flex-col items-start">
                  <span>{location.name}</span>
                  <span className="text-sm text-muted-foreground">
                    {getLocationItemCount(location)} items
                  </span>
                </div>
              </div>
            </AccordionTrigger>
            <AccordionContent>
              <div className="space-y-4">
                {isPhysicalLocation(location) ? (
                  /* For physical locations, render sublocations as cards */
                  <>
                    {/* Render sublocation cards with selection indicator */}
                    {location.sublocations?.map((sublocation, cardIndex) => {
                      const cardData: LocationCardData = {
                        id: sublocation.id,
                        name: sublocation.name,
                        description: sublocation.description,
                        locationType: sublocation.type,
                        bgColor: sublocation.metadata?.bgColor,
                        items: (sublocation.items || []) as GameItem[],
                        sublocations: [],
                        mapCoordinates: sublocation.metadata?.notes,
                        createdAt: sublocation.createdAt,
                        updatedAt: sublocation.updatedAt
                      };
                      return (
                        <div
                          key={`${location.name}-${cardIndex}`}
                          className={`relative rounded-lg transition-all ${
                            activeCard?.id === cardData.id && isActiveCardFromCurrentLocation
                              ? 'ring-2 ring-primary ring-offset-2'
                              : ''
                          }`}
                        >
                          <MediaStoragePageAccordionCard
                            card={cardData}
                            id={`${location.name}-${cardIndex}`}
                            setActive={(card) => handleSetActive(card, index)}
                            isDigital={false}
                          />
                          {activeCard?.id === cardData.id && isActiveCardFromCurrentLocation && (
                            <div className="absolute -top-2 -right-2 bg-primary text-primary-foreground rounded-full p-1">
                              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                <polyline points="20 6 9 17 4 12"></polyline>
                              </svg>
                            </div>
                          )}
                        </div>
                      );
                    })}

                    {/* Add Sublocation button - moved outside the loop */}
                    {type === 'physical' && (
                      <Button
                        variant="default"
                        onClick={handleAddSublocationClick}
                        className="my-4"
                      >
                        Add Sublocation to {location.name}
                      </Button>
                    )}

                    {/* Helper text when a sublocation is auto-selected */}
                    {isPhysicalLocation(location) &&
                     location.sublocations?.length === 1 &&
                     activeCard?.id === location.sublocations[0].id &&
                     isActiveCardFromCurrentLocation && (
                      <div className="text-sm text-muted-foreground italic mb-2">
                        <span className="font-medium">{activeCard?.name}</span> selected. You can edit it using the button below.
                      </div>
                    )}

                    {/* Action buttons for Add/Edit */}
                    <div className={`flex space-x-2 mt-4 ${location.sublocations?.length === 0 ? 'hidden' : ''}`}>
                      {/* Edit button with enhanced visual cues */}
                      <Button
                        variant={activeCard && isActiveCardFromCurrentLocation ? "default" : "outline"}
                        onClick={handleEditClick}
                        disabled={!activeCard || !isActiveCardFromCurrentLocation}
                        className="flex items-center relative"
                      >
                        <PencilIcon className="h-4 w-4 mr-2" />
                        Edit Sublocation
                        {activeCard && isActiveCardFromCurrentLocation && (
                          <span className="absolute -top-2 -right-2 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[10px] text-primary-foreground">
                            ✓
                          </span>
                        )}
                      </Button>

                      {/* Delete button */}
                      <Button
                        variant="destructive"
                        onClick={handleDeleteClick}
                        disabled={!activeCard || !isActiveCardFromCurrentLocation}
                        className="flex items-center"
                      >
                        <Trash2Icon className="h-4 w-4 mr-2" />
                        Delete Sublocation
                      </Button>

                      {/* Edit Sublocation Drawer - without a trigger */}
                      <DrawerContainer
                        open={openEditDrawer}
                        onOpenChange={setOpenEditDrawer}
                        title="Edit Sublocation"
                        description="Update details for this sublocation"
                      >
                        {activeCard && isActiveCardFromCurrentLocation && activeLocationIndex !== null && (
                          <MediaPageSublocationForm
                            sublocationData={{
                              ...activeCard,
                              locationType: activeCard.locationType
                            }}
                            parentLocationId={locationData[activeLocationIndex].id}
                            isEditing={true}
                            buttonText="Update Sublocation"
                            onSuccess={() => {
                              setOpenEditDrawer(false);
                              setActiveCard(null);
                              setActiveCardLocationIndex(null);
                            }}
                          />
                        )}
                      </DrawerContainer>
                    </div>

                    {/* Delete Confirmation Dialog */}
                    <Dialog open={openDeleteDialog} onOpenChange={setOpenDeleteDialog}>
                      <DialogContent>
                        <DialogHeader>
                          <DialogTitle>Delete Sublocation</DialogTitle>
                          <DialogDescription>
                            Are you sure you want to delete the sublocation "{activeCard?.name}"? This action cannot be undone.
                          </DialogDescription>
                        </DialogHeader>
                        <DialogFooter>
                          <Button
                            variant="outline"
                            onClick={() => setOpenDeleteDialog(false)}
                          >
                            Cancel
                          </Button>
                          <Button
                            variant="destructive"
                            onClick={handleConfirmDelete}
                          >
                            Delete
                          </Button>
                        </DialogFooter>
                      </DialogContent>
                    </Dialog>
                  </>
                ) : (
                  /* For digital locations, show items directly */
                  <div className="space-y-2 pl-8">
                    {(location as DigitalLocation).items?.map((item, index) => (
                      <div
                        key={`${item.label}-${index}`}
                        className="flex items-center justify-between py-2"
                      >
                        <div>
                          <p className="font-medium">{item.name}</p>
                          <p className="text-sm text-muted-foreground">
                            {item.platform.charAt(0).toUpperCase() + item.platform.slice(1)} {item.platformVersion}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </AccordionContent>
          </AccordionItem>
        )) : null}
      </Accordion>
      )}
    </div>
  );
}