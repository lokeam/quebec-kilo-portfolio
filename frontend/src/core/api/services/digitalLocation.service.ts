/**
 * Digital Location Service
 *
 * Provides functions for managing digital locations through the backend API.
 */

import { axiosInstance } from '@/core/api/client/axios-instance';
import { apiRequest } from '@/core/api/utils/apiRequest';
import type { DigitalLocation, CreateDigitalLocationRequest } from '@/types/domain/digital-location';

interface DigitalLocationResponseWrapper {
  success: boolean;
  digital: DigitalLocation;
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

interface DigitalLocationsResponseWrapper {
  success: boolean;
  digital: DigitalLocation[];
  metadata: {
    timestamp: string;
    request_id: string;
  };
}

const DIGITAL_LOCATION_ENDPOINT = '/v1/locations/digital';

/**
 * Fetches all digital locations for the current user
 */
export const getAllDigitalLocations = (): Promise<DigitalLocation[]> =>
  apiRequest('getAllLocations', () =>
    axiosInstance
      .get<DigitalLocationsResponseWrapper>(DIGITAL_LOCATION_ENDPOINT)
      .then(response => response.data.digital)
  );

/**
 * Fetches a specific digital location by ID
 */
export const getSingleDigitalLocation = (id: string): Promise<DigitalLocation> =>
  apiRequest(`getLocationById(${id})`, () =>
    axiosInstance
      .get<DigitalLocationResponseWrapper>(`${DIGITAL_LOCATION_ENDPOINT}/${id}`)
      .then(response => response.data.digital)
  );

/**
 * Creates a new digital location
 */
export const createDigitalLocation = (input: CreateDigitalLocationRequest): Promise<DigitalLocation> =>
  apiRequest('createLocation', () =>
    axiosInstance
      .post<DigitalLocationResponseWrapper>(DIGITAL_LOCATION_ENDPOINT, input)
      .then(response => response.data.digital)
  );

/**
 * Updates an existing digital location
 */
export const updateDigitalLocation = (id: string, input: Partial<CreateDigitalLocationRequest>): Promise<DigitalLocation> =>
  apiRequest(`updateLocation(${id})`, () =>
    axiosInstance
      .put<DigitalLocationResponseWrapper>(`${DIGITAL_LOCATION_ENDPOINT}/${id}`, input)
      .then(response => response.data.digital)
  );

/**
 * Deletes a digital location
 */
export const deleteDigitalLocation = (id: string): Promise<void> =>
  apiRequest(`deleteLocation(${id})`, () =>
    axiosInstance
      .delete<DigitalLocationResponseWrapper>(`${DIGITAL_LOCATION_ENDPOINT}/${id}`)
      .then(() => {})  // Return void
  );