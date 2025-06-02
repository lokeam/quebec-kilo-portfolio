import type { FormValues } from '@/features/dashboard/components/organisms/OnlineServicesPage/OnlineServiceForm/OnlineServiceForm';

export const DEFAULT_FORM_VALUES: FormValues = {
  name: "",
  isActive: true,
  url: "",
  billingCycle: "",
  costPerCycle: 0,
  nextPaymentDate: new Date(),
  paymentMethod: "",
  isSubscriptionService: false
};