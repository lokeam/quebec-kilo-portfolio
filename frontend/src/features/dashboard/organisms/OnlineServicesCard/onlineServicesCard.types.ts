export interface OnlineService {
  name: string;
  url: string;
  monthlyFee: string;
  annualFee: string;
  plan: "" | "1 Month" | "3 Month" | "12 Month";
  renewalMonth: string;
  avatar: string;
};

export interface OnlineServicesData {
  totalAnnual: string;
  renewsThisMonth: string[];
  totalServices: number;
  services: OnlineService[];
};

export interface ServiceListProps {
  services: OnlineService[];
};