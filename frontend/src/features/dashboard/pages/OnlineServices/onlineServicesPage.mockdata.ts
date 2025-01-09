
export type PaymentMethodType = "Alipay" | "Amex" | "Code" | "CodeFront" | "Diners" | "Discover" |
"Elo" | "Generic" | "Hiper" | "Hipercard" | "Jcb" | "Maestro" | "Mastercard" |
"Mir" | "Paypal" | "Unionpay" | "Visa";

export interface OnlineService {
  name: string;
  label: string;
  logo: string;
  tierName: string;
  billingCycle: string;
  url: string;
  monthlyFee: string;
  quarterlyFee: string;
  annualFee: string;
  renewalDay: string;
  renewalMonth: string;
  isActive: boolean;
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
  },
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
  },
  {
    name: 'humble',
    label: 'Humble Bundle',
    logo: 'humble',
    tierName: '',
    url: 'https://www.humblebundle.com/',
    billingCycle: '1 yr', // 1 mo - $11.99, 1 yr - $129
    monthlyFee: '$11.99',
    quarterlyFee: '$35.97',
    annualFee: '$129',
    renewalDay: '1',
    renewalMonth: 'January',
    isActive: false,
    paymentMethod: 'Paypal',
  },
  {
    name: 'greenman',
    label: 'Green Man Gaming',
    logo: 'greenman',
    tierName: '', // Unknown - requires research
    url: 'https://www.greenmangaming.com/',
    billingCycle: 'NA',
    monthlyFee: 'FREE',
    quarterlyFee: 'FREE',
    annualFee: 'FREE',
    renewalDay: 'NA',
    renewalMonth: 'NA',
    isActive: false,
    paymentMethod: 'Visa',
  },
  {
    name: 'fanatical',
    label: 'Fanatical',
    logo: 'fanatical',
    tierName: '', // Unknown - requires research
    url: 'https://www.fanatical.com/',
    billingCycle: 'NA',
    monthlyFee: 'FREE',
    quarterlyFee: 'FREE',
    annualFee: 'FREE',
    renewalDay: 'NA',
    renewalMonth: 'NA',
    isActive: false,
    paymentMethod: 'Amex',
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
  },
  {
    name: 'geforce',
    label: 'GeForce Now',
    logo: 'nvidia',
    tierName: 'priority', // 1-month: free (add supported), Performance - 9.99, Ultimate - 19.99 // 6-month: Free, Performance - 49.99, Ultimate - 99.99
    url: 'https://www.nvidia.com/en-us/geforce/products/geforce-now/',
    billingCycle: '1 yr',
    monthlyFee: '$4.99',
    quarterlyFee: '$17.97',
    annualFee: '$59.94',
    renewalDay: '1',
    renewalMonth: 'February',
    isActive: false,
    paymentMethod: 'Alipay',
  },
  {
    name: 'primegaming',
    label: 'Prime Gaming',
    logo: 'prime',
    tierName: '', // prime - 14.99mo/139yr
    url: 'https://www.amazongames.com/en-us/support/prime-gaming/articles/what-is-prime-gaming',
    billingCycle: '1 yr',
    monthlyFee: '$14.99',
    quarterlyFee: '$44.97',
    annualFee: '$139.92',
    renewalDay: '1',
    renewalMonth: 'March',
    isActive: true,
    paymentMethod: 'Visa',
  },
  {
    name: 'playpass',
    label: 'Google Play Pass',
    logo: 'playpass',
    tierName: '', // prime - 5.99mo / 29.99yr
    url: 'https://play.google.com/intl/en_us/about/play-pass/',
    billingCycle: '1 mo',
    monthlyFee: '$5.99',
    quarterlyFee: '$17.97',
    annualFee: '$59.94',
    renewalDay: '1',
    renewalMonth: 'August',
    isActive: false,
    paymentMethod: 'Discover',
  },
  {
    name: 'eaplay',
    label: 'EA Play',
    logo: 'ea',
    tierName: '', // play - 5.99mo/39.99yr, playpro - 16.99mo/119.99yr
    url: 'https://www.ea.com/en-us/games/ea-play',
    billingCycle: '1 yr', // 1 mo, 1 yr
    monthlyFee: '$5.99',
    quarterlyFee: '$17.97',
    annualFee: '$59.94',
    renewalDay: '1',
    renewalMonth: 'August',
    isActive: false,
    paymentMethod: 'Mastercard',
  },
  {
    name: 'quest+',
    label: 'Meta Quest+',
    logo: 'meta',
    tierName: 'monthly', // monthly -  7.99, yearly - 59.99
    url: 'https://www.meta.com/questplus',
    billingCycle: '1 yr',
    monthlyFee: '$7.99',
    quarterlyFee: '$23.97',
    annualFee: '$59.94',
    renewalDay: '1',
    renewalMonth: 'August',
    isActive: false,
    paymentMethod: 'Amex',
  },
];


