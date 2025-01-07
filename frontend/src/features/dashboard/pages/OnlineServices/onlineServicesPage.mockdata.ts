
export type PaymentMethodType = "Alipay" | "Amex" | "Code" | "CodeFront" | "Diners" | "Discover" |
"Elo" | "Generic" | "Hiper" | "Hipercard" | "Jcb" | "Maestro" | "Mastercard" |
"Mir" | "Paypal" | "Unionpay" | "Visa";

export interface OnlineService {
  name: string;
  label: string;
  logo: string;
  tier: string;
  billingCycle: string; // FREE, 1 mo, 3 mo, 1yr
  currency: string;
  price: string;
  paymentMethod?: PaymentMethodType;
}

export interface OnlineServicesList {
  services: OnlineService[];
};


export const onlineServicesPageMockData: OnlineService[] = [
  {
    name: 'playstation',
    label: 'Playstation Network',
    logo: 'playstation',
    tier: 'essential', // Essential, Extra, Premium
    billingCycle: '1 yr', // 1mo, 3mo, 1yr
    currency: 'USD',
    price: '$6.66',
    paymentMethod: 'Visa',
  },
  {
    name: 'xbox',
    label: 'Xbox Network',
    logo: 'xbox',
    tier: 'livegold', // Live Gold, Game Pass, Game Pass Ultimate
    billingCycle: '3 mo',
    currency: 'USD',
    price: '$5',
    paymentMethod: 'Mastercard',
  },
  {
    name: 'nintendo',
    label: 'Nintendo Switch Online',
    logo: 'nintendo',
    tier: 'online', // Switch Online, Switch Online + Expansion Pack
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$100',
    paymentMethod: 'Visa',
  },
  {
    name: 'steam',
    label: 'Steam',
    logo: 'steam',
    tier: '', // No tier
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
    paymentMethod: 'Paypal',
  },
  {
    name: 'epicgames',
    label: 'Epic Games Store',
    logo: 'epic',
    tier: '', // Unknown - requires research
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
    paymentMethod: 'Alipay',
  },
  {
    name: 'gog',
    label: 'GOG.com',
    logo: 'gog',
    tier: '', // Unknown - requires research
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
    paymentMethod: 'Paypal',
  },
  {
    name: 'humble',
    label: 'Humble Bundle',
    logo: 'humble',
    tier: '', // Unknown - requires research
    billingCycle: '1 yr',
    currency: 'USD',
    price: '11.99',
    paymentMethod: 'Visa',
  },
  {
    name: 'greenman',
    label: 'Green Man Gaming',
    logo: 'greenmanlogo',
    tier: '', // Unknown - requires research
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
    paymentMethod: 'Mir',
  },
  {
    name: 'fanatical',
    label: 'Fanatical',
    logo: 'fanatical',
    tier: '', // Unknown - requires research
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
    paymentMethod: 'Mastercard',
  },
  {
    name: 'applearcade',
    label: 'Apple Arcade',
    logo: 'apple',
    tier: '', // monthly, yearly
    billingCycle: '1 mo',
    currency: 'USD',
    price: '$6.99',
    paymentMethod: 'Visa',
  },
  {
    name: 'netflixgames',
    label: 'Netflix Games',
    logo: 'netflix',
    tier: '', // standardads, standard, premium
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$15.49',
    paymentMethod: 'Jcb',
  },
  {
    name: 'geforce',
    label: 'GeForce Now',
    logo: 'nvidia',
    tier: 'priority', // priority, premium
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$5.99',
    paymentMethod: 'Amex',
  },
  {
    name: 'primegaming',
    label: 'Prime Gaming',
    logo: 'prime',
    tier: '', // prime - 14.99mo/139yr
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$100',
    paymentMethod: 'Visa',
  },
  {
    name: 'playpass',
    label: 'Google Play Pass',
    logo: 'playpass',
    tier: '', // prime - 5.99mo / 29.99yr
    billingCycle: '1 mo',
    currency: 'USD',
    price: '$29.99',
    paymentMethod: 'Jcb',
  },
  {
    name: 'eaplay',
    label: 'EA Play',
    logo: 'ea',
    tier: '', // play - 5.99mo/39.99yr, playpro - 16.99mo/119.99yr
    billingCycle: '1 yr', // 1 mo, 1 yr
    currency: 'USD',
    price: '$5.99',
    paymentMethod: 'Paypal',
  },
  {
    name: 'quest+',
    label: 'Meta Quest+',
    logo: 'meta',
    tier: 'monthly', // monthly, yearly
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$7.99', // monthly - 7.99, yearly - 59.99
    paymentMethod: 'Alipay',
  },
]