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
    physicalLocations: [
      {
        id: 'condostudy-1',
        name: 'Condo study',
        label: 'condostudy',
        type: 'physical',
        locationType: PhysicalLocationType.APARTMENT,
        mapCoordinates: '40°44\'53.2"N 73°59\'05.6"W',
        createdAt: new Date('2024-01-01'),
        updatedAt: new Date('2024-01-27'),
        sublocations: [
          {
            id: 'bookshelf-a-1',
            name: 'Bookshelf A',
            description: 'Bookshelf A',
            locationType: SublocationType.SHELF,
            items: [
              {
                id: 'gradius-v-1',
                imageUrl: 'https://i.imgur.com/example1.jpg',
                title: 'Gradius V',
                favorite: false,
                dateAdded: '2023-01-15',
                type: 'physical',
                platform: {
                  category: 'console',
                  model: 'PlayStation 2'
                },
                location: {
                  name: 'Condo study',
                  category: 'apartment',
                  subname: 'Bookshelf A',
                  sublocation: 'shelf'
                }
              },
              {
                id: 'streetfighter-iii-ts-1',
                imageUrl: 'https://i.imgur.com/example2.jpg',
                title: 'Street Fighter III: Third Strike',
                favorite: true,
                dateAdded: '2023-01-16',
                type: 'physical',
                platform: {
                  category: 'console',
                  model: 'PlayStation 2'
                },
                location: {
                  name: 'Condo study',
                  category: 'apartment',
                  subname: 'Bookshelf A',
                  sublocation: 'shelf'
                }
              },
              {
                id: 'zelda-wind-waker-1',
                imageUrl: 'https://i.imgur.com/example5.jpg',
                title: 'The Legend of Zelda: Wind Waker',
                favorite: false,
                dateAdded: '2023-01-17',
                type: 'physical',
                platform: {
                  category: 'console',
                  model: 'Gamecube'
                },
                location: {
                  name: 'Condo study',
                  category: 'apartment',
                  subname: 'Bookshelf A',
                  sublocation: 'shelf'
                }
              },
              {
                id: 'metroid-prime-1',
                imageUrl: 'https://i.imgur.com/example6.jpg',
                title: 'Metroid Prime',
                favorite: true,
                dateAdded: '2023-01-18',
                type: 'physical',
                platform: {
                  category: 'console',
                  model: 'Gamecube'
                },
                location: {
                  name: 'Condo study',
                  category: 'apartment',
                  subname: 'Bookshelf A',
                  sublocation: 'shelf'
                }
              },
              {
                id: 'metal-gear-solid-3-snake-eater-1',
                imageUrl: 'https://i.imgur.com/example7.jpg',
                title: 'Metal Gear Solid 3: Snake Eater',
                favorite: false,
                dateAdded: '2023-01-19',
                type: 'physical',
                platform: {
                  category: 'console',
                  model: 'PlayStation 3'
                },
                location: {
                  name: 'Condo study',
                  category: 'apartment',
                  subname: 'Bookshelf A',
                  sublocation: 'shelf'
                }
              },
              {
                id: 'dragon-age-origins-1',
                imageUrl: 'https://i.imgur.com/example8.jpg',
                title: 'Dragon Age: Origins',
                favorite: false,
                dateAdded: '2023-01-20',
                type: 'physical',
                platform: {
                  category: 'console',
                  model: 'PlayStation 3'
                },
                location: {
                  name: 'Condo study',
                  category: 'apartment',
                  subname: 'Bookshelf A',
                  sublocation: 'shelf'
                }
              },
              {
                id: 'shadow-of-the-colossus-1',
                imageUrl: 'https://i.imgur.com/example9.jpg',
                title: 'Shadow of the Colossus',
                favorite: true,
                dateAdded: '2023-01-21',
                type: 'physical',
                platform: {
                  category: 'console',
                  model: 'PlayStation 3'
                },
                location: {
                  name: 'Condo study',
                  category: 'apartment',
                  subname: 'Bookshelf A',
                  sublocation: 'shelf'
                }
              }
            ],
          },
        ]
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
            locationType: SublocationType.CONSOLE,
            items: [
              {
                id: 'super-mario-party-jamboree-1',
                imageUrl: 'https://i.imgur.com/example10.jpg',
                title: 'Super Mario Party Jamboree',
                favorite: false,
                dateAdded: '2023-01-22',
                type: 'physical',
                platform: {
                  category: 'console',
                  model: 'Nintendo Switch'
                },
                location: {
                  name: 'Sister\'s house',
                  category: 'house',
                  subname: 'Living room media cabinet',
                  sublocation: 'console'
                }
              },
              {
                id: 'luigis-mansion-3-1',
                imageUrl: 'https://i.imgur.com/example11.jpg',
                title: 'Luigi\'s Mansion 3',
                favorite: true,
                dateAdded: '2023-01-22',
                type: 'physical',
                platform: {
                  category: 'console',
                  model: 'Nintendo Switch'
                },
                location: {
                  name: 'Sister\'s house',
                  category: 'house',
                  subname: 'Living room media cabinet',
                  sublocation: 'console'
                }
              }
            ]
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
            id: 'elden-ring-1',
            imageUrl: 'https://i.imgur.com/example12.jpg',
            title: 'ELDEN RING',
            favorite: true,
            dateAdded: '2023-01-25',
            type: 'digital',
            platform: {
              category: 'pc',
              model: 'PC'
            },
            location: {
              service: 'steam',
              diskSize: {
                value: '60',
                unit: 'GB'
              }
            }
          },
          {
            id: 'helldivders-2-1',
            imageUrl: 'https://i.imgur.com/example13.jpg',
            title: 'HELLDIVERS 2',
            favorite: false,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'pc',
              model: 'PC'
            },
            location: {
              service: 'steam',
              diskSize: {
                value: '100',
                unit: 'GB'
              }
            }
          },
          {
            id: 'yakuza-0-1',
            imageUrl: 'https://i.imgur.com/example14.jpg',
            title: 'Yakuza 0',
            favorite: false,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'pc',
              model: 'PC'
            },
            location: {
              service: 'steam',
              diskSize: {
                value: '40',
                unit: 'GB'
              }
            }
          },
          {
            id: 'yakuza-kiwami-1',
            imageUrl: 'https://i.imgur.com/example15.jpg',
            title: 'Yakuza Kiwami',
            favorite: false,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'pc',
              model: 'PC'
            },
            location: {
              service: 'steam',
              diskSize: {
                value: '35',
                unit: 'GB'
              }
            }
          },
          {
            id: 'bloodstained-curse-of-the-moon-1',
            imageUrl: 'https://i.imgur.com/example16.jpg',
            title: 'Bloodstained: Curse of the Moon',
            favorite: false,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'pc',
              model: 'PC'
            },
            location: {
              service: 'steam',
              diskSize: {
                value: '2',
                unit: 'GB'
              }
            }
          },
          {
            id: 'bloodstained-curse-of-the-moon-2-1',
            imageUrl: 'https://i.imgur.com/example17.jpg',
            title: 'Bloodstained: Curse of the Moon 2',
            favorite: false,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'pc',
              model: 'PC'
            },
            location: {
              service: 'steam',
              diskSize: {
                value: '3',
                unit: 'GB'
              }
            }
          },
          {
            id: 'resident-evil-2-1',
            imageUrl: 'https://i.imgur.com/example18.jpg',
            title: 'Resident Evil 2',
            favorite: true,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'pc',
              model: 'PC'
            },
            location: {
              service: 'steam',
              diskSize: {
                value: '26',
                unit: 'GB'
              }
            }
          },
          {
            id: 'resident-evil-3-1',
            imageUrl: 'https://i.imgur.com/example19.jpg',
            title: 'Resident Evil 3',
            favorite: false,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'pc',
              model: 'PC'
            },
            location: {
              service: 'steam',
              diskSize: {
                value: '45',
                unit: 'GB'
              }
            }
          },
          {
            id: 'resident-evil-4-1',
            imageUrl: 'https://i.imgur.com/example20.jpg',
            title: 'Resident Evil 4',
            favorite: true,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'pc',
              model: 'PC'
            },
            location: {
              service: 'steam',
              diskSize: {
                value: '60',
                unit: 'GB'
              }
            }
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
            imageUrl: 'https://i.imgur.com/example21.jpg',
            title: 'Bloodbourne',
            favorite: true,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'console',
              model: 'PlayStation 4'
            },
            location: {
              service: 'playstation',
              diskSize: {
                value: '40',
                unit: 'GB'
              }
            }
          },
          {
            id: 'wanted-dead-1',
            imageUrl: 'https://i.imgur.com/example22.jpg',
            title: 'Wanted: Dead',
            favorite: false,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'console',
              model: 'PlayStation 5'
            },
            location: {
              service: 'playstation',
              diskSize: {
                value: '50',
                unit: 'GB'
              }
            }
          },
          {
            id: 'returnal-1',
            imageUrl: 'https://i.imgur.com/example23.jpg',
            title: 'Returnal',
            favorite: true,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'console',
              model: 'PlayStation 5'
            },
            location: {
              service: 'playstation',
              diskSize: {
                value: '60',
                unit: 'GB'
              }
            }
          },
          {
            id: 'elden-ring-1',
            imageUrl: 'https://i.imgur.com/example24.jpg',
            title: 'Elden Ring',
            favorite: true,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'console',
              model: 'PlayStation 5'
            },
            location: {
              service: 'playstation',
              diskSize: {
                value: '60',
                unit: 'GB'
              }
            }
          },
          {
            id: "hogwart's-legacy-1",
            imageUrl: 'https://i.imgur.com/example25.jpg',
            title: "Hogwart's Legacy",
            favorite: false,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'console',
              model: 'PlayStation 5'
            },
            location: {
              service: 'playstation',
              diskSize: {
                value: '85',
                unit: 'GB'
              }
            }
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
            imageUrl: 'https://i.imgur.com/example26.jpg',
            title: 'Limbo',
            favorite: false,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'console',
              model: 'Nintendo Switch'
            },
            location: {
              service: 'nintendo',
              diskSize: {
                value: '1',
                unit: 'GB'
              }
            }
          },
          {
            id: 'the-legend-of-zelda-breath-of-the-wild-1',
            imageUrl: 'https://i.imgur.com/example27.jpg',
            title: 'The Legend of Zelda: Breath of the Wild',
            favorite: true,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'console',
              model: 'Nintendo Switch'
            },
            location: {
              service: 'nintendo',
              diskSize: {
                value: '18',
                unit: 'GB'
              }
            }
          },
          {
            id: 'the-legend-of-zelda-tears-of-the-kingdom-1',
            imageUrl: 'https://i.imgur.com/example28.jpg',
            title: 'The Legend of Zelda: Tears of the Kingdom',
            favorite: true,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'console',
              model: 'Nintendo Switch'
            },
            location: {
              service: 'nintendo',
              diskSize: {
                value: '18',
                unit: 'GB'
              }
            }
          },
          {
            id: 'stardew-valley-1',
            imageUrl: 'https://i.imgur.com/example29.jpg',
            title: 'Stardew Valley',
            favorite: false,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'console',
              model: 'Nintendo Switch'
            },
            location: {
              service: 'nintendo',
              diskSize: {
                value: '1',
                unit: 'GB'
              }
            }
          },
          {
            id: 'pikmin-4-1',
            imageUrl: 'https://i.imgur.com/example30.jpg',
            title: 'Pikmin 4',
            favorite: true,
            dateAdded: '2023-01-26',
            type: 'digital',
            platform: {
              category: 'console',
              model: 'Nintendo Switch'
            },
            location: {
              service: 'nintendo',
              diskSize: {
                value: '10',
                unit: 'GB'
              }
            }
          }
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