export type AvailableService = {
  name: string;
  label: string;
  logo: string;
  tierName: string;
  billingCycle: string;
  url: string;
  monthlyFee: string;
  quarterlyFee: string;
  annualFee: string;
  renewalDay: string;
  renewalMonth: string;
  isActive: boolean;
};

export const mockAvailableServices: OnlineService[] = [
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
  },
  {
    name: 'steam',
    label: 'Steam',
    logo: 'steam',
    tierName: '', // No tier
    url: 'https://store.steampowered.com',
    billingCycle: 'FREE',
    monthlyFee: 'FREE',
    quarterlyFee: 'FREE',
    annualFee: 'FREE',
    renewalDay: 'NA',
    renewalMonth: 'NA',
    isActive: true,
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
  },
  {
    name: 'humble',
    label: 'Humble Bundle',
    logo: 'humble',
    tierName: '',
    url: 'https://www.humblebundle.com/',
    billingCycle: '1 yr', // 1 mo - $11.99, 1 yr - $129
    monthlyFee: '$11.99',
    quarterlyFee: '$35.97',
    annualFee: '$129',
    renewalDay: '1',
    renewalMonth: 'January',
    isActive: false,
  },
  {
    name: 'greenman',
    label: 'Green Man Gaming',
    logo: 'greenmanlogo',
    tierName: '', // Unknown - requires research
    url: 'https://www.greenmangaming.com/',
    billingCycle: 'NA',
    monthlyFee: 'FREE',
    quarterlyFee: 'FREE',
    annualFee: 'FREE',
    renewalDay: 'NA',
    renewalMonth: 'NA',
    isActive: false,
  },
  {
    name: 'fanatical',
    label: 'Fanatical',
    logo: 'fanatical',
    tierName: '', // Unknown - requires research
    url: 'https://www.fanatical.com/',
    billingCycle: 'NA',
    monthlyFee: 'FREE',
    quarterlyFee: 'FREE',
    annualFee: 'FREE',
    renewalDay: 'NA',
    renewalMonth: 'NA',
    isActive: false,
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
  },
  {
    name: 'geforce',
    label: 'GeForce Now',
    logo: 'nvidia',
    tierName: 'priority', // 1-month: free (add supported), Performance - 9.99, Ultimate - 19.99 // 6-month: Free, Performance - 49.99, Ultimate - 99.99
    url: 'https://www.nvidia.com/en-us/geforce/products/geforce-now/',
    billingCycle: '1 yr',
    monthlyFee: '$4.99',
    quarterlyFee: '$17.97',
    annualFee: '$59.94',
    renewalDay: '1',
    renewalMonth: 'February',
    isActive: false,
  },
  {
    name: 'primegaming',
    label: 'Prime Gaming',
    logo: 'prime',
    tierName: '', // prime - 14.99mo/139yr
    url: 'https://www.amazongames.com/en-us/support/prime-gaming/articles/what-is-prime-gaming',
    billingCycle: '1 yr',
    monthlyFee: '$14.99',
    quarterlyFee: '$44.97',
    annualFee: '$139.92',
    renewalDay: '1',
    renewalMonth: 'March',
    isActive: true,
  },
  {
    name: 'playpass',
    label: 'Google Play Pass',
    logo: 'playpass',
    tierName: '', // prime - 5.99mo / 29.99yr
    url: 'https://play.google.com/intl/en_us/about/play-pass/',
    billingCycle: '1 mo',
    monthlyFee: '$5.99',
    quarterlyFee: '$17.97',
    annualFee: '$59.94',
    renewalDay: '1',
    renewalMonth: 'August',
    isActive: false,
  },
  {
    name: 'eaplay',
    label: 'EA Play',
    logo: 'ea',
    tierName: '', // play - 5.99mo/39.99yr, playpro - 16.99mo/119.99yr
    url: 'https://www.ea.com/en-us/games/ea-play',
    billingCycle: '1 yr', // 1 mo, 1 yr
    monthlyFee: '$5.99',
    quarterlyFee: '$17.97',
    annualFee: '$59.94',
    renewalDay: '1',
    renewalMonth: 'August',
    isActive: false,
  },
  {
    name: 'quest+',
    label: 'Meta Quest+',
    logo: 'meta',
    tierName: 'monthly', // monthly -  7.99, yearly - 59.99
    url: 'https://www.meta.com/questplus',
    billingCycle: '1 yr',
    monthlyFee: '$7.99',
    quarterlyFee: '$23.97',
    annualFee: '$59.94',
    renewalDay: '1',
    renewalMonth: 'August',
    isActive: false,
  },
]
