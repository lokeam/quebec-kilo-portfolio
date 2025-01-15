import { useCallback, useEffect } from 'react';
import { Button } from '@/shared/components/ui/button';
import { LayoutGrid, LayoutList, Sheet } from 'lucide-react';
import { Input } from '@/shared/components/ui/input';

// Components
import { FilterDropdown } from '@/shared/components/ui/FilterDropdown/FilterDropdown';

// Hooks
import { useFilterCheckboxes } from '@/shared/components/ui/FilterDropdown/useFilterCheckboxes';
import { useOnlineServicesStore, ViewModes } from '@/features/dashboard/lib/stores/onlineServicesStore';

// Constants
import { BILLING_CYCLE_OPTIONS, PAYMENT_METHOD_OPTIONS } from '@/shared/components/ui/FilterDropdown/filterOptions.consts';

export function OnlineServicesToolbar() {
  const {
    viewMode,
    setViewMode,
    setSearchQuery,
    setBillingCycleFilters,
    setPaymentMethodFilters,
  } = useOnlineServicesStore();
  const billingCycleFilter = useFilterCheckboxes(
    BILLING_CYCLE_OPTIONS.map(option => option.key)
  );
  const paymentMethodFilter = useFilterCheckboxes(
    PAYMENT_METHOD_OPTIONS.map(option => option.key)
  );

  const handleSearchChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setSearchQuery(event.target.value.toLowerCase());
    },
    [setSearchQuery]
  );

  useEffect(() => {
    const selectedBillingCycles = Object.entries(billingCycleFilter.checkboxes)
      .filter(([_, isChecked]) => isChecked)
      .map(([key]) => key);

    console.log('Selected Billing Cycles:', selectedBillingCycles); // Debug log

    setBillingCycleFilters(selectedBillingCycles);
  }, [billingCycleFilter.checkboxes, setBillingCycleFilters]);

  useEffect(() => {
    const selectedPaymentMethods = Object.entries(paymentMethodFilter.checkboxes)
      .filter(([_, isChecked]) => isChecked)
      .map(([key]) => key);

    console.log('Selected Payment Methods:', selectedPaymentMethods); // Debug log

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
          options={BILLING_CYCLE_OPTIONS}
          width="140px"
          {...billingCycleFilter}
        />

        <FilterDropdown
          label="Payment Method"
          options={PAYMENT_METHOD_OPTIONS}
          width="180px"
          {...paymentMethodFilter}
        />
      </div>

      <div className="flex items-center gap-3">
        <span className="text-sm text-gray-500">View</span>
        <div className="flex bg-black rounded-md p-1 gap-1">
          <Button
            variant={viewMode === ViewModes.GRID ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setViewMode(ViewModes.GRID)}
          >
            <LayoutGrid className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === ViewModes.LIST ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setViewMode(ViewModes.LIST)}
          >
            <LayoutList className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === ViewModes.TABLE ? 'default' : 'ghost'}
            size="icon"
            className="h-8 w-8"
            onClick={() => setViewMode(ViewModes.TABLE)}
          >
            <Sheet className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
