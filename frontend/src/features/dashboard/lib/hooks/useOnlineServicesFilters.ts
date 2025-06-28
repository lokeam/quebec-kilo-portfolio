import { useMemo } from 'react';
import type { DigitalLocationBFFResponseItem } from '@/types/domain/digital-location';

interface FilterOption {
  key: string;
  label: string;
}

interface OnlineServicesFilterOptions {
  paymentMethods: FilterOption[];
  billingCycles: FilterOption[];
}

// Helper function to format billing cycle labels
const formatBillingCycleLabel = (cycle: string): string => {
  if (!cycle || cycle.trim() === '') {
    return 'Free';
  }

  // Map common billing cycle formats to user-friendly labels
  const cycleMap: Record<string, string> = {
    '1 month': 'Monthly',
    '3 month': 'Quarterly',
    '6 month': 'Semi-Annual',
    '12 month': 'Annual',
    '1 year': 'Annual',
    '': 'Free'
  };

  return cycleMap[cycle] || cycle;
};

export function useOnlineServicesFilters(services: DigitalLocationBFFResponseItem[]): OnlineServicesFilterOptions {
  return useMemo(() => {
    if (!services || services.length === 0) {
      return { paymentMethods: [], billingCycles: [] };
    }

    // Extract unique payment methods from services
    const uniquePaymentMethods = Array.from(new Set(
      services
        .map(service => service.paymentMethod)
        .filter(method => method && method.trim() !== '') // Filter out empty/null values
    ))
    .sort()
    .map(method => ({
      key: method,
      label: method.charAt(0).toUpperCase() + method.slice(1).toLowerCase() // Capitalize first letter
    }));

    // Extract unique billing cycles from services
    const uniqueBillingCycles = Array.from(new Set(
      services
        .map(service => service.billingCycle)
        .filter(cycle => cycle !== undefined && cycle !== null) // Include empty strings for "Free" services
    ))
    .sort((a, b) => {
      // Sort with empty strings (Free) first, then alphabetically
      if (a === '' && b !== '') return -1;
      if (a !== '' && b === '') return 1;
      return a.localeCompare(b);
    })
    .map(cycle => ({
      key: cycle,
      label: formatBillingCycleLabel(cycle)
    }));

    return {
      paymentMethods: uniquePaymentMethods,
      billingCycles: uniqueBillingCycles
    };
  }, [services]);
}