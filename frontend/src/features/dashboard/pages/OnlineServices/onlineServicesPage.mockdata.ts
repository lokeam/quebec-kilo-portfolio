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


export const mockdata = {
  currentMonth: "Jan",
  previousMonth: "Dec",
  currentMonthTotal: 100,
  previousMonthTotal: 90,
  services: [
    {
      name: 'playstation',
      label: 'Playstation Network',
      logo: 'psnlogo',
      tier: 'essential', // Essential, Extra, Premium
      billingCycle: '1 yr', // 1mo, 3mo, 1yr
      currency: 'USD',
      price: '$6.66',
      paymentMethod: 'visa',
    },
    {
      name: 'xbox',
      label: 'Xbox Network',
      logo: 'xboxlogo',
      tier: 'livegold', // Live Gold, Game Pass, Game Pass Ultimate
      billingCycle: '1 yr',
      currency: 'USD',
      price: '$5',
      paymentMethod: 'visa',
    },
    {
      name: 'nintendo',
      label: 'Nintendo Switch Online',
      logo: 'switchlogo',
      tier: 'online', // Switch Online, Switch Online + Expansion Pack
      billingCycle: '1 yr',
      currency: 'USD',
      price: '$100',
      paymentMethod: 'paypal',
    },
    {
      name: 'steam',
      label: 'Steam',
      logo: 'steamlogo',
      tier: '', // No tier
      billingCycle: 'NA',
      currency: 'USD',
      price: 'FREE',
      paymentMethod: 'paypal',
    },
    {
      name: 'epicgames',
      label: 'Epic Games Store',
      logo: 'epicgameslogo',
      tier: '', // Unknown - requires research
      billingCycle: 'NA',
      currency: 'USD',
      price: 'FREE',
      paymentMethod: 'mastercard',
    },
    {
      name: 'gog',
      label: 'GOG.com',
      logo: 'goglogo',
      tier: '', // Unknown - requires research
      billingCycle: 'NA',
      currency: 'USD',
      price: 'FREE',
      paymentMethod: 'visa',
    },
    {
      name: 'humble',
      label: 'Humble Bundle',
      logo: 'humblelogo',
      tier: '', // Unknown - requires research
      billingCycle: '1 yr',
      currency: 'USD',
      price: '11.99',
      paymentMethod: 'mastercard',
    },
    {
      name: 'greenman',
      label: 'Green Man Gaming',
      logo: 'greenmanlogo',
      tier: '', // Unknown - requires research
      billingCycle: 'NA',
      currency: 'USD',
      price: 'FREE',
      paymentMethod: 'discover',
    },
    {
      name: 'fanatical',
      label: 'Fanatical',
      logo: 'fanaticallogo',
      tier: '', // Unknown - requires research
      billingCycle: 'NA',
      currency: 'USD',
      price: 'FREE',
      paymentMethod: 'googlepay',
    },
    {
      name: 'ubisoft',
      label: 'Ubisoft Connect',
      logo: 'ubisoftlogo',
      tier: '', // classic - 7.99mo, premium - 17.99mo
      billingCycle: '1 mo',
      currency: 'USD',
      price: '$18',
      paymentMethod: 'mastercard',
    },
    {
      name: 'googleplaypass',
      label: 'Google Play Pass',
      logo: 'googleplaypasslogo',
      tier: '', // monthly -  4.99, yearly - 29.99
      billingCycle: '1 yr',
      currency: 'USD',
      price: '$100',
      paymentMethod: 'googlepay',
    },
    {
      name: 'applearcade',
      label: 'Apple Arcade',
      logo: 'applearcadelogo',
      tier: '', // monthly, yearly
      billingCycle: '1 yr',
      currency: 'USD',
      price: '$6.99',
      paymentMethod: 'applepay',
    },
    {
      name: 'netflixgames',
      label: 'Netflix Games',
      logo: 'netflixgameslogo',
      tier: '', // standardads, standard, premium
      billingCycle: '1 yr',
      currency: 'USD',
      price: '$15.49',
      paymentMethod: 'visa',
    },
    {
      name: 'geforce',
      label: 'GeForce Now',
      logo: 'geforcelogo',
      tier: 'priority', // priority, premium
      billingCycle: '1 yr',
      currency: 'USD',
      price: '$5.99',
      paymentMethod: 'amex',
    },
    {
      name: 'primegaming',
      label: 'Prime Gaming',
      logo: 'primegaminglogo',
      tier: '', // prime - 14.99mo/139yr
      billingCycle: '1 yr',
      currency: 'USD',
      price: '$100',
      paymentMethod: 'discover',
    },
    {
      name: 'amazonluna',
      label: 'Amazon Luna',
      logo: 'amazonlunalogo',
      tier: 'monthly', // monthly
      billingCycle: '1 mo',
      currency: 'USD',
      price: '$9.99',
      paymentMethod: 'samsungpay',
    },
    {
      name: 'eaplay',
      label: 'EA Play',
      logo: 'eaplaylogo',
      tier: 'play', // play - 5.99mo/39.99yr, playpro - 16.99mo/119.99yr
      billingCycle: '1 yr', // 1 mo, 1 yr
      currency: 'USD',
      price: '$5.99',
      paymentMethod: 'mastercard',
    },
    {
      name: 'videogamesmonthly',
      label: 'Video Games Monthly',
      logo: 'vg monthly',
      tier: '3up›', // 3up, 4up, 5up, pwpak, megabox
      billingCycle: '1 yr',
      currency: 'USD',
      price: '$34.99', // 3up - 34.99, 4up - 39.99, 5up - 44.99, pwpak - 79.99, megabox - 152.99
      paymentMethod: 'discover',
    },
    {
      name: 'retrogametreasure',
      label: 'Retro Game Treasure',
      logo: 'retrogametreasurelogo',
      tier: '', // monthtomonth, 3monthprepay, 6monthprepay
      billingCycle: '1 yr',
      currency: 'USD',
      price: '$100', // m2m - 39.99, 3mo - 116.97, 6mo - 227.97
      paymentMethod: 'visa',
    },
    {
      name: 'quest+',
      label: 'Meta Quest+',
      logo: 'metaquestlogo',
      tier: 'monthly', // monthly, yearly
      billingCycle: '1 yr',
      currency: 'USD',
      price: '$7.99', // monthly - 7.99, yearly - 59.99
      paymentMethod: 'amex',
    },
  ]
}

