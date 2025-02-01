
export const onlineServicesPageMockData = {
  services: [
    {
      id: 'psn-001',
      name: 'playstation',
      status: 'active',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: '3 month',
        fees: {
          monthly: '$6.66',
          quarterly: '$18.99',
          annual: '$66.66',
        },
        renewalDate: {
          day: '1',
          month: 'January',
        },
        paymentMethod: 'Visa',
      },
      tier: {
        name: 'essential', // from "tierName"
        features: ['Online Play', 'Monthly Games', 'Cloud Storage'],
      },
      features: ['Cross-Platform Play', 'Cloud Saves', 'Member Discounts'],
      label: 'Playstation Network',
      logo: 'playstation',
      url: 'https://www.playstation.com/en-us/playstation-network/',
    },
    {
      id: 'xbox-002',
      name: 'xbox',
      status: 'active',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: '1 month',
        fees: {
          monthly: '$14.99',
          quarterly: '$44.97',
          annual: '$139.92',
        },
        renewalDate: {
          day: '15',
          month: 'March',
        },
        paymentMethod: 'Mastercard',
      },
      tier: {
        name: 'Game Pass Standard',
        features: ['Game Pass Library', 'Xbox Live Gold', 'Exclusive Discounts'],
      },
      features: ['Cloud Gaming', 'Cross-Platform Play', 'Member Perks'],
      label: 'Xbox Network',
      logo: 'xbox',
      url: 'https://www.xbox.com/en-US/xbox-game-pass',
    },
    {
      id: 'nintendo-003',
      name: 'nintendo',
      status: 'active',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: '1 year',
        fees: {
          monthly: '$3.99',
          quarterly: '$11.97',
          annual: '$34.96',
        },
        renewalDate: {
          day: '1',
          month: 'January',
        },
        paymentMethod: 'Visa',
      },
      tier: {
        name: 'Switch Online',
        features: ['Online Play', 'Classic Game Library'],
      },
      features: ['Cloud Saves', 'Smartphone App', 'Special Offers'],
      label: 'Nintendo Switch Online',
      logo: 'nintendo',
      url: 'https://www.nintendo.com/us/switch/online/',
    },
    {
      id: 'steam-004',
      name: 'steam',
      status: 'active',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: 'NA',
        fees: {
          monthly: 'FREE',
          quarterly: 'FREE',
          annual: 'FREE',
        },
        renewalDate: {
          day: 'NA',
          month: 'NA',
        },
        paymentMethod: 'Paypal',
      },
      tier: {
        name: '',
        features: [],
      },
      features: ['Achievements', 'Cloud Saves', 'Workshop Support'],
      label: 'Steam',
      logo: 'steam',
      url: 'https://store.steampowered.com',
    },
    {
      id: 'epic-005',
      name: 'epicgames',
      status: 'inactive',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: 'NA',
        fees: {
          monthly: 'FREE',
          quarterly: 'FREE',
          annual: 'FREE',
        },
        renewalDate: {
          day: 'NA',
          month: 'NA',
        },
        paymentMethod: 'Alipay',
      },
      tier: {
        name: '',
        features: [],
      },
      features: ['Free Weekly Games', 'Cloud Saves', 'Cross-Platform Support'],
      label: 'Epic Games Store',
      logo: 'epic',
      url: 'https://store.epicgames.com/en-US/',
    },
    {
      id: 'gog-006',
      name: 'gog',
      status: 'inactive',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: 'NA',
        fees: {
          monthly: 'FREE',
          quarterly: 'FREE',
          annual: 'FREE',
        },
        renewalDate: {
          day: 'NA',
          month: 'NA',
        },
        paymentMethod: 'Jcb',
      },
      tier: {
        name: '',
        features: [],
      },
      features: ['DRM-Free Games', 'Cloud Saves', 'Game Backups'],
      label: 'GOG.com',
      logo: 'gog',
      url: 'https://www.gog.com/games',
    },
    {
      id: 'humble-007',
      name: 'humble',
      status: 'inactive',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: '1 year',
        fees: {
          monthly: '$11.99',
          quarterly: '$35.97',
          annual: '$129',
        },
        renewalDate: {
          day: '1',
          month: 'January',
        },
        paymentMethod: 'Paypal',
      },
      tier: {
        name: '',
        features: [],
      },
      features: ['Monthly Bundles', 'Store Discounts', 'Charity Donations'],
      label: 'Humble Bundle',
      logo: 'humble',
      url: 'https://www.humblebundle.com/',
    },
    {
      id: 'greenman-008',
      name: 'greenman',
      status: 'inactive',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: 'NA',
        fees: {
          monthly: 'FREE',
          quarterly: 'FREE',
          annual: 'FREE',
        },
        renewalDate: {
          day: 'NA',
          month: 'NA',
        },
        paymentMethod: 'Visa',
      },
      tier: {
        name: '',
        features: [],
      },
      features: ['PC Game Store', 'Deals and Discounts'],
      label: 'Green Man Gaming',
      logo: 'greenman',
      url: 'https://www.greenmangaming.com/',
    },
    {
      id: 'fanatical-009',
      name: 'fanatical',
      status: 'inactive',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: 'NA',
        fees: {
          monthly: 'FREE',
          quarterly: 'FREE',
          annual: 'FREE',
        },
        renewalDate: {
          day: 'NA',
          month: 'NA',
        },
        paymentMethod: 'Amex',
      },
      tier: {
        name: '',
        features: [],
      },
      features: ['Bundled Deals', 'Storefront Discounts'],
      label: 'Fanatical',
      logo: 'fanatical',
      url: 'https://www.fanatical.com/',
    },
    {
      id: 'apple-arcade-010',
      name: 'applearcade',
      status: 'active',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: '1 month',
        fees: {
          monthly: '$6.99',
          quarterly: '$20.97',
          annual: '$49.99',
        },
        renewalDate: {
          day: '1',
          month: 'August',
        },
        paymentMethod: 'Mastercard',
      },
      tier: {
        name: '',
        features: [],
      },
      features: ['No Ads', 'Family Sharing', 'Offline Play'],
      label: 'Apple Arcade',
      logo: 'apple',
      url: 'https://www.apple.com/apple-arcade/',
    },
    {
      id: 'netflix-games-011',
      name: 'netflixgames',
      status: 'active',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: '1 year',
        fees: {
          monthly: '$15.49',
          quarterly: '$46.47',
          annual: '$139.41',
        },
        renewalDate: {
          day: '1',
          month: 'January',
        },
        paymentMethod: 'Visa',
      },
      tier: {
        name: '',
        features: [],
      },
      features: ['Mobile Games Catalog', 'Ad-Free Gaming'],
      label: 'Netflix Games',
      logo: 'netflix',
      url: 'https://www.netflix.com/',
    },
    {
      id: 'geforce-012',
      name: 'geforce',
      status: 'inactive',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: '1 year',
        fees: {
          monthly: '$4.99',
          quarterly: '$17.97',
          annual: '$59.94',
        },
        renewalDate: {
          day: '1',
          month: 'February',
        },
        paymentMethod: 'Alipay',
      },
      tier: {
        name: 'priority',
        features: ['Priority Queue', 'Extended Session Length'],
      },
      features: ['Cloud Gaming', 'RTX On', 'Cross-Platform'],
      label: 'GeForce Now',
      logo: 'nvidia',
      url: 'https://www.nvidia.com/en-us/geforce/products/geforce-now/',
    },
    {
      id: 'prime-gaming-013',
      name: 'primegaming',
      status: 'active',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: '1 year',
        fees: {
          monthly: '$14.99',
          quarterly: '$44.97',
          annual: '$139.92',
        },
        renewalDate: {
          day: '1',
          month: 'March',
        },
        paymentMethod: 'Visa',
      },
      tier: {
        name: '',
        features: [],
      },
      features: ['Free PC Games', 'In-Game Loot', 'Twitch Subscription'],
      label: 'Prime Gaming',
      logo: 'prime',
      url: 'https://www.amazongames.com/en-us/support/prime-gaming/articles/what-is-prime-gaming',
    },
    {
      id: 'play-pass-014',
      name: 'playpass',
      status: 'inactive',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: '1 month',
        fees: {
          monthly: '$5.99',
          quarterly: '$17.97',
          annual: '$59.94',
        },
        renewalDate: {
          day: '1',
          month: 'August',
        },
        paymentMethod: 'Discover',
      },
      tier: {
        name: '',
        features: [],
      },
      features: ['Ad-Free Apps', 'Family Sharing'],
      label: 'Google Play Pass',
      logo: 'playpass',
      url: 'https://play.google.com/intl/en_us/about/play-pass/',
    },
    {
      id: 'ea-play-015',
      name: 'eaplay',
      status: 'inactive',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: '1 year',
        fees: {
          monthly: '$5.99',
          quarterly: '$17.97',
          annual: '$59.94',
        },
        renewalDate: {
          day: '1',
          month: 'August',
        },
        paymentMethod: 'Mastercard',
      },
      tier: {
        name: '',
        features: [],
      },
      features: ['Early Trials', 'Game Library', 'Member Discounts'],
      label: 'EA Play',
      logo: 'ea',
      url: 'https://www.ea.com/en-us/games/ea-play',
    },
    {
      id: 'quest-plus-016',
      name: 'quest+',
      status: 'inactive',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      type: 'GAMING',
      billing: {
        cycle: '1 year',
        fees: {
          monthly: '$7.99',
          quarterly: '$23.97',
          annual: '$59.94',
        },
        renewalDate: {
          day: '1',
          month: 'August',
        },
        paymentMethod: 'Amex',
      },
      tier: {
        name: 'monthly',
        features: ['Two Free VR Games / Month'],
      },
      features: ['VR Game Discounts', 'App Lab Access'],
      label: 'Meta Quest+',
      logo: 'meta',
      url: 'https://www.meta.com/questplus',
    },
  ],
  totalServices: 16,
};



export type AvailableService = {
  name: string;
  label: string;
  logo: string;
  plan?: string;
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
