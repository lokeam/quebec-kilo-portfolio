export interface Sublocation {
  id: string;
  name: string;
  type: string;
  parentLocationId: string;
  createdAt: Date;
  updatedAt: Date;
  [key: string]: unknown;
}