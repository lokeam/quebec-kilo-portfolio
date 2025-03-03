// Shadcn UI Components
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/shared/components/ui/accordion';
import {
  Card,
  CardHeader,
} from '@/shared/components/ui/card';

// Utils
import { motion } from 'framer-motion';
import { useDomainMaps } from '@/features/dashboard/lib/hooks/useDomainMaps';

// Types
import type { SubLocationType } from '@/features/dashboard/lib/types/media-storage/constants';
import type { SubLocation } from '@/features/dashboard/lib/types/media-storage/physical';
import type { GameItem } from '@/features/dashboard/lib/types/media-storage/items';

// Icons
import { IconCloudDataConnection } from '@tabler/icons-react';

export interface LocationCardData {
  name: string;
  description: string;
  locationType: string;
  bgColor?: string;
  items?: GameItem[];
  sublocations?: SubLocation[];
}

interface CardProps {
  card: LocationCardData;
  id: string;
  setActive: (card: LocationCardData) => void;
  isDigital: boolean;
}

export function MediaStoragePageAccordionCard({
  card,
  id,
  setActive,
}: CardProps) {
  const { sublocation: sublocationIcons } = useDomainMaps();

  // Get the icon based on the sublocation type
  const getSubLocationIcon = () => {
    const IconComponent = sublocationIcons[card.locationType.toLowerCase() as Lowercase<SubLocationType>];
    return IconComponent ? (
      <IconComponent className="w-full h-full" />
    ) : (
      <IconCloudDataConnection className="w-full h-full" />
    );
  };

  return (
    <motion.div
      layoutId={`card-${card.name}-${id}`}
      key={`card-${card.name}-${id}`}
      onClick={() => setActive(card)}
    >
      <Card className="cursor-pointer">
        <Accordion type="single" collapsible className="w-full">
          <AccordionItem value="games">
            <AccordionTrigger>
              <CardHeader className="flex flex-row items-center gap-4">
                <motion.div layoutId={`image-${card.name}-${id}`} className="w-14 h-14">
                  <div className="w-10 h-10 shrink-0 text-white flex items-center justify-center">
                    {getSubLocationIcon()}
                  </div>
                </motion.div>
                <div>
                  <motion.h3
                    layoutId={`name-${card.name}-${id}`}
                    className="font-medium"
                  >
                    {card.name}
                  </motion.h3>
                  {/* <motion.p
                    layoutId={`description-${card.description}-${id}`}
                    className="text-muted-foreground"
                  >
                    {card.description}
                  </motion.p> */}
                </div>
              </CardHeader>
            </AccordionTrigger>

            <AccordionContent>
            <div className="space-y-2 pl-24">
              {card.items?.map((item, index) => (
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
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </Card>
    </motion.div>
  )
}
