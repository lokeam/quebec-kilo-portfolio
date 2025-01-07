export interface OnlineService {
  name: string;
  label: string;
  logo: string;
  tier: string;
  billingCycle: string; // FREE, 1 mo, 3 mo, 1yr
  currency: string;
  price: string;
  paymentMethod?: string;
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
  },
  {
    name: 'xbox',
    label: 'Xbox Network',
    logo: 'xbox',
    tier: 'livegold', // Live Gold, Game Pass, Game Pass Ultimate
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$5',
  },
  {
    name: 'nintendo',
    label: 'Nintendo Switch Online',
    logo: 'nintendo',
    tier: 'online', // Switch Online, Switch Online + Expansion Pack
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$100',
  },
  {
    name: 'steam',
    label: 'Steam',
    logo: 'steam',
    tier: '', // No tier
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
  },
  {
    name: 'epicgames',
    label: 'Epic Games Store',
    logo: 'epic',
    tier: '', // Unknown - requires research
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
  },
  {
    name: 'gog',
    label: 'GOG.com',
    logo: 'gog',
    tier: '', // Unknown - requires research
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
  },
  {
    name: 'humble',
    label: 'Humble Bundle',
    logo: 'humble',
    tier: '', // Unknown - requires research
    billingCycle: '1 yr',
    currency: 'USD',
    price: '11.99',
  },
  {
    name: 'greenman',
    label: 'Green Man Gaming',
    logo: 'greenmanlogo',
    tier: '', // Unknown - requires research
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
  },
  {
    name: 'fanatical',
    label: 'Fanatical',
    logo: 'fanatical',
    tier: '', // Unknown - requires research
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
  },
  {
    name: 'applearcade',
    label: 'Apple Arcade',
    logo: 'apple',
    tier: '', // monthly, yearly
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$6.99',
  },
  {
    name: 'netflixgames',
    label: 'Netflix Games',
    logo: 'netflix',
    tier: '', // standardads, standard, premium
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$15.49',
  },
  {
    name: 'geforce',
    label: 'GeForce Now',
    logo: 'nvidia',
    tier: 'priority', // priority, premium
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$5.99',
  },
  {
    name: 'primegaming',
    label: 'Prime Gaming',
    logo: 'prime',
    tier: '', // prime - 14.99mo/139yr
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$100',
  },
  {
    name: 'playpass',
    label: 'Google Play Pass',
    logo: 'playpass',
    tier: '', // prime - 5.99mo / 29.99yr
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$29.99',
  },
  {
    name: 'eaplay',
    label: 'EA Play',
    logo: 'ea',
    tier: '', // play - 5.99mo/39.99yr, playpro - 16.99mo/119.99yr
    billingCycle: '1 yr', // 1 mo, 1 yr
    currency: 'USD',
    price: '$5.99',
  },
  {
    name: 'quest+',
    label: 'Meta Quest+',
    logo: 'meta',
    tier: 'monthly', // monthly, yearly
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$7.99', // monthly - 7.99, yearly - 59.99
  },
]
