import type { LocationIconBgColor } from '@/types/domain/location-types';

const VALID_COLORS: LocationIconBgColor[] = [
  'red',
  'green',
  'blue',
  'orange',
  'gold',
  'purple',
  'brown',
  'gray',
  'pink'
];

/**
 * Type guard to validate if a string is a valid LocationIconBgColor
 * @param color - The color string to validate
 * @returns True if the color is valid, false otherwise
 */
export function isValidLocationBgColor(color: string | undefined): color is LocationIconBgColor {
  if (!color) return false;
  return VALID_COLORS.includes(color as LocationIconBgColor);
}