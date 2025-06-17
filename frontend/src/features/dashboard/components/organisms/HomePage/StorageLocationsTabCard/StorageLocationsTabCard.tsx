import { useState, useMemo } from 'react';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs';
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/shared/components/ui/pagination"
//import type { DigitalStorageService, PhysicalStorageLocation } from './storageLocationsTabCard.mockdata';
import { ITEMS_PER_PAGE } from '@/features/dashboard/lib/constants/dashboard.constants';
import { StorageLocationList } from './StorageLocationList';
import type { LocationIconBgColor } from '@/types/domain/location-types';

// Types needed for StorageLocationList component
export type StorageLocationTabCardPhysicalItem = {
  sublocationId: string;
  sublocationName: string;
  sublocationType: string;
  storedItems: number;
  parentLocationId: string;
  parentLocationName: string;
  parentLocationType: string;
  parentLocationBgColor?: LocationIconBgColor;
};

export type StorageLocationTabCardDigitalItem = {
  logo: string;
  name: string;
  url: string;
  billingCycle: string;
  monthlyFee: number;
  storedItems: number;
}


type StorageLocationsTabCardProps = {
  digitalLocations: StorageLocationTabCardDigitalItem[];
  sublocations: StorageLocationTabCardPhysicalItem[];
}

export function StorageLocationsTabCard({
  digitalLocations,
  sublocations,
}: StorageLocationsTabCardProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const [activeTab, setActiveTab] = useState<'physical' | 'digital'>('physical');

  const services = useMemo(() =>
    activeTab === 'physical' ? (sublocations || []) : (digitalLocations || []),
    [activeTab, sublocations, digitalLocations]
  );
  const totalPages = Math.ceil(services.length / ITEMS_PER_PAGE);
  const shouldShowPagination = services.length > ITEMS_PER_PAGE;

  const paginatedServices = useMemo(() =>
    services.slice(
      (currentPage - 1) * ITEMS_PER_PAGE,
      currentPage * ITEMS_PER_PAGE
    ),
    [currentPage, services]
  );

  const totalPhysicalLocations = sublocations.length;
  const totalDigitalLocations = digitalLocations.length;

  return (
    <Card className="col-span-full lg:col-span-2">
      <CardHeader>
        <CardTitle>Storage Locations</CardTitle>
        <CardDescription>
          {`${totalPhysicalLocations} physical storage locations | ${totalDigitalLocations} online storage locations`}
        </CardDescription>
      </CardHeader>

      <Tabs
        defaultValue="physical"
        className="flex-grow flex flex-col"
        onValueChange={(value) => {
          setActiveTab(value as "physical" | "digital")
          setCurrentPage(1)
        }}
      >
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="physical">Physical Storage</TabsTrigger>
          <TabsTrigger value="digital">Online Storage</TabsTrigger>
        </TabsList>

          {/* Physical Storage */}
          <TabsContent
            value="physical"
            className="flex-grow flex flex-col data-[state=inactive]:hidden"
          >
            <CardContent className="flex-grow overflow-auto pt-4">
              <StorageLocationList services={paginatedServices} isPhysical={true} />
            </CardContent>
          </TabsContent>

          {/* Digital Storage */}
          <TabsContent
            value="digital"
            className="flex-grow flex flex-col data-[state=inactive]:hidden"
          >
            <CardContent className="flex-grow overflow-auto pt-4">
              <StorageLocationList services={paginatedServices} isPhysical={false} />
            </CardContent>
          </TabsContent>

       </Tabs>
      { shouldShowPagination && (
          <CardFooter>
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious
                    onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                    isActive={currentPage === 1}
                    aria-disabled={currentPage === 1}
                    className={currentPage === 1 ? 'pointer-events-none opacity-50' : 'cursor-pointer'}
                  />
                </PaginationItem>
                <PaginationItem>
                  <PaginationNext
                    onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
                    isActive={currentPage === totalPages}
                    aria-disabled={currentPage === totalPages}
                    className={currentPage === totalPages ? 'pointer-events-none opacity-50' : 'cursor-pointer'}
                  />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </CardFooter>
      )}

    </Card>
  );
};
