
import type { DigitalLocation } from '@/types/domain/online-service';

import { SERVICE_STATUS_CODES, type ServiceStatusCode } from '@/shared/constants/service.constants';

export function getServiceStatusColor(status: ServiceStatusCode): string {
  const statusColors = {
    [SERVICE_STATUS_CODES.ACTIVE]: 'text-green-500',
    [SERVICE_STATUS_CODES.INACTIVE]: 'text-gray-400',
    [SERVICE_STATUS_CODES.ERROR]: 'text-red-500'
  };

  return statusColors[status] ?? statusColors[SERVICE_STATUS_CODES.INACTIVE];
}

export function isPaidService(service: DigitalLocation): boolean {

  const hasMonthlyCost = service?.monthlyCost !== 0;
  const hasBillingCycle = service?.billingCycle !== "";
  const hasPaymentMethod = service?.paymentMethod !== "";

  console.log('testing isPaidService', {hasMonthlyCost, hasBillingCycle, hasPaymentMethod});

  return hasMonthlyCost && hasBillingCycle && hasPaymentMethod;
}

export function formatCurrency(amount: string | null | undefined): string {
  if (!amount || amount === 'FREE' || amount === '0') {
    return 'FREE';
  }
  // Ensure consistent currency formatting
  return amount.startsWith('$') ? amount : `$${amount}`;
}

export function isRenewalMonth(service: DigitalLocation): boolean {
  if (!service?.nextPaymentDate) {
    return false; // If no renewal date, it's not renewing this month
  }

  const currentMonth = new Date().toLocaleString('default', { month: 'long' });
  return currentMonth === service.nextPaymentDate.toLocaleString();
}
