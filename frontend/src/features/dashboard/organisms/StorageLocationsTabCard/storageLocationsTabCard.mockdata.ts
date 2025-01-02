export interface PhysicalStorageLocation {
  name: string;
  itemsStored: string;
  avatar: string;
}

export interface DigitalStorageService {
  name: string;
  url: string;
  price: string;
  itemsStored: string;
  avatar: string;
}

export interface StorageLocationsData {
  totalDigitalLocations: string;
  totalPhysicalLocations: string;
  digitalStorage: DigitalStorageService[];
  physicalStorage: PhysicalStorageLocation[];
}

const physicalStorageLocations: PhysicalStorageLocation[] = [
  {
    name: "Study bookshelf",
    itemsStored: "7",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "Study Midtower PC",
    itemsStored: "4",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "Living room media cabinet",
    itemsStored: "11",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "Office break room",
    itemsStored: "3",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "Public Storage",
    itemsStored: "9",
    avatar: "/placeholder.svg?height=36&width=36"
  },
];

const digitalStorageServices: DigitalStorageService[] = [
  {
    name: "Google Drive",
    url: "https://drive.google.com/drive/home",
    price: "FREE",
    itemsStored: "8",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "Playstation Network",
    url: "https://www.playstation.com/en-us/playstation-network/",
    price: "$5.99/month",
    itemsStored: "5",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "XBOX Game Pass",
    url: "https://www.xbox.com/en-US/xbox-game-pass",
    price: "$7.99/month",
    itemsStored: "12",
    avatar: "/placeholder.svg?height=36&width=36"
  },
  {
    name: "Switch Online",
    url: "https://accounts.nintendo.com/",
    price: "$8.99/month",
    itemsStored: "6",
    avatar: "/placeholder.svg?height=36&width=36"
  },
];

export const storageLocationsData: StorageLocationsData = {
  totalDigitalLocations: "4",
  totalPhysicalLocations: "5",
  digitalStorage: digitalStorageServices,
  physicalStorage: physicalStorageLocations
};