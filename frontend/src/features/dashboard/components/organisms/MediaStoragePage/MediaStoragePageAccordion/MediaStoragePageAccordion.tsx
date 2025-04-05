import { useState } from 'react';

// Shadcn UI Components
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/shared/components/ui/accordion';
import { Avatar } from '@/shared/components/ui/avatar';

// Components
import { DrawerContainer } from '@/features/dashboard/components/templates/DrawerContainer' ;
import {
  MediaStoragePageAccordionCard,
} from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/MediaStoragePageAccordionCard';
import { MediaPageSublocationForm } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaPageSublocationForm/MediaPageSublocationForm';

// Hooks
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';

// Type
import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';
import type { DigitalLocation } from '@/features/dashboard/lib/types/media-storage/digital';
import type { LocationCardData } from '@/features/dashboard/components/organisms/MediaStoragePage/MediaStoragePageAccordion/MediaStoragePageAccordionCard';
import type { MediaStorageMetadata } from '@/features/dashboard/lib/types/media-storage/metadata';

// Guards
import { isPhysicalLocation } from '@/features/dashboard/lib/types/media-storage/guards';

interface MediaStoragePageAccordionProps {
  locationData: PhysicalLocation[] | DigitalLocation[];
  title: string;
  meta: MediaStorageMetadata;
  type: 'physical' | 'digital';
}

interface MediaStoragePageAccordionProps {
  locationData: PhysicalLocation[] | DigitalLocation[];
  title: string;
  meta: MediaStorageMetadata;
  type: 'physical' | 'digital';
}

export function MediaStoragePageAccordion({
  locationData,
  title,
  meta,
  type
}: MediaStoragePageAccordionProps) {
  const [activeCard, setActiveCard] = useState<LocationCardData | null>(null);
  const [open, setOpen] = useState<boolean>(false);
  const { games, location: locationIcons } = useDomainMaps();

  const handleSetActive = (card: LocationCardData) => {
    setActiveCard(card);
    console.log("Active card:", card);
  }

  // Get relevant counts based on type
  const itemCount = type === 'physical' ? meta.counts.items.physical : meta.counts.items.digital;
  const locationCount = type === 'physical' ? meta.counts.locations.physical : meta.counts.locations.digital;

  // Helper function to get total items for a location
  const getLocationItemCount = (location: PhysicalLocation | DigitalLocation) => {
    const locationStats = meta.counts.items.byLocation[location.label];
    return locationStats?.total ?? 0;
  }

  // Helper function to get the appropriate icon component
  const getLocationIcon = (location: PhysicalLocation | DigitalLocation) => {
    if (type === 'physical') {
      const physicalLocation = location as PhysicalLocation;
      const IconComponent = locationIcons[physicalLocation.locationType];
      return IconComponent ? <IconComponent className="h-4 w-4" /> : null;
    } else {
      const digitalLocation = location as DigitalLocation;
      const LogoComponent = games[digitalLocation.label];
      return LogoComponent ? <LogoComponent className="h-4 w-4" /> : null;
    }
  }

  return (
    <div className="flex flex-col gap-4 border rounded-md p-4 mb-10">
      <div className="space-y-2">
        <h2 className='text-2xl font-bold tracking-tight'>{title}</h2>
        <div className="text-sm text-muted-foreground">
          <p>{itemCount} items</p>
          <p>{locationCount} locations</p>
        </div>
      </div>

      <Accordion type="single" collapsible className="w-full">
        {locationData.map((location, index) => (
          <AccordionItem key={`item-${index + 1}`} value={`item-${index + 1}`}>
            <AccordionTrigger>
              <div className="flex items-center space-x-3">
                <Avatar className="h-8 w-8 bg-muted flex items-center justify-center">
                  <div className="flex items-center justify-center w-full h-full">
                    {getLocationIcon(location)}
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
                    {location.subLocations?.map((sublocation, cardIndex) => (
                      <MediaStoragePageAccordionCard
                        key={`${location.name}-${cardIndex}`}
                        card={sublocation}
                        id={`${location.name}-${cardIndex}`}
                        setActive={handleSetActive}
                        isDigital={false}
                      />
                    ))}
                    <DrawerContainer
                      open={open}
                      onOpenChange={setOpen}
                      triggerAddLocation="Add Sublocation"
                      title="SubLocation (Unit Storage)"
                      description="Where in your physical location do you keep your games?"
                    >
                      <MediaPageSublocationForm onSuccess={() => setOpen(false)} />
                    </DrawerContainer>
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
        ))}
      </Accordion>
    </div>
  );
}
