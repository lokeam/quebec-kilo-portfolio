
import type { OnlineService, OnlineServicesData } from './onlineServicesCard.types';

export const onlineServices: OnlineService[] = [
  {
    name: 'steam',
    label: 'Steam',
    logo: 'steam',
    tierName: '', // No tier
    url: 'https://store.steampowered.com',
    billingCycle: 'NA',
    monthlyFee: 'FREE',
    quarterlyFee: 'FREE',
    annualFee: 'FREE',
    renewalDay: 'NA',
    renewalMonth: 'NA',
    isActive: true,
    paymentMethod: 'Paypal',
    plan: '',
  },
  {
    name: 'gog',
    label: 'GOG.com',
    logo: 'gog',
    tierName: '',
    url: 'https://www.gog.com/games',
    billingCycle: 'NA',
    monthlyFee: 'FREE',
    quarterlyFee: 'FREE',
    annualFee: 'FREE',
    renewalDay: 'NA',
    renewalMonth: 'NA',
    isActive: false,
    paymentMethod: 'Jcb',
    plan: '12 Month',
  },
  {
    name: 'nintendo',
    label: 'Nintendo Switch Online',
    logo: 'nintendo',
    tierName: 'Switch Online', // Switch Online, Switch Online + Expansion Pack
    url: 'https://www.nintendo.com/us/switch/online/',
    billingCycle: '1 yr',
    monthlyFee: '$3.99',
    quarterlyFee: '$11.97',
    annualFee: '$34.96',
    renewalDay: '1',
    renewalMonth: 'January',
    isActive: true,
    paymentMethod: 'Visa',
    plan: '1 Month',
  },
  {
    name: 'playstation',
    label: 'Playstation Network',
    logo: 'playstation',
    tierName: 'essential', // Essential, Extra, Premium
    url: 'https://www.playstation.com/en-us/playstation-network/',
    billingCycle: '1 yr', // 1mo, 3mo, 1yr
    monthlyFee: '$6.66',
    quarterlyFee: '$18.99',
    annualFee: '$66.66',
    renewalDay: '1',
    renewalMonth: 'January',
    isActive: true,
    paymentMethod: 'Visa',
    plan: '3 Month',
  },
  {
    name: 'xbox',
    label: 'Xbox Network',
    logo: 'xbox',
    tierName: 'Game Pass Standard', // PC Game Pass, Game Pass Core, Game Pass Standard, Game Pass Ultimate
    url: 'https://www.xbox.com/en-US/xbox-game-pass',
    billingCycle: '1 mo',
    monthlyFee: '$14.99',
    quarterlyFee: '$44.97',
    annualFee: '$139.92',
    renewalDay: '15',
    renewalMonth: 'March',
    isActive: true,
    paymentMethod: 'Mastercard',
    plan: '1 Month',
  },
  {
    name: 'epicgames',
    label: 'Epic Games Store',
    logo: 'epic',
    tierName: '',
    url: 'https://store.epicgames.com/en-US/',
    billingCycle: 'NA',
    monthlyFee: 'FREE',
    quarterlyFee: 'FREE',
    annualFee: 'FREE',
    renewalDay: 'NA',
    renewalMonth: 'NA',
    isActive: false,
    paymentMethod: 'Alipay',
    plan: '',
  },
  {
    name: 'applearcade',
    label: 'Apple Arcade',
    logo: 'apple',
    tierName: '', // monthly, yearly
    url: 'https://www.apple.com/apple-arcade/',
    billingCycle: '1 mo',
    monthlyFee: '$6.99',
    quarterlyFee: '$20.97',
    annualFee: '$49.99',
    renewalDay: '1',
    renewalMonth: 'August',
    isActive: true,
    paymentMethod: 'Mastercard',
    plan: '1 Month',
  },
  {
    name: 'netflixgames',
    label: 'Netflix Games',
    logo: 'netflix',
    tierName: '', // standardads - $6.99, standard - $15.49, premium - $22.99
    url: 'https://www.netflix.com/',
    billingCycle: '1 yr',
    monthlyFee: '$15.49',
    quarterlyFee: '$46.47',
    annualFee: '$139.41',
    renewalDay: '1',
    renewalMonth: 'January',
    isActive: true,
    paymentMethod: 'Visa',
    plan: '1 Month',
  },
];

export const onlineServicesData: OnlineServicesData = {
  // Calculate total annual fees from non-free services
  totalAnnual: "$380.04", // Sum of Nintendo ($80.04), PSN ($180), and Xbox ($120)

  // Assuming current month is March 2024, these services renew this month
  renewsThisMonth: ["Nintendo Switch Online", "XBOX Live"],

  // Total count of services
  totalServices: onlineServices.length,

  // All services
  services: onlineServices
};
