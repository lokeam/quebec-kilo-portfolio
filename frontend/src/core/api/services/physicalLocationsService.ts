import { axiosInstance } from '@/core/api/client/axios-instance';
import { transformResponse } from '@/core/api/utils/transformResponse';
import type { SnakeCaseVersion } from '@/core/api/utils/apiTypes';

/**
 * Types for the frontend app using camelCase
 */
export interface PhysicalLocation {
  id: string;
  userId: string;
  name: string;
  label: string;
  locationType: string;
  mapCoordinates: string;
  createdAt: Date;
  updatedAt: Date;
  sublocations?: Sublocation[];
}

export interface Sublocation {
  id: string;
  name: string;
  description?: string;
  // Add other camelCase fields as needed
}

/**
 * Types for the API responses using snake_case (as returned by backend)
 * Generated automatically from the camelCase interfaces
 */
export type PhysicalLocationSnakeCase = SnakeCaseVersion<PhysicalLocation>;
export type SublocationSnakeCase = SnakeCaseVersion<Sublocation>;

/**
 * Get all physical locations
 *
 * This demonstrates all three approaches:
 *
 * 1. Legacy: Just use the original response and let components handle snake_case
 * 2. Manual: Explicitly transform the response
 * 3. Typed: Transform and cast to the camelCase type
 */

// 1. Legacy approach - keep backward compatibility
export async function getPhysicalLocationsLegacy() {
  return axiosInstance.get<PhysicalLocationSnakeCase[]>('/physical-locations');
}

// 2. Manual transformation - opt-in to camelCase
export async function getPhysicalLocations() {
  const response = await axiosInstance.get<PhysicalLocationSnakeCase[]>('/physical-locations');
  return transformResponse(response) as unknown as PhysicalLocation[];
}

// 3. Add enhanced functions that handle specific locations
export async function getPhysicalLocation(id: string): Promise<PhysicalLocation> {
  const response = await axiosInstance.get<PhysicalLocationSnakeCase>(`/physical-locations/${id}`);
  return transformResponse(response) as unknown as PhysicalLocation;
}

export async function createPhysicalLocation(
  locationData: Omit<PhysicalLocation, 'id' | 'createdAt' | 'updatedAt'>
): Promise<PhysicalLocation> {
  // Convert camelCase back to snake_case for the API request
  const snakeCaseData: Partial<PhysicalLocationSnakeCase> = {
    name: locationData.name,
    label: locationData.label,
    location_type: locationData.locationType,
    map_coordinates: locationData.mapCoordinates,
    // Add other fields as needed
  };

  const response = await axiosInstance.post<PhysicalLocationSnakeCase>('/physical-locations', snakeCaseData);
  return transformResponse(response) as unknown as PhysicalLocation;
}