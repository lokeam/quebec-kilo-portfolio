import {
  PhysicalLocationType,
  SublocationType,
  GamePlatform
} from '@/features/dashboard/lib/types/media-storage/constants';
import type { PhysicalLocation } from '@/features/dashboard/lib/types/media-storage/physical';
import type { DigitalLocation } from '@/features/dashboard/lib/types/media-storage/digital';
import type { MediaStorageMetadata } from '@/features/dashboard/lib/types/media-storage/metadata';

interface MediaStorageResponse {
  data: {
    physicalLocations: PhysicalLocation[];
    digitalLocations: DigitalLocation[];
  };
  meta: MediaStorageMetadata;
}

export const mediaStoragePageMockData: MediaStorageResponse = {
  data: {
    name: '',
    label: '',
    physicalLocations: [
      {
        id: 'condostudy-1', // Required id
        name: 'Condo study',
        label: 'condostudy',
        type: 'physical', // Required discriminator
        locationType: PhysicalLocationType.APARTMENT,
        mapCoordinates: '40°44\'53.2"N 73°59\'05.6"W',
        createdAt: new Date('2024-01-01'), // Required timestamp
        updatedAt: new Date('2024-01-27'), // Required timestamp
        sublocations: [
          {
            id: 'bookshelf-a-1', // Required id
            name: 'Bookshelf A',
            description: 'Bookshelf A',
            locationType: SublocationType.SHELF,
            items: [
              {
                id: 'gradius-v-1',
                name: 'Gradius V',
                label: 'gradiusv',
                type: 'physical', // NOTE: Backend must add this discriminator
                platform: GamePlatform.PLAYSTATION,
                platformVersion: 'PS2',
                acquiredDate: new Date('2023-01-15'),
                // @ts-expect-error: Backend must add this discriminator
                condition: 'excellent',
                hasOriginalCase: true,
              },
              {
                id: 'streetfighter-iii-ts-1',
                name: 'Street Fighter III: Third Strike',
                label: 'streetfighteriiithirdstrike',
                platform: GamePlatform.PLAYSTATION,
                platformVersion: 'PS2',
                acquiredDate: new Date('2023-01-16'),
                // @ts-expect-error: Backend must add this discriminator
                condition: 'good',
                hasOriginalCase: false,
                hasManual: false
              },
              {
                id: 'zelda-wind-waker-1',
                name: 'The Legend of Zelda: Wind Waker',
                label: 'legendofzeldawindwaker',
                platform: GamePlatform.NINTENDO,
                platformVersion: 'Gamecube',
                acquiredDate: new Date('2023-01-17'),
                // @ts-expect-error: Backend must add this discriminator
                condition: 'good',
                hasOriginalCase: false,
                hasManual: false
              },
              {
                id: 'metroid-prime-1',
                name: 'Metroid Prime',
                label: 'metroidprime',
                platform: GamePlatform.NINTENDO,
                platformVersion: 'Gamecube',
                acquiredDate: new Date('2023-01-18'),
                // @ts-expect-error: Backend must add this discriminator
                condition: 'good',
                hasOriginalCase: false,
                hasManual: false
              },
              {
                id: 'metal-gear-solid-3-snake-eater-1',
                name: 'Metal Gear Solid 3: Snake Eater',
                label: 'metalgearsolid3snakeeater',
                platform: GamePlatform.PLAYSTATION,
                platformVersion: 'PS3',
                acquiredDate: new Date('2023-01-19'),
                // @ts-expect-error: Backend must add this discriminator
                condition: 'good',
                hasOriginalCase: false,
                hasManual: false
              },
              {
                id: 'dragon-age-origins-1',
                name: 'Dragon Age: Origins',
                label: 'dragonageorigins',
                platform: GamePlatform.PLAYSTATION,
                platformVersion: 'PS3',
                acquiredDate: new Date('2023-01-20'),
                // @ts-expect-error: Backend must add this discriminator
                condition: 'good',
                hasOriginalCase: false,
                hasManual: false
              },
              {
                id: 'shadow-of-the-colossus-1',
                name: 'Shadow of the Colossus',
                label: 'shadowofthecolossus',
                platform: GamePlatform.PLAYSTATION,
                platformVersion: 'PS3',
                acquiredDate: new Date('2023-01-21'),
                // @ts-expect-error: Backend must add this discriminator
                condition: 'good',
                hasOriginalCase: false,
                hasManual: false
              }
            ],
          },
        ],
        items: [],
      },
      {
        id: 'sisters-house-1',
        name: 'Sister\'s house',
        label: 'sistershouse',
        type: 'physical',
        locationType: PhysicalLocationType.HOUSE,
        mapCoordinates: '28°33\'01.1"N 81°29\'30.4"W',
        createdAt: new Date('2024-01-01'),
        updatedAt: new Date('2024-01-27'),
        sublocations: [
          {
            name: 'Living room media cabinet',
            description: 'Living room media cabinet',
            items: [
              {
                id: 'super-mario-party-jamboree-1',
                name: 'Super Mario Party Jamboree',
                label: 'supermariopartyjamboree',
                platform: GamePlatform.NINTENDO,
                platformVersion: 'Switch',
                acquiredDate: new Date('2023-01-22'),
                // @ts-expect-error: Backend must add this discriminator
                condition: 'good',
                hasOriginalCase: false,
                hasManual: false
              },
              {
                id: 'luigis-mansion-3-1',
                name: 'Luigi\'s Mansion 3',
                label: 'luigisman3',
                platform: GamePlatform.NINTENDO,
                platformVersion: 'Switch',
                acquiredDate: new Date('2023-01-22'),
                // @ts-expect-error: Backend must add this discriminator
                condition: 'good',
                hasOriginalCase: false,
                hasManual: false
              }
            ],
            locationType: SublocationType.CONSOLE,
          }
        ],
        items: [],
      }
    ],
    digitalLocations: [
      {
        name: 'Steam',
        label: GamePlatform.STEAM,
        url: 'https://store.steampowered.com/',
        isActive: true,
        isFree: true,
        monthlyFee: '0',
        locationImage: '',
        mapCoordinates: '',
        items: [
          {
            id: 'dark-souls-3-1',
            name: 'Dark Souls 3',
            label: 'darksouls3',
            platform: GamePlatform.STEAM,
            platformVersion: 'PC',
            acquiredDate: new Date('2023-01-23'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'dark-souls-remastered-1',
            name: 'Dark Souls: REMASTERED',
            label: 'darksoulsremastered',
            platform: GamePlatform.STEAM,
            platformVersion: 'PC',
            acquiredDate: new Date('2023-01-24'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'elden-ring-1',
            name: 'ELDEN RING',
            label: 'eldenring',
            platform: GamePlatform.STEAM,
            platformVersion: 'PC',
            acquiredDate: new Date('2023-01-25'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'helldivders-2-1',
            name: 'HELLDIVERS 2',
            label: 'helldivders2',
            platform: GamePlatform.STEAM,
            platformVersion: 'PC',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'yakuza-0-1',
            name: 'Yakuza 0',
            label: 'yakuzazero',
            platform: GamePlatform.STEAM,
            platformVersion: 'PC',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'yakuza-kiwami-1',
            name: 'Yakuza Kiwami',
            label: 'yakuzakiwami',
            platform: GamePlatform.STEAM,
            platformVersion: 'PC',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'bloodstained-curse-of-the-moon-1',
            name: 'Bloodstained: Curse of the Moon',
            label: 'bloodstainedcurseofthemoon',
            platform: GamePlatform.STEAM,
            platformVersion: 'PC',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'bloodstained-curse-of-the-moon-2-1',
            name: 'Bloodstained: Curse of the Moon 2',
            label: 'bloodstainedcurseofthemoon2',
            platform: GamePlatform.STEAM,
            platformVersion: 'PC',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'resident-evil-2-1',
            name: 'Resident Evil 2',
            label: 'residentevil2',
            platform: GamePlatform.STEAM,
            platformVersion: 'PC',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'resident-evil-3-1',
            name: 'Resident Evil 3',
            label: 'residentevil3',
            platform: GamePlatform.STEAM,
            platformVersion: 'PC',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'resident-evil-4-1',
            name: 'Resident Evil 4',
            label: 'residentevil4',
            platform: GamePlatform.STEAM,
            platformVersion: 'PC',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          }
        ],
      },
      {
        name: 'Playstation Network',
        label: GamePlatform.PLAYSTATION,
        url: 'https://www.playstation.com/en-us/playstation-plus/',
        isActive: true,
        isFree: false,
        monthlyFee: '10',
        locationImage: '',
        mapCoordinates: '',
        items: [
          {
            id: 'bloodbourne-1',
            name: 'Bloodbourne',
            label: 'bloodbourne',
            platform: GamePlatform.PLAYSTATION,
            platformVersion: 'PS3',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'wanted-dead-1',
            name: 'Wanted: Dead',
            label: 'wanteddead',
            platform: GamePlatform.PLAYSTATION,
            platformVersion: 'PS5',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'returnal-1',
            name: 'Returnal',
            label: 'returnal',
            platform: GamePlatform.PLAYSTATION,
            platformVersion: 'PS5',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'elden-ring-1',
            name: 'Elden Ring',
            label: 'eldenring',
            platform: GamePlatform.PLAYSTATION,
            platformVersion: 'PS5',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: "hogwart's-legacy-1",
            name: "Hogwart's Legacy",
            label: 'hogwartslegacy',
            platform: GamePlatform.PLAYSTATION,
            platformVersion: 'PS5',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          }
        ],
      },
      {
        name: 'Nintendo Switch Online',
        label: GamePlatform.NINTENDO,
        url: 'https://www.nintendo.com/switch/nintendo-switch-online/',
        isActive: true,
        isFree: false,
        monthlyFee: '5',
        locationImage: '',
        mapCoordinates: '',
        items: [
          {
            id: 'limbo-1',
            name: 'Limbo',
            label: 'limbo',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'the-legend-of-zelda-breath-of-the-wild-1',
            name: 'The Legend of Zelda: Breath of the Wild',
            label: 'legendofzeldabreathofthewild',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'the-legend-of-zelda-tears-of-the-kingdom-1',
            name: 'The Legend of Zelda: Tears of the Kingdom',
            label: 'legendofzeldatearsofthekingdom',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'stardew-valley-1',
            name: 'Stardew Valley',
            label: 'stardewvalley',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'pikmin-4-1',
            name: 'Pikmin 4',
            label: 'pikmin4',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'capcom-beat\'em-up-bundle-1',
            name: 'Capcom Beat\'em Up Bundle',
            label: 'capcombeatemupbundle',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'arcade-archives-in-the-hunt-1',
            name: 'Arcade Archives: In the Hunt',
            label: 'arcadearchivesinthehunt',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'manifold-garden-1',
            name: 'Manifold Garden',
            label: 'manifoldgarden',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'okami-hd-1',
            name: 'Okami HD',
            label: 'okamihd',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'carto-1',
            name: 'Carto',
            label: 'carto',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'overcooked-2-1',
            name: 'Overcooked! 2',
            label: 'overcooked2',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'lovers-in-a-dangerous-spacetime-1',
            name: 'Lovers in a Dangerous Spacetime',
            label: 'loversinadangerousspacetime',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'mega-man-legacy-collection-1',
            name: 'Mega Man Legacy Collection',
            label: 'megamanlegacycollection',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'crawl-1',
            name: 'Crawl',
            label: 'crawl',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'tetris-99-1',
            name: 'Tetris 99',
            label: 'tetris99',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
          {
            id: 'arcade-archives-puzzle-bobble-1',
            name: 'Arcade Archives: Puzzle Bobble',
            label: 'arcadearchivespuzzlebobble',
            platform: GamePlatform.NINTENDO,
            platformVersion: 'Switch',
            acquiredDate: new Date('2023-01-26'),
            // @ts-expect-error: Backend must add this discriminator
            condition: 'good',
            hasOriginalCase: false,
            hasManual: false
          },
        ],
      }
    ],
  },
  meta: {
    counts: {
      locations: {
        total: 5,
        physical: 2,
        digital: 3,
      },
      items: {
        total: 41,
        physical: 9,
        digital: 32,
        byLocation: {
          'condostudy': {
            total: 7,
            inSublocations: 7,
          },
          'sistershouse': {
            total: 2,
            inSublocations: 2,
          },
          'steam': {
            total: 11,
            inSublocations: 0,
          },
          'sony': {
            total: 5,
            inSublocations: 0,
          },
          'nintendo': {
            total: 16,
            inSublocations: 0,
          },
        },
      },
    },
  },
};

// interface GameItem {
//   itemName: string;
//   itemLabel: string;
//   itemPlatform: string;
//   itemPlatformVersion: string;
// }

// export interface BaseLocation {
//   name: string;
//   label: string;
//   subLocations?: SubLocation[];  // Use SubLocation instead of LocationCardData
// }

// export interface SubLocation {
//   name: string;
//   description: string;
//   locationType: SubLocationType;
//   items?: GameItem[];
// }

// interface MediaStorageMetadata {
//   counts: {
//     locations: {
//       total: number;
//       physical: number;
//       digital: number;
//     };
//     items: {
//       total: number;
//       physical: number;
//       digital: number;
//       byLocation: Record<string, {
//         total: number;
//         inSublocations: number;
//       }>;
//     };
//   };
// }

// interface MediaStorageResponse {
//   data: {
//     name: string;
//     label: string;
//     physicalLocations: PhysicalLocation[];
//     digitalLocations: DigitalLocation[];
//   };
//   meta: MediaStorageMetadata;
// }

// export const mediaStoragePageMockData: MediaStorageResponse = {
//   data: {
//     name: '',
//     label: '',
//     physicalLocations: [
//       {
//         name: 'Condo study',
//         label: 'condostudy',
//         locationType: PhysicalLocationType.APARTMENT,
//         mapCoordinates: '40°44\'53.2"N 73°59\'05.6"W',
//         sublocations: [
//           {
//             name: 'Bookshelf A',
//             description: 'Bookshelf A',
//             items: [
//               {
//                 itemName: 'Gradius V',
//                 itemLabel: 'gradiusv',
//                 itemPlatform: 'Playstation',
//                 itemPlatformVersion: 'PS2',
//               },
//               {
//                 itemName: 'Street Fighter III: Third Strike',
//                 itemLabel: 'streetfighteriiithirdstrike',
//                 itemPlatform: 'Playstation',
//                 itemPlatformVersion: 'PS2',
//               },
//               {
//                 itemName: 'The Legend of Zelda: Wind Waker',
//                 itemLabel: 'legendofzeldawindwaker',
//                 itemPlatform: 'Nintendo',
//                 itemPlatformVersion: 'Gamecube',
//               },
//               {
//                 itemName: 'Metroid Prime',
//                 itemLabel: 'metroidprime',
//                 itemPlatform: 'Nintendo',
//                 itemPlatformVersion: 'Gamecube',
//               },
//               {
//                 itemName: 'Metal Gear Solid 3: Snake Eater',
//                 itemLabel: 'metalgearsolid3snakeeater',
//                 itemPlatform: 'Playstation',
//                 itemPlatformVersion: 'PS3',
//               },
//               {
//                 itemName: 'Dragon Age: Origins',
//                 itemLabel: 'dragonageorigins',
//                 itemPlatform: 'Playstation',
//                 itemPlatformVersion: 'PS3',
//               },
//               {
//                 itemName: 'Shadow of the Colossus',
//                 itemLabel: 'shadowofthecolossus',
//                 itemPlatform: 'Playstation',
//                 itemPlatformVersion: 'PS3',
//               }
//             ],
//             locationType: SubLocationType.SHELF,
//           },
//         ],
//         items: [],
//       },
//       {
//         name: 'Sister\'s house',
//         label: 'sistershouse',
//         locationType: PhysicalLocationType.HOUSE,
//         mapCoordinates: '28°33\'01.1"N 81°29\'30.4"W',
//         sublocations: [
//           {
//             name: 'Living room media cabinet',
//             description: 'Living room media cabinet',
//             items: [
//               {
//                 itemName: 'Super Mario Party Jamboree',
//                 itemLabel: 'supermariopartyjamboree',
//                 itemPlatform: 'Nintendo',
//                 itemPlatformVersion: 'Switch',
//               },
//               {
//                 itemName: 'Luigi\'s Mansion 3',
//                 itemLabel: 'luigisman3',
//                 itemPlatform: 'Nintendo',
//                 itemPlatformVersion: 'Switch',
//               }
//             ],
//             locationType: SubLocationType.CONSOLE,
//           }
//         ],
//         items: [],
//       }
//     ],
//     digitalLocations: [
//       {
//         name: 'Steam',
//         label: GamePlatform.STEAM,
//         url: 'https://store.steampowered.com/',
//         isActive: true,
//         isFree: true,
//         monthlyFee: '0',
//         locationImage: '',
//         mapCoordinates: '',
//         items: [
//           {
//             itemName: 'Dark Souls 3',
//             itemLabel: 'darksouls3',
//             itemPlatform: 'Steam',
//             itemPlatformVersion: 'PC',
//           },
//           {
//             itemName: 'Dark Souls: REMASTERED',
//             itemLabel: 'darksoulsremastered',
//             itemPlatform: 'Steam',
//             itemPlatformVersion: 'PC',
//           },
//           {
//             itemName: 'ELDEN RING',
//             itemLabel: 'eldenring',
//             itemPlatform: 'Steam',
//             itemPlatformVersion: 'PC',
//           },
//           {
//             itemName: 'HELLDIVERS 2',
//             itemLabel: 'helldivders2',
//             itemPlatform: 'Steam',
//             itemPlatformVersion: 'PC',
//           },
//           {
//             itemName: 'Yakuza 0',
//             itemLabel: 'yakuzazero',
//             itemPlatform: 'Steam',
//             itemPlatformVersion: 'PC',
//           },
//           {
//             itemName: 'Yakuza Kiwami',
//             itemLabel: 'yakuzakiwami',
//             itemPlatform: 'Steam',
//             itemPlatformVersion: 'PC',
//           },
//           {
//             itemName: 'Bloodstained: Curse of the Moon',
//             itemLabel: 'bloodstainedcurseofthemoon',
//             itemPlatform: 'Steam',
//             itemPlatformVersion: 'PC',
//           },
//           {
//             itemName: 'Bloodstained: Curse of the Moon 2',
//             itemLabel: 'bloodstainedcurseofthemoon2',
//             itemPlatform: 'Steam',
//             itemPlatformVersion: 'PC',
//           },
//           {
//             itemName: 'Resident Evil 2',
//             itemLabel: 'residentevil2',
//             itemPlatform: 'Steam',
//             itemPlatformVersion: 'PC',
//           },
//           {
//             itemName: 'Resident Evil 3',
//             itemLabel: 'residentevil3',
//             itemPlatform: 'Steam',
//             itemPlatformVersion: 'PC',
//           },
//           {
//             itemName: 'Resident Evil 4',
//             itemLabel: 'residentevil4',
//             itemPlatform: 'Steam',
//             itemPlatformVersion: 'PC',
//           }
//         ],
//       },
//       {
//         name: 'Playstation Network',
//         label: GamePlatform.SONY,
//         url: 'https://www.playstation.com/en-us/playstation-plus/',
//         isActive: true,
//         isFree: false,
//         monthlyFee: '10',
//         locationImage: '',
//         mapCoordinates: '',
//         items: [
//           {
//             itemName: 'Bloodbourne',
//             itemLabel: 'bloodbourne',
//             itemPlatform: 'Playstation',
//             itemPlatformVersion: 'PS3',
//           },
//           {
//             itemName: 'Wanted: Dead',
//             itemLabel: 'wanteddead',
//             itemPlatform: 'Playstation',
//             itemPlatformVersion: 'PS5',
//           },
//           {
//             itemName: 'Returnal',
//             itemLabel: 'returnal',
//             itemPlatform: 'Playstation',
//             itemPlatformVersion: 'PS5',
//           },
//           {
//             itemName: 'Elden Ring',
//             itemLabel: 'eldenring',
//             itemPlatform: 'Playstation',
//             itemPlatformVersion: 'PS5',
//           },
//           {
//             itemName: "Hogwart's Legacy",
//             itemLabel: 'hogwartslegacy',
//             itemPlatform: 'Playstation',
//             itemPlatformVersion: 'PS5',
//           }
//         ],
//       },
//       {
//         name: 'Nintendo Switch Online',
//         label: GamePlatform.NINTENDO,
//         url: 'https://www.nintendo.com/switch/nintendo-switch-online/',
//         isActive: true,
//         isFree: false,
//         monthlyFee: '5',
//         locationImage: '',
//         mapCoordinates: '',
//         items: [
//           {
//             itemName: 'Limbo',
//             itemLabel: 'limbo',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'The Legend of Zelda: Breath of the Wild',
//             itemLabel: 'legendofzeldabreathofthewild',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'The Legend of Zelda: Tears of the Kingdom',
//             itemLabel: 'legendofzeldatearsofthekingdom',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Stardew Valley',
//             itemLabel: 'stardewvalley',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Pikmin 4',
//             itemLabel: 'pikmin4',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Capcom Beat\'em Up Bundle',
//             itemLabel: 'capcombeatemupbundle',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Arcade Archives: In the Hunt',
//             itemLabel: 'arcadearchivesinthehunt',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Manifold Garden',
//             itemLabel: 'manifoldgarden',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Okami HD',
//             itemLabel: 'okamihd',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Carto',
//             itemLabel: 'carto',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Overcooked! 2',
//             itemLabel: 'overcooked2',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Lovers in a Dangerous Spacetime',
//             itemLabel: 'loversinadangerousspacetime',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Mega Man Legacy Collection',
//             itemLabel: 'megamanlegacycollection',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Crawl',
//             itemLabel: 'crawl',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Tetris 99',
//             itemLabel: 'tetris99',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//           {
//             itemName: 'Arcade Archives: Puzzle Bobble',
//             itemLabel: 'arcadearchivespuzzlebobble',
//             itemPlatform: 'Nintendo',
//             itemPlatformVersion: 'Switch',
//           },
//         ],
//       }
//     ],
//   },
//   meta: {
//     counts: {
//       locations: {
//         total: 5,
//         physical: 2,
//         digital: 3,
//       },
//       items: {
//         total: 41,
//         physical: 9,
//         digital: 32,
//         byLocation: {
//           'condostudy': {
//             total: 7,
//             inSublocations: 7,
//           },
//           'sistershouse': {
//             total: 2,
//             inSublocations: 2,
//           },
//           'steam': {
//             total: 11,
//             inSublocations: 0,
//           },
//           'sony': {
//             total: 5,
//             inSublocations: 0,
//           },
//           'nintendo': {
//             total: 16,
//             inSublocations: 0,
//           },
//         },
//       },
//     },
//   },
// };
