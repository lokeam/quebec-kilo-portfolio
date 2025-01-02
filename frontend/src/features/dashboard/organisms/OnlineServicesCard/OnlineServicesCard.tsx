import { useState, useMemo } from 'react';
import { Card, CardContent, CardDescription,CardFooter, CardHeader, CardTitle } from '@/shared/components/ui/card';
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/shared/components/ui/pagination"
import { ServiceList } from './OnlineServicesList';
import type { OnlineService } from './onlineServicesCard.types';
import { ITEMS_PER_PAGE } from '@/features/dashboard/constants/dashboard.constants';

type OnlineServicesCardProps = {
  totalAnnual: string;
  renewsThisMonth: string[];
  totalServices: number;
  services: OnlineService[];
};

export function OnlineServicesCard({
  totalAnnual,
  renewsThisMonth,
  totalServices,
  services,
}: OnlineServicesCardProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const totalPages = Math.ceil(services.length / ITEMS_PER_PAGE);
  const showShowPagination = services.length > ITEMS_PER_PAGE;

  const paginatedServices = useMemo(() =>
    services.slice(
      (currentPage - 1) * ITEMS_PER_PAGE,
      currentPage * ITEMS_PER_PAGE
    ),
    [currentPage, services]
  );

  return (
    <Card className="col-span-full lg:col-span-2 flex flex-col h-full">
      <CardHeader>
        <CardTitle>Online Gaming Services</CardTitle>
        <CardDescription>
          {totalAnnual} total annual fees | {totalServices} total services | {renewsThisMonth.length} renews this month
        </CardDescription>
      </CardHeader>

      <CardContent className="flex-1 overflow-auto">
        <ServiceList services={paginatedServices} />
      </CardContent>

      { showShowPagination && (
        <CardFooter>
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                  isActive={currentPage === 1}
                  className={currentPage === 1 ? 'pointer-events-none opacity-50' : 'cursor-pointer'}
                />
              </PaginationItem>
              <PaginationItem>
                <PaginationNext
                  onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
                  isActive={currentPage === totalPages}
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
