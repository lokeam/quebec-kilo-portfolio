
type MediaTypeDomain = "games" | "movies" | "oneTimePurchase" | "hardware" | "dlc" | "inGamePurchase" | "subscription";

// @ts-nocheck

// Dashboard
export const homePageMockData = {
  // First row, all single statistics cards
  totalGames: 72,
  gameStats: {
    title: 'Games',
    icon: 'games',
    value: 72,
    lastUpdated: 1736398800000,
  },
  subscriptionStats: {
    title: 'Monthly Online Services Costs',
    icon: 'coin',
    value: 120,
    lastUpdated: 1748750400000,
  },
  digitalLocationStats: {
    title: 'Digital Storage Locations',
    icon: 'onlineServices',
    value: 3,
    lastUpdated: 1744084800000,
  },
  physicalLocationStats: {
    title: 'Physical Storage Locations',
    icon: 'package',
    value: 5,
    lastUpdated: 1693886400000,
  },
  // Online Gaming Services Card
  subscriptionTotal: 380.04,
  subscriptionRecurringNextMonth: 2,
  digitalLocations: [
    {
      logo: 'steam',
      name: 'Steam',
      url: 'https://store.steampowered.com/',
      billingCycle: '',
      monthlyFee: 0,
      storedItems: 35,
    },
    {
      logo: 'gog',
      name: 'GOG.com',
      url: 'https://www.gog.com/games',
      billingCycle: '',
      monthlyFee: 0,
      storedItems: 5,
    },
    {
      logo: 'playstation',
      name: 'Playstation Network',
      url: 'https://www.playstation.com/en-us/psn/',
      billingCycle: '1 month',
      monthlyFee: 15.00,
      storedItems: 5,
    },
    {
      logo: 'nintendo',
      name: 'Nintendo Switch Online',
      url: 'https://www.nintendo.com/switch/online/',
      billingCycle: '1 month',
      monthlyFee: 3.99,
      storedItems: 12,
    },
    {
      logo: 'apple',
      name: 'Apple Arcade',
      url: 'https://www.apple.com/apple-arcade/',
      billingCycle: '1 month',
      monthlyFee: 6.99,
      storedItems: 10,
    },
    {
      logo: 'xbox',
      name: 'Xbox Game Pass',
      url: 'https://www.xbox.com/en-us/xbox-game-pass',
      billingCycle: '12 months',
      monthlyFee: 6.24,
      storedItems: 10,
    },
    {
      logo: 'playpass',
      name: 'Google Play Pass',
      url: 'https://www.google.com/playpass',
      billingCycle: '1 month',
      monthlyFee: 4.99,
      storedItems: 10,
    },
  ],
  // Storage Locations Tab Card
  sublocations: [
    {
      sublocationId: '58',
      sublocationName: 'Study bookshelf',
      sublocationType: 'shelf',
      storedItems: 7,
      parentLocationId: '25',
      parentLocationName: 'Chris Redfield\'s House',
      parentLocationType: 'house',
      parentLocationBgColor: 'red' as const,
    },
    {
      sublocationId: '59',
      sublocationName: 'Living room media console',
      sublocationType: 'console',
      storedItems: 11,
      parentLocationId: '32',
      parentLocationName: 'Jill Valentine\'s Condo',
      parentLocationType: 'apartment',
      parentLocationBgColor: 'blue' as const,
    },
    {
      sublocationId: '60',
      sublocationName: 'Closet by front door',
      sublocationType: 'closet',
      storedItems: 9,
      parentLocationId: '32',
      parentLocationName: 'Jill Valentine\'s Condo',
      parentLocationType: 'apartment',
      parentLocationBgColor: 'blue' as const,
    },
    {
      sublocationId: '60',
      sublocationName: 'Overlanding storage bin',
      sublocationType: 'box',
      storedItems: 4,
      parentLocationId: '32',
      parentLocationName: 'Chris Redfield\'s Truck',
      parentLocationType: 'vehicle',
      parentLocationBgColor: 'green' as const,
    }
  ],
  // Items by Platform Card
  newItemsThisMonth: 2,
  // NOTE: platformList item needs to be dynamically created in the ItemsByPlatformCard in order to associate a fill color for the pie chart
  platformList: [
    {
      platform: 'ps4',
      itemCount: 4,
    },
    {
      platform: 'pc',
      itemCount: 38,
    },
    {
      platform: 'xbox',
      itemCount: 16,
    },
    {
      platform: 'switch',
      itemCount: 11,
    },
    {
      platform: 'mobile',
      itemCount: 3,
    }
  ],
  // Monthly Spending Card
  mediaTypeDomains: ["oneTimePurchase", "hardware", "dlc", "inGamePurchase", "subscription"] as MediaTypeDomain[],
  monthlyExpenditures: [
    {
      date: "2025-01-01",
      oneTimePurchase: 40.00,
      hardware: 150.00,
      dlc: 0,
      inGamePurchase: 6.65,
      subscription: 5,
    },
    {
      date: "2025-02-01",
      oneTimePurchase: 5.99,
      hardware: 0,
      dlc: 0,
      inGamePurchase: 0,
      subscription: 5,
    },
    {
      date: "2025-03-01",
      oneTimePurchase: 0,
      hardware: 0,
      dlc: 0,
      inGamePurchase: 5.99,
      subscription: 20.00,
    },
    {
      date: "2025-04-01",
      oneTimePurchase: 390.1,
      hardware: 359.37,
      dlc: 29.99,
      inGamePurchase: 20.00,
      subscription: 100,
    },
    {
      date: "2025-06-01",
      oneTimePurchase: 480,
      hardware: 534.04,
      dlc: 267.02,
      inGamePurchase: 178.01,
      subscription: 356.02,
    },
    {
      date: "2025-07-01",
      oneTimePurchase: 0,
      hardware: 0,
      dlc: 0,
      inGamePurchase: 0,
      subscription: 0,
    },
    {
      date: "2025-08-01",
      oneTimePurchase: 0,
      hardware: 0,
      dlc: 0,
      inGamePurchase: 0,
      subscription: 0,
    },
    {
      date: "2025-09-01",
      oneTimePurchase: 0,
      hardware: 0,
      dlc: 0,
      inGamePurchase: 0,
      subscription: 0,
    },
    {
      date: "2025-10-01",
      oneTimePurchase: 0,
      hardware: 0,
      dlc: 0,
      inGamePurchase: 0,
      subscription: 0,
    },
    {
      date: "2025-11-01",
      oneTimePurchase: 0,
      hardware: 0,
      dlc: 0,
      inGamePurchase: 0,
      subscription: 0,
    },
    {
      date: "2025-12-01",
      oneTimePurchase: 0,
      hardware: 0,
      dlc: 0,
      inGamePurchase: 0,
      subscription: 0,
    },
  ],
}