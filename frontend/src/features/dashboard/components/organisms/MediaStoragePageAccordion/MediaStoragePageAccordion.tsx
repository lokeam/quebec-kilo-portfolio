import { useState } from 'react';

// Shadcn UI Components
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/shared/components/ui/accordion';
import { Button } from '@/shared/components/ui/button';
import { Avatar } from '@/shared/components/ui/avatar';

// Components
import {
  MediaStoragePageAccordionCard,
} from '@/features/dashboard/components/organisms/MediaStoragePageAccordion/MediaStoragePageAccordionCard';

// Hooks
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';

// Type
import type {
  PhysicalLocation,
  DigitalLocation,
  LocationCardData,
  MediaStorageMetadata
} from '@/features/dashboard/types/media-storage.types';

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
                {type === 'physical' ? (
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
                    <Button>Add Sublocation</Button>
                  </>
                ) : (
                  /* For digital locations, show items directly */
                  <div className="space-y-2 pl-8">
                    {(location as DigitalLocation).items?.map((item, index) => (
                      <div
                        key={`${item.itemLabel}-${index}`}
                        className="flex items-center justify-between py-2"
                      >
                        <div>
                          <p className="font-medium">{item.itemName}</p>
                          <p className="text-sm text-muted-foreground">
                            {item.itemPlatform} {item.itemPlatformVersion}
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
