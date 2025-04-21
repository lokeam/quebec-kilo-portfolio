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
  // Check if billing exists and has fees
  if (!service?.billing?.fees) {
    return true; // If no billing info or fees, consider it free
  }

  // Now safely check the fees
  const { fees } = service.billing;

  // Only check monthly fee since it's the only one guaranteed in the type
  return fees.monthly === 'FREE' || fees.monthly === '0';
}

export function formatCurrency(amount: string | null | undefined): string {
  if (!amount || amount === 'FREE' || amount === '0') {
    return 'FREE';
  }
  // Ensure consistent currency formatting
  return amount.startsWith('$') ? amount : `$${amount}`;
}

export function isRenewalMonth(service: OnlineService): boolean {
  if (!service?.billing?.renewalDate?.month) {
    return false; // If no renewal date, it's not renewing this month
  }

  const currentMonth = new Date().toLocaleString('default', { month: 'long' });
  return currentMonth === service.billing.renewalDate.month;
}
