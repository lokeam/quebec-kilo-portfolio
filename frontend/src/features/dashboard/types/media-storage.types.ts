
export interface GameItem {
  itemName: string;
  itemLabel: string;
  itemPlatform: string;
  itemPlatformVersion: string;
}

export enum GamePlatform {
  STEAM = 'steam',
  EPIC = 'epic',
  GOG = 'gog',
  PLAYSTATION = 'playstation',
  SONY = 'sony',
  XBOX = 'xbox',
  NINTENDO = 'nintendo'
}

export interface Location {
  id: string;
  type: 'physical' | 'digital';
  name: string;
  coordinates?: string;
  subLocations?: string[];
}

export interface LocationCardData {
  name: string;
  description: string;
  locationType: string;
  bgColor?: string;
  items?: GameItem[];
}

export interface BaseLocation {
  name: string;
  label: string;
  subLocations?: LocationCardData[];
}

export interface SubLocation {
  name: string;
  description: string;
  locationType: SubLocationType;
  items?: GameItem[];
}

export enum PhysicalLocationType {
  HOUSE = 'house',
  APARTMENT = 'apartment',
  OFFICE = 'office',
  WAREHOUSE = 'warehouse'
}

export enum SubLocationType {
  SHELF = 'shelf',
  CONSOLE = 'console',
  CABINET = 'cabinet',
  CLOSET = 'closet',
  DRAWER = 'drawer',
  BOX = 'box'
}

export interface PhysicalLocation extends BaseLocation {
  locationType: PhysicalLocationType;
  mapCoordinates?: string;
  items?: GameItem[];
}

export interface DigitalLocation extends BaseLocation {
  label: GamePlatform;
  url: string;
  isActive: boolean;
  isFree: boolean;
  monthlyFee?: string;
  locationImage?: string;
  mapCoordinates?: string;
  items?: GameItem[];
}

export interface MediaStorageMetadata {
  counts: {
    items: {
      physical: number;
      digital: number;
      byLocation: Record<string, { total: number }>;
    };
    locations: {
      physical: number;
      digital: number;
    };
  };
}