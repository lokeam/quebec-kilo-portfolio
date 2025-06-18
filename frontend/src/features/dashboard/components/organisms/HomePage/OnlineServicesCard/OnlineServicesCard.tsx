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
import { ITEMS_PER_PAGE } from '@/features/dashboard/lib/constants/dashboard.constants';

type OnlineServicesCardProps = {
  subscriptionTotal: number;
  digitalLocations: {
    logo: string;
    name: string;
    url: string;
    billingCycle: string;
    monthlyFee: number;
    storedItems: number;
    renewsNextMonth: boolean; // NOTE IMPLEMENT FEATURE TO SHOW BADGE IF NOT MONTHLY BILLING CYCLE AND RENEWS NEXT MONTH
  }[];
};

export function OnlineServicesCard({
  subscriptionTotal,
  digitalLocations,
}: OnlineServicesCardProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const totalPages = Math.ceil(digitalLocations.length / ITEMS_PER_PAGE);
  const showShowPagination = digitalLocations.length > ITEMS_PER_PAGE;

  const totaldigitalLocations = digitalLocations.length;
  const subscriptionRecurringNextMonth = digitalLocations.filter(location => location.renewsNextMonth).length;

  const paginatedServices = useMemo(() =>
    digitalLocations.slice(
      (currentPage - 1) * ITEMS_PER_PAGE,
      currentPage * ITEMS_PER_PAGE
    ),
    [currentPage, digitalLocations]
  );

  return (
    <Card className="col-span-full lg:col-span-2 flex flex-col h-full">
      <CardHeader>
        <CardTitle>Online Gaming Services</CardTitle>
        <CardDescription>
          {subscriptionTotal} total annual fees | {totaldigitalLocations} total services | {subscriptionRecurringNextMonth} renews this month
        </CardDescription>
      </CardHeader>

      <CardContent className="flex-1 overflow-auto">
        <ServiceList digitalLocations={paginatedServices} />
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
