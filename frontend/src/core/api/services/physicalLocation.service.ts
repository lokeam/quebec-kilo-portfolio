/**
 * Physical Location Service
 *
 * Provides functions for managing physical locations and their sublocations through the backend API.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';
import type { PhysicalLocation, CreatePhysicalLocationRequest, LocationsBFFResponse } from '@/types/domain/physical-location';
import type { Sublocation, CreateSublocationRequest } from '@/types/domain/sublocation';
import { logger } from '@/core/utils/logger/logger';

// Response wrappers for physical locations
interface PhysicalLocationResponseWrapper {
  success: boolean;
  physical: PhysicalLocation;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

interface PhysicalLocationsBFFResponseWrapper {
  success: boolean;
  physical: LocationsBFFResponse;
  metadata: { timestamp: string; request_id: string; };
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
const PHYSICAL_LOCATION_BFF_ENDPOINT = '/v1/locations/physical/bff';
const PHYSICAL_LOCATION_CRUD_ENDPOINT = '/v1/locations/physical';
const SUBLOCATION_CRUD_ENDPOINT = '/v1/locations/sublocations';

// Physical Location Operations

/**
 * Fetches all physical locations for the current user on physical location page
 */
export const getPhysicalLocationsBFFResponse = (): Promise<LocationsBFFResponse> =>
  apiRequest('getPhysicalLocationsBFFResponse', async () => {
    console.log('[DEBUG] getPhysicalLocationsBFFResponse: Making API request');
    const response = await axiosInstance.get<{ physical: LocationsBFFResponse }>(PHYSICAL_LOCATION_BFF_ENDPOINT);
    console.log('[DEBUG] getPhysicalLocationsBFFResponse: Raw API response:', response.data);

    if (!response.data.physical) {
      console.error('[DEBUG] getPhysicalLocationsBFFResponse: No physical data in response:', response.data);
      throw new Error('No physical data in response');
    }

    console.log('[DEBUG] getPhysicalLocationsBFFResponse: Successfully extracted physical data:', response.data.physical);
    return response.data.physical;
  });

/**
 * Fetches a specific physical location by ID
 */
export const getSinglePhysicalLocation = (id: string): Promise<PhysicalLocation> =>
  apiRequest(`getPhysicalLocationById(${id})`, () =>
    axiosInstance
      .get<PhysicalLocationResponseWrapper>(`${PHYSICAL_LOCATION_CRUD_ENDPOINT}/${id}`)
      .then(response => response.data.physical)
  );

/**
 * Creates a new physical location
 */
export const createPhysicalLocation = (input: CreatePhysicalLocationRequest): Promise<LocationsBFFResponse> =>
  apiRequest('createPhysicalLocation', () =>
    axiosInstance
      .post<PhysicalLocationsBFFResponseWrapper>(PHYSICAL_LOCATION_CRUD_ENDPOINT, input)
      .then(response => response.data.physical)
  );

/**
 * Updates an existing physical location
 */
export const updatePhysicalLocation = (id: string, input: Partial<CreatePhysicalLocationRequest>): Promise<PhysicalLocation> =>
  apiRequest(`updatePhysicalLocation(${id})`, () =>
    axiosInstance
      .put<PhysicalLocationResponseWrapper>(`${PHYSICAL_LOCATION_CRUD_ENDPOINT}/${id}`, input)
      .then(response => response.data.physical)
  );

/**
 * Deletes a physical location
 */
export const deletePhysicalLocation = (ids: string | string[]): Promise<DeletePhysicalLocationResponse['physical']> => {
  const idParam = Array.isArray(ids) ? ids.join(',') : ids;
  logger.debug('Making delete request for physical location(s):', { idParam });

  return apiRequest(`deletePhysicalLocation(${idParam})`, () => {
    logger.debug('Executing delete request to:', { url: `${PHYSICAL_LOCATION_CRUD_ENDPOINT}?ids=${idParam}` });
    return axiosInstance
      .delete<DeletePhysicalLocationResponse>(`${PHYSICAL_LOCATION_CRUD_ENDPOINT}?ids=${idParam}`)
      .then((response) => {
        logger.debug('Delete request successful', { response: response.data });

        if (!response.data.physical) {
          throw new Error('Invalid response: missing physical data');
        }

        if (!response.data.physical.success) {
          throw new Error('Delete operation was not successful');
        }

        return response.data.physical;
      })
      .catch((error) => {
        logger.error('Delete request failed:', error);
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
      .get<SublocationsResponseWrapper>(`${PHYSICAL_LOCATION_CRUD_ENDPOINT}/${physicalLocationId}/sublocations`)
      .then(response => response.data.sublocations)
  );

/**
 * Fetches a specific sublocation by ID
 */
export const getSingleSublocation = (id: string): Promise<Sublocation> =>
  apiRequest(`getSublocationById(${id})`, () =>
    axiosInstance
      .get<SublocationResponseWrapper>(`${SUBLOCATION_CRUD_ENDPOINT}/${id}`)
      .then(response => response.data.sublocation)
  );

/**
 * Creates a new sublocation
 */
export const createSublocation = (input: CreateSublocationRequest): Promise<Sublocation> =>
  apiRequest('createSublocation', () =>
    axiosInstance
      .post<SublocationResponseWrapper>(SUBLOCATION_CRUD_ENDPOINT, input)
      .then(response => response.data.sublocation)
  );

/**
 * Updates an existing sublocation
 */
export const updateSublocation = (id: string, input: Partial<CreateSublocationRequest>): Promise<Sublocation> =>
  apiRequest(`updateSublocation(${id})`, () =>
    axiosInstance
      .put<SublocationResponseWrapper>(`${SUBLOCATION_CRUD_ENDPOINT}/${id}`, input)
      .then(response => response.data.sublocation)
  );

/**
 * Deletes a sublocation
 */
export const deleteSublocation = (ids: string | string[]): Promise<DeleteSublocationResponse['sublocation']> => {
  const idParam = Array.isArray(ids) ? ids.join(',') : ids;
  console.log('Making delete request for sublocation(s):', idParam);

  return apiRequest(`deleteSublocation(${idParam})`, () => {
    console.log('Executing delete request to:', `${SUBLOCATION_CRUD_ENDPOINT}?ids=${idParam}`);
    return axiosInstance
      .delete<DeleteSublocationResponse>(`${SUBLOCATION_CRUD_ENDPOINT}?ids=${idParam}`)
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