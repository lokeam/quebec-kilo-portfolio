import { useCallback, useEffect } from 'react';
import { Button } from '@/shared/components/ui/button';
import { LayoutGrid, LayoutList, Sheet } from 'lucide-react';
import { Input } from '@/shared/components/ui/input';

// Components
import { FilterDropdown } from '@/shared/components/ui/FilterDropdown/FilterDropdown';

// Hooks
import { useFilterCheckboxes } from '@/shared/components/ui/FilterDropdown/useFilterCheckboxes';
import { useOnlineServicesStore } from '@/features/dashboard/lib/stores/onlineServicesStore';

interface FilterOption {
  key: string;
  label: string;
}

interface OnlineServicesToolbarProps {
  paymentMethods: FilterOption[];
  billingCycles: FilterOption[];
}

export function OnlineServicesToolbar({ paymentMethods, billingCycles }: OnlineServicesToolbarProps) {
  const {
    viewMode,
    setViewMode,
    setSearchQuery,
    setBillingCycleFilters,
    setPaymentMethodFilters,
  } = useOnlineServicesStore();

  const billingCycleFilter = useFilterCheckboxes(
    billingCycles.map(option => option.key)
  );
  const paymentMethodFilter = useFilterCheckboxes(
    paymentMethods.map(option => option.key)
  );

  const handleSearchChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setSearchQuery(event.target.value.toLowerCase());
    },
    [setSearchQuery]
  );

  useEffect(() => {
    const selectedBillingCycles = Object.entries(billingCycleFilter.checkboxes)
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    .filter(([_, isChecked]) => isChecked === true)
      .map(([key]) => key);

    setBillingCycleFilters(selectedBillingCycles);
  }, [billingCycleFilter.checkboxes, setBillingCycleFilters]);

  useEffect(() => {
    const selectedPaymentMethods = Object.entries(paymentMethodFilter.checkboxes)
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    .filter(([_, isChecked]) => isChecked === true)
    .map(([key]) => key);

    setPaymentMethodFilters(selectedPaymentMethods);
  }, [paymentMethodFilter.checkboxes, setPaymentMethodFilters]);

  return (
    <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
      <div className="flex flex-wrap items-center gap-3">
        <Input
          className="w-[300px]"
          id="filterServices"
          placeholder="Filter Services..."
          onChange={handleSearchChange}
        />
        <span className="text-sm text-gray-500">Filter by</span>

        <FilterDropdown
          label="Billing Cycle"
          options={billingCycles}
          width="140px"
          {...billingCycleFilter}
        />

        <FilterDropdown
          label="Payment Method"
          options={paymentMethods}
          width="180px"
          {...paymentMethodFilter}
        />
      </div>

      <div className="flex items-center gap-3">
        <span className="text-sm text-gray-500">View</span>
        <div className="flex bg-black rounded-md p-1 gap-1">
          <Button
            variant={viewMode === 'grid' ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setViewMode('grid')}
          >
            <LayoutGrid className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === 'list' ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setViewMode('list')}
          >
            <LayoutList className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === 'table' ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setViewMode('table')}
          >
            <Sheet className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
