export interface OnlineService {
  name: string;
  label?: string;
  logo?: string;
  tierName?: string;
  billingCycle?: string;
  quarterlyFee?: string;
  renewalDay?: string;
  isActive?: boolean;
  paymentMethod?: string;
  url: string;
  monthlyFee: string;
  annualFee: string;
  plan: "" | "1 Month" | "3 Month" | "12 Month";
  renewalMonth: string;
  avatar?: string;
};

export interface OnlineServicesData {
  totalAnnual: string;
  renewsThisMonth: string[];
  totalServices: number;
  services: OnlineService[];
};

export interface ServiceListProps {
  digitalLocations: {
    logo: string;
    name: string;
    url: string;
    billingCycle: string;
    monthlyFee: number;
    storedItems: number;
  }[];
};