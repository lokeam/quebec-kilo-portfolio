import type { ServiceStatusCode, ServiceType } from '@/shared/constants/service.constants';

/**
 * Core properties shared by all online services.
 * Provides the foundation for more specific service types.
 */
export interface BaseOnlineService {
  id: string;
  name: string;
  label: string;
  logo: string;
  url: string;
  status: ServiceStatusCode;
  createdAt: Date;
  updatedAt: Date;
  isSubscriptionService: boolean;
  serviceType: ServiceType;
  isActive: boolean;
};

/**
 * Detailed status information for service health monitoring.
 * Used for service availability tracking and notifications.
 */
export type ServiceStatus = {
  code: ServiceStatusCode;
  message?: string;
  lastChecked: string;
};
