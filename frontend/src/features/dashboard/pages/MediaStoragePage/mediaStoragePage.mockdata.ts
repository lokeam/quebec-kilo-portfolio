//import type { MediaStorageData } from '@/features/dashboard/lib/types/service.types';
import { PhysicalLocationType, SubLocationType, GamePlatform } from '@/features/dashboard/types/media-storage.types';
import type { PhysicalLocation, DigitalLocation } from '@/features/dashboard/types/media-storage.types';

interface GameItem {
  itemName: string;
  itemLabel: string;
  itemPlatform: string;
  itemPlatformVersion: string;
}

export interface BaseLocation {
  name: string;
  label: string;
  subLocations?: SubLocation[];  // Use SubLocation instead of LocationCardData
}

export interface SubLocation {
  name: string;
  description: string;
  locationType: SubLocationType;
  items?: GameItem[];
}

interface MediaStorageMetadata {
  counts: {
    locations: {
      total: number;
      physical: number;
      digital: number;
    };
    items: {
      total: number;
      physical: number;
      digital: number;
      byLocation: Record<string, {
        total: number;
        inSubLocations: number;
      }>;
    };
  };
}

interface MediaStorageResponse {
  data: {
    name: string;
    label: string;
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
        name: 'Condo study',
        label: 'condostudy',
        locationType: PhysicalLocationType.APARTMENT,
        mapCoordinates: '40°44\'53.2"N 73°59\'05.6"W',
        subLocations: [
          {
            name: 'Bookshelf A',
            description: 'Bookshelf A',
            items: [
              {
                itemName: 'Gradius V',
                itemLabel: 'gradiusv',
                itemPlatform: 'Playstation',
                itemPlatformVersion: 'PS2',
              },
              {
                itemName: 'Street Fighter III: Third Strike',
                itemLabel: 'streetfighteriiithirdstrike',
                itemPlatform: 'Playstation',
                itemPlatformVersion: 'PS2',
              },
              {
                itemName: 'The Legend of Zelda: Wind Waker',
                itemLabel: 'legendofzeldawindwaker',
                itemPlatform: 'Nintendo',
                itemPlatformVersion: 'Gamecube',
              },
              {
                itemName: 'Metroid Prime',
                itemLabel: 'metroidprime',
                itemPlatform: 'Nintendo',
                itemPlatformVersion: 'Gamecube',
              },
              {
                itemName: 'Metal Gear Solid 3: Snake Eater',
                itemLabel: 'metalgearsolid3snakeeater',
                itemPlatform: 'Playstation',
                itemPlatformVersion: 'PS3',
              },
              {
                itemName: 'Dragon Age: Origins',
                itemLabel: 'dragonageorigins',
                itemPlatform: 'Playstation',
                itemPlatformVersion: 'PS3',
              },
              {
                itemName: 'Shadow of the Colossus',
                itemLabel: 'shadowofthecolossus',
                itemPlatform: 'Playstation',
                itemPlatformVersion: 'PS3',
              }
            ],
            locationType: SubLocationType.SHELF,
          },
        ],
        items: [],
      },
      {
        name: 'Sister\'s house',
        label: 'sistershouse',
        locationType: PhysicalLocationType.HOUSE,
        mapCoordinates: '28°33\'01.1"N 81°29\'30.4"W',
        subLocations: [
          {
            name: 'Living room media cabinet',
            description: 'Living room media cabinet',
            items: [
              {
                itemName: 'Super Mario Party Jamboree',
                itemLabel: 'supermariopartyjamboree',
                itemPlatform: 'Nintendo',
                itemPlatformVersion: 'Switch',
              },
              {
                itemName: 'Luigi\'s Mansion 3',
                itemLabel: 'luigisman3',
                itemPlatform: 'Nintendo',
                itemPlatformVersion: 'Switch',
              }
            ],
            locationType: SubLocationType.CONSOLE,
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
            itemName: 'Dark Souls 3',
            itemLabel: 'darksouls3',
            itemPlatform: 'Steam',
            itemPlatformVersion: 'PC',
          },
          {
            itemName: 'Dark Souls: REMASTERED',
            itemLabel: 'darksoulsremastered',
            itemPlatform: 'Steam',
            itemPlatformVersion: 'PC',
          },
          {
            itemName: 'ELDEN RING',
            itemLabel: 'eldenring',
            itemPlatform: 'Steam',
            itemPlatformVersion: 'PC',
          },
          {
            itemName: 'HELLDIVERS 2',
            itemLabel: 'helldivders2',
            itemPlatform: 'Steam',
            itemPlatformVersion: 'PC',
          },
          {
            itemName: 'Yakuza 0',
            itemLabel: 'yakuzazero',
            itemPlatform: 'Steam',
            itemPlatformVersion: 'PC',
          },
          {
            itemName: 'Yakuza Kiwami',
            itemLabel: 'yakuzakiwami',
            itemPlatform: 'Steam',
            itemPlatformVersion: 'PC',
          },
          {
            itemName: 'Bloodstained: Curse of the Moon',
            itemLabel: 'bloodstainedcurseofthemoon',
            itemPlatform: 'Steam',
            itemPlatformVersion: 'PC',
          },
          {
            itemName: 'Bloodstained: Curse of the Moon 2',
            itemLabel: 'bloodstainedcurseofthemoon2',
            itemPlatform: 'Steam',
            itemPlatformVersion: 'PC',
          },
          {
            itemName: 'Resident Evil 2',
            itemLabel: 'residentevil2',
            itemPlatform: 'Steam',
            itemPlatformVersion: 'PC',
          },
          {
            itemName: 'Resident Evil 3',
            itemLabel: 'residentevil3',
            itemPlatform: 'Steam',
            itemPlatformVersion: 'PC',
          },
          {
            itemName: 'Resident Evil 4',
            itemLabel: 'residentevil4',
            itemPlatform: 'Steam',
            itemPlatformVersion: 'PC',
          }
        ],
      },
      {
        name: 'Playstation Network',
        label: GamePlatform.SONY,
        url: 'https://www.playstation.com/en-us/playstation-plus/',
        isActive: true,
        isFree: false,
        monthlyFee: '10',
        locationImage: '',
        mapCoordinates: '',
        items: [
          {
            itemName: 'Bloodbourne',
            itemLabel: 'bloodbourne',
            itemPlatform: 'Playstation',
            itemPlatformVersion: 'PS3',
          },
          {
            itemName: 'Wanted: Dead',
            itemLabel: 'wanteddead',
            itemPlatform: 'Playstation',
            itemPlatformVersion: 'PS5',
          },
          {
            itemName: 'Returnal',
            itemLabel: 'returnal',
            itemPlatform: 'Playstation',
            itemPlatformVersion: 'PS5',
          },
          {
            itemName: 'Elden Ring',
            itemLabel: 'eldenring',
            itemPlatform: 'Playstation',
            itemPlatformVersion: 'PS5',
          },
          {
            itemName: "Hogwart's Legacy",
            itemLabel: 'hogwartslegacy',
            itemPlatform: 'Playstation',
            itemPlatformVersion: 'PS5',
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
            itemName: 'Limbo',
            itemLabel: 'limbo',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'The Legend of Zelda: Breath of the Wild',
            itemLabel: 'legendofzeldabreathofthewild',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'The Legend of Zelda: Tears of the Kingdom',
            itemLabel: 'legendofzeldatearsofthekingdom',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Stardew Valley',
            itemLabel: 'stardewvalley',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Pikmin 4',
            itemLabel: 'pikmin4',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Capcom Beat\'em Up Bundle',
            itemLabel: 'capcombeatemupbundle',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Arcade Archives: In the Hunt',
            itemLabel: 'arcadearchivesinthehunt',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Manifold Garden',
            itemLabel: 'manifoldgarden',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Okami HD',
            itemLabel: 'okamihd',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Carto',
            itemLabel: 'carto',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Overcooked! 2',
            itemLabel: 'overcooked2',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Lovers in a Dangerous Spacetime',
            itemLabel: 'loversinadangerousspacetime',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Mega Man Legacy Collection',
            itemLabel: 'megamanlegacycollection',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Crawl',
            itemLabel: 'crawl',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Tetris 99',
            itemLabel: 'tetris99',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
          },
          {
            itemName: 'Arcade Archives: Puzzle Bobble',
            itemLabel: 'arcadearchivespuzzlebobble',
            itemPlatform: 'Nintendo',
            itemPlatformVersion: 'Switch',
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
            inSubLocations: 7,
          },
          'sistershouse': {
            total: 2,
            inSubLocations: 2,
          },
          'steam': {
            total: 11,
            inSubLocations: 0,
          },
          'sony': {
            total: 5,
            inSubLocations: 0,
          },
          'nintendo': {
            total: 16,
            inSubLocations: 0,
          },
        },
      },
    },
  },
};
