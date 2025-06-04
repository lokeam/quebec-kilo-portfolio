/**
 * Physical Location Service
 *
 * Provides functions for managing physical locations and their sublocations through the backend API.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';
import type { PhysicalLocation, CreatePhysicalLocationRequest } from '@/types/domain/physical-location';
import type { Sublocation, CreateSublocationRequest } from '@/types/domain/sublocation';

// Response wrappers for physical locations
interface PhysicalLocationResponseWrapper {
  success: boolean;
  physical: PhysicalLocation;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

interface PhysicalLocationsResponseWrapper {
  success: boolean;
  physical: PhysicalLocation[];
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

interface SublocationResponseWrapper {
  success: boolean;
  sublocation: Sublocation;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

interface SublocationsResponseWrapper {
  success: boolean;
  sublocations: Sublocation[];
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

export interface DeletePhysicalLocationResponse {
  physical: {
    success: boolean;
    deleted_count: number;
    location_ids: string[];
  };
}

export interface DeleteSublocationResponse {
  sublocation: {
    success: boolean;
    deleted_count: number;
    sublocation_ids: string[];
  };
}

const PHYSICAL_LOCATION_ENDPOINT = '/v1/locations/physical';
const SUBLOCATION_ENDPOINT = '/v1/locations/sublocations';

// Physical Location Operations

/**
 * Fetches all physical locations for the current user
 */
export const getAllPhysicalLocations = (): Promise<PhysicalLocation[]> =>
  apiRequest('getAllPhysicalLocations', () =>
    axiosInstance
      .get<PhysicalLocationsResponseWrapper>(PHYSICAL_LOCATION_ENDPOINT)
      .then(response => response.data.physical)
  );

/**
 * Fetches a specific physical location by ID
 */
export const getSinglePhysicalLocation = (id: string): Promise<PhysicalLocation> =>
  apiRequest(`getPhysicalLocationById(${id})`, () =>
    axiosInstance
      .get<PhysicalLocationResponseWrapper>(`${PHYSICAL_LOCATION_ENDPOINT}/${id}`)
      .then(response => response.data.physical)
  );

/**
 * Creates a new physical location
 */
export const createPhysicalLocation = (input: CreatePhysicalLocationRequest): Promise<PhysicalLocation> =>
  apiRequest('createPhysicalLocation', () =>
    axiosInstance
      .post<PhysicalLocationResponseWrapper>(PHYSICAL_LOCATION_ENDPOINT, input)
      .then(response => response.data.physical)
  );

/**
 * Updates an existing physical location
 */
export const updatePhysicalLocation = (id: string, input: Partial<CreatePhysicalLocationRequest>): Promise<PhysicalLocation> =>
  apiRequest(`updatePhysicalLocation(${id})`, () =>
    axiosInstance
      .put<PhysicalLocationResponseWrapper>(`${PHYSICAL_LOCATION_ENDPOINT}/${id}`, input)
      .then(response => response.data.physical)
  );

/**
 * Deletes a physical location
 */
export const deletePhysicalLocation = (ids: string | string[]): Promise<DeletePhysicalLocationResponse['physical']> => {
  const idParam = Array.isArray(ids) ? ids.join(',') : ids;
  console.log('Making delete request for physical location(s):', idParam);

  return apiRequest(`deletePhysicalLocation(${idParam})`, () => {
    console.log('Executing delete request to:', `${PHYSICAL_LOCATION_ENDPOINT}?ids=${idParam}`);
    return axiosInstance
      .delete<DeletePhysicalLocationResponse>(`${PHYSICAL_LOCATION_ENDPOINT}?ids=${idParam}`)
      .then((response) => {
        console.log('Delete request successful');
        return response.data.physical;
      })
      .catch((error) => {
        console.error('Delete request failed:', error);
        throw error;
      });
  });
};

// Sublocation Operations

/**
 * Fetches all sublocations for a physical location
 */
export const getAllSublocations = (physicalLocationId: string): Promise<Sublocation[]> =>
  apiRequest(`getAllSublocations(${physicalLocationId})`, () =>
    axiosInstance
      .get<SublocationsResponseWrapper>(`${PHYSICAL_LOCATION_ENDPOINT}/${physicalLocationId}/sublocations`)
      .then(response => response.data.sublocations)
  );

/**
 * Fetches a specific sublocation by ID
 */
export const getSingleSublocation = (id: string): Promise<Sublocation> =>
  apiRequest(`getSublocationById(${id})`, () =>
    axiosInstance
      .get<SublocationResponseWrapper>(`${SUBLOCATION_ENDPOINT}/${id}`)
      .then(response => response.data.sublocation)
  );

/**
 * Creates a new sublocation
 */
export const createSublocation = (input: CreateSublocationRequest): Promise<Sublocation> =>
  apiRequest('createSublocation', () =>
    axiosInstance
      .post<SublocationResponseWrapper>(SUBLOCATION_ENDPOINT, input)
      .then(response => response.data.sublocation)
  );

/**
 * Updates an existing sublocation
 */
export const updateSublocation = (id: string, input: Partial<CreateSublocationRequest>): Promise<Sublocation> =>
  apiRequest(`updateSublocation(${id})`, () =>
    axiosInstance
      .put<SublocationResponseWrapper>(`${SUBLOCATION_ENDPOINT}/${id}`, input)
      .then(response => response.data.sublocation)
  );

/**
 * Deletes a sublocation
 */
export const deleteSublocation = (ids: string | string[]): Promise<DeleteSublocationResponse['sublocation']> => {
  const idParam = Array.isArray(ids) ? ids.join(',') : ids;
  console.log('Making delete request for sublocation(s):', idParam);

  return apiRequest(`deleteSublocation(${idParam})`, () => {
    console.log('Executing delete request to:', `${SUBLOCATION_ENDPOINT}?ids=${idParam}`);
    return axiosInstance
      .delete<DeleteSublocationResponse>(`${SUBLOCATION_ENDPOINT}?ids=${idParam}`)
      .then((response) => {
        console.log('Delete request successful');
        return response.data.sublocation;
      })
      .catch((error) => {
        console.error('Delete request failed:', error);
        throw error;
      });
  });
};