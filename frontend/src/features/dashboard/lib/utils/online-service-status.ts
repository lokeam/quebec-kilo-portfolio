
import type { DigitalLocation } from '@/types/domain/online-service';

import { SERVICE_STATUS_CODES, type ServiceStatusCode } from '@/shared/constants/service.constants';

import { formatCurrency } from '@/features/dashboard/lib/utils/formatCurrency';


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

  return hasMonthlyCost && hasBillingCycle && hasPaymentMethod;
}

export function formatServicePrice(amount: string | number | null | undefined): string {
  if (amount == null || amount === 0 || amount === '0' || amount === 'FREE') {
    return 'FREE';
  }

  const numericAmount = typeof amount === 'string' ? parseFloat(amount) : amount;
  return formatCurrency(numericAmount);
}

export function isRenewalMonth(service: DigitalLocation): boolean {
  if (!service?.nextPaymentDate) {
    return false; // If no renewal date, it's not renewing this month
  }

  const currentMonth = new Date().toLocaleString('default', { month: 'long' });
  return currentMonth === service.nextPaymentDate.toLocaleString();
}
