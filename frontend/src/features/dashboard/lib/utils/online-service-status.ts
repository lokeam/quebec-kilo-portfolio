import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import { SERVICE_STATUS_CODES, type ServiceStatusCode } from '@/shared/constants/service.constants';

export function getServiceStatusColor(status: ServiceStatusCode): string {
  const statusColors = {
    [SERVICE_STATUS_CODES.ACTIVE]: 'text-green-500',
    [SERVICE_STATUS_CODES.INACTIVE]: 'text-gray-400',
    [SERVICE_STATUS_CODES.ERROR]: 'text-red-500'
  };

  return statusColors[status] ?? statusColors[SERVICE_STATUS_CODES.INACTIVE];
}

export function isServiceFree(service: OnlineService): boolean {
  return service.billing.fees.monthly === 'FREE' &&
         service.billing.fees.quarterly === 'FREE' &&
         service.billing.fees.annual === 'FREE';
}

export function formatCurrency(amount: string): string {
  if (amount === 'FREE') return amount;
  // Ensure consistent currency formatting
  return amount.startsWith('$') ? amount : `$${amount}`;
}

export function isRenewalMonth(service: OnlineService): boolean {
  const currentMonth = new Date().toLocaleString('default', { month: 'long' });
  return currentMonth === service.billing.renewalDate.month;
}
