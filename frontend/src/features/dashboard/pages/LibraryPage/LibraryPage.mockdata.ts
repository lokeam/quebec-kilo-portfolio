import type { LibraryGameItem } from '@/types/domain/library-types';

export const libraryPageMockData: LibraryGameItem[] = [
  {
    id: 718670,
    name: "Nier Automata - Game of the Year Edition",
    coverUrl: "https://i.imgur.com/dwswpNq.jpg",
    firstReleaseDate: 1487894400,
    rating: 4.5,
    isInLibrary: true,
    isInWishlist: false,
    gameType: {
      displayText: "Action RPG",
      normalizedText: "action-rpg"
    },
    favorite: true,
    gamesByPlatformAndLocation: [
      {
        gameId: 718670,
        platformId: 48,
        platformName: "PlayStation 4",
        type: "physical",
        locationId: "condo-study-1",
        locationName: "Condo Study",
        locationType: "apartment",
        sublocationId: "bookshelf-a-1",
        sublocationName: "Study bookshelf",
        sublocationType: "shelf",
        sublocationBgColor: "red"
      }
    ]
  },
  {
    id: 427520,
    name: "Factorio",
    coverUrl: "https://i.imgur.com/lDgAyT9.jpg",
    firstReleaseDate: 1451606400,
    rating: 4.8,
    isInLibrary: true,
    isInWishlist: false,
    gameType: {
      displayText: "Strategy",
      normalizedText: "strategy"
    },
    favorite: false,
    gamesByPlatformAndLocation: [
      {
        gameId: 427520,
        platformId: 1,
        platformName: "PC",
        type: "digital",
        locationId: "steam-1",
        locationName: "Steam",
        locationType: "digital"
      }
    ]
  },
  {
    id: 332330,
    name: "Don't Starve Together",
    coverUrl: "https://i.imgur.com/mnP42vi.jpg",
    firstReleaseDate: 1420070400,
    rating: 4.2,
    isInLibrary: true,
    isInWishlist: false,
    gameType: {
      displayText: "Survival",
      normalizedText: "survival"
    },
    favorite: false,
    gamesByPlatformAndLocation: [
      {
        gameId: 332330,
        platformId: 1,
        platformName: "PC",
        type: "digital",
        locationId: "steam-1",
        locationName: "Steam",
        locationType: "digital"
      }
    ]
  },
  {
    id: 646570,
    name: "Slay the Spire",
    coverUrl: "https://i.imgur.com/iSnOrlw.jpg",
    firstReleaseDate: 1514764800,
    rating: 4.7,
    isInLibrary: true,
    isInWishlist: false,
    gameType: {
      displayText: "Roguelike",
      normalizedText: "roguelike"
    },
    favorite: false,
    gamesByPlatformAndLocation: [
      {
        gameId: 646570,
        platformId: 1,
        platformName: "PC",
        type: "digital",
        locationId: "steam-1",
        locationName: "Steam",
        locationType: "digital"
      }
    ]
  },
  {
    id: 684410,
    name: "Bridge Constructor Portal",
    coverUrl: "https://i.imgur.com/3w5Q5PL.jpg",
    firstReleaseDate: 1514764800,
    rating: 4.0,
    isInLibrary: true,
    isInWishlist: false,
    gameType: {
      displayText: "Puzzle",
      normalizedText: "puzzle"
    },
    favorite: false,
    gamesByPlatformAndLocation: [
      {
        gameId: 684410,
        platformId: 1,
        platformName: "PC",
        type: "digital",
        locationId: "steam-1",
        locationName: "Steam",
        locationType: "digital"
      }
    ]
  },
  {
    id: 555150,
    name: "The First Tree",
    coverUrl: "https://i.imgur.com/F2aSEb4.jpg",
    firstReleaseDate: 1504224000,
    rating: 4.3,
    isInLibrary: true,
    isInWishlist: false,
    gameType: {
      displayText: "Adventure",
      normalizedText: "adventure"
    },
    favorite: false,
    gamesByPlatformAndLocation: [
      {
        gameId: 555150,
        platformId: 1,
        platformName: "PC",
        type: "digital",
        locationId: "steam-1",
        locationName: "Steam",
        locationType: "digital"
      }
    ]
  }
];
