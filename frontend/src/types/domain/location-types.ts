/**
 * Primary location type discriminator.
 * Used to distinguish between physical and digital storage locations.
 */
export const LocationType = {
  PHYSICAL: 'physical',
  DIGITAL: 'digital'
} as const;

export type LocationType = typeof LocationType[keyof typeof LocationType];

/**
 * Types of physical storage locations.
 * Represents different real-world storage environments.
 */
export const PhysicalLocationType = {
  HOUSE: 'house',
  APARTMENT: 'apartment',
  OFFICE: 'office',
  WAREHOUSE: 'warehouse',
} as const;

export type PhysicalLocationType = typeof PhysicalLocationType[keyof typeof PhysicalLocationType];

/**
 * Types of storage subdivisions within physical locations.
 * Represents specific storage units or furniture.
 */
export const SublocationType = {
  shelf: 'shelf',
  console: 'console',
  cabinet: 'cabinet',
  closet: 'closet',
  drawer: 'drawer',
  box: 'box'
} as const;

export type SublocationType = typeof SublocationType[keyof typeof SublocationType];