export const onlineServicesPageMockData: OnlineService[] = [
  {
    name: 'playstation',
    label: 'Playstation Network',
    logo: 'psnlogo',
    tier: 'essential', // Essential, Extra, Premium
    billingCycle: '1 yr', // 1mo, 3mo, 1yr
    currency: 'USD',
    price: '$6.66',
  },
  {
    name: 'xbox',
    label: 'Xbox Network',
    logo: 'xboxlogo',
    tier: 'livegold', // Live Gold, Game Pass, Game Pass Ultimate
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$5',
  },
  {
    name: 'nintendo',
    label: 'Nintendo Switch Online',
    logo: 'switchlogo',
    tier: 'online', // Switch Online, Switch Online + Expansion Pack
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$100',
  },
  {
    name: 'steam',
    label: 'Steam',
    logo: 'steamlogo',
    tier: '', // No tier
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
  },
  {
    name: 'epicgames',
    label: 'Epic Games Store',
    logo: 'epicgameslogo',
    tier: '', // Unknown - requires research
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
  },
  {
    name: 'gog',
    label: 'GOG.com',
    logo: 'goglogo',
    tier: '', // Unknown - requires research
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
  },
  {
    name: 'humble',
    label: 'Humble Bundle',
    logo: 'humblelogo',
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
    logo: 'fanaticallogo',
    tier: '', // Unknown - requires research
    billingCycle: 'NA',
    currency: 'USD',
    price: 'FREE',
  },
  {
    name: 'ubisoft',
    label: 'Ubisoft Connect',
    logo: 'ubisoftlogo',
    tier: '', // classic - 7.99mo, premium - 17.99mo
    billingCycle: '1 mo',
    currency: 'USD',
    price: '$18',
  },
  {
    name: 'googleplaypass',
    label: 'googlepay Play Pass',
    logo: 'googleplaypasslogo',
    tier: '', // monthly -  4.99, yearly - 29.99
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$100',
  },
  {
    name: 'applearcade',
    label: 'Apple Arcade',
    logo: 'applearcadelogo',
    tier: '', // monthly, yearly
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$6.99',
  },
  {
    name: 'netflixgames',
    label: 'Netflix Games',
    logo: 'netflixgameslogo',
    tier: '', // standardads, standard, premium
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$15.49',
  },
  {
    name: 'geforce',
    label: 'GeForce Now',
    logo: 'geforcelogo',
    tier: 'priority', // priority, premium
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$5.99',
  },
  {
    name: 'primegaming',
    label: 'Prime Gaming',
    logo: 'primegaminglogo',
    tier: '', // prime - 14.99mo/139yr
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$100',
  },
  {
    name: 'amazonluna',
    label: 'Amazon Luna',
    logo: 'amazonlunalogo',
    tier: 'monthly', // monthly
    billingCycle: '1 mo',
    currency: 'USD',
    price: '$9.99',
  },
  {
    name: 'eaplay',
    label: 'EA Play',
    logo: 'eaplaylogo',
    tier: '', // play - 5.99mo/39.99yr, playpro - 16.99mo/119.99yr
    billingCycle: '1 yr', // 1 mo, 1 yr
    currency: 'USD',
    price: '$5.99',
  },
  {
    name: 'videogamesmonthly',
    label: 'Video Games Monthly',
    logo: 'vg monthly',
    tier: '3up›', // 3up, 4up, 5up, pwpak, megabox
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$34.99', // 3up - 34.99, 4up - 39.99, 5up - 44.99, pwpak - 79.99, megabox - 152.99
  },
  {
    name: 'retrogametreasure',
    label: 'Retro Game Treasure',
    logo: 'retrogametreasurelogo',
    tier: '', // monthtomonth, 3monthprepay, 6monthprepay
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$100', // m2m - 39.99, 3mo - 116.97, 6mo - 227.97
  },
  {
    name: 'quest+',
    label: 'Meta Quest+',
    logo: 'metaquestlogo',
    tier: 'monthly', // monthly, yearly
    billingCycle: '1 yr',
    currency: 'USD',
    price: '$7.99', // monthly - 7.99, yearly - 59.99
  },
]
