export interface PlatformData {
  domain: 'games' | 'books' | 'movies'
  platform: string
  itemCount: number
};

export type PlatformItem = {
  platform: string;
  itemCount: number;
  fill: string;
};
