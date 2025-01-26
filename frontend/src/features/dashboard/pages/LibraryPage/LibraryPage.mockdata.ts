import type { LibraryItem } from "@/features/dashboard/lib/types/library/items";

export const libraryPageMockData: LibraryItem[] = [
  {
    id: "718670",
    imageUrl: "https://i.imgur.com/dwswpNq.jpg",
    title: "Nier Automata - Game of the Year Edition",
    favorite: true,
    dateAdded: "2024-01-01",
    type: "physical",
    platform: {
      category: "Console",
      model: "PlayStation 4",
    },
    location: {
      name: "Condo",
      category: "apartment",
      subname: "Study bookshelf",
      sublocation: "shelf",
    }
  },
  {
    id: "427520",
    imageUrl: "https://i.imgur.com/lDgAyT9.jpg",
    title: "Factorio",
    favorite: false,
    dateAdded: "2024-01-01",
    type: "digital",
    platform: {
      category: "PC",
    },
    location: {
      service: "steam",
      diskSize: {
        value: "5",
        unit: "GB",
      }
    }
  },
  {
    id: "332330",
    imageUrl: "https://i.imgur.com/mnP42vi.jpg",
    title: "Don't Starve Together",
    favorite: false,
    dateAdded: "2024-01-01",
    type: "digital",
    platform: {
      category: "PC",
    },
    location: {
      service: "steam",
      diskSize: {
        value: "685",
        unit: "MB",
      }
    }
  },
  {
    id: "646570",
    imageUrl: "https://i.imgur.com/iSnOrlw.jpg",
    title: "Slay the Spire",
    favorite: false,
    dateAdded: "2024-01-01",
    type: "digital",
    platform: {
      category: "PC",
    },
    location: {
      service: "steam",
      diskSize: {
        value: "1.2",
        unit: "GB",
      }
    }
  },
  {
    id: "684410",
    imageUrl: "https://i.imgur.com/3w5Q5PL.jpg",
    title: "Bridge Constructor Portal",
    favorite: false,
    dateAdded: "2024-01-01",
    type: "digital",
    platform: {
      category: "PC",
    },
    location: {
      service: "steam",
      diskSize: {
        value: "200",
        unit: "MB",
      }
    }
  },
  {
    id: "555150",
    imageUrl: "https://i.imgur.com/F2aSEb4.jpg",
    title: "The First Tree",
    favorite: false,
    dateAdded: "2024-01-01",
    type: "digital",
    platform: {
      category: "PC",
    },
    location: {
      service: "steam",
      diskSize: {
        value: "2",
        unit: "GB",
      }
    }
  },
];
