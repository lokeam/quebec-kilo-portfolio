// frontend/src/core/api/adapters/dashboard.adapter.ts

import type { DashboardResponse, MonthlyExpenditure, MediaTypeDomain } from '@/core/api/services/dashboard.service';
import type { ChartConfig } from '@/shared/components/ui/chart';
import { logger } from '@/core/utils/logger/logger';

export interface AdaptedMonthlyExpenditure {
  date: string;
  oneTimePurchase: number;
  hardware: number;
  dlc: number;
  inGamePurchase: number;
  subscription: number;
}

export interface AdaptedDashboardData extends Omit<DashboardResponse, 'monthlyExpenditures' | 'mediaTypeDomains'> {
  monthlyExpenditures: AdaptedMonthlyExpenditure[];
  mediaTypeDomains: MediaTypeDomain[];
}

// Constants
// Config for dashboard chart:
export const dashboardChartConfig = {
  oneTimePurchase: {
    label: "One Time Purchase",
    color: "hsl(var(--chart-1))",
  },
  dlc: {
    label: "DLC",
    color: "hsl(var(--chart-2))",
  },
  inGamePurchase: {
    label: "In Game Purchase",
    color: "hsl(var(--chart-3))",
  },
  hardware: {
    label: "Hardware",
    color: "hsl(var(--chart-4))",
  },
  subscription: {
    label: "Subscription",
    color: "hsl(var(--chart-5))",
  },
} satisfies ChartConfig;

// Main adapter object
export const dashboardAdapter = {
  /**
   * Filters mediaTypeDomains to only include valid chart domains
   */
  getValidMediaTypeDomains: (mediaTypeDomains: MediaTypeDomain[], chartConfig: ChartConfig): MediaTypeDomain[] => {
    const validDomains = mediaTypeDomains.filter((domain) => chartConfig[domain]);

    if (process.env.NODE_ENV === 'development') {
      const invalidDomains = mediaTypeDomains.filter((domain) => !chartConfig[domain]);
      if (invalidDomains.length > 0) {
        logger.warn('Some mediaTypeDomains are not present in chartConfig:', invalidDomains);
      }
    }

    return validDomains;
  },

  /**
   * Main adapter function that transforms dashboard response
   */
  toDisplayModel: (response: DashboardResponse, chartConfig: ChartConfig): AdaptedDashboardData => {
    const adaptedExpenditures = dashboardAdapter.adaptMonthlyExpenditures(response.monthlyExpenditures);
    const validDomains = dashboardAdapter.getValidMediaTypeDomains(response.mediaTypeDomains, chartConfig);

    logger.debug('Dashboard adapter: Transformed data', {
      originalExpendituresCount: response.monthlyExpenditures.length,
      adaptedExpendituresCount: adaptedExpenditures.length,
      validDomains,
    });

    return {
      ...response,
      monthlyExpenditures: adaptedExpenditures,
      mediaTypeDomains: validDomains,
    };
  },

  // No formatting, just pass through the numbers
  adaptMonthlyExpenditures: (monthlyExpenditures: MonthlyExpenditure[]): AdaptedMonthlyExpenditure[] => {
    return monthlyExpenditures.map(expenditure => ({
      date: expenditure.date,
      oneTimePurchase: expenditure.oneTimePurchase,
      hardware: expenditure.hardware,
      dlc: expenditure.dlc,
      inGamePurchase: expenditure.inGamePurchase,
      subscription: expenditure.subscription,
    }));
  },
};