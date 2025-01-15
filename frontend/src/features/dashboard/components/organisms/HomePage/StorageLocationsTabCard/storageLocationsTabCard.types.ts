export type PhysicalStorageLocation = {
  name: string;
  itemsStored: string;
  avatar: string;
};

export type DigitalStorageService = {
  name: string;
  url: string;
  price: string;
  itemsStored: string;
  avatar: string;
};

export type StorageLocationsData = {
  totalDigitalLocations: string;
  totalPhysicalLocations: string;
  digitalStorage: DigitalStorageService[];
  physicalStorage: PhysicalStorageLocation[];
